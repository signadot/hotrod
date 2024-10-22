# Hot R.O.D. - Rides on Demand

This demo is based on the Jaeger HotROD demo but has 
been modified considerably to showcase Signadot & Sandboxes.

![image](/docs/graph.png)

## Running

First, [install Signadot Operator](https://www.signadot.com/docs/installation/signadot-operator)
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
