package employerprofile

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateEmployerProfileRequest_JSONParsing(t *testing.T) {
	tests := []struct {
		name        string
		requestBody string
		expectError bool
	}{
		{
			name: "Minimal request with only company_name",
			requestBody: `{
				"company_name": "kisanlink"
			}`,
			expectError: false,
		},
		{
			name: "Request with empty strings",
			requestBody: `{
				"company_name": "kisanlink",
				"industry": "",
				"company_size": "",
				"website_url": "",
				"recruiter_name": ""
			}`,
			expectError: false,
		},
		{
			name: "Request with arrays",
			requestBody: `{
				"company_name": "kisanlink",
				"job_categories": ["AgriTech Development"],
				"hiring_locations": ["Hyderabad"],
				"hiring_types": []
			}`,
			expectError: false,
		},
		{
			name: "Full request body",
			requestBody: `{
				"company_name": "kisanlink",
				"website_url": "",
				"industry": "AgriTech / Smart Farming",
				"company_size": "11-50 employees",
				"company_description": "",
				"recruiter_name": "Meher Prasad",
				"designation": "Recruiter",
				"official_email": "meher@gmail.com",
				"phone_number": "7386727005",
				"linkedin_profile": "",
				"job_categories": ["AgriTech Development"],
				"hiring_locations": ["Hyderabad"],
				"hiring_types": []
			}`,
			expectError: false,
		},
		{
			name:        "Empty request body",
			requestBody: `{}`,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req UpdateEmployerProfileRequest
			err := json.Unmarshal([]byte(tt.requestBody), &req)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUpdateEmployerProfileRequest_FieldValidation(t *testing.T) {
	// Test that fields are properly parsed
	requestBody := `{
		"company_name": "kisanlink",
		"industry": "AgriTech / Smart Farming",
		"company_size": "11-50 employees",
		"recruiter_name": "Meher Prasad",
		"job_categories": ["AgriTech Development"],
		"hiring_locations": ["Hyderabad"],
		"hiring_types": []
	}`

	var req UpdateEmployerProfileRequest
	err := json.Unmarshal([]byte(requestBody), &req)
	assert.NoError(t, err)

	// Verify fields are parsed correctly
	assert.Equal(t, "kisanlink", req.CompanyName)
	assert.Equal(t, "AgriTech / Smart Farming", req.Industry)
	assert.Equal(t, "11-50 employees", req.CompanySize)
	assert.Equal(t, "Meher Prasad", req.RecruiterName)
	assert.Equal(t, []string{"AgriTech Development"}, req.JobCategories)
	assert.Equal(t, []string{"Hyderabad"}, req.HiringLocations)
	assert.Equal(t, []string{}, req.HiringTypes)
}

func TestUpdateEmployerProfileRequest_EmptyFields(t *testing.T) {
	// Test that empty fields are handled correctly
	requestBody := `{
		"company_name": "kisanlink",
		"industry": "",
		"company_size": "",
		"website_url": "",
		"recruiter_name": ""
	}`

	var req UpdateEmployerProfileRequest
	err := json.Unmarshal([]byte(requestBody), &req)
	assert.NoError(t, err)

	// Verify non-empty fields are parsed
	assert.Equal(t, "kisanlink", req.CompanyName)

	// Verify empty fields are empty strings
	assert.Equal(t, "", req.Industry)
	assert.Equal(t, "", req.CompanySize)
	assert.Equal(t, "", req.WebsiteURL)
	assert.Equal(t, "", req.RecruiterName)
}
