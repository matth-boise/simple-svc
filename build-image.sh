#!/bin/bash
echo build simple-svc and update minikube docker registry
mkdir -p bin
CGO_ENABLED=0 GOOS=linux go build -o bin/simple-svc -a src/simple-svc.go
eval $(minikube docker-env)
docker build -t simple-svc:v1 .
