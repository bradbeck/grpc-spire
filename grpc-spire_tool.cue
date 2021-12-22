package grpcspire

import "encoding/json"

all: [ for x in [service, deployment] for y in x {y}]

command: {
	apply: task: print: {
		kind: "print"
		text: json.MarshalStream(all)
	}
}
