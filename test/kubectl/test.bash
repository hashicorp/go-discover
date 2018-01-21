#!/usr/bin/env bash

function valid_ip() {
    local  ip=$1
    local  stat=1

    if [[ $ip =~ ^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}$ ]]; then
        OIFS=$IFS
        IFS='.'
        ip=($ip)
        IFS=$OIFS
        [[ ${ip[0]} -le 255 && ${ip[1]} -le 255 \
            && ${ip[2]} -le 255 && ${ip[3]} -le 255 ]]
        stat=$?
    fi
    return $stat
}

tagged_ips=$(kubectl --namespace hashicorp logs discover-nginx | tail -1)

for ip in $tagged_ips; do
	if valid_ip $ip; then stat='good'; else stat='bad'; fi
	printf "%-20s: %s\n" "$ip" "$stat"
	if [ "$stat" != "good" ] ; then
		exit 1
	fi
done

echo "OK"
exit 0
