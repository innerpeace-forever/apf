package apf

import "github.com/cihub/seelog"

type Logger struct {
	loggers map[string]seelog.LoggerInterface
}
