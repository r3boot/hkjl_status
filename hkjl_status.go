package main

import (
	"errors"
	"flag"
	"github.com/r3boot/rlib/logger"
	"gopkg.in/redis.v3"
	"net/http"
	"os"
	"strconv"
	"text/template"
	"time"
)

// URLs to monitor status for
const MAINSITE_URL string = "https://www.hackenkunjeleren.nl/"
const COMMUNITY_URL string = "https://community.hackenkunjeleren.nl/"

var MONITORED_URLS []string = []string{
	"https://www.hackenkunjeleren.nl/",
	"https://community.hackenkunjeleren.nl/",
}

// Poll timeout
const D_TIMEOUT string = "30s"

// Maximum number of buffered items in polling queue
const D_MAX_BUFFERED int = 30

// Default values for CLI arguments
const D_DEBUG bool = false
const D_TIMESTAMP bool = false
const D_OUTPUT string = "/srv/www/hkjl_status/htdocs"
const D_TEMPLATES string = "/usr/share/hkjl_status/templates"
const D_REDIS string = "localhost:6379"

// Constants used to denote the global status
const D_GREEN int = 0
const D_ORANGE int = 1
const D_RED int = 2

// Datastructure used to capture details of a polled website
type Response struct {
	Url   string
	Code  int
	Time  float64
	Error string
}

// Datastructure used to write details into a template
type TemplateData struct {
	Status    int
	Timestamp string
	Responses []Response
}

// Command-line arguments supported by this application
var debug = flag.Bool("D", D_DEBUG, "Enable debugging output")
var timestamp = flag.Bool("T", D_TIMESTAMP, "Enable timestamps in output")
var output_dir = flag.String("o", D_OUTPUT, "Directory in which to write output")
var templates = flag.String("t", D_TEMPLATES, "Directory containing templates")
var redis_addr = flag.String("redis", D_REDIS, "Host:port on which redis is running")

// Logging framework
var Log logger.Log

// Redis client
var Redis redis.Client

func ConnectToRedis(uri string) {
	Redis = *redis.NewClient(&redis.Options{
		Addr: uri,
	})
}

func RedisReachable() bool {
	var pong string
	var err error

	pong, err = Redis.Ping().Result()
	return (pong == "PONG" && err != nil)
}

/* LoadTemplate -- Loads the template 'name' and returns a byte array with
 * the content of the template. Err will be non-nil when any error occurs.
 */
func LoadTemplate(fname string) (data []byte, err error) {
	var fs os.FileInfo
	var fd *os.File
	var size int

	if fs, err = os.Stat(fname); err != nil {
		return
	}

	if fs.IsDir() {
		err = errors.New(fname + " is a directory")
		return
	}

	if fd, err = os.Open(fname); err != nil {
		return
	}
	defer fd.Close()

	data = make([]byte, fs.Size())
	if size, err = fd.Read(data); err != nil {
		data = []byte{}
		return
	}

	if int64(size) != fs.Size() {
		data = []byte{}
		err = errors.New("Incorrect number of bytes read: " + strconv.Itoa(size))
	}

	return
}

/* PollSite -- Polls a website and returns the http code, the response time
 * and any error that occurred during polling the website
 */
func PollSite(url string, rc chan Response) (response Response) {
	var client *http.Client
	var res *http.Response
	var err error
	var t_start time.Time
	var timeout time.Duration

	t_start = time.Now()

	response.Url = url
	response.Code = 666
	response.Error = "Unknown error"

	timeout, err = time.ParseDuration(D_TIMEOUT)
	if err != nil {
		response.Error = err.Error()
		rc <- response
		return
	}

	client = &http.Client{Timeout: timeout}

	res, err = client.Get(url)
	if err != nil {
		response.Error = err.Error()
		rc <- response
		return
	}

	response.Url = url
	response.Code = res.StatusCode
	response.Time = time.Since(t_start).Seconds() * 1000

	if response.Code != 200 && response.Code != 666 {
		response.Error = "Received status code " + strconv.Itoa(response.Code) + " during download of page"
	}

	rc <- response

	return
}

/* WriteStatusPage -- Generates the actual status page
 */
func WriteStatusPage(templates string, output_dir string) (err error) {
	var fd *os.File
	var template_data []byte
	var url string
	var fname string
	var output string
	var output_new string
	var data TemplateData
	var tmpl *template.Template
	var response_chan chan Response
	var response Response
	var num_monitored_urls int
	var num_200s int

	// Poll websites concurrently to retrieve status info
	Log.Debug("Poll websites")
	num_monitored_urls = len(MONITORED_URLS)
	response_chan = make(chan Response, num_monitored_urls)

	// Start polling goroutines
	for _, url = range MONITORED_URLS {
		Log.Debug("Starting polling routine for " + url)
		go PollSite(url, response_chan)
	}

	// Pull results from the polling response channel
	for i := 0; i < num_monitored_urls; i++ {
		response = <-response_chan
		data.Responses = append(data.Responses, response)
	}

	// Calculate global status
	num_200s = 0
	for _, response = range data.Responses {
		if response.Code == 200 {
			num_200s += 1
		}
	}

	if num_200s == num_monitored_urls {
		data.Status = D_GREEN
	} else if num_200s > 0 && num_200s < num_monitored_urls {
		data.Status = D_ORANGE
	} else {
		data.Status = D_RED
	}

	data.Timestamp = time.Now().Format(time.RFC850)

	Log.Debug("Loading template file data")
	fname = templates + "/index.html"
	if template_data, err = LoadTemplate(fname); err != nil {
		err = errors.New("Failed to read template data: " + err.Error())
		return
	}

	Log.Debug("Parsing template file data")
	tmpl, err = template.New("statusPage").Parse(string(template_data))
	if err != nil {
		err = errors.New("Failed to load template: " + err.Error())
		return
	}

	Log.Debug("Opening output file")
	output = output_dir + "/index.html"
	output_new = output_dir + "/index.html.new"
	fd, err = os.OpenFile(output_new, (os.O_CREATE | os.O_WRONLY), 0644)
	if err != nil {
		err = errors.New("Failed to open output file: " + err.Error())
		return
	}
	defer fd.Close()

	Log.Debug("Write output file")
	if err = tmpl.Execute(fd, data); err != nil {
		fd.Close()
		os.Remove(output_new)
		err = errors.New("Failed to render template: " + err.Error())
		return
	}
	fd.Close()

	Log.Debug("Renaming output file")
	os.Rename(output_new, output)

	return
}

func init() {
	flag.Parse()

	Log.UseDebug = *debug
	Log.UseVerbose = *debug
	Log.UseTimestamp = *timestamp
	Log.Debug("Logging initialized")

	ConnectToRedis(*redis_addr)
	if RedisReachable() {
		Log.Debug("Redis initialized")
	} else {
		Log.Warning("Unable to connect to redis, not enabling graphs")
	}
}

func main() {
	var err error

	Log.Debug("Polling and updating status information")

	if err = WriteStatusPage(*templates, *output_dir); err != nil {
		Log.Fatal("Failed to render status page: " + err.Error())
	}
}
