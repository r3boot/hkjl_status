package main

import (
	"flag"
	"github.com/r3boot/rlib/logger"
	"net/http"
	"time"
)

// URLs to monitor status for
const MAINSITE_URL string = "https://www.hackenkunjeleren.nl/"
const COMMUNITY_URL string = "https://community.hackenkunjeleren.nl/"

const D_DEBUG bool = false
const D_TIMESTAMP bool = false

type Response struct {
	Code  int
	Time  float64
	Error error
}

var debug = flag.Bool("D", D_DEBUG, "Enable debugging output")
var timestamp = flag.Bool("T", D_TIMESTAMP, "Enable timestamps in output")

var Log logger.Log

func PollSite(url string) (response Response) {
	var res *http.Response
	var err error
	var t_start time.Time

	t_start = time.Now()

	res, err = http.Get(url)
	if err != nil {
		response.Error = err
		return
	}

	response.Code = res.StatusCode
	response.Time = time.Since(t_start).Seconds()

	return
}

func init() {
	flag.Parse()

	Log.UseDebug = *debug
	Log.UseVerbose = *debug
	Log.UseTimestamp = *timestamp
	Log.Debug("Logging initialized")
}

func main() {
	Log.Debug("Polling and updating status information")

	Log.Debug(PollSite(MAINSITE_URL))
	Log.Debug(PollSite(COMMUNITY_URL))
}
