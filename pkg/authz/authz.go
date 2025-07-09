package authz

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

func CheckAAAPermission(username, resource, action, resourceID, jwtToken string) (bool, error) {
	payload := map[string]interface{}{
		"user_id":       resourceID, // Use the user's UUID as user_id
		"resource_name": resource,
		"action":        action,
		"principal_id":  resourceID, // Use the same ID as principal_id for user's own resources
	}

	body, _ := json.Marshal(payload)

	log.Printf("=== AAA AUTHORIZATION DEBUG ===")
	log.Printf("Calling AAA service: %s/check-permission", os.Getenv("AAA_SERVICE_URL"))
	log.Printf("Payload: %s", string(body))
	log.Printf("UserID: %s", resourceID)
	log.Printf("ResourceName: %s", resource)
	log.Printf("Action: %s", action)
	log.Printf("PrincipalID: %s", resourceID)
	log.Printf("Username: %s", username)
	log.Printf("JWT Token: %s", jwtToken[:50]+"...") // Show first 50 chars

	req, err := http.NewRequest("POST", os.Getenv("AAA_SERVICE_URL")+"/check-permission", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return false, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwtToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error calling AAA service: %v", err)
		return false, err
	}
	defer resp.Body.Close()

	responseBody, _ := io.ReadAll(resp.Body)
	log.Printf("AAA Response Status: %d", resp.StatusCode)
	log.Printf("AAA Response Body: %s", string(responseBody))
	log.Printf("AAA Request Headers: %v", req.Header)
	log.Printf("AAA Request URL: %s", req.URL.String())

	var out struct {
		Status  bool     `json:"status"`
		Message string   `json:"message"`
		Data    bool     `json:"data"`
		Errors  []string `json:"errors"`
	}
	json.Unmarshal(responseBody, &out)

	log.Printf("AAA Authorization Result: Status=%v, Data=%v, Message=%s", out.Status, out.Data, out.Message)
	log.Printf("=== END AAA AUTHORIZATION DEBUG ===")

	return out.Data, nil
}
