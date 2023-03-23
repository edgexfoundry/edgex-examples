#!/bin/sh
#
# Copyright (C) 2022 Intel Corporation
#
# SPDX-License-Identifier: Apache-2.0
#

# server (CA, needed for curl)

openssl ecparam -genkey -name secp384r1 -noout -out server-ca.key
openssl req -x509 -new -key server-ca.key -subj "/CN=ServerCA" -sha384 -out server-ca.pem

# server (leaf, needed for nginx)

openssl ecparam -genkey -name secp384r1 -noout -out server.key
openssl req -subj "/CN=edgex" -addext "subjectAltName = DNS:edgex" -config openssl.conf -key server.key -sha384 -new -out server.req
openssl x509 -sha384 -extfile openssl.conf -extensions server_ext -extensions edgex_san -CA server-ca.pem -CAkey server-ca.key -CAcreateserial -req -in server.req -days 365 -out server.pem

# client (CA, needed for ngnix mutual auth)

openssl ecparam -genkey -name secp384r1 -noout -out client-ca.key
openssl req -x509 -new -key client-ca.key -subj "/CN=ClientCA" -sha384 -out client-ca.pem

# client (leaf, needed for curl)

openssl ecparam -genkey -name secp384r1 -noout -out client.key
openssl req -subj "/CN=client" -config openssl.conf -key client.key -sha384 -new -out client.req
openssl x509 -sha384 -extfile openssl.conf -extensions client_ext -CA client-ca.pem -CAkey client-ca.key -CAcreateserial -req -in client.req -days 365 -out client.pem
