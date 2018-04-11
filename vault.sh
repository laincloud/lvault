#!/bin/sh

cd /lain/app

sleep 5

ip=$(ip addr | grep 'state UP' -A2 | tail -n1 | awk '{print $2}' | cut -f1  -d'/')

echo "storage \"consul\" {
    path = \"vault/\"
	address = \"consul.lain:8500\"
}

listener \"tcp\" {
	address = \"0.0.0.0:8200\"
	tls_disable = 1
}
disable_mlock=true
api_addr = \"http://$ip:8200\"
cluster_addr = \"http://$ip:8201\"
">/lain/app/vault.conf

echo "start"
exec vault server -config=/lain/app/vault.conf
