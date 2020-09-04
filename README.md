# google-secret-k8s
Kubernetes based CRD and controller for synchronizing google secret manager data to k8s application

# Installation

```
$ kubectl create namespace secret-operator-system
$ helm install gsecret-k8s charts/gsecret-k8s --namespace secret-operator-system
```