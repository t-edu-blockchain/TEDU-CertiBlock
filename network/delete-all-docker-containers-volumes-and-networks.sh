#!/bin/bash
set -x
set -e
set -o pipefail

# 1. Stop and Remove All Containers
docker ps -q | xargs -r docker stop
docker ps -aq | xargs -r docker rm

# 2. Remove All Volumes
docker volume prune -f
docker volume ls -q | xargs -r docker volume rm

# 3. Remove All Networks
docker network prune -f
docker network ls --filter type=custom -q | xargs -r docker network rm

echo "DONE"
