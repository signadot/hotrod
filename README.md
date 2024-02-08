# Hot R.O.D. - Rides on Demand

This demo is based on the Jaeger HotROD demo but has 
been modified considerably to showcase Signadot & Sandboxes.

![image](/docs/graph.png)

## Running

First, [install Signadot Operator](https://docs.signadot.com/docs/installation)
if you haven't already.

### Run everything in Kubernetes

Decide on a namespace in which to install HotROD and then run:

```sh
kubectl create ns "${NAMESPACE}"
kubectl -n "${NAMESPACE}" apply -k k8s/overlays/prod/<devmesh | istio>
```

To uninstall:

```bash
kubectl delete ns "${NAMESPACE}"
```

### Release

If you want to set a newly pushed set of images as the default for the application, you can
run the following:

```bash
cd k8s/base && kustomize edit set image signadot/hotrod:$(RELEASE_TAG)
```
