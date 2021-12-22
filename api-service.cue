package grpcspire

service: "api-service": spec: {
	ports: [{
		name:     "http"
		nodePort: 30080
	}]
	type: "NodePort"
}
