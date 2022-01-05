# gRPC SPIFFE/SPIRE Example Project

## Starting Minikube

```bash
minikube start \
    --extra-config=apiserver.service-account-signing-key-file=/var/lib/minikube/certs/sa.key \
    --extra-config=apiserver.service-account-key-file=/var/lib/minikube/certs/sa.pub \
    --extra-config=apiserver.service-account-issuer=api \
    --extra-config=apiserver.service-account-api-audiences=api,spire-server \
    --extra-config=apiserver.authorization-mode=Node,RBAC
```

## Building

In a dedicated shell, configure docker to use the minikube registry and then build the
images and push them into it.

```bash
eval $(minikube docker-env)
make all-images
```

## Setup Spire

Setup Spire using the quickstart manifests and then wait for the spire server to
be ready.

```bash
k apply -k spire/quickstart
k wait --for=condition=ready pod/spire-server-0 -n spire --timeout=5m
```

Create registration entries to assign SPIFFE IDs to the add server and the api client.

```bash
k exec -n spire spire-server-0 -- \
    /opt/spire/bin/spire-server entry create \
    -spiffeID spiffe://example.org/ns/spire/sa/spire-agent \
    -selector k8s_sat:cluster:demo-cluster \
    -selector k8s_sat:agent_ns:spire \
    -selector k8s_sat:agent_sa:spire-agent \
    -node

k exec -n spire spire-server-0 -- \
    /opt/spire/bin/spire-server entry create \
    -spiffeID spiffe://example.org/ns/default/add \
    -parentID spiffe://example.org/ns/spire/sa/spire-agent \
    -selector k8s:ns:default \
    -selector k8s:sa:default \
    -selector k8s:pod-label:app:add

k exec -n spire spire-server-0 -- \
    /opt/spire/bin/spire-server entry create \
    -spiffeID spiffe://example.org/ns/default/api \
    -parentID spiffe://example.org/ns/spire/sa/spire-agent \
    -selector k8s:ns:default \
    -selector k8s:sa:default \
    -selector k8s:pod-label:app:api
```

## Running

Deply the add service and api service.

```bash
cue apply | k apply -f -
```

In a dedicated shell, configure a tunnel to the api service.

```bash
minikube service api-service --url
```

## Validating

```bash
k describe pod -l app=add
k describe pod -l app=api

k logs -l app=add -f
k logs -l app=api -f
```

Use the port for the api service tunnel above to exercise the add endpoint.

```bash
http :<port>/add/3/5
```

## Cleanup

Cleanup the add/api deployments/services.

```bash
cue apply | k delete -f -
```

## To Do

* Find an alternate method for exposing the api service endpoint.

## References

* <https://grpc.io/docs/languages/go/quickstart/>
* <https://medium.com/hackernoon/how-to-develop-go-grpc-microservices-and-deploy-in-kubernates-5eace0425bf8>
* <https://konghq.com/blog/kubernetes-ingress-grpc-example/>
* <https://kubernetes.github.io/ingress-nginx/examples/grpc/>
* <https://github.com/fullstorydev/grpcurl>
* <https://dev.to/techschoolguru/the-complete-grpc-course-protobuf-go-java-2af6>
* <https://istio.io/latest/blog/2021/proxyless-grpc/>
* <https://www.youtube.com/watch?v=cSxHGt7tc88>
* <https://github.com/grpc/grpc-go/tree/master/examples/features/xds>
* <https://www.youtube.com/watch?v=IbcJ8kNmsrE>
* <https://github.com/salrashid123/k8s_grpc_xds>
* <https://github.com/envoyproxy/go-control-plane>
* <https://thebottomturtle.io/Solving-the-bottom-turtle-SPIFFE-SPIRE-Book.pdf>
* <https://letsencrypt.org>
* <https://spiffe.io/docs/latest/try/getting-started-k8s/>
* <https://spiffe.io/docs/latest/deploying/spire_agent/>
* <https://github.com/spiffe/go-spiffe/tree/main/v2/examples/spiffe-grpc>
* <https://github.com/spiffe/spire/blob/main/doc/plugin_agent_workloadattestor_k8s.md>
* <https://github.com/spiffe/spire-tutorials/tree/master/k8s/quickstart>
* <https://github.com/ryoya-fujimoto/grpc-testing>
* <https://github.com/google/addlicense>
