package main

import "github.com/opensibyl/git-connector/pkg"

func main() {
	engine := pkg.CreateServer()
	err := engine.Run(":9448")
	if err != nil {
		return
	}
}
