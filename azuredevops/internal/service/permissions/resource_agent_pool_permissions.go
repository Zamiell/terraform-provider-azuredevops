package permissions

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	securityhelper "github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service/permissions/utils"
)

const agentPoolsSecurityToken = "AgentPools"

// ResourceAgentPoolPermissions schema and implementation for agent pool permission resource
func ResourceAgentPoolPermissions() *schema.Resource {
	return &schema.Resource{
		Create: resourceAgentPoolPermissionsCreateOrUpdate,
		Read:   resourceAgentPoolPermissionsRead,
		Update: resourceAgentPoolPermissionsCreateOrUpdate,
		Delete: resourceAgentPoolPermissionsDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: securityhelper.CreatePermissionResourceSchema(map[string]*schema.Schema{
			"agent_pool_id": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringIsNotWhiteSpace,
				Required:     true,
				ForceNew:     true,
			},
		}),
	}
}

func resourceAgentPoolPermissionsCreateOrUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	sn, err := securityhelper.NewSecurityNamespace(d, clients, securityhelper.SecurityNamespaceIDValues.DistributedTask, createAgentPoolToken)
	if err != nil {
		return err
	}

	if err := securityhelper.SetPrincipalPermissions(d, sn, nil, false); err != nil {
		return err
	}

	return resourceAgentPoolPermissionsRead(d, m)
}

func resourceAgentPoolPermissionsRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	sn, err := securityhelper.NewSecurityNamespace(d, clients, securityhelper.SecurityNamespaceIDValues.DistributedTask, createAgentPoolToken)
	if err != nil {
		return err
	}

	principalPermissions, err := securityhelper.GetPrincipalPermissions(d, sn)
	if err != nil {
		return err
	}
	if principalPermissions == nil {
		d.SetId("")
		log.Printf("[INFO] Permissions for ACL token %q not found. Removing from state", sn.GetToken())
		return nil
	}

	d.Set("permissions", principalPermissions.Permissions)
	return nil
}

func resourceAgentPoolPermissionsDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	sn, err := securityhelper.NewSecurityNamespace(d, clients, securityhelper.SecurityNamespaceIDValues.DistributedTask, createAgentPoolToken)
	if err != nil {
		return err
	}

	if err := securityhelper.SetPrincipalPermissions(d, sn, &securityhelper.PermissionTypeValues.NotSet, true); err != nil {
		return err
	}
	return nil
}

func createAgentPoolToken(d *schema.ResourceData, clients *client.AggregatedClient) (string, error) {
	agentPoolID, ok := d.GetOk("agent_pool_id")
	if !ok {
		return "", fmt.Errorf("Failed to get 'agent_pool_id' from schema")
	}

	poolID, err := strconv.Atoi(agentPoolID.(string))
	if err != nil {
		return "", fmt.Errorf("parsing agent_pool_id: %w", err)
	}

	return fmt.Sprintf("%s/%d", agentPoolsSecurityToken, poolID), nil
}
