#!/bin/sh


#(echo "begin"; sleep 15; sh start.sh; echo "ok")&
sleep 5
echo 'backend "etcd" {
	path = "vault/"
	address = "http://etcd.lain:4001"
	advertise_addr ="http://'$LAIN_PROCNAME-$DEPLOYD_POD_INSTANCE_NO.$LAIN_APPNAME'.lain:8200"
}

listener "tcp" {
	address = "0.0.0.0:8200"
	tls_disable = 0
	tls_cert_file = "./certs/server.pem"
	tls_key_file = "./certs/server.key"
}
disable_mlock=false' >vaultetcd.conf

echo "start"
exec vault server -config=vaultetcd.conf
