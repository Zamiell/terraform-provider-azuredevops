//go:build (all || permissions || resource_agent_pool_permissions) && (!exclude_permissions || !resource_agent_pool_permissions)
// +build all permissions resource_agent_pool_permissions
// +build !exclude_permissions !resource_agent_pool_permissions

package permissions

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestAgentPoolPermissions_CreateAgentPoolToken(t *testing.T) {
	d := getAgentPoolPermissionsResource(t, "5")
	token, err := createAgentPoolToken(d, nil)
	assert.NoError(t, err)
	assert.Equal(t, "AgentPools/5", token)

	d = getAgentPoolPermissionsResource(t, "abc")
	token, err = createAgentPoolToken(d, nil)
	assert.Empty(t, token)
	assert.Error(t, err)

	d = getAgentPoolPermissionsResource(t, "")
	token, err = createAgentPoolToken(d, nil)
	assert.Empty(t, token)
	assert.Error(t, err)
}

func getAgentPoolPermissionsResource(t *testing.T, agentPoolID string) *schema.ResourceData {
	d := schema.TestResourceDataRaw(t, ResourceAgentPoolPermissions().Schema, nil)
	if agentPoolID != "" {
		d.Set("agent_pool_id", agentPoolID)
	}
	return d
}
