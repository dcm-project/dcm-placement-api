package opa

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Validator struct {
	server string
}

func NewValidator(server string) *Validator {
	return &Validator{server: server}
}

func (v *Validator) EvalPolicy(ctx context.Context, policy string, input interface{}) (bool, error) {
	requestBody := map[string]interface{}{
		"input": input,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return false, err
	}

	url := fmt.Sprintf("%s/v1/data/%s/allow", v.server, policy)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, err
	}

	// Extract result from OPA response
	if allow, ok := result["result"].(bool); ok {
		return allow, nil
	}

	return false, nil
}
