#!/bin/sh
#
# Copyright (C) 2022 Intel Corporation
#
# SPDX-License-Identifier: Apache-2.0
#
kubectl create namespace edgex

kubectl delete secret edgex-tls -n edgex
kubectl create secret tls -n edgex edgex-tls --cert=server.pem --key=server.key

kubectl delete secret edgex-client-ca -n edgex
kubectl create secret generic -n edgex edgex-client-ca --from-file=ca.crt=client-ca.pem

