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


## Development

### Frontend

To run frontend you could easily run with `air` that helps with hot-reload. 

Before running `air` or manual steps you have to set up the following env
```shell
export KAFKA_BROKER=kafka-headless.namespace.svc:9092
export REDIS_ADDR=redis.namespace.svc:6379
export FRONTEND_LOCATION_ADDR=location.namespace.svc:8081
```

Now let's run the frontend
```shell
air
```

That will listen for the changes and restart the server every change.

If no want to use this approach, you could
```shell
make build-frontend-app
go run ./cmd/hotrod/main.go frontend
```
