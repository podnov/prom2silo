package main

import (
	"encoding/json"
	"os"
	"testing"
)

func Test_convertPrometheusAlertToScienceLogicAlert_status_firing(t *testing.T) {
	givenAlignedResource := "/device/42"

	oldAlignedResource := os.Getenv(alignedResourceEnvVarName)
	os.Setenv(alignedResourceEnvVarName, givenAlignedResource)
	defer func() {
		os.Setenv(alignedResourceEnvVarName, oldAlignedResource)
	}()

	givenAlertAnnotations := map[string]string{
		"description": "Pod /silo-etl-config-dynamic-apps-config-etl-stage-1516893600-6czrl is was restarted 5.042016806722689 times within the last hour",
		"summary":     "my-summary",
	}

	givenAlertLabels := map[string]string{
		"instance": "10.233.96.177:8080",
		"severity": "warning",
	}

	givenAlert := PrometheusAlert{}
	givenAlert.Alertname = "PodFrequentlyRestarting"
	givenAlert.Annotations = givenAlertAnnotations
	givenAlert.Labels = givenAlertLabels
	givenAlert.Status = "firing"

	actualAlert := convertPrometheusAlertToScienceLogicAlert(givenAlert)

	actualBytes, err := json.MarshalIndent(actualAlert, "", "    ")
	if err != nil {
		t.Error(err)
	}

	actualAlertJson := string(actualBytes)

	expectedAlertJson := `{
    "aligned_resource": "/device/42",
    "message": "Firing: PodFrequentlyRestarting\nSeverity: Warning\nInstance: 10.233.96.177:8080\n\nPod /silo-etl-config-dynamic-apps-config-etl-stage-1516893600-6czrl is was restarted 5.042016806722689 times within the last hour"
}`

	if actualAlertJson != expectedAlertJson {
		t.Errorf("got: %s, want: %s", actualAlertJson, expectedAlertJson)
	}
}

func Test_convertPrometheusAlertToScienceLogicAlert_status_resolved(t *testing.T) {
	givenAlignedResource := "/device/42"

	oldAlignedResource := os.Getenv(alignedResourceEnvVarName)
	os.Setenv(alignedResourceEnvVarName, givenAlignedResource)
	defer func() {
		os.Setenv(alignedResourceEnvVarName, oldAlignedResource)
	}()

	givenAlertAnnotations := map[string]string{
		"description": "Pod /silo-etl-config-dynamic-apps-config-etl-stage-1516893600-6czrl is was restarted 5.042016806722689 times within the last hour",
		"summary":     "my-summary",
	}

	givenAlertLabels := map[string]string{
		"instance": "10.233.96.177:8080",
		"severity": "warning",
	}

	givenAlert := PrometheusAlert{}
	givenAlert.Alertname = "PodFrequentlyRestarting"
	givenAlert.Annotations = givenAlertAnnotations
	givenAlert.Labels = givenAlertLabels
	givenAlert.Status = "resolved"

	actualAlert := convertPrometheusAlertToScienceLogicAlert(givenAlert)

	actualBytes, err := json.MarshalIndent(actualAlert, "", "    ")
	if err != nil {
		t.Error(err)
	}

	actualAlertJson := string(actualBytes)

	expectedAlertJson := `{
    "aligned_resource": "/device/42",
    "message": "Resolved: PodFrequentlyRestarting\nSeverity: Warning\nInstance: 10.233.96.177:8080\n\nPod /silo-etl-config-dynamic-apps-config-etl-stage-1516893600-6czrl is was restarted 5.042016806722689 times within the last hour"
}`

	if actualAlertJson != expectedAlertJson {
		t.Errorf("got: %s, want: %s", actualAlertJson, expectedAlertJson)
	}
}
