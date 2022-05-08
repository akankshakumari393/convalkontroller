# Validation webhook controller

This webhook controller validates the `Depkon` resource on create and update request and check if configmap and all the Deployments are present in the mentioned namespace.

### Install using helm
```
helm repo add akankshakumari393 https://akankshakumari393.github.io/helm-charts
kubectl create namespace valcontroller
helm install depkon akankshakumari393/convalkontroller -n valcontroller
```