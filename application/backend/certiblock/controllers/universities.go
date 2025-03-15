package controllers

import (
	"CertiBlock/application/backend/certiblock/base"
	"CertiBlock/application/backend/certiblock/base/data"
	"CertiBlock/application/backend/certiblock/services/universities"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UniversitiesAPI(context *base.ApplicationContext, r *gin.RouterGroup) {
	r.POST("/certificate-file", AddCertificateFile(context))
}

// POST /api/universities/certificate-file
// @Tags universities
// @Summary Add a certificate file
// @Description Add a certificate file
// @Accept json
// @Produce json
// @Param certificateFile body data.CertificateFileUpload true "Certificate file data"
// @Success 201 {object} gin.H
// @Failure 400 {object} gin.H
// @Router /api/universities/certificate-file [post]
func AddCertificateFile(context *base.ApplicationContext) func(c *gin.Context) {
	return func(c *gin.Context) {
		var certificateFileUpload data.CertificateFileUpload
		if err := c.ShouldBindJSON(&certificateFileUpload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		_, err := universities.SaveCertificateFile(context, &certificateFileUpload)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusCreated, gin.H{})
	}
}
