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

func ExtractNodeIDsFromAnnotation(ds *v1.DaemonSet, annotation string) ([]string, error) {
    nodeIDs := []string{}

	if _, ok := ds.Annotations[annotation]; !ok {
		return []string{}, nil
	}

    multiLock := &multiLockAnnotationValue{}
    err := json.Unmarshal([]byte(ds.Annotations[annotation]), multiLock)
    if err != nil {
		return nil, err
    }

	for _, lock := range multiLock.LockAnnotations {
		nodeIDs = append(nodeIDs, lock.NodeID)
	}

	if len(nodeIDs) > 0 {
		return nodeIDs, nil
	}

    singleLock := &lockAnnotationValue{}
    err = json.Unmarshal([]byte(ds.Annotations[annotation]), singleLock)
    if err != nil {
		return nil, err
    }

	nodeIDs = append(nodeIDs, singleLock.NodeID)
	if nodeIDs[0] != "manual" {
		return nodeIDs, nil
	}
	return []string{}, nil
}
