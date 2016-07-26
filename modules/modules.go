package modules

import (
	"github.com/r3boot/rlib/logger"
)

type Response struct {
	Code  int
	Time  float64
	Error error
}

var Log logger.Log

func Setup(l logger.Log) {
	Log = l
}
