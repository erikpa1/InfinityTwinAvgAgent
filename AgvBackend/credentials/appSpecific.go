package credentials

import "fmt"

func IsSlaveApplication() bool {
	val := GetEnvOrDefault(
		fmt.Sprintf("TURTLE_%s_MASTER", APP_NAME),
		"0")
	return val == "1" || val == "2"
}

func IsMasterApplication() bool {
	val := GetEnvOrDefault(
		fmt.Sprintf("TURTLE_%s_MASTER", APP_NAME),
		"0")

	return val == "0" || val == "2"
}
