package crsctl

import (
	"os/exec"
	"time"

	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
)

const (
	PluginName    = "crsctl"
	PluginVersion = 1
	PluginVedor   = "mfms"
	crsctlPath    = "crsctl_path"
	resourceName  = "resource"
	crsName       = "crs"
)

type Plugin struct {
	initialized  bool
	crsctlPath   string
	crsctlTarget map[string]string
}

func NewCollector() *Plugin {
	return &Plugin{initialized: false, crsctlTarget: map[string]string{}}
}

func (p *Plugin) GetConfigPolicy() (plugin.ConfigPolicy, error) {
	policy := plugin.NewConfigPolicy()
	policy.AddNewStringRule([]string{PluginVedor, PluginName}, crsctlPath, true)
	return *policy, nil
}

func (p *Plugin) GetMetricTypes(config plugin.Config) ([]plugin.Metric, error) {
	mts := []plugin.Metric{}
	for _, group := range []string{resourceName, crsName} {
		m := plugin.Metric{Namespace: createNamespace(group)}
		mts = append(mts, m)
	}
	return mts, nil
}

func (p *Plugin) CollectMetrics(metrics []plugin.Metric) ([]plugin.Metric, error) {
	var err error
	mts := []plugin.Metric{}
	if !p.initialized {
		config := metrics[0].Config
		p.crsctlPath, err = config.GetString(crsctlPath)
		if err != nil {
			return nil, err
		}
		for k := range config {
			p.crsctlTarget[k], _ = config.GetString(k)
		}
		delete(p.crsctlTarget, crsctlPath)
		p.initialized = true
	}

	statusOutput, err := p.getStatusResourceOutput()
	if err != nil {
		return nil, err
	}
	statusResource := parseStatusResource(statusOutput)
	statusResourceResult := checkStatusResource(statusResource, p.crsctlTarget)

	checkOutput, err := p.getCrsCheckOutput()
	if err != nil {
		return nil, err
	}
	crsCheck := parseCrsCheck(checkOutput)
	crsCheckResult := checkCrsCheck(crsCheck)

	allResults := map[string]CheckResult{resourceName: statusResourceResult, crsName: crsCheckResult}

	ts := time.Now()
	for _, metric := range metrics {
		group := metric.Namespace[2].Value
		results := allResults[group]
		for _, svcName := range results.OK {
			ns := createNamespace(group)
			ns[3].Value = svcName
			mt := plugin.Metric{
				Namespace: ns,
				Data:      0,
				Timestamp: ts,
			}
			mts = append(mts, mt)
		}
		for _, svcName := range results.NotOK {
			ns := createNamespace(group)
			ns[3].Value = svcName
			mt := plugin.Metric{
				Namespace: ns,
				Data:      1,
				Timestamp: ts,
			}
			mts = append(mts, mt)
		}
	}

	return mts, nil
}

func (p *Plugin) getStatusResourceOutput() (string, error) {
	cmd := exec.Command(p.crsctlPath, "status", "resource")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func (p *Plugin) getCrsCheckOutput() (string, error) {
	cmd := exec.Command(p.crsctlPath, "check", "crs")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func createNamespace(group string) plugin.Namespace {
	namespace := plugin.NewNamespace(PluginVedor, PluginName, group)
	namespace = namespace.AddDynamicElement("resource", "resource name")
	namespace = namespace.AddStaticElement("available")
	return namespace
}
