package main

import ()

type PrometheusPayload struct {
	Alerts []PrometheusAlert `json:"alerts"`
}

type PrometheusAlert struct {
	Alertname   string            `json:"alertname,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Status      string            `json:"status,omitempty"`
}

type ScienceLogicAlert struct {
	AlignedResource string `json:"aligned_resource,omitempty"`
	ForceYId        string `json:"force_yid,omitempty"`
	ForceYName      string `json:"force_yname,omitempty"`
	ForceYType      string `json:"force_ytype,omitempty"`
	Message         string `json:"message,omitempty"`
	MessageTime     string `json:"message_time,omitempty"`
	Threshold       string `json:"threshold,omitempty"`
	Value           string `json:"value,omitempty"`
}
