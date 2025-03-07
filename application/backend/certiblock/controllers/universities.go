package controllers

import (
	"certiblock/base"
	"certiblock/base/data"
	"certiblock/services/certificates"
	"certiblock/services/enrollment_certificates"
	"certiblock/services/universities"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UniversitiesAPI(context *base.ApplicationContext, r *gin.RouterGroup) {
	r.POST("", RegisterUniversity(context))
	r.POST("/info", GetInfo(context))
	r.POST("/enroll", EnrollStudent(context))
	r.POST("/certificate", IssueCertificate(context))
}

type PrivateKey struct {
	PrivateKey string `json:"private_key"`
}

// POST /api/universities/info
// @Tags universities
// @Summary Get university information
// @Description Get information about a university using its private key
// @Accept json
// @Produce json
// @Param private_key body PrivateKey true "Private key"
// @Success 200 {object} data.BCUniversityOutput
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Router /api/universities/info [post]
func GetInfo(context *base.ApplicationContext) func(c *gin.Context) {
	return func(c *gin.Context) {
		var pr PrivateKey
		if err := c.ShouldBindJSON(&pr); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		uni, err := universities.GetByPrivateKey(context, pr.PrivateKey)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, uni)
	}
}

// POST /api/universities
// @Tags universities
// @Summary Register a university
// @Description Register a university
// @Accept json
// @Produce json
// @Param university body data.BCUniversityInput true "University data"
// @Success 201 {object} data.BCUniversityOutput
// @Failure 400 {object} gin.H
// @Router /api/universities [post]
func RegisterUniversity(context *base.ApplicationContext) func(c *gin.Context) {
	return func(c *gin.Context) {
		var uni data.BCUniversityInput
		if err := c.ShouldBindJSON(&uni); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		uniOutput, err := universities.Register(context, uni)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, uniOutput)
	}
}

// POST /api/universities/enroll
// @Tags universities
// @Summary Enroll a student to a university
// @Description Enroll a student to a university
// @Accept json
// @Produce json
// @Param student_public_key body string true "Student public key"
// @Param hash body string true "CCCD hash"
// @Success 201 {object} data.BCEnrollmentCertificateOutput
// @Failure 400 {object} gin.H
// @Router /api/universities/enroll [post]
// EnrollStudent handles the enrollment of a student by processing the input JSON
// and issuing an enrollment certificate. It returns a handler function for a Gin
// HTTP request.
//
// The handler function performs the following steps:
// 1. Binds the incoming JSON request body to a BCEnrollmentCertificateInput struct.
// 2. If the binding fails, it responds with a 400 Bad Request status and an error message.
// 3. Calls the Issue function to generate an enrollment certificate.
// 4. If the Issue function returns an error, it responds with a 400 Bad Request status and an error message.
// 5. If successful, it responds with a 200 OK status and the issued enrollment certificate.
//
// Parameters:
// - context: A pointer to the ApplicationContext which provides necessary context for issuing the certificate.
//
// Returns:
// - A Gin handler function that processes the enrollment request.
func EnrollStudent(context *base.ApplicationContext) func(c *gin.Context) {
	return func(c *gin.Context) {
		var enrollmentCertificateInput data.BCEnrollmentCertificateInput
		if err := c.ShouldBindJSON(&enrollmentCertificateInput); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		ec, err := enrollment_certificates.Issue(context, enrollmentCertificateInput)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, ec)
	}
}

// POST /api/universities/certificate
// @Tags universities
// @Summary Issue a certificate
// @Description Issue a certificate
// @Accept json
// @Produce json
// @Param certificate body data.BCCertificateInput true "Certificate data"
// @Success 201 {object} data.BCCertificateOutput
// @Failure 400 {object} gin.H
// @Router /api/universities/certificate [post]
func IssueCertificate(context *base.ApplicationContext) func(c *gin.Context) {
	return func(c *gin.Context) {
		var certificateInput data.BCCertificateInput
		if err := c.ShouldBindJSON(&certificateInput); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		certificate, err := certificates.Issue(context, certificateInput)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, certificate)
	}
}
