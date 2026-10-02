#!/bin/bash
docker stack services cmm_service_api_metric -q | while read -r line; do
   docker service update --force $line
done
