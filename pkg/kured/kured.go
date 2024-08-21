package kured

import (
	"encoding/json"
	"time"

	v1 "k8s.io/api/apps/v1"
)

// this can be extracted from "github.com/kubereboot/kured/pkg/daemonsetlock"
type lockAnnotationValue struct {
	NodeID   string        `json:"nodeID"`
	Metadata interface{}   `json:"metadata,omitempty"`
	Created  time.Time     `json:"created"`
	TTL      time.Duration `json:"TTL"`
}

// this can be extracted from "github.com/kubereboot/kured/pkg/daemonsetlock"
type multiLockAnnotationValue struct {
	MaxOwners       int                   `json:"maxOwners"`
	LockAnnotations []lockAnnotationValue `json:"locks"`
}

type SilenceNode struct {
	NodeID     string
	SilenceEnd time.Time
}

type TimeProvider func() time.Time

func ExtractNodeIDsFromAnnotation(ds *v1.DaemonSet, annotation string, silenceDuration time.Duration, nowProvider TimeProvider) ([]SilenceNode, error) {
	now := nowProvider()
	silencerArray := []SilenceNode{}

	if _, ok := ds.Annotations[annotation]; !ok {
		return silencerArray, nil
	}

	multiLock := &multiLockAnnotationValue{}
	err := json.Unmarshal([]byte(ds.Annotations[annotation]), multiLock)
	if err != nil {
		return nil, err
	}

	for _, lock := range multiLock.LockAnnotations {
		// TODO: silence just for silenceEnd.Sub(now) duration
		silenceEnd := lock.Created.Add(silenceDuration)
		if silenceEnd.After(now) {
			silencerArray = append(silencerArray, SilenceNode{
				NodeID:     lock.NodeID,
				SilenceEnd: silenceEnd,
			})
		}
	}

	if len(silencerArray) > 0 {
		return silencerArray, nil
	}

	singleLock := &lockAnnotationValue{}
	err = json.Unmarshal([]byte(ds.Annotations[annotation]), singleLock)
	if err != nil {
		return nil, err
	}

	if singleLock.NodeID != "manual" {
		// TODO: silence just for silenceEnd.Sub(now) duration
		silenceEnd := singleLock.Created.Add(silenceDuration)
		if silenceEnd.After(now) {
			silencerArray = append(silencerArray, SilenceNode{
				NodeID:     singleLock.NodeID,
				SilenceEnd: silenceEnd,
			})
		}
		return silencerArray, nil
	} else {
		return silencerArray, nil
	}
}
