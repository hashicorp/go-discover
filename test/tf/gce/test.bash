#!/usr/bin/env bash

public_ips=$(terraform output public_ips | tr ',' ' ')
tagged_ips=$(terraform output tagged_ips | tr ',' ' ')

for h in $public_ips; do
	ip=$(ssh -q ubuntu@$h \
		-i tf_rsa \
		-o UserKnownHostsFile=/dev/null \
		-o StrictHostKeyChecking=no \
		/tmp/discover -q addrs \
			provider=gce \
			credentials_file=/tmp/gce.json \
			tag_value=consul-0)
	if [ "$ip" != "$tagged_ips" ] ; then
		echo "got $ip on $h want $tagged_ips"
		exit 1
	fi
done

echo "OK"
exit 0
