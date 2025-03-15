package controllers

import (
	"CertiBlock/application/backend/certiblock/base"
	"CertiBlock/application/backend/certiblock/base/data"
	"CertiBlock/application/backend/certiblock/services/students"
	"net/http"

	"github.com/gin-gonic/gin"
)

func StudentsAPI(context *base.ApplicationContext, r *gin.RouterGroup) {
	r.POST("", RegisterStudent(context))
	r.POST("/login", LoginStudent(context))
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
