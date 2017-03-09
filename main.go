package main

import (
	config "plugin-to-SONA/config"
	log "plugin-to-SONA/log"
	plugin "plugin-to-SONA/v1"
)

func main() {
	log.Logging()
	log.Info("Starting plugin for SONA")
	config.ParseCommandLine()
	log.Info("Terminating with code: ", plugin.Run())
}
