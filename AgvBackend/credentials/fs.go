package credentials

import (
	"fmt"
	"os/user"
	"path/filepath"
	"turtle/core/lgr"
)

func GetDarwinWorkspace() string {

	fromEnv := GetEnvOrDefault(fmt.Sprintf("TURTLE_%s_MASTER", APP_NAME), "")

	if fromEnv == "" {
		usr, err := user.Current()
		appSupportPath := filepath.Join(usr.HomeDir, "Library", "Application Support")
		if err == nil {
			return appSupportPath + "/" + GetAppName()
		} else {
			lgr.Error(err.Error())
		}

	} else {
		return fromEnv + "/" + GetAppName()
	}

	return "../TurtleIntelligenceStorage" + "/" + GetAppName()

}

func GetStoragePath() string {
	return GetEnvOrDefault(fmt.Sprintf("TURTLE_%s_STORAGE", APP_NAME), "LOCALAPPDATA")
}
