#!/bin/bash
service=$1
env=$2
port_line=`grep 'Port=' $GOPATH/src/go-microservice-structure/$service/config.$env.ini`
arrIN=(${port_line//=/ })
echo ${arrIN[1]}
