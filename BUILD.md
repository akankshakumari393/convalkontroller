# Build locally

`go build`

# Build with docker and push image 

```
docker build -t akankshakumari393/convalkontroller:0.0.1 .
docker push akankshakumari393/convalkontroller:0.0.1
```

# Generate keys to be used by webhook controller

```
openssl req -new -x509 -subj "/CN=convalkontroller.default.svc"  -addext "subjectAltName = DNS:convalkontroller.default.svc" -nodes -newkey rsa:4096 -keyout ./manifests/certs/tls.key -out ./manifests/certs/tls.crt -days 365
```
# Add them to be used by secrets

```
kubectl create secret generic certs --from-file=manifests/certs/tls.crt --from-file=manifests/certs/tls.key
```

# Create ClusterRole and ClusterRolebinding as depkon resource can be created in any namespace
```
kubectl create clusterrole convalkon-role --verb=get --resource=deployments,configmaps

kubectl create clusterrolebinding convalkon-role-binding --clusterrole=convalkon-role --user=system:serviceaccount:default:default
```

# Create Validation webhook configuration 
```
# update the ca bundle with `certs` secret tls.crt data
kubectl create -f manifests/validating-webhook.yaml
```

# Create the deployment

```
kubectl create -f mainfests/deployment.yaml
```

# Alternate Deploy from local chart

```
kubectl create namespace {releaseNamespace}
helm install {releaseName} ./convalkontroller/ -n {releaseNamespace}
```

# create the depkon Resource for validation

```
kubectl create -f https://raw.githubusercontent.com/akankshakumari393/depkon/master/manifests/depkon-cr.yaml
```
