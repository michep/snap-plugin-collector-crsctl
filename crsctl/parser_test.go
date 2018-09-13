package crsctl

import (
	"fmt"
	"testing"
)

func Test_checkStatusResource(t *testing.T) {
	rs := parseStatusResource(resourceStatusOutput)
	ch := checkStatusResource(rs, resourceTarget)
	fmt.Printf("checkStatusResource: %+v\n", ch)
}

func Test_checkCrsCheck(t *testing.T) {
	rs := parseCrsCheck(statusCheckOutput)
	ch := checkCrsCheck(rs)
	fmt.Printf("checkCrsCheck: %+v\n", ch)
}

var (
	resourceStatusOutput = `NAME=both_good
TYPE=ora.diskgroup.type
TARGET=ONLINE       , ONLINE
STATE=ONLINE on lux, ONLINE on nox

NAME=lux_good
TYPE=ora.diskgroup.type
TARGET=ONLINE , ONLINE
STATE=ONLINE on lux, OFFLINE

NAME=nox_good
TYPE=ora.diskgroup.type
TARGET=ONLINE , ONLINE
STATE=OFFLINE, ONLINE on nox

NAME=both_no_good
TYPE=ora.scan_listener.type
TARGET=ONLINE
STATE=OFFLINE , OFFLINE

NAME=one_no_good
TYPE=ora.scan_listener.type
TARGET=ONLINE
STATE=OFFLINE

NAME=running_on_lux
TYPE=ora.ons.type
TARGET=ONLINE       , ONLINE
STATE=ONLINE on lux

NAME=running_on_nox
TYPE=ora.ons.type
TARGET=ONLINE       , ONLINE
STATE=ONLINE on nox
`
	statusCheckOutput = `CRS-4638: Oracle High Availability Services is online
CRS-4535: Cannot communicate with Cluster Ready Services
CRS-4529: Cluster Synchronization Services is online
CRS-4533: Event Manager is online
`
	resourceTarget = map[string]string{
		"nodes":          "lux,nox",
		"both_good":      "any",  // OK
		"lux_good":       "both", // NotOK
		"nox_good":       "nox",  // OK
		"both_no_good":   "any",  // NotOK
		"one_no_good":    "both", // NotOK
		"running_on_lux": "lux",  // OK
		"running_on_nox": "any",  // OK
		"unknown":        "any",  // NotOK
	}
)
