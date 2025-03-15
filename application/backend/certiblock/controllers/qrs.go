package controllers

import (
	"CertiBlock/application/backend/certiblock/base"
	"CertiBlock/application/backend/certiblock/base/data"
	"CertiBlock/application/backend/certiblock/services/qrs"
	"net/http"

	"github.com/gin-gonic/gin"
)

func QRsAPI(context *base.ApplicationContext, r *gin.RouterGroup) {
	r.POST("/create", CreateQR(context))
	r.GET("/verify/:token", VerifyQR(context))
}

// POST /api/qrs/create
// @Tags qrs
// @Summary Create a QR code
// @Description Create a QR code
// @Accept json
// @Produce json
// @Param qr body data.QRInput true "QR data"
// @Success 201 {object} data.QROutput
// @Failure 400 {object} gin.H
// @Router /api/qrs/create [post]
func CreateQR(context *base.ApplicationContext) func(c *gin.Context) {
	return func(c *gin.Context) {
		var qrInput data.QRInput
		if err := c.ShouldBindJSON(&qrInput); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		qr, err := qrs.CreateQR(context, &qrInput)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, qr)
	}
}

// GET /api/qrs/verify/:token
// @Tags qrs
// @Summary Verify a QR code
// @Description Verify a QR code
// @Accept json
// @Produce json
// @Param qr path string true "QR token"
// @Success 200 {object} string "The file content in base64"
// @Failure 400 {object} gin.H
// @Router /api/qrs/verify [post]
func VerifyQR(context *base.ApplicationContext) func(c *gin.Context) {
	return func(c *gin.Context) {
		token := c.Param("token")
		file, err := qrs.VerifyAndDecodeQR(context, token)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, file)
	}
}
