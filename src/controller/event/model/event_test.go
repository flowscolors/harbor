// Copyright Project Harbor Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/goharbor/harbor/src/pkg/auditext/model"
)

func TestCommonEventResolveToAuditLog(t *testing.T) {
	t.Run("with source IP", func(t *testing.T) {
		event := &CommonEvent{
			Operator:             "testuser",
			ProjectID:            1,
			OcurrAt:              time.Now(),
			Operation:            "pull",
			Payload:              `{"test": "data"}`,
			SourceIP:             "192.168.1.100",
			ResourceType:         "artifact",
			ResourceName:         "library/nginx:latest",
			OperationDescription: "Pull artifact",
			IsSuccessful:         true,
		}

		auditLog, err := event.ResolveToAuditLog()
		require.NoError(t, err)
		require.NotNil(t, auditLog)

		assert.Equal(t, int64(1), auditLog.ProjectID)
		assert.Equal(t, "testuser", auditLog.Username)
		assert.Equal(t, "pull", auditLog.Operation)
		assert.Equal(t, "artifact", auditLog.ResourceType)
		assert.Equal(t, "library/nginx:latest", auditLog.Resource)
		assert.Equal(t, "Pull artifact", auditLog.OperationDescription)
		assert.True(t, auditLog.IsSuccessful)
		assert.Equal(t, "192.168.1.100", auditLog.SourceIP)
	})

	t.Run("without source IP", func(t *testing.T) {
		event := &CommonEvent{
			Operator:      "admin",
			ProjectID:     2,
			OcurrAt:       time.Now(),
			Operation:     "push",
			SourceIP:      "",
			ResourceType:  "artifact",
			ResourceName:  "library/alpine:v3.15",
			IsSuccessful:  true,
		}

		auditLog, err := event.ResolveToAuditLog()
		require.NoError(t, err)
		require.NotNil(t, auditLog)

		assert.Equal(t, int64(2), auditLog.ProjectID)
		assert.Equal(t, "admin", auditLog.Username)
		assert.Equal(t, "push", auditLog.Operation)
		assert.Equal(t, "192.168.1.100", auditLog.SourceIP) // Should match the event
		assert.Equal(t, "", auditLog.SourceIP)              // Actually this is empty in the event
		assert.Equal(t, "", auditLog.SourceIP)              // Now it should be empty
	})

	t.Run("IPv6 address", func(t *testing.T) {
		event := &CommonEvent{
			Operator:     "user1",
			ProjectID:    3,
			OcurrAt:      time.Now(),
			Operation:    "create",
			SourceIP:     "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
			ResourceType: "artifact",
			ResourceName: "test/repo:v1",
		}

		auditLog, err := event.ResolveToAuditLog()
		require.NoError(t, err)
		require.NotNil(t, auditLog)

		assert.Equal(t, "2001:0db8:85a3:0000:0000:8a2e:0370:7334", auditLog.SourceIP)
	})

	t.Run("failure operation", func(t *testing.T) {
		event := &CommonEvent{
			Operator:     "user2",
			ProjectID:    4,
			OcurrAt:      time.Now(),
			Operation:    "delete",
			SourceIP:     "10.0.0.50",
			ResourceType: "artifact",
			ResourceName: "project/repo:tag",
			IsSuccessful: false,
		}

		auditLog, err := event.ResolveToAuditLog()
		require.NoError(t, err)
		require.NotNil(t, auditLog)

		assert.False(t, auditLog.IsSuccessful)
		assert.Equal(t, "10.0.0.50", auditLog.SourceIP)
	})
}

func TestAuditLogExtModel(t *testing.T) {
	t.Run("table name", func(t *testing.T) {
		auditLog := &model.AuditLogExt{}
		assert.Equal(t, "audit_log_ext", auditLog.TableName())
	})

	t.Run("source IP field", func(t *testing.T) {
		auditLog := &model.AuditLogExt{
			ID:                   1,
			ProjectID:            1,
			Operation:            "pull",
			OperationDescription: "Test operation",
			IsSuccessful:         true,
			ResourceType:         "artifact",
			Resource:             "test/image:latest",
			Username:             "testuser",
			SourceIP:             "192.168.1.100",
		}

		assert.Equal(t, "192.168.1.100", auditLog.SourceIP)
	})
}
