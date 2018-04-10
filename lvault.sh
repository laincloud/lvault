#!/bin/sh

sleep 3

cd /lain/app
exec /lain/app/lvault-0.1.linux.amd64 -ssoserver=$SSO_SERVER_NAME -ssoid=$SSO_CLIENT_ID -ssosecret=$SSO_CLIENT_SECRET -https=$HTTPS -debug=$DEBUG