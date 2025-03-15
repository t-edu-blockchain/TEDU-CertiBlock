package controllers

import (
	"CertiBlock/application/backend/certiblock/base"
	"CertiBlock/application/backend/certiblock/base/data"
	"CertiBlock/application/backend/certiblock/services/certificates"
	"CertiBlock/application/backend/certiblock/services/students"
	"net/http"

	"github.com/gin-gonic/gin"
)

func StudentsAPI(context *base.ApplicationContext, r *gin.RouterGroup) {
	r.POST("", RegisterStudent(context))
	r.POST("/login", LoginStudent(context))
	r.POST("/certificates", GetAllCertificatesByStudent(context))
	r.POST("/certificates/:certUUID", GetOneCertificateOfStudent(context))
}

// POST /api/students
// @Tags students
// @Summary Register a student
// @Description Register a student
// @Accept json
// @Produce json
// @Param student body data.StudentInput true "Student data"
// @Success 201 {object} data.StudentOutput
// @Failure 400 {object} gin.H
// @Router /api/students [post]
func RegisterStudent(context *base.ApplicationContext) func(c *gin.Context) {
	return func(c *gin.Context) {
		var studentInput data.StudentInput
		if err := c.ShouldBindJSON(&studentInput); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		student, err := students.RegisterStudent(context, &studentInput)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, student)
	}
}

// POST /api/students/login
// @Tags students
// @Summary Login a student
// @Description Login a student
// @Accept json
// @Produce json
// @Param student body data.StudentInput true "Student data"
// @Success 200 {object} data.StudentOutput
// @Failure 400 {object} gin.H
// @Router /api/students/login [post]
func LoginStudent(context *base.ApplicationContext) func(c *gin.Context) {
	return func(c *gin.Context) {
		var studentInput data.StudentInput
		if err := c.ShouldBindJSON(&studentInput); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		student, err := students.LoginStudent(context, &studentInput)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, student)
	}
}

// POST /api/students/certificates
// @Tags students
// @Summary Get all certificates of a student
// @Description Get all certificates of a student
// @Accept json
// @Produce json
// @Param student body data.StudentAuth true "Student data"
// @Success 200 {object} []data.CertificateOutput
// @Failure 400 {object} gin.H
// @Router /api/students/certificates [post]
func GetAllCertificatesByStudent(context *base.ApplicationContext) func(c *gin.Context) {
	return func(c *gin.Context) {
		var studentAuth data.StudentAuth
		if err := c.ShouldBindJSON(&studentAuth); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		certificates, err := certificates.GetAllCertificatesByStudent(context, studentAuth.PrivateKey)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, certificates)
	}
}

// POST /api/students/certificates/:certUUID
// @Tags students
// @Summary Get one certificate of a student
// @Description Get one certificate of a student
// @Accept json
// @Produce json
// @Param student body data.StudentAuth true "Student data"
// @Param certUUID path string true "Certificate UUID"
// @Success 200 {object} data.CertificateOutputFull
// @Failure 400 {object} gin.H
// @Router /api/students/certificates/{certUUID} [post]
func GetOneCertificateOfStudent(context *base.ApplicationContext) func(c *gin.Context) {
	return func(c *gin.Context) {
		var studentAuth data.StudentAuth
		if err := c.ShouldBindJSON(&studentAuth); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		certUUID := c.Param("certUUID")
		certificate, err := certificates.GetOneCertificate(context, studentAuth.PrivateKey, certUUID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, certificate)
	}
}
