# gRPC SPIFFE/SPIRE Example Project

## Starting Minikube

```bash
minikube start
```

## Building

In a dedicated shell, configure docker to use the minikube registry and then build the
images and push them into it.

```bash
eval $(minikube docker-env)
make all-images
```

## Setup Spire

Setup Spire using the spiffe-csi manifests and then wait for the spire server and agent
to be ready.

```bash
k apply -f spiffe-csi/spire-namespace.yaml
k apply -f spiffe-csi/spire-csi-driver.yaml
k apply -f spiffe-csi/spire-server.yaml
k rollout status -n spire deployment/spire-server
k apply -f spiffe-csi/spire-agent.yaml
k rollout status -n spire daemonset/spire-agent
```

Create registration entries to assign SPIFFE IDs to the add server and the api client.

```bash
k exec -it -n spire deployment/spire-server -- \
    /opt/spire/bin/spire-server entry create \
    -spiffeID spiffe://example.org/ns/spire/sa/spire-agent \
    -selector k8s_psat:cluster:demo-cluster \
    -selector k8s_psat:agent_ns:spire \
    -selector k8s_psat:agent_sa:spire-agent \
    -node

k exec -it -n spire deployment/spire-server -- \
    /opt/spire/bin/spire-server entry create \
    -spiffeID spiffe://example.org/ns/default/add \
    -parentID spiffe://example.org/ns/spire/sa/spire-agent \
    -selector k8s:ns:default \
    -selector k8s:sa:default \
    -selector k8s:pod-label:app:add

k exec -it -n spire deployment/spire-server -- \
    /opt/spire/bin/spire-server entry create \
    -spiffeID spiffe://example.org/ns/default/api \
    -parentID spiffe://example.org/ns/spire/sa/spire-agent \
    -selector k8s:ns:default \
    -selector k8s:sa:default \
    -selector k8s:pod-label:app:api
```

## Running

Deploy the add service and api service.

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

## Debugging

If something goes wrong we can take a look at the spire-agent logs. If needed
we can deploy the spire-agent image as a client workload and run spire-agent
commands directly on it.

```bash
k logs -l app=spire-agent -n spire -c spire-agent -f

k exec -it -n spire deployment/spire-server -- \
    /opt/spire/bin/spire-server entry create \
    -spiffeID spiffe://example.org/ns/default/client \
    -parentID spiffe://example.org/ns/spire/sa/spire-agent \
    -selector k8s:ns:default \
    -selector k8s:sa:default \
    -selector k8s:pod-label:app:client

k exec -it -n spire deployment/spire-server -- /opt/spire/bin/spire-server entry show

k apply -f spiffe-csi/client-deployment.yaml
k rollout status deployment/client

k exec -it deployment/client -- /opt/spire/bin/spire-agent api fetch -socketPath /spiffe-workload-api/spire-agent.sock
```

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
* <https://github.com/spiffe/spire/blob/main/doc/spire_server.md>
* <https://github.com/spiffe/spire-tutorials/tree/master/k8s/quickstart>
* <https://github.com/ryoya-fujimoto/grpc-testing>
* <https://github.com/google/addlicense>
* <https://github.com/Nordix/Meridio/blob/master/pkg/nsm/client.go>
* <https://github.com/HewlettPackard/k8s-sigstore-attestor/blob/main/pkg/server/plugin/upstreamauthority/spire/spire_server_client.go>
* <https://github.com/nishantapatil3/spire-federation-kind/blob/main/src/broker-webapp/main.go>
* <https://github.com/spiffe/spiffe-csi/tree/main/example>
* <https://github.com/nicktrav/spire-envoy-k8s-demo/blob/master/scripts/4-delete-spire-entries.sh>
