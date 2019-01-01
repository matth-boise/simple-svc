#!/bin/bash
echo build simple-envoy and update minikube docker registry
eval $(minikube docker-env)
docker build -t simple-envoy:v1 .
