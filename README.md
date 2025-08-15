# Hot R.O.D. - Rides on Demand

This demo is based on the Jaeger HotROD demo but has 
been modified considerably to showcase Signadot & Sandboxes.

![image](/docs/graph.png)


## Running

First, [install Signadot Operator](https://www.signadot.com/docs/installation/signadot-operator)
if you haven't already.
Next, install and configure the [Signadot CLI](https://www.signadot.com/docs/getting-started/installation/signadot-cli).

### Install the HotROD Application

Decide on a namespace in which to install HotROD and then run:

```bash
kubectl create ns "${NAMESPACE}"
kubectl -n "${NAMESPACE}" apply -k k8s/overlays/prod/<devmesh | istio>
```

To uninstall:

```bash
kubectl delete ns "${NAMESPACE}"
```


## Development

You can easily run any of the HotROD services on your workstation by leveraging
Signadot sandboxes. Let’s walk through an example of how to do this for the
frontend application.

### Frontend

Let's start by connecting the cluster and creating a sandbox with a local
workload. Be sure to replace `<cluster-name>` and `<namespace>` with the
appropriate values for your setup:

```bash
# Connect to the cluster
#
signadot local connect


# Create a local sandbox for the frontend app
#
signadot sandbox apply -f - <<EOF
name: local-hotrod-frontend
spec:
  cluster: <cluster-name>
  local:
  - name: local-frontend
    from:
      kind: Deployment
      namespace: <namespace>
      name: frontend
    mappings:
    - port: 8080
      toLocal: localhost:8080
EOF
```

The above sandbox ensures that any traffic within the cluster directed to
`http://frontend.<namespace>:8080` is routed to your local machine at
`localhost:8080`.

Now it's time to run the local version of the frontend application, which will
bind to `localhost:8080`.

Since the frontend depends on several other services, such as the location
service, Redis, Kafka, and others, we'll continue using the existing in-cluster
instances of these dependencies. To ensure the local frontend can seamlessly
communicate with them, we'll reuse the environment configuration from the
in-cluster baseline workload.

To load these environment variables into your local shell, run:

```bash
eval "$(signadot sandbox get-env local-hotrod-frontend)"
```

This command fetches and sets the necessary environment variables so your local
frontend can interact with the in-cluster services just as it would within the
cluster.

To run the frontend application locally, you have a couple of options.

If you prefer automatic hot-reloading during development, you can use `air`:

```bash
air
```

Alternatively, you can manually build and run the frontend using the `make` and
`go run` commands:

```bash
make build-frontend-app
go run ./cmd/hotrod frontend
```

Both approaches will launch the frontend bound to `localhost:8080`, ready to
receive traffic from the Signadot sandbox.

Now you can even run Signadot Smart Tests against your local frontend to
validate its behavior in a realistic environment. For example:

```bash
signadot st run --sandbox=local-hotrod-frontend --publish
Created test run with ID "eu3o6mpz6mtc" in cluster "xrc-test".

Test run status:
✅   ...rt-tests/frontend/post-frontend-dispatch.star   [ID: eu3o6mpz6mtc-1, STATUS: completed]
✅   ...hotrod/smart-tests/location/get-location.star   [ID: eu3o6mpz6mtc-2, STATUS: completed]

Test run summary:
* Executions
   ✅ 2/2 tests completed
* Diffs
   ✅ No high/medium relevance differences found
* Checks
   ✅ 2 checks passed
```


## Deployment

### Release 

To build and push new images, we can leverage by using the `make release`.

For the case of releasing latest images we can do `RELEASE_TAG=latest make release`.
Note that you can replace the `RELEASE_TAG` with the value you need.

### Considerations

You have to make sure you have the rights to write in signadot/hotrod.
