/*
Copyright 2020 The Rook Authors. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package pool to manage a rook pool.
package pool

import (
	"testing"

	cephv1 "github.com/rook/rook/pkg/apis/ceph.rook.io/v1"
	opcontroller "github.com/rook/rook/pkg/operator/ceph/controller"
	"github.com/stretchr/testify/assert"
)

func TestToCustomResourceStatus(t *testing.T) {
	mirroringStatus := &cephv1.PoolRadosNamespaceMirroringStatusSummarySpec{}
	mirroringStatus.Health = "HEALTH_OK"
	mirroringInfo := &cephv1.PoolRadosNamespaceMirroringInfo{
		Mode:     "pool",
		SiteName: "rook-ceph-emea",
		Peers: []cephv1.PeersSpec{
			{UUID: "82656994-3314-4996-ac4c-263c2c9fd081"},
		},
	}

	// Test 1: Empty so it's disabled
	{
		newMirroringStatus, newMirroringInfo, _ := toCustomResourceStatus(&cephv1.MirroringStatusSpec{}, mirroringStatus, &cephv1.MirroringInfoSpec{}, mirroringInfo, &cephv1.SnapshotScheduleStatusSpec{}, []cephv1.SnapshotSchedulesSpec{}, "")
		assert.NotEmpty(t, newMirroringStatus.PoolRadosNamespaceMirroringStatus)
		assert.Equal(t, "HEALTH_OK", newMirroringStatus.PoolRadosNamespaceMirroringStatus.Health)
		assert.Equal(t, "pool", newMirroringInfo.Mode)
	}

	// Test 2: snap sched
	{
		snapSchedStatus := []cephv1.SnapshotSchedulesSpec{
			{
				Pool:  "my-pool",
				Image: "pool/image",
			},
		}
		newMirroringStatus, newMirroringInfo, newSnapshotScheduleStatus := toCustomResourceStatus(&cephv1.MirroringStatusSpec{}, mirroringStatus, &cephv1.MirroringInfoSpec{}, mirroringInfo, &cephv1.SnapshotScheduleStatusSpec{}, snapSchedStatus, "")
		assert.NotEmpty(t, newMirroringStatus.PoolRadosNamespaceMirroringStatus)
		assert.Equal(t, "HEALTH_OK", newMirroringStatus.PoolRadosNamespaceMirroringStatus.Health)
		assert.NotEmpty(t, newMirroringInfo.Mode, "pool")
		assert.NotEmpty(t, newSnapshotScheduleStatus)
	}
}

func TestUpdateStatusInfo(t *testing.T) {

	cephBlockPoolReplicated := &cephv1.CephBlockPool{
		Spec: cephv1.NamedBlockPoolSpec{
			Name: "test-pool-replicated",
			PoolSpec: cephv1.PoolSpec{
				Replicated: cephv1.ReplicatedSpec{
					Size: 3,
				},
			},
		},
		Status: &cephv1.CephBlockPoolStatus{
			Phase: cephv1.ConditionProgressing,
		},
	}
	updateStatusInfo(cephBlockPoolReplicated)
	statusInfo := cephBlockPoolReplicated.Status.Info
	assert.Equal(t, "Replicated", statusInfo["type"])
	assert.Equal(t, cephv1.DefaultFailureDomain, statusInfo["failureDomain"])
	assert.Empty(t, statusInfo[opcontroller.RBDMirrorBootstrapPeerSecretName])

	cephBlockPoolErasureCoded := &cephv1.CephBlockPool{
		Spec: cephv1.NamedBlockPoolSpec{
			Name: "test-pool-erasure-coded",
			PoolSpec: cephv1.PoolSpec{
				FailureDomain: "osd",
				ErasureCoded: cephv1.ErasureCodedSpec{
					CodingChunks: 6,
					DataChunks:   2,
				},
			},
		},
		Status: &cephv1.CephBlockPoolStatus{
			Phase: cephv1.ConditionProgressing,
		},
	}
	updateStatusInfo(cephBlockPoolErasureCoded)
	statusInfo = cephBlockPoolErasureCoded.Status.Info
	assert.Equal(t, "Erasure Coded", statusInfo["type"])
	assert.Equal(t, "osd", statusInfo["failureDomain"])
	assert.Empty(t, statusInfo[opcontroller.RBDMirrorBootstrapPeerSecretName])

	cephBlockPoolErasureCoded.Spec.PoolSpec.Mirroring = cephv1.MirroringSpec{
		Enabled: true,
	}

	updateStatusInfo(cephBlockPoolErasureCoded)
	statusInfo = cephBlockPoolErasureCoded.Status.Info
	assert.Empty(t, statusInfo[opcontroller.RBDMirrorBootstrapPeerSecretName])

	cephBlockPoolErasureCoded.Status.Phase = cephv1.ConditionReady
	updateStatusInfo(cephBlockPoolErasureCoded)
	statusInfo = cephBlockPoolErasureCoded.Status.Info
	assert.NotEmpty(t, statusInfo[opcontroller.RBDMirrorBootstrapPeerSecretName])
}
