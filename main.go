package main

import (
	"os"
	"plugin-to-SONA"
	"plugin-to-SONA/config"
)

func main() {
	config.ParseCommandLine()
	os.Exit(plugin.Run())
}
