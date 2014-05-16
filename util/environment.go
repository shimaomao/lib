package util

import (
	"os"
	"strings"
)

type EnvironmentInfo string

var (
	EILocal    = EnvironmentInfo("local")
	EIDev      = EnvironmentInfo("dev")
	EIProd     = EnvironmentInfo("prod")
	EILocalGAE = EnvironmentInfo("local-gae")
	EIDevGAE   = EnvironmentInfo("dev-gae")
	EIProdGAE  = EnvironmentInfo("prod-gae")
)
var CurrentEnvironment EnvironmentInfo

func init() {
	resolveEnvironment()
}

func IsProduction() bool {
	return CurrentEnvironment == EIProd || CurrentEnvironment == EIProdGAE
}

func IsOnGAE() bool {
	// return false for LocalGAE since it is controllable.
	return CurrentEnvironment == EIDevGAE || CurrentEnvironment == EIProdGAE
}

func resolveEnvironment() {
	if os.Getenv("SPEEDLAND_ENV") != "" {
		CurrentEnvironment = EnvironmentInfo(os.Getenv("SPEEDLAND_ENV"))
		return
	} else {
		// We could not configure envvar in GAE, so need to detect
		if os.Getenv("RUN_WITH_DEVAPPSERVER") == "1" {
			CurrentEnvironment = EILocalGAE
			return
		} else {
			if os.Getenv("USER") == "" {
				// GAE environment, check the directory name
				pwd, _ := os.Getwd()
				for _, s := range strings.Split(pwd, "/") {
					if strings.HasSuffix(s, "-dev") {
						CurrentEnvironment = EIDevGAE
						return
					}
				}
				CurrentEnvironment = EIProdGAE
			}
		}
	}
}
