package tools

import (
	"strconv"
	"turtle/core/lgr"
)

func StringToInt32(value string, defval int32) int32 {
	status, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		lgr.Error(err.Error())
		return defval
	} else {
		return int32(status)
	}

}

func StringToInt64(value string, defval int64) int64 {
	status, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		lgr.Error(err.Error())
		return defval
	} else {
		return status
	}

}
