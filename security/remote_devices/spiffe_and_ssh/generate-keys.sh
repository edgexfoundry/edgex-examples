#!/bin/sh

test -f id_rsa || ssh-keygen -N '' -C device-ssh-proxy -t rsa -b 4096 -f id_rsa
test -d local/ssh_keys || mkdir local/ssh_keys
cp -f id_rsa* local/ssh_keys
cp -f id_rsa.pub remote/sshd-remote/authorized_keys
