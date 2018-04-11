#!/bin/sh

cd /lain/app

sleep 5
echo 'storage "consul" {
    path = "vault/"
	address = "consul.lain:8500"
	api_addr = "http://lvault.lain.local"
}

listener "tcp" {
	address = "0.0.0.0:8200"
	tls_disable = 1
}
disable_mlock=true' >/lain/app/vaultetcd.conf

echo "start"
exec vault server -config=/lain/app/vaultetcd.conf
