package main

import (
	"flag"
	"github.com/josuehennemann/logger"
	"net/http"
	"runtime/debug"
)

var (
	ServiceName      string
	fileConf         = flag.String("fileConf", "", "Endereco do arquivo de configuração")
	config           *Config
	Logger           *logger.Logger
	Access           *logger.Logger
	listHandlersHttp map[string]*MyHandler
)

const (
	FORMAT_DATETIME = "2006-01-02 15:04:05"
)

func recoverPanic() {
	_recoverPanic(nil)
}

func recoverPanicHttp(w http.ResponseWriter) {
	_recoverPanic(w)
}

func _recoverPanic(w http.ResponseWriter) {
	if rec := recover(); rec != nil {
		if w != nil {
			responseInternalError(w)
		}
		stack := debug.Stack()
		Logger.WritePanic(rec, stack)
	}
}
