package http

import "fmt"

type IntegrationError struct {
	StatusCode  int    `json:"status_code,omitempty"`
	Description string `json:"description,omitempty"`
	Message     string `json:"message,omitempty"`
}

func (ie *IntegrationError) Error() string {
	return fmt.Sprintf("Status %d: Integration Error: %s", ie.StatusCode, ie.Message)
}
