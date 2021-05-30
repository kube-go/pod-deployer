package main

import (
	"encoding/json"
	"os"

	v1 "k8s.io/api/core/v1"
)

func main() {
	pod := &v1.Pod{}
	pod.Name = "nginx"
	// edit deployment spec

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "    ")
	enc.Encode(pod)
}
