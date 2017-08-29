#!/usr/bin/env bash

public_ips=$(terraform output public_ips | tr ',' ' ')
tagged_ips=$(terraform output tagged_ips | tr ',' ' ')

for h in $public_ips; do
	ip=$(ssh -q root@$h \
		-i tf_rsa \
		-o UserKnownHostsFile=/dev/null \
		-o StrictHostKeyChecking=no \
		/tmp/discover -q addrs \
			provider=digitalocean \
			region=nyc3 \
			tag_name=go-discover-test-tag \
			api_token=$TF_VAR_digitalocean_token)
	if [ "$ip" != "$tagged_ips" ] ; then
		echo "got $ip on $h want $tagged_ips"
		exit 1
	fi
done

echo "OK"
exit 0
