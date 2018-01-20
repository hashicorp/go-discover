#!/usr/bin/env bash
set -e
. ./utils.bash

mkdir -p cli-outputs

# Test if ARM_TENANT_ID and ARM_SUBSCRIPTION_ID are set,
# otherwise try and find sane defaults.
if [ -z "${ARM_SUBSCRIPTION_ID}" -o -z "${ARM_TENANT_ID}" ] ; then
    echo "--- ARM_{TENANT_ID,SUBSCRIPTION_ID} not found"
    echo "--- Going to try and auto configure"
    accounts=$(az account list)

    export ARM_SUBSCRIPTION_ID=$(echo ${accounts} | jq -r .[0].id)
    if [ -z "${ARM_SUBSCRIPTION_ID}" ]; then
        echo "!!! Could not figure out subscription id."
        echo "!!! Make sure you're logged into an Aszure account."
    fi
    echo "--- Using subscription: ${ARM_SUBSCRIPTION_ID}"

    # required by envBuilder
    export ARM_TENANT_ID=$(echo ${accounts} | jq -r .[0].tenantId)
    if [ -z "${ARM_TENANT_ID}" ]; then
        echo "!!! Could not figure out tenant id."
        echo "!!! Make sure you're logged into an Azure account."
    fi
    echo "--- Using tenant id: ${ARM_TENANT_ID}"
fi

# Create discover role & credentials
if [ ! -f cli-outputs/role.json ]; then
    envsubst < ./fixtures/role.json.tmpl > cli-outputs/role.json
    echo "creating role"
    az role definition create --role-definition @cli-outputs/role.json > cli-outputs/roleDefinition.json
    echo "sleeping for 30 seconds to let custom role become eventually consistent"
    sleep 30
fi
if [ ! -f cli-outputs/discoverAuth.json ]; then
    echo "creating discover-test-discover user"
    az ad sp create-for-rbac --name discover-test-discover --role=$(roleName) --scopes="/subscriptions/${ARM_SUBSCRIPTION_ID}" > cli-outputs/discoverAuth.json
    envBuilder "discover"
fi

echo "Done setting up vmss account."
