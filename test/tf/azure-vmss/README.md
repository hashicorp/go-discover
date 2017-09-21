# Azure Tests For Virtual Machine Scale Sets

## Prerequisites
- [Azure CLI](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli)
- [jq](https://stedolan.github.io/jq/download/)
- [envsubst](https://www.gnu.org/software/gettext/)
- Get an Azure account
    - `az login` to the account

## Running the tests

- `make build`
    - Should create a role with minimum read-only action
    - Should create a service principal with that role
    - Should create the following resources:
        - network
        - public load balancer with nat forwarding rules
        - virtual machine scale set
        - copy the discover binary to a host
- `make test`
    - Should test that the discover binary, with minimal credentials, generates a list of IPs that are associated with the virtual machine scale set.
- `make destroy`
    - Destroys all the azure resources

## Known issues

The Azure Roles API is quite frequently inconsistent when it comes to deleting the role.

ex.

```
08:52 $ az role definition list --custom-role-only
[]
08:52 $ az role definition list --custom-role-only
[
  {
    "name": "6b8b3e02-82e4-4d9d-86bf-b54787b9a274",
    ...
    "type": "Microsoft.Authorization/roleDefinitions"
  }
]
08:53 $ az role definition list --custom-role-only
[]
08:53 $ az role definition list --custom-role-only
[]
08:53 $ az role definition list --custom-role-only
[]
08:53 $ az role definition list --custom-role-only
[
  {
    "name": "6b8b3e02-82e4-4d9d-86bf-b54787b9a274",
    ...
    "type": "Microsoft.Authorization/roleDefinitions"
  }
]
```

This means if you're developing on this test you will quite often need to run:

```bash
az role definition delete --name <role id>
az role definition list --custom-role-only
```

many, many times until the API is consistent about the state of the role.
