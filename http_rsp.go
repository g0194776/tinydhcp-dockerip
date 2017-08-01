package main

type HttpResponse struct {
	ErrorID  int    `json:"error-id"`
	Reason   string `json:"reason,omitempty"`
	DockerIP string `json:"docker-ip,omitempty"`
}
