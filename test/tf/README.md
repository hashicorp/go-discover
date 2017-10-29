# Terraform Setups

This folder contains Terraform setups for the supported cloud providers.

### Amazon AWS


### Azure

#### Credentials

https://www.terraform.io/docs/providers/azurerm/index.html#creating-credentials

```shell
# install Azure CLI (https://github.com/Azure/azure-cli)
curl -L https://aka.ms/InstallAzureCli | bash

# 1. Login
$ az login

# 2. Get SubscriptionID
$ az account list
[
  {
    "cloudName": "AzureCloud",
    "id": "subscription_id",
    "isDefault": true,
    "name": "Gratis versie",
    "state": "Enabled",
    "tenantId": "tenant_id",
    "user": {
      "name": "user@email.com",
      "type": "user"
    }
  }
]


# 3. Switch to subscription
$ az account set --subscription="subscription_id"

# 4. Create ClientID and Secret
$ az ad sp create-for-rbac --role="Contributor" --scopes="/subscriptions/subscription_id"
{
  "appId": "client_id",
  "displayName": "azure-cli-2017-07-18-16-51-43",
  "name": "http://azure-cli-2017-07-18-16-51-43",
  "password": "client_secret",
  "tenant": "tenant_id"
}

# 5. Export the Credentials for the client
export ARM_CLIENT_ID=client_id
export ARM_CLIENT_SECRET=client_secret
export ARM_TENANT_ID=tenant_id
export ARM_SUBSCRIPTION_ID=subscription_id

# 6. Test the credentials
$ az vm list-sizes
```

#### Run the test

```shell

# create environment
make build

# run test (should return OK)
make test

# cleanup
make destroy

```

### Google Cloud


### Scaleway

todo

### Softlayer

tbd
