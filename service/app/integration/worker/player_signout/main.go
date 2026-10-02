package main

import (
	"github.com/tnminh2412/go-microservice-structure/service/app/integration/worker/player_signout/cmd"
)

var revision string

func main() {
	cmd.SetRevision(revision)
	cmd.Execute()
}
