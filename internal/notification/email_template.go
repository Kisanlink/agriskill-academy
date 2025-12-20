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
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Application Status Update</title>
</head>
<body style="margin: 0; padding: 0; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; line-height: 1.6; background-color: #f4f6f8; color: #333333;">
    <!-- Main Container -->
    <div style="max-width: 600px; margin: 0 auto; background-color: #ffffff; overflow: hidden; box-shadow: 0 4px 6px rgba(0,0,0,0.05);">
        
        <!-- Header with AgriJobs Logo -->
        <div style="background: linear-gradient(135deg, #2196F3 0%, #1976D2 100%); padding: 35px 20px; text-align: center;">
            {{if .LogoURL}}
            <img src="{{.LogoURL}}" alt="AgriJobs" style="max-width: 180px; height: auto; display: block; margin: 0 auto 15px auto; background-color: rgba(255,255,255,0.1); padding: 10px; border-radius: 8px;">
            {{else}}
            <h1 style="color: #ffffff; margin: 0; font-size: 28px; font-weight: 600;">AgriJobs</h1>
            {{end}}
            <p style="color: rgba(255,255,255,0.9); margin: 10px 0 0 0; font-size: 14px;">Application Status Update</p>
        </div>

        <!-- Content -->
        <div style="padding: 40px 30px;">
            
            <h2 style="color: #2c3e50; margin-top: 0; font-size: 24px; text-align: center; font-weight: 600;">Hi {{.StudentName}}!</h2>
            
            <p style="font-size: 16px; color: #555555; text-align: center; margin-bottom: 30px;">There has been an update to your job application.</p>

            <!-- Company Card with Large Logo -->
            <div style="background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%); border: 2px solid #dee2e6; border-radius: 12px; padding: 30px 20px; margin: 30px 0; text-align: center; box-shadow: 0 2px 8px rgba(0,0,0,0.08);">
                <!-- Company Logo - Larger and More Prominent -->
                {{if .CompanyLogo}}
                <div style="margin-bottom: 20px;">
                    <img src="{{.CompanyLogo}}" alt="{{.CompanyName}}" style="max-width: 120px; max-height: 120px; width: auto; height: auto; border-radius: 12px; margin: 0 auto; display: block; background-color: white; padding: 8px; box-shadow: 0 4px 12px rgba(0,0,0,0.15); object-fit: contain;">
                </div>
                {{end}}
                
                <h3 style="margin: 0 0 8px 0; color: #2c3e50; font-size: 20px; font-weight: 600;">{{.JobTitle}}</h3>
                <p style="margin: 0 0 15px 0; color: #495057; font-size: 16px;">
                    <strong style="color: #2196F3;">{{.CompanyName}}</strong>
                </p>
                
                {{if .CompanyWebsite}}
                <div style="margin-top: 15px;">
                    <a href="{{.CompanyWebsite}}" target="_blank" style="display: inline-block; color: #2196F3; text-decoration: none; font-size: 14px; font-weight: 500; padding: 8px 16px; border: 1px solid #2196F3; border-radius: 6px; background-color: white;">
                        🌐 Visit Company Website →
                    </a>
                </div>
                {{end}}
            </div>

            <!-- Status Section -->
            <div style="text-align: center; margin: 35px 0;">
                <p style="margin-bottom: 12px; font-size: 13px; text-transform: uppercase; letter-spacing: 1.5px; color: #6c757d; font-weight: 600;">Current Status</p>
                <div style="display: inline-block; padding: 14px 35px; background: linear-gradient(135deg, #e3f2fd 0%, #bbdefb 100%); color: #1565c0; border-radius: 50px; font-weight: bold; font-size: 20px; text-transform: capitalize; border: 2px solid #90caf9; box-shadow: 0 2px 8px rgba(33, 150, 243, 0.2);">
                    {{.Status}}
                </div>
            </div>

            <!-- Message (if available) -->
            {{if .StatusMessage}}
            <div style="background-color: #f8f9fa; border-left: 5px solid #2196F3; padding: 20px; margin: 30px 0; border-radius: 4px;">
                <p style="margin: 0; font-style: italic; color: #495057; font-size: 15px; line-height: 1.6;">"{{.StatusMessage}}"</p>
            </div>
            {{end}}

            <!-- CTA Button -->
            <div style="text-align: center; margin-top: 40px;">
                <a href="{{.ApplicationLink}}" style="display: inline-block; background: linear-gradient(135deg, #2196F3 0%, #1976D2 100%); color: #ffffff; padding: 16px 40px; border-radius: 8px; text-decoration: none; font-weight: bold; font-size: 16px; box-shadow: 0 4px 12px rgba(33, 150, 243, 0.4); transition: all 0.3s ease;">
                    View Application Details →
                </a>
            </div>

        </div>

        <!-- Footer -->
        <div style="background-color: #2c3e50; padding: 35px 30px; text-align: center; color: #95a5a6; font-size: 13px;">
            <p style="margin: 0 0 12px 0; color: #bdc3c7;">&copy; {{.CurrentYear}} Agriskill Academy. All rights reserved.</p>
            <p style="margin: 0 0 20px 0; color: #95a5a6; font-size: 12px;">You received this email because you applied for a job on Agriskill Academy.</p>
            <div style="margin-top: 20px; padding-top: 20px; border-top: 1px solid #34495e;">
                <a href="#" style="color: #7f8c8d; text-decoration: none; margin: 0 15px; font-size: 12px;">Privacy Policy</a>
                <span style="color: #34495e;">|</span>
                <a href="#" style="color: #7f8c8d; text-decoration: none; margin: 0 15px; font-size: 12px;">Help Center</a>
                <span style="color: #34495e;">|</span>
                <a href="#" style="color: #7f8c8d; text-decoration: none; margin: 0 15px; font-size: 12px;">Contact Us</a>
            </div>
        </div>
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
