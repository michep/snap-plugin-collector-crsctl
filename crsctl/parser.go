package crsctl

import (
	"strings"
)

var (
	crsServices = map[string]string{
		"Oracle High Availability Services": "HAS",
		"Cluster Ready Services":            "CRS",
		"Cluster Synchronization Services":  "CSS",
		"Event Manager":                     "EM",
	}
)

type ResourceStatus struct {
	Name        string
	OnlineNodes []string
}

type ResourceStatuses []*ResourceStatus

type CrsCheck struct {
	Available   []string
	Unavailable []string
}

type CheckResult struct {
	OK    []string
	NotOK []string
}

func parseStatusResource(output string) (result []*ResourceStatus) {
	for _, part := range strings.Split(output, newLine+newLine) {
		partLines := strings.Split(part, newLine)
		if strings.HasPrefix(partLines[0], "NAME") {
			name := strings.Split(partLines[0], "=")[1]
			statesLine := strings.Split(partLines[3], "=")[1]
			states := strings.Split(statesLine, ",")

			r := ResourceStatus{Name: name}
			for _, state := range states {
				if strings.Index(state, "ONLINE") != -1 {
					node := strings.Split(state, "on")[1]
					r.OnlineNodes = append(r.OnlineNodes, strings.TrimSpace(node))
				}
			}
			result = append(result, &r)
		}
	}
	return result
}

func parseCrsCheck(output string) (result CrsCheck) {
	for _, line := range strings.Split(output, newLine) {
		for svc, svcabbr := range crsServices {
			if strings.Index(line, svc) != -1 {
				if strings.Index(line, "is online") == -1 {
					result.Unavailable = append(result.Unavailable, svcabbr)
				} else {
					result.Available = append(result.Available, svcabbr)
				}
			}
		}
	}
	return result
}

func checkCrsCheck(crsCheck CrsCheck) (result CheckResult) {
	result.OK = crsCheck.Available
	result.NotOK = crsCheck.Unavailable
	return result
}

func checkStatusResource(statusResource ResourceStatuses, target map[string]string) (result CheckResult) {
	targetNodes := target["nodes"]
	for targetSvc, opt := range target {
		if targetSvc == "nodes" {
			continue
		}

		if !statusResource.hasName(targetSvc) {
			result.NotOK = append(result.NotOK, targetSvc)
			continue
		}

		online := statusResource.byName(targetSvc).OnlineNodes
		if len(online) == 0 {
			result.NotOK = append(result.NotOK, targetSvc)
			continue
		}

		checkres := false
		if opt == "all" {
			checkres = checkAll(targetNodes, online)
		} else if opt == "any" {
			checkres = checkAny(targetNodes, online)
		} else {
			checkres = checkAll(opt, online)
		}

		if checkres {
			result.OK = append(result.OK, targetSvc)
		} else {
			result.NotOK = append(result.NotOK, targetSvc)
		}

	}
	return result

}

func checkAll(nodes_str string, online []string) bool {
	nodes := strings.Split(nodes_str, ",")
	if len(nodes) != len(online) {
		return false
	}
	for _, node := range nodes {
		if !has(online, node) {
			return false
		}
	}
	return true
}

func checkAny(nodes_str string, online []string) bool {
	nodes := strings.Split(nodes_str, ",")
	for _, node := range online {
		if !has(nodes, node) {
			return false
		}
	}
	return true
}

func has(strs []string, str string) bool {
	for _, s := range strs {
		if s == str {
			return true
		}
	}
	return false
}

func (rs ResourceStatuses) hasName(name string) bool {
	for _, node := range rs {
		if node.Name == name {
			return true
		}
	}
	return false
}

func (rs ResourceStatuses) byName(name string) *ResourceStatus {
	for _, node := range rs {
		if node.Name == name {
			return node
		}
	}
	return nil
}
