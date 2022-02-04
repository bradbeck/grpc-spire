package grpcspire

import (
	k8sCoreV1 "k8s.io/api/core/v1"
	k8sAppsV1 "k8s.io/api/apps/v1"
)

service: [Name=_]: k8sCoreV1.#Service & {
	apiVersion: "v1"
	kind:       "Service"
	metadata: name: Name
}

deployment: [Name=_]: k8sAppsV1.#Deployment & {
	apiVersion: "apps/v1"
	kind:       "Deployment"
	metadata: name: Name
}

#service: {
	name: string
	port: number
    deployment: "\(name)-deployment"
    service: "\(name)-service"
    image: "bradbeck/\(name)-service"
}

#serviceList: [...#service]

#services: #serviceList & [{
	name: "add"
	port: 50051
}, {
	name: "api"
	port: 8080
}]

for s in #services {
	deployment: "\(s.deployment)": {
		metadata: labels: app: s.name
		spec: {
			selector: matchLabels: app: s.name
			replicas: 1
			template: {
				metadata: labels: app: s.name
				spec: {
					containers: [{
						name:            s.name
						image:           s.image
						imagePullPolicy: "Never"
						ports: [{
							containerPort: s.port
							name:          s.service
						}]
						volumeMounts: [{
							name: "spiffe-workload-api"
							mountPath: "/spiffe-workload-api"
							readOnly: true
						}]
						env: [{
							name: "SPIFFE_ENDPOINT_SOCKET"
							value: "unix:///spiffe-workload-api/spire-agent.sock"
						}]
					}]
					volumes: [{
						name: "spiffe-workload-api"
						csi: driver: "csi.spiffe.io"
					}]
				}
			}
		}
	}
	service: "\(s.service)": spec: {
		selector: app: s.name
		ports: [{
			port:        s.port
			targetPort?: s.service
		}]
	}
}
