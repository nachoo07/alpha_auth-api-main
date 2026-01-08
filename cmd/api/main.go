package main

import (
	"github.com/AlphaCodinggroup/alpha_auth-api/cmd/api/modules"
)

func main() {
	app := modules.NewApp()
	app.Run()
}
