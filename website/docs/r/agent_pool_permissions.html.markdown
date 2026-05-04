---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_agent_pool_permissions"
description: |-
  Manages permissions for an Azure DevOps agent pool
---

# azuredevops_agent_pool_permissions

Manages permissions for an Azure DevOps agent pool.

## Example Usage

```hcl
resource "azuredevops_agent_pool" "example" {
  name           = "Example-pool"
  auto_provision = false
  auto_update    = false
}

data "azuredevops_project" "example" {
  name = "Example Project"
}

data "azuredevops_group" "example-build-administrators" {
  project_id = data.azuredevops_project.example.id
  name       = "Build Administrators"
}

resource "azuredevops_agent_pool_permissions" "example" {
  agent_pool_id = azuredevops_agent_pool.example.id
  principal     = data.azuredevops_group.example-build-administrators.id
  permissions = {
    View                   = "allow"
    Use                    = "allow"
    Create                 = "allow"
    Manage                 = "allow"
    AdministerPermissions  = "allow"
  }
}
```

## Roles

The Azure DevOps UI uses roles to assign permissions for agent pools.

| Role          | Allowed Permissions       |
|---------------|---------------------------|
| Reader        | View                      |
| User          | View, Use, Create         |
| Creator       | View, Use, Create         |
| Administrator | All permissions           |

## Argument Reference

The following arguments are supported:

* `agent_pool_id` - (Required) The ID of the agent pool.

* `principal` - (Required) The **group** principal to assign the permissions.

* `permissions` - (Required) the permissions to assign. The following permissions are available.

  | Permission             | Description                         |
  |------------------------|-------------------------------------|
  | View                   | View agent pool resources           |
  | Manage                 | Manage agent pool resources         |
  | Listen                 | Listen for agent pool events        |
  | AdministerPermissions  | Administer agent pool permissions   |
  | Use                    | Use agent pool resources            |
  | Create                 | Create agent pool resources         |

---

* `replace` - (Optional) Replace (`true`) or merge (`false`) the permissions. Default: `true`

## Relevant Links

* [Azure DevOps Service REST API 7.1 - Security](https://learn.microsoft.com/en-us/rest/api/azure/devops/security/?view=azure-devops-rest-7.1)
* [Azure DevOps Security namespace and permission reference](https://learn.microsoft.com/en-us/azure/devops/organizations/security/namespace-reference?view=azure-devops#distributedtask)
* [Azure DevOps Agent pool security](https://learn.microsoft.com/en-us/azure/devops/pipelines/agents/pools-queues?view=azure-devops#security)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 minutes) Used when creating the Agent Pool Permission.
* `read` - (Defaults to 5 minute) Used when retrieving the Agent Pool Permission.
* `update` - (Defaults to 10 minutes) Used when updating the Agent Pool Permission.
* `delete` - (Defaults to 10 minutes) Used when deleting the Agent Pool Permission.

## Import

The resource does not support import.

## PAT Permissions Required

- **Security**: vso.security_manage - Grants the ability to read, write, and manage security permissions.
