package config

import (
	"flag"
	"os"
	"time"
)

const (
	EXITNORMAL     = 0
	EXITNONROOT    = 1
	EXITLOG        = 2
	EXITKUBECONFIG = 3
	EXITSERVICEENV = 4
	EXITKUBEINIT   = 5
)

var EventHandlingInterval = time.Duration(5)
var LogFormat = "text"
var LogLevel = "info"
var KubeConfig = "" // ex) /etc/kubernetes/admin.conf
const (
	DeamonSetName = "plugin2sona"
)

func init() {
}

func ParseCommandLine() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flag.StringVar(&LogFormat, "logformat", LogFormat, "log format = {text, json}")
	flag.StringVar(&LogLevel, "loglevel", LogLevel, "log level = {info, warning, fatal, error, panic, debug}")
	flag.StringVar(&KubeConfig, "kubeconfig", KubeConfig, "absolute path to kubernetes config file")
	flag.Parse()
}
