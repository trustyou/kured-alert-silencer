package kured_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trustyou/kured-alert-silencer/pkg/kured"

	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestExtractNodeIDsFromAnnotation(t *testing.T) {
    const KuredNodeLockAnnotation string = "weave.works/kured-node-lock"

    tests := []struct {
        name     string
        annotationValue string
        want     []string
    }{
        {
            name: "single node ID",
            annotationValue: `{"nodeID":"kind-control-plane2","metadata":{"unschedulable":false},"created":"2024-05-31T06:31:37.441623522Z","TTL":0}`,
            want: []string{"kind-control-plane2"},
        },
        {
            name: "multiple node IDs",
            annotationValue: `{"maxOwners":2,"locks":[{"nodeID":"kind-worker2","metadata":{"unschedulable":false},"created":"2024-05-30T11:41:32.735905893Z","TTL":0},{"nodeID":"kind-control-plane","metadata":{"unschedulable":false},"created":"2024-05-30T11:41:49.868231413Z","TTL":0}]}`,
            want: []string{"kind-worker2", "kind-control-plane"},
        },
        {
            name: "manual lock",
            annotationValue: `{"nodeID":"manual"}`,
            want: []string{},
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            ds := &v1.DaemonSet{
                ObjectMeta: metav1.ObjectMeta{
                    Annotations: map[string]string{
                        KuredNodeLockAnnotation: tt.annotationValue,
                    },
                },
            }

            result, err := kured.ExtractNodeIDsFromAnnotation(ds, KuredNodeLockAnnotation)
			require.NoError(t, err)
			require.Equal(t, result, tt.want)

        })
    }
}

func TestExtractNodeIDsFromAnnotationNoAnnotation(t *testing.T) {
	ds := &v1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{},
		},
	}

	result, err := kured.ExtractNodeIDsFromAnnotation(ds, "weave.works/kured-node-lock")
	require.NoError(t, err)
	require.Equal(t, result, []string{})

}
