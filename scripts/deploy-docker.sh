#!/bin/bash

set -e

echo "CGO disabled for production build"
export CGO_ENABLED=0

go generate github.com/AscensionBlockchain/LandingPage/cmd/server

# go get github.com/google/gops
# cp $GOPATH/bin/gops ./services/getascension/executables

#  maydo: call upx(1) on each compiled exe before copying them in

mkdir -p ./services/getascension/executables
go build -mod vendor -v $VENDOR -o ./services/getascension/executables/server \
	github.com/AscensionBlockchain/LandingPage/cmd/server
	

# deploy
echo "resetting all services"
docker-compose build
# docker-compose down --remove-orphans
docker-compose down --volumes
docker-compose up --remove-orphans -d

## watch logs
# rm -f /tmp/getascension.log
# docker-compose logs -f >/tmp/getascension.log 2>&1 & disown
# sleep 3
# tail -f /tmp/getascension.log | go run social/tools/neat -v

echo "Deployment complete"
