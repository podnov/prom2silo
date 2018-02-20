package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty"
	"log"
	"os"
	"strings"
)

const siloAlertApiUri = "api/alert"

func main() {
	listenAddress := os.Getenv("PROM2SILO_LISTEN_ADDRESS")

	if listenAddress == "" {
		listenAddress = ":8080"
	}

	r := gin.Default()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.POST("/v1", handleV1Post)

	r.Run(listenAddress)
}

func convertPrometheusAlertToScienceLogicAlert(prometheusAlert PrometheusAlert) ScienceLogicAlert {
	alertname := prometheusAlert.Alertname
	description := prometheusAlert.Annotations["description"]
	instance := prometheusAlert.Labels["instance"]
	severity := strings.Title(prometheusAlert.Labels["severity"])
	status := strings.Title(prometheusAlert.Status)

	message := fmt.Sprintf(`%v: %v
Severity: %v
Instance: %v

%v`, status, alertname,
		severity,
		instance,
		description)

	result := ScienceLogicAlert{}

	result.AlignedResource = os.Getenv("PROM2SILO_SILO_ALIGNED_RESOURCE")
	result.Message = message

	return result
}

func handleV1Post(c *gin.Context) {
	log.Printf("Handling v1 POST\n")

	prometheusPayload := PrometheusPayload{}
	err := c.BindJSON(&prometheusPayload)
	if err != nil {
		panic(err.Error())
	}

	alerts := prometheusPayload.Alerts
	alertCount := len(alerts)
	log.Printf("Received Prometheus payload with [%v] alerts\n", alertCount)

	for _, alert := range alerts {
		sendScienceLogicAlert(alert)
	}

	log.Printf("Sent [%v] ScienceLogic alerts\n", alertCount)
}

func sendScienceLogicAlert(prometheusAlert PrometheusAlert) {
	scienceLogicAlert := convertPrometheusAlertToScienceLogicAlert(prometheusAlert)

	scienceLogicAlertJson, err := json.Marshal(scienceLogicAlert)
	if err == nil {
		message := fmt.Sprintf("Sending ScienceLogic alert [%s]\n", scienceLogicAlertJson)
		log.Printf(message)
	} else {
		log.Printf("Could not serialize converted ScienceLogic alert: %s\n", err)
	}

	siloBaseUrl := os.Getenv("PROM2SILO_SILO_BASE_URL")

	url := fmt.Sprintf("%s/%s", siloBaseUrl, siloAlertApiUri)

	createdScienceLogicAlert := ScienceLogicAlert{}

	siloUsername := os.Getenv("PROM2SILO_SILO_USERNAME")
	siloPassword := os.Getenv("PROM2SILO_SILO_PASSWORD")

	response, err := resty.R().
		SetBasicAuth(siloUsername, siloPassword).
		SetBody(scienceLogicAlert).
		SetResult(&createdScienceLogicAlert).
		Post(url)

	if err != nil {
		panic(err.Error())
	}

	createdScienceLogicAlertJson, err := json.Marshal(createdScienceLogicAlert)

	if err == nil {
		message := fmt.Sprintf("Received status [%v] and response body [%s] from ScienceLogic\n", response.StatusCode(), createdScienceLogicAlertJson)
		log.Printf(message)
	} else {
		log.Printf("Could not serialize ScienceLogic alert response: %v\n", err)
	}
}
