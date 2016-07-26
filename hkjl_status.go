package main

import (
	"flag"
	"github.com/r3boot/hkjl_status/modules"
	"github.com/r3boot/rlib/logger"
)

const D_DEBUG bool = false
const D_TIMESTAMP bool = false

const MAINSITE_URL string = "https://www.hackenkunjeleren.nl/"
const COMMUNITY_URL string = "https://community.hackenkunjeleren.nl/"

var debug = flag.Bool("D", D_DEBUG, "Enable debugging output")
var timestamp = flag.Bool("T", D_TIMESTAMP, "Enable timestamps in output")

var Log logger.Log

func init() {
	flag.Parse()

	Log.UseDebug = *debug
	Log.UseVerbose = *debug
	Log.UseTimestamp = *timestamp
	Log.Debug("Logging initialized")

	modules.Setup(Log)
	Log.Debug("Modules initialized")
}

func main() {
	Log.Debug("Polling and updating status information")

	Log.Debug(modules.PollSite(MAINSITE_URL))
	Log.Debug(modules.PollSite(COMMUNITY_URL))
}
