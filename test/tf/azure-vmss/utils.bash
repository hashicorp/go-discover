#!/usr/bin/env bash
envBuilder() {
    NAME=${1}
    export ARM_CLIENT_ID=$(jq -r .appId < cli-outputs/${NAME}Auth.json)
    export ARM_CLIENT_SECRET=$(jq -r .password < cli-outputs/${NAME}Auth.json)
    envsubst > cli-outputs/${NAME}.env < ./fixtures/account.env.tmpl
    echo "Built ${NAME}.env"
}

roleName() {
    if [ -f cli-outputs/roleDefinition.json ]; then
        jq -r .name < cli-outputs/roleDefinition.json
    fi
}

clientID() {
    CLIENT=${1}
    if [ -f cli-outputs/${CLIENT}.json ]; then
        jq -r .appId < cli-outputs/${CLIENT}.json
    fi
}
