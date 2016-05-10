#!/bin/bash

#export VAULT_ADDR='http://127.0.0.1:8200'

vault unseal 8490ca4cac392f560e39a801e6ba0aa68d4d403b259f664b627411f5ec53e821

vault auth 4429e1bc-5d1c-34a6-b1d2-0623a46635c4

if [ $DEPLOYD_POD_INSTANCE_NO -eq 4 ]
then
	vault audit-enable file path=audit.log
fi
