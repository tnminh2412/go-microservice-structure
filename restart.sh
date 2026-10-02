#!/bin/bash
while read -r line ; do
    docker rm -f $line
done < <(docker ps | grep "$1_runner" | awk '{print $1}')