package main

import (
	"log"
	"os"

	"github.com/imrenagi/google-secret-k8s/agent-inject/cmd"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func main() {
	c := cmd.NewRootCommand(os.Args[1:])
	err := c.Execute()
	if err != nil {
		log.Fatalf("Unable to execute application command")
	}
}
