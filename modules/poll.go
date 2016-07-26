package modules

import (
	"net/http"
	"time"
)

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
