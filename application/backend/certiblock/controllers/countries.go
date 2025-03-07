package controllers

import (
	"certiblock/base"
	"certiblock/base/data"
	"certiblock/services/countries"
	"net/http"
	"strconv"

	"github.com/samber/lo"

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
// @Success 200 {array} data.CountryOutput
// @Failure 500 {object} gin.H
// @Router /countries [get]
func GetCountries(context *base.ApplicationContext) func(c *gin.Context) {
	return func(c *gin.Context) {
		countries, err := countries.GetAll(context)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, lo.Map(countries, func(a data.Country, _ int) data.CountryOutput {
			return data.CountryOutputResponse(a)
		}))
	}
}

// GET /api/countries/:id
// @Tags countries
// @Summary Get a country by ID
// @Description Get a country by ID
// @Produce json
// @Param id path int true "Country ID"
// @Success 200 {object} data.CountryOutput
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Router /countries/{id} [get]
// GetCountryById handles the request to retrieve a country by its ID.
// It expects an ID parameter in the URL, converts it to an integer, and
// fetches the corresponding country from the database. If the ID is invalid
// or the country is not found, it returns an appropriate error response.
//
// Parameters:
// - context: A pointer to the base.ApplicationContext which holds the application context.
//
// Returns:
// - A function that takes a gin.Context and processes the request to get a country by ID.
//
// Responses:
// - 400 Bad Request: If the ID parameter is not a valid integer.
// - 404 Not Found: If no country is found with the given ID.
// - 200 OK: If the country is found, returns the country data in JSON format.
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

		c.JSON(http.StatusOK, data.CountryOutputResponse(*country))
	}
}
