#!/usr/bin/env bash
set -e
. ./utils.bash

if [ -f cli-outputs/discoverAuth.json ]; then
    armclientid=$(clientID discoverAuth)
    echo "Attempting to remove service principal: ${armclientid}"
    az ad sp delete --id ${armclientid}
    if [ $? -eq 0 ]; then
        rm cli-outputs/discover.env
        rm cli-outputs/discoverAuth.json
        echo "Successfully removed  service principal"
    fi
fi

if [ -f cli-outputs/roleDefinition.json ]; then
    roleid=$(roleName)
    echo "Attempting to remove role: ${roleid}"
    az role definition delete --name ${roleid}
    if [ $? -eq 0 ]; then
        rm cli-outputs/roleDefinition.json
        rm cli-outputs/role.json
        echo "Successfully removed role"
    fi
fi
