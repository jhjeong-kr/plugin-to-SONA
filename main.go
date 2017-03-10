package main

import (
	config "plugin-to-SONA/config"
	log "plugin-to-SONA/log"
	plugin "plugin-to-SONA/v1"
)

func main() {
	log.Info("starting plugin for SONA")
	config.ParseCommandLine()
	log.Info("terminating with code: ", plugin.Run(), ", bye, bye~")
}
