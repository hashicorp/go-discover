#!/usr/bin/env bash
set -e

public_ip=$(terraform output public_ip)
vm_scale_set=$(terraform output vm_scale_set)
resource_group=$(terraform output resource_group)

official_ips=$(az vmss nic list --resource-group "${resource_group}" --vmss-name "${vm_scale_set}" \
               | jq -r '[.[].ipConfigurations[].privateIpAddress] | join(" ")')

. cli-outputs/discover.env

ips=$(ssh -q ubuntu@${public_ip} \
    -i tf_rsa \
    -p 50000 \
    -o UserKnownHostsFile=/dev/null \
    -o StrictHostKeyChecking=no \
    ./discover -q addrs \
        provider=azure \
        tenant_id=${ARM_TENANT_ID} \
        subscription_id=${ARM_SUBSCRIPTION_ID} \
        client_id=${ARM_CLIENT_ID} \
        secret_access_key=${ARM_CLIENT_SECRET} \
        resource_group=${resource_group} \
        vm_scale_set=${vm_scale_set})

if [ "${ips}" != "${official_ips}" ] ; then
    echo "got ${ips} on ${public_ip} want ${official_ips}"
    exit 1
fi

echo "OK"
exit 0
