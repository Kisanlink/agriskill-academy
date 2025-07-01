package employerapplication

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, handler *EmployerApplicationHandler) {
	rg.GET("/employer/jobs/:jobId/applications", handler.GetApplicationsForJob)
	rg.PUT("/employer/applications/:applicationId/status", handler.UpdateStatus)
	rg.GET("/employer/applicants/:studentId/profile", handler.GetApplicantProfile)
	rg.POST("/employer/applications/:applicationId/message", handler.SendMessage)
	rg.GET("/employer/applications/:applicationId/messages", handler.GetMessages)

	// Student-side
	rg.GET("/student/applications", handler.GetApplicationsByStudent)
}
