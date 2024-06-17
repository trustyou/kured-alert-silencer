package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/prometheus/common/version"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/trustyou/kured-alert-silencer/pkg/kured"
	"github.com/trustyou/kured-alert-silencer/pkg/silence"

	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	// Command line flags
	dsName              string
	dsNamespace         string
	lockAnnotation      string
	logFormat           string
	logLevel            string
	alertmanagerURL     string
	silenceDuration     string
	silenceMatchersJSON string
	showVersion         bool
)

const (
	// KuredNodeLockAnnotation is the canonical string value for the kured node-lock annotation
	KuredNodeLockAnnotation string = "weave.works/kured-node-lock"
	EnvPrefix                      = "KURED_ALERT_SILENCER"
)

// flagToEnvVar converts command flag name to equivalent environment variable name
func flagToEnvVar(flag string) string {
	envVarSuffix := strings.ToUpper(strings.ReplaceAll(flag, "-", "_"))
	return fmt.Sprintf("%s_%s", EnvPrefix, envVarSuffix)
}

// bindFlags binds each cobra flag to its associated viper configuration (environment variable)
func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		// Environment variables can't have dashes in them, so bind them to their equivalent keys with underscores
		if strings.Contains(f.Name, "-") {
			v.BindEnv(f.Name, flagToEnvVar(f.Name))
		}

		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !f.Changed && v.IsSet(f.Name) {
			val := v.Get(f.Name)
			log.Infof("Binding %s command flag to environment variable: %s", f.Name, flagToEnvVar(f.Name))
			cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}

// bindViper initializes viper and binds command flags with environment variables
func bindViper(cmd *cobra.Command, args []string) error {
	v := viper.New()

	v.SetEnvPrefix(EnvPrefix)
	v.AutomaticEnv()
	bindFlags(cmd, v)

	return nil
}

// NewRootCommand construct the Cobra root command
func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:               "kured-alert-silencer",
		Short:             "An opinionated way of silencing alerts while Kured reboot k8s nodes.",
		PersistentPreRunE: bindViper,
		Run:               root,
	}

	rootCmd.PersistentFlags().StringVar(&dsNamespace, "ds-namespace", "kube-system",
		"namespace containing daemonset on which Kured place the lock")
	rootCmd.PersistentFlags().StringVar(&dsName, "ds-name", "kured",
		"name of daemonset on which to place lock")
	rootCmd.PersistentFlags().StringVar(&lockAnnotation, "lock-annotation", KuredNodeLockAnnotation,
		"annotation in which to record locking node")
	rootCmd.PersistentFlags().StringVar(&logFormat, "log-format", "text",
		"use text or json log format")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info",
		"debug, info, warn, error, fatal, or panic")
	rootCmd.PersistentFlags().StringVar(&alertmanagerURL, "alertmanager-url", "http://localhost:9093",
		"Alertmanager URL to silence alerts")
	rootCmd.PersistentFlags().StringVar(&silenceDuration, "silence-duration", "10m",
		"Silence duration for alerts in Go duration format (e.g. 10m, 1h, 2h30m)")
	rootCmd.PersistentFlags().StringVar(
		&silenceMatchersJSON,
		"silence-matchers-json",
		`[{"name": "instance", "value": "{{.NodeName}}", "isRegex": false}]`,
		`JSON string with format [{"name": "instance", "value": "{{.NodeName}}", "isRegex": false}, {"name": "alertname", "value": "node_reboot", "isRegex": false}]`)
	rootCmd.PersistentFlags().BoolVar(&showVersion, "version", false, "Show version and exit")
	return rootCmd
}

func main() {
	cmd := NewRootCommand()

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func root(cmd *cobra.Command, args []string) {
	if showVersion {
		fmt.Print(version.Print("kured-alert-silencer"))
		os.Exit(0)
	}

	level, err := log.ParseLevel(logLevel)
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(level)

	if logFormat == "json" {
		log.SetFormatter(&log.JSONFormatter{})
	}

	log.Info("Kured Alert Silencer starting")

	ctx := context.Background()
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Create a new clientset which includes our clientset scheme
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("Kured daemon set namespace: %s", dsNamespace)
	log.Infof("Kured daemon set name: %s", dsName)
	log.Infof("lock annotation: %s", lockAnnotation)
	log.Infof("silence duration: %s", silenceDuration)
	log.Infof("silence matchers JSON: %s", silenceMatchersJSON)

	log.Info("Watching DaemonSet")
	watcher, err := client.AppsV1().DaemonSets(dsNamespace).Watch(ctx, metav1.ListOptions{
		FieldSelector: fields.OneTermEqualSelector("metadata.name", dsName).String(),
	})
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)

	// Process events
	go func() {
		alertmanager, err := silence.NewAlertmanagerClient(alertmanagerURL)
		if err != nil {
			log.Fatal(err)
		}
		for event := range watcher.ResultChan() {
			switch event.Type {
			case watch.Added, watch.Modified:
				ds := event.Object.(*v1.DaemonSet)
				nodeIDs, err := kured.ExtractNodeIDsFromAnnotation(ds, lockAnnotation)
				if err != nil {
					log.Fatal(err)
				}

				for _, nodeID := range nodeIDs {
					log.Infof("Silencing alerts for node %s", nodeID)
					err = silence.SilenceAlerts(alertmanager, silenceMatchersJSON, nodeID, silenceDuration)
					if err != nil {
						log.Fatal(err)
					}
				}
			case watch.Deleted:
				log.Info("DaemonSet deleted")
			case watch.Error:
				log.Error("Error watching DaemonSet")
			}
		}
		wg.Done()
	}()

	wg.Wait()
}
