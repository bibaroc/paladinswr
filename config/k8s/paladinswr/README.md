# Deploying in Kubernetes

## Adding Kubernetes Secrets
```sh
kubectl create secret generic paladinswr-injector --from-env-file=./config/k8s/paladinswr/secrets.injector.env -n paladinswr
```