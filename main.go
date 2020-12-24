package main

import (
	"fmt"
	"os"

	build "github.com/DhruvBisla/ogolb/pkg/build"
	serve "github.com/DhruvBisla/ogolb/pkg/serve"
	setup "github.com/DhruvBisla/ogolb/pkg/setup"
)

func main() {
	if len(os.Args) > 1 {
		command := os.Args[1]
		switch command {
		case "setup":
			fmt.Println("Setup requested")
			setup.Setup()
		case "build":
			fmt.Println("Build requested")
			build.Build()
		case "serve":
			fmt.Println("Serve reqeusted")
			serve.Serve()
		case "help":
			fmt.Println("Try 'init' to get started or 'build' to build your project")
		default:
			fmt.Println("Nothing given")
		}
	} else {
		fmt.Println("Try 'help' to learn about things you can do")
	}
}
