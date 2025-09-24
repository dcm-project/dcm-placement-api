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

func (v *Validator) EvalPolicy(ctx context.Context, policy string, input interface{}) (map[string]interface{}, error) {
	requestBody := map[string]interface{}{
		"input": input,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/v1/data/%s", v.server, policy)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result["result"].(map[string]interface{}), nil
}

func (v *Validator) IsInputValid(result map[string]interface{}) bool {
	return result["inputvalid"].(bool)
}

func (v *Validator) IsOutputValid(result map[string]interface{}) bool {
	return result["outputvalid"].(bool)
}

func (v *Validator) GetOutputZones(result map[string]interface{}) []string {
	res := result["result"].([]interface{})[0]
	zones := res.(map[string]interface{})["output"].(map[string]interface{})["zones"].([]interface{})
	zonesList := make([]string, len(zones))
	for i, zone := range zones {
		zonesList[i] = zone.(string)
	}
	return zonesList
}
