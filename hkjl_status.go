package main

import (
	"errors"
	"flag"
	"github.com/r3boot/rlib/logger"
	"net/http"
	"os"
	"strconv"
	"text/template"
	"time"
)

// URLs to monitor status for
const MAINSITE_URL string = "https://www.hackenkunjeleren.nl/"
const COMMUNITY_URL string = "https://community.hackenkunjeleren.nl/"

const D_TIMEOUT string = "30s"

const D_DEBUG bool = false
const D_TIMESTAMP bool = false
const D_OUTPUT string = "/srv/www/hkjl_status/htdocs"
const D_TEMPLATES string = "/usr/share/hkjl_status/templates"

const D_GREEN int = 0
const D_ORANGE int = 1
const D_RED int = 2

type Response struct {
	Url   string
	Code  int
	Time  float64
	Error string
}

type TemplateData struct {
	Status    int
	Timestamp string
	Mainsite  Response
	Community Response
}

var debug = flag.Bool("D", D_DEBUG, "Enable debugging output")
var timestamp = flag.Bool("T", D_TIMESTAMP, "Enable timestamps in output")
var output_dir = flag.String("o", D_OUTPUT, "Directory in which to write output")
var templates = flag.String("t", D_TEMPLATES, "Directory containing templates")

var Log logger.Log

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

	rc <- response

	return
}

/* WriteStatusPage -- Generates the actual status page
 */
func WriteStatusPage(templates string, output_dir string) (err error) {
	var fd *os.File
	var template_data []byte
	var fname string
	var output string
	var output_new string
	var data TemplateData
	var tmpl *template.Template
	var ms_chan chan Response
	var cs_chan chan Response

	// Poll websites to retrieve status info
	Log.Debug("Poll websites")
	ms_chan = make(chan Response, 1)
	cs_chan = make(chan Response, 1)
	go PollSite(MAINSITE_URL, ms_chan)
	go PollSite(COMMUNITY_URL, cs_chan)
	data.Mainsite = <-ms_chan
	data.Community = <-cs_chan

	// Perform magic to fill template data
	data.Status = D_GREEN
	if data.Mainsite.Code != 200 && data.Community.Code != 200 {
		data.Status = D_RED
	} else if data.Mainsite.Code != 200 {
		data.Status = D_ORANGE
	} else if data.Community.Code != 200 {
		data.Status = D_ORANGE
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
}

func main() {
	var err error

	Log.Debug("Polling and updating status information")

	if err = WriteStatusPage(*templates, *output_dir); err != nil {
		Log.Fatal("Failed to render status page: " + err.Error())
	}
}
