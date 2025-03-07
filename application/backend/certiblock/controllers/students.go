package controllers

import (
	"certiblock/base"
	"certiblock/base/data"
	"certiblock/services/students"
	"net/http"

	"github.com/gin-gonic/gin"
)

func StudentsAPI(context *base.ApplicationContext, r *gin.RouterGroup) {
	r.GET("/:public_key", GetStudentByPublicKey(context))
	r.POST("", RegisterStudent(context))
}

// GET /api/students/:public_key
// @Tags students
// @Summary Get a student by public key
// @Description Get a student by public key
// @Produce json
// @Param public_key path string true "Student public key"
// @Success 200 {object} data.BCStudentOutput
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Router /students/{id} [get]
func GetStudentByPublicKey(context *base.ApplicationContext) func(c *gin.Context) {
	return func(c *gin.Context) {
		publicKey := c.Param("public_key")

		student, err := students.GetByPublicKey(context, publicKey)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Student not found",
			})
			return
		}

		c.JSON(http.StatusOK, data.BCStudentOutputResponse(student))
	}
}

// POST /api/students
// @Tags students
// @Summary Register a student
// @Description Register a student
// @Accept json
// @Produce json
// @Param student body data.BCStudentInput true "Student data"
// @Success 201 {object} data.BCStudentOutput
// @Failure 400 {object} gin.H
// @Router /students [post]
// RegisterStudent handles the registration of a new student.
// It expects a JSON payload representing the student input data.
// If the input data is invalid, it responds with a 400 Bad Request status and an error message.
// If the registration is successful, it responds with a 201 Created status and the student output data.
//
// Parameters:
// - context: A pointer to the base.ApplicationContext.
//
// Returns:
// - A Gin handler function that processes the student registration request.
func RegisterStudent(context *base.ApplicationContext) func(c *gin.Context) {
	return func(c *gin.Context) {
		var studentInput data.BCStudentInput
		if err := c.ShouldBindJSON(&studentInput); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		student, err := students.Register(context, studentInput)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, data.BCStudentOutputResponse(student))
	}
}
