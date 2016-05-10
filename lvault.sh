#!/bin/sh

sleep 3
environ=$LAIN_DOMAIN
DEBUG="true"

source ./config

export SSO_CLIENT_ID=$clientid
export SSO_CLIENT_SECRET=$clientsec
export SSO_SERVER_NAME=$ssoserver

exec ./lvault-0.1.linux.amd64 -ssoserver=$SSO_SERVER_NAME -ssoid=$SSO_CLIENT_ID -ssosecret=$SSO_CLIENT_SECRET -debug=$DEBUG
