# Hot R.O.D. - Rides on Demand

A simple ride-sharing application that allows the end users to request rides to
one of 4 locations and have a nearby driver assigned along with an ETA.

It consists of 4 services: `frontend`, `location`, `driver` and `route`, as well as some
stateful components like `Kafka`, `Redis` and `MySQL`.


![image](https://www.signadot.com/docs/img/hotrod-arch.png)


## Running

First, [install Signadot Operator](https://docs.signadot.com/docs/installation)
if you haven't already.

### Run everything in Kubernetes

Decide on a namespace in which to install HotROD and then run:

```sh
kubectl create ns "${NAMESPACE}"
kubectl -n "${NAMESPACE}" apply -k https://github.com/signadot/hotrod/k8s/overlays/prod/devmesh
```

To uninstall:

```sh
kubectl -n "${NAMESPACE}" delete -k https://github.com/signadot/hotrod/k8s/overlays/prod/devmesh
```
