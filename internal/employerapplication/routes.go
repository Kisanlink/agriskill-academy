package employerapplication

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *EmployerApplicationHandler) {
	fmt.Printf("DEBUG: Registering employer application routes\n")

	rg.GET("/employer/jobs/:jobId/applications", handler.GetApplicationsForJob)
	rg.GET("/employer/jobs/:jobId/applications/debug", handler.DebugApplications)
	rg.PUT("/employer/applications/:applicationId/status", handler.UpdateStatus)
	rg.GET("/employer/applicants/:studentId/profile", handler.GetApplicantProfile)
	rg.POST("/employer/applications/:applicationId/message", handler.SendMessage)
	rg.GET("/employer/applications/:applicationId/messages", handler.GetMessages)

	// Student-side
	rg.GET("/student/applications", handler.GetApplicationsByStudent)
	rg.GET("/student/applications/:applicationId/messages", handler.GetMessages)
	rg.POST("/student/applications/:applicationId/message", handler.SendMessage)

	fmt.Printf("DEBUG: Employer application routes registered successfully\n")
}
