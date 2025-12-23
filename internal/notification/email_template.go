package notification

import (
	"bytes"
	"fmt"
	"html/template"
	"time"
)

type EmailTemplateService struct{}

func NewEmailTemplateService() *EmailTemplateService {
	return &EmailTemplateService{}
}

// RenderNewJobEmail generates a visually engaging email for new job alerts
func (s *EmailTemplateService) RenderNewJobEmail(jobData map[string]interface{}) (string, error) {
	// Add current year for copyright if not present
	if _, ok := jobData["CurrentYear"]; !ok {
		jobData["CurrentYear"] = time.Now().Year()
	}

	tmpl := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>New Job Opportunity - AgriSkill Academy</title>
    <style>
        /* Reset and Client Specific Styles */
        body { margin: 0; padding: 0; -webkit-text-size-adjust: 100%; -ms-text-size-adjust: 100%; background-color: #f4f7f6; }
        table, td { border-collapse: collapse; mso-table-lspace: 0pt; mso-table-rspace: 0pt; }
        img { border: 0; height: auto; line-height: 100%; outline: none; text-decoration: none; -ms-interpolation-mode: bicubic; }
        
        /* Typography */
        body, td, th { font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; color: #333333; }
        
        /* AgriSkill Palette */
        .primary-color { color: #1b5e20; } /* Deep Green */
        .bg-primary { background-color: #1b5e20; }
        .accent-bg { background-color: #e8f5e9; } /* Light Green */
        .button-color { background-color: #2e7d32; }
    </style>
</head>
<body style="margin: 0; padding: 0; background-color: #f4f7f6;">
    <div style="background-color: #f4f7f6;">
        <div style="display: none; max-height: 0px; overflow: hidden;">
            A new {{.JobTitle}} role at {{.Company}} matches your profile.
        </div>

        <table role="presentation" width="100%" border="0" cellspacing="0" cellpadding="0">
            <tr>
                <td align="center" style="padding: 20px 0;">
                    
                    <table role="presentation" width="600" border="0" cellspacing="0" cellpadding="0" style="width: 100%; max-width: 600px; background-color: #ffffff; border-radius: 8px; box-shadow: 0 4px 6px rgba(0,0,0,0.05); overflow: hidden;">
                        
                        <tr>
                            <td align="center" style="background-color: #1b5e20; padding: 30px 20px; background-image: linear-gradient(135deg, #1b5e20 0%, #2e7d32 100%);">
                                {{if .LogoURL}}
                                <img src="{{.LogoURL}}" alt="AgriSkill Academy" width="150" style="display: block; width: 150px; max-width: 100%; margin-bottom: 10px;">
                                {{else}}
                                <h1 style="color: #ffffff; margin: 0; font-size: 24px; font-weight: 700; letter-spacing: 1px;">AGRISKILL <span style="font-weight: 300;">ACADEMY</span></h1>
                                {{end}}
                                <p style="color: #e8f5e9; margin: 10px 0 0 0; font-size: 14px; text-transform: uppercase; letter-spacing: 1px;">New Job Alert</p>
                            </td>
                        </tr>

                        <tr>
                            <td style="padding: 40px 30px 20px 30px;">
                                <h2 style="color: #333333; margin: 0 0 15px 0; font-size: 22px;">Hello {{.StudentName}},</h2>
                                <p style="font-size: 16px; line-height: 1.6; color: #555555; margin: 0;">
                                    We found a new opportunity that matches your skills in the agricultural sector.
                                </p>
                            </td>
                        </tr>

                        <tr>
                            <td style="padding: 0 30px;">
                                <table role="presentation" width="100%" border="0" cellspacing="0" cellpadding="0" style="background-color: #f8fcf8; border: 1px solid #c8e6c9; border-radius: 8px;">
                                    <tr>
                                        <td style="padding: 25px;">
                                            <h3 style="margin: 0 0 5px 0; color: #1b5e20; font-size: 20px;">{{.JobTitle}}</h3>
                                            <p style="margin: 0 0 20px 0; color: #455a64; font-size: 16px; font-weight: 600;">{{.Company}}</p>
                                            
                                            <table role="presentation" width="100%" border="0" cellspacing="0" cellpadding="0">
                                                <tr>
                                                    <td width="50%" style="padding-bottom: 15px; vertical-align: top;">
                                                        <p style="margin: 0; font-size: 12px; color: #888888; text-transform: uppercase;">Location</p>
                                                        <p style="margin: 3px 0 0 0; font-size: 14px; color: #333333;"><strong>{{.Location}}</strong></p>
                                                    </td>
                                                    <td width="50%" style="padding-bottom: 15px; vertical-align: top;">
                                                        <p style="margin: 0; font-size: 12px; color: #888888; text-transform: uppercase;">Job Type</p>
                                                        <p style="margin: 3px 0 0 0; font-size: 14px; color: #333333;"><strong>{{.JobType}}</strong></p>
                                                    </td>
                                                </tr>
                                                <tr>
                                                    <td width="50%" style="padding-bottom: 15px; vertical-align: top;">
                                                        <p style="margin: 0; font-size: 12px; color: #888888; text-transform: uppercase;">Experience</p>
                                                        <p style="margin: 3px 0 0 0; font-size: 14px; color: #333333;"><strong>{{.Experience}}</strong></p>
                                                    </td>
                                                    {{if .Salary}}
                                                    <td width="50%" style="padding-bottom: 15px; vertical-align: top;">
                                                        <p style="margin: 0; font-size: 12px; color: #888888; text-transform: uppercase;">Salary</p>
                                                        <p style="margin: 3px 0 0 0; font-size: 14px; color: #333333;"><strong>{{.Salary}}</strong></p>
                                                    </td>
                                                    {{end}}
                                                </tr>
                                            </table>

                                            {{if .Description}}
                                            <div style="border-top: 1px solid #e0e0e0; margin-top: 10px; padding-top: 15px;">
                                                <p style="margin: 0; font-size: 14px; line-height: 1.6; color: #666666;">
                                                    {{.Description}}
                                                </p>
                                            </div>
                                            {{end}}
                                        </td>
                                    </tr>
                                </table>
                            </td>
                        </tr>

                        <tr>
                            <td align="center" style="padding: 35px 30px;">
                                <table role="presentation" border="0" cellspacing="0" cellpadding="0">
                                    <tr>
                                        <td align="center" bgcolor="#2e7d32" style="border-radius: 6px;">
                                            <a href="{{.JobLink}}" target="_blank" style="font-size: 16px; font-weight: bold; color: #ffffff; text-decoration: none; padding: 14px 40px; border: 1px solid #2e7d32; display: inline-block; border-radius: 6px; font-family: sans-serif;">Apply Now</a>
                                        </td>
                                    </tr>
                                </table>
                            </td>
                        </tr>

                        <tr>
                            <td style="background-color: #f4f7f6; padding: 30px; text-align: center; border-top: 1px solid #e0e0e0;">
                                <p style="margin: 0 0 10px 0; font-size: 12px; color: #888888;">
                                    &copy; {{.CurrentYear}} AgriSkill Academy. Cultivating Careers.
                                </p>
                                <p style="margin: 0; font-size: 12px; color: #888888;">
                                    You received this email because you subscribed to job alerts.<br>
                                    <a href="{{.UnsubscribeURL}}" style="color: #2e7d32; text-decoration: underline;">Unsubscribe</a> | <a href="{{.ManagePreferencesURL}}" style="color: #2e7d32; text-decoration: underline;">Manage Preferences</a>
                                </p>
                            </td>
                        </tr>

                    </table>
                    </td>
            </tr>
        </table>
    </div>
</body>
</html>
`

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

// RenderStatusUpdateEmail generates a professional status update with company branding
func (s *EmailTemplateService) RenderStatusUpdateEmail(appData map[string]interface{}) (string, error) {
	// Add current year for copyright if not present
	if _, ok := appData["CurrentYear"]; !ok {
		appData["CurrentYear"] = time.Now().Year()
	}

	tmpl := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Application Status Update</title>
</head>
<body style="margin: 0; padding: 0; background-color: #f0f2f5; font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;">
    
    <table role="presentation" width="100%" border="0" cellspacing="0" cellpadding="0">
        <tr>
            <td align="center" style="padding: 20px 0;">
                
                <table role="presentation" width="600" border="0" cellspacing="0" cellpadding="0" style="width: 100%; max-width: 600px; background-color: #ffffff; border-radius: 12px; box-shadow: 0 8px 24px rgba(0,0,0,0.08); overflow: hidden;">
                    
                    <tr>
                        <td style="padding: 20px 30px; border-bottom: 2px solid #f0f0f0;">
                            <table role="presentation" width="100%" border="0" cellspacing="0" cellpadding="0">
                                <tr>
                                    <td align="left">
                                        {{if .LogoURL}}
                                        <img src="{{.LogoURL}}" alt="AgriSkill Academy" width="120" style="display: block; width: 120px;">
                                        {{else}}
                                        <span style="color: #1b5e20; font-weight: bold; font-size: 18px;">AGRISKILL</span>
                                        {{end}}
                                    </td>
                                    <td align="right" style="color: #888888; font-size: 12px; font-weight: 500; text-transform: uppercase;">
                                        Application Update
                                    </td>
                                </tr>
                            </table>
                        </td>
                    </tr>

                    <tr>
                        <td align="center" style="padding: 40px 30px;">
                            <h1 style="margin: 0 0 10px 0; color: #333333; font-size: 24px;">Hi {{.StudentName}},</h1>
                            <p style="margin: 0 0 30px 0; color: #666666; font-size: 16px;">The status of your application has changed.</p>

                            <div style="background-color: #e3f2fd; border: 1px solid #bbdefb; color: #1565c0; display: inline-block; padding: 12px 30px; border-radius: 50px; font-size: 18px; font-weight: 700; letter-spacing: 0.5px; box-shadow: 0 2px 4px rgba(0,0,0,0.05);">
                                {{.Status}}
                            </div>
                        </td>
                    </tr>

                    <tr>
                        <td style="padding: 0 30px 40px 30px;">
                            <table role="presentation" width="100%" border="0" cellspacing="0" cellpadding="0" style="background-color: #fafafa; border-radius: 8px; border: 1px solid #eeeeee;">
                                <tr>
                                    <td style="padding: 25px;">
                                        
                                        <table role="presentation" width="100%" border="0" cellspacing="0" cellpadding="0">
                                            <tr>
                                                {{if .CompanyLogo}}
                                                <td width="60" style="padding-right: 20px; vertical-align: middle;">
                                                    <img src="{{.CompanyLogo}}" alt="{{.CompanyName}}" width="60" height="60" style="width: 60px; height: 60px; border-radius: 8px; object-fit: cover; border: 1px solid #dddddd;">
                                                </td>
                                                {{end}}
                                                <td style="vertical-align: middle;">
                                                    <h3 style="margin: 0 0 5px 0; color: #333333; font-size: 18px;">{{.JobTitle}}</h3>
                                                    <p style="margin: 0; color: #1b5e20; font-weight: 600;">{{.CompanyName}}</p>
                                                    {{if .CompanyWebsite}}
                                                    <a href="{{.CompanyWebsite}}" style="font-size: 12px; color: #666666; text-decoration: none; margin-top: 5px; display: inline-block;">Visit Website &rarr;</a>
                                                    {{end}}
                                                </td>
                                            </tr>
                                        </table>

                                        {{if .StatusMessage}}
                                        <div style="margin-top: 20px; padding-top: 20px; border-top: 1px dashed #dddddd;">
                                            <p style="margin: 0 0 5px 0; font-size: 12px; color: #888888; text-transform: uppercase; font-weight: 700;">Message from Employer:</p>
                                            <p style="margin: 0; font-size: 15px; color: #444444; font-style: italic; line-height: 1.5;">
                                                "{{.StatusMessage}}"
                                            </p>
                                        </div>
                                        {{end}}

                                    </td>
                                </tr>
                            </table>
                        </td>
                    </tr>

                    <tr>
                        <td align="center" style="padding-bottom: 40px;">
                            <a href="{{.ApplicationLink}}" style="background-color: #1b5e20; color: #ffffff; padding: 15px 35px; text-decoration: none; border-radius: 6px; font-weight: bold; font-size: 16px; display: inline-block; transition: background 0.3s;">View Application Status</a>
                        </td>
                    </tr>

                    <tr>
                        <td style="background-color: #37474f; padding: 30px; text-align: center; color: #cfd8dc; font-size: 13px;">
                            <p style="margin: 0 0 10px 0;">&copy; {{.CurrentYear}} AgriSkill Academy</p>
                            <p style="margin: 0 0 10px 0; font-size: 12px;">
                                <a href="{{.UnsubscribeURL}}" style="color: #eceff1; text-decoration: underline;">Unsubscribe</a> | 
                                <a href="{{.ManagePreferencesURL}}" style="color: #eceff1; text-decoration: underline;">Manage Preferences</a>
                            </p>
                            <div style="margin-top: 10px;">
                                <a href="#" style="color: #eceff1; text-decoration: none; margin: 0 10px;">Help Center</a>
                                <a href="#" style="color: #eceff1; text-decoration: none; margin: 0 10px;">Privacy Policy</a>
                            </div>
                        </td>
                    </tr>

                </table>
            </td>
        </tr>
    </table>
</body>
</html>
`

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
