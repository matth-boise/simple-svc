# simple-svc
kubernetes simple service with envoy sidecar - work in progress

intended to be used with minikube

steps:
- start minikube
- `./build-image.sh` to build docker image from src/simple-svc.go and publish to minikube repository
- `(cd simple-envoy; ./build-image.sh)` to build basic envoy docker image and publish to minikube repository
- install helm on minikube
- `helm install helm-chart --set name=service-a --name service-a`
- `helm install helm-chart --set name=service-b --name service-b`
- `minikube service --url service-a ; minikube service --url service-b`
- assuming the service url if service-a is http://192.168.99.100:30001
  - `curl http://192.168.99.100:30001/xxx` will call service-a, which will echo back the path "/xxx" and the pod ip
  -  `curl http://192.168.99.100:30001/service-b/xxx` will call service-a, which will extract service-b from the path and call it with "/xxx"
  -  `curl http://192.168.99.100:30001/service-b-envoy/xxx` will call service-a, which will extract service-b-envoy from the path and call it (which actually calls localhost:80 - which is the envoy sidecar port for egress, which is maps to service-b) with "/xxx"
