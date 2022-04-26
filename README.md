# Hot R.O.D. - Rides on Demand

This demo is based on the Jaeger Hotrod demo [hotrod-tutorial]

![image](https://user-images.githubusercontent.com/906471/151587572-56d39bc2-c20f-4d87-85b8-7bc7859ac52f.png)


## Running

First, [install Signadot Operator](https://docs.signadot.com/docs/installation)
if you haven't already.

### Run everything in Kubernetes

Decide on a namespace in which to install HotROD and then run:

```sh
kubectl create ns "${NAMESPACE}"
kubectl -n "${NAMESPACE}" apply -f k8s/pieces
```

To try the demo Signadot Resource Plugin, you must also install the HotROD
reource plugin into the `signadot` namespace where Signadot Operator runs:

```sh
kubectl -n signadot apply -f resource-plugin
```

To uninstall:

```sh
kubectl -n "${NAMESPACE}" delete -f k8s/pieces
kubectl -n signadot delete -f resource-plugin
```
