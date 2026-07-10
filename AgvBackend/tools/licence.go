package tools

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"os/exec"
	"strings"
	"turtle/core/lgr"
	"turtle/credentials"
	"turtle/vfs"
)

type LicType string

const LIC_TYPE_DISC LicType = ""
const LIC_TYPE_DISC1 LicType = "disc"
const LIC_TYPE_ENV LicType = "env"
const LIC_TYPE_URL LicType = "url"

type License struct {
	Type     string            `json:"type"`
	Machine  string            `json:"machine"`
	IpAdress string            `json:"ip_address"`
	Keys     map[string]string `json:"keys"`
}

func (self *License) HasLicenseFor(key string) bool {
	_, ok := self.Keys[key]
	return ok
}

func GetLicenseString() (string, error) {
	source := LicType(credentials.LicenseSource())

	if source == LIC_TYPE_DISC1 || source == LIC_TYPE_DISC {
		dataStr, err := vfs.GetFileStringFromWDNew("licences/licence.json")
		return dataStr, err
	} else if source == LIC_TYPE_ENV {
		return credentials.GetLicense(), nil
	} else if source == LIC_TYPE_URL {
		return credentials.GetLicense(), nil
	}

	return "", errors.New("license source not supported")
}

func GetMachineID() (string, error) {
	cmd := exec.Command("wmic", "csproduct", "get", "uuid")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(output), "\n")
	if len(lines) < 2 {
		return "", fmt.Errorf("unexpected output from wmic command")
	}
	return strings.TrimSpace(lines[1]), nil
}

func CheckLicence(licenseKey string) bool {
	machineID, err := GetMachineID()

	if credentials.LicenceDisabled() {
		return true
	}

	if err != nil {
		lgr.Info("Failed to get machine ID:", err)
		return false
	}
	lgr.Info("Machine ID:", machineID)

	dataStr, err := GetLicenseString()
	if err != nil {
		lgr.Info("Failed to read licence file:", err)
		return false
	}

	var claims jwt.MapClaims
	_, err = jwt.ParseWithClaims(dataStr, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(licenseKey), nil
	})
	if err != nil {
		lgr.Info("Failed to parse JWT:", err)
		return false
	}

	if machine, ok := claims["machine"].(string); ok && machine == machineID {
		return true
	}

	lgr.Info("Machine ID and licence ID mismatch")
	return false
}
