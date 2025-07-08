package authz

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)

func CheckAAAPermission(username, resource, action, resourceID, jwtToken string) (bool, error) {
	payload := map[string]interface{}{
		"username": username,
		"resource": resource,
		"action":   action,
	}
	if resourceID != "" {
		payload["resource_id"] = resourceID
	}

	body, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", os.Getenv("AAA_SERVICE_URL")+"/check-permission", bytes.NewBuffer(body))
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwtToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var out struct {
		Success bool `json:"success"`
		Data    bool `json:"data"`
	}
	json.NewDecoder(resp.Body).Decode(&out)
	return out.Data, nil
}
