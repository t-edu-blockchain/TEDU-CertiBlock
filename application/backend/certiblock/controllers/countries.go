package controllers

import (
	"CertiBlock/application/backend/certiblock/base"
	"CertiBlock/application/backend/certiblock/services/countries"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CountriesAPI(context *base.ApplicationContext, r *gin.RouterGroup) {
	r.GET("", GetCountries(context))
}

// GET /api/countries
// @Tags countries
// @Summary Get all countries
// @Description Get all countries
// @Produce json
// @Success 200 {array} data.Country
// @Failure 500 {object} gin.H
// @Router /api/countries [get]
func GetCountries(context *base.ApplicationContext) func(c *gin.Context) {
	return func(c *gin.Context) {
		countries, err := countries.GetAll(context)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, countries)
	}
}

// GET /api/countries/:id
// @Tags countries
// @Summary Get a country by ID
// @Description Get a country by ID
// @Produce json
// @Param id path int true "Country ID"
// @Success 200 {object} data.Country
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Router /api/countries/{id} [get]
func GetCountryById(context *base.ApplicationContext) func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		ID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid ID",
			})
			return
		}

		country, err := countries.GetById(context, ID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Country not found",
			})
			return
		}

		c.JSON(http.StatusOK, country)
	}
}
