package notification

import (
	"bytes"
	"fmt"
	"html/template"
)

type EmailTemplateService struct{}

func NewEmailTemplateService() *EmailTemplateService {
	return &EmailTemplateService{}
}

func (s *EmailTemplateService) RenderNewJobEmail(jobData map[string]interface{}) (string, error) {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>New Job Posted</title>
	<style>
		body {
			font-family: Arial, sans-serif;
			line-height: 1.6;
			color: #333;
			max-width: 600px;
			margin: 0 auto;
			padding: 20px;
		}
		.header {
			background-color: #4CAF50;
			color: white;
			padding: 20px;
			text-align: center;
			border-radius: 5px 5px 0 0;
		}
		.content {
			background-color: #f9f9f9;
			padding: 20px;
			border: 1px solid #ddd;
		}
		.job-details {
			background-color: white;
			padding: 15px;
			margin: 15px 0;
			border-left: 4px solid #4CAF50;
		}
		.button {
			display: inline-block;
			padding: 12px 24px;
			background-color: #4CAF50;
			color: white;
			text-decoration: none;
			border-radius: 5px;
			margin-top: 20px;
		}
		.footer {
			text-align: center;
			padding: 20px;
			color: #666;
			font-size: 12px;
		}
	</style>
</head>
<body>
	<div class="header">
		{{if .LogoURL}}
		<div style="text-align: center; margin-bottom: 15px;">
			<img src="{{.LogoURL}}" alt="AgriJobs Logo" style="max-width: 200px; height: auto; display: block; margin: 0 auto;" />
		</div>
		{{end}}
		<h2>New Job Opportunity!</h2>
	</div>
	<div class="content">
		<p>Hi {{.StudentName}},</p>
		<p>A new job has been posted that might interest you:</p>
		<div class="job-details">
			<h3>{{.JobTitle}}</h3>
			<p><strong>Company:</strong> {{.Company}}</p>
			<p><strong>Location:</strong> {{.Location}}</p>
			<p><strong>Job Type:</strong> {{.JobType}}</p>
			<p><strong>Experience:</strong> {{.Experience}}</p>
			{{if .Salary}}
			<p><strong>Salary:</strong> {{.Salary}}</p>
			{{end}}
			{{if .Description}}
			<p><strong>Description:</strong> {{.Description}}</p>
			{{end}}
		</div>
		<a href="{{.JobLink}}" class="button">View Full Details</a>
	</div>
	<div class="footer">
		<p>You're receiving this email because you have job alerts enabled in your notification preferences.</p>
		<p>You can manage your preferences in your account settings.</p>
	</div>
</body>
</html>`

	t, err := template.New("newjob").Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, jobData); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

func (s *EmailTemplateService) RenderStatusUpdateEmail(appData map[string]interface{}) (string, error) {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>Application Status Update</title>
	<style>
		body {
			font-family: Arial, sans-serif;
			line-height: 1.6;
			color: #333;
			max-width: 600px;
			margin: 0 auto;
			padding: 20px;
		}
		.header {
			background-color: #2196F3;
			color: white;
			padding: 20px;
			text-align: center;
			border-radius: 5px 5px 0 0;
		}
		.content {
			background-color: #f9f9f9;
			padding: 20px;
			border: 1px solid #ddd;
		}
		.status-box {
			background-color: white;
			padding: 15px;
			margin: 15px 0;
			border-left: 4px solid #2196F3;
		}
		.status {
			font-size: 18px;
			font-weight: bold;
			color: #2196F3;
			margin: 10px 0;
		}
		.button {
			display: inline-block;
			padding: 12px 24px;
			background-color: #2196F3;
			color: white;
			text-decoration: none;
			border-radius: 5px;
			margin-top: 20px;
		}
		.footer {
			text-align: center;
			padding: 20px;
			color: #666;
			font-size: 12px;
		}
	</style>
</head>
<body>
	<div class="header">
		{{if .LogoURL}}
		<div style="text-align: center; margin-bottom: 15px;">
			<img src="{{.LogoURL}}" alt="AgriJobs Logo" style="max-width: 200px; height: auto; display: block; margin: 0 auto;" />
		</div>
		{{end}}
		<h2>Application Status Update</h2>
	</div>
	<div class="content">
		<p>Hi {{.StudentName}},</p>
		<p>Your application status has been updated:</p>
		<div class="status-box">
			<p><strong>Job:</strong> {{.JobTitle}} at {{.Company}}</p>
			<p class="status">New Status: {{.Status}}</p>
			{{if .StatusMessage}}
			<p>{{.StatusMessage}}</p>
			{{end}}
		</div>
		<a href="{{.ApplicationLink}}" class="button">View Application</a>
	</div>
	<div class="footer">
		<p>You're receiving this email because you have application updates enabled in your notification preferences.</p>
		<p>You can manage your preferences in your account settings.</p>
	</div>
</body>
</html>`

	t, err := template.New("statusupdate").Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, appData); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}
