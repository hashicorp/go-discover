#!/usr/bin/env bash

public_ips=$(terraform output public_ips | tr ',' ' ')
private_ips=$(terraform output private_ips | tr ',' ' ')
tagged_ips=$(terraform output tagged_ips | tr ',' ' ')

for h in $public_ips; do
	ip=$(ssh -q ubuntu@$h \
		-i tf_rsa \
		-o UserKnownHostsFile=/dev/null \
		-o StrictHostKeyChecking=no \
		./discover -q addrs \
			provider=azure \
			tenant_id=$ARM_TENANT_ID \
			subscription_id=$ARM_SUBSCRIPTION_ID \
			client_id=$ARM_CLIENT_ID \
			secret_access_key=$ARM_CLIENT_SECRET \
			tag_name=consul \
			tag_value=server)
	if [ "$ip" != "$tagged_ips" ] ; then
		echo "got $ip on $h want $tagged_ips"
		exit 1
	fi
done

echo "OK"
exit 0
