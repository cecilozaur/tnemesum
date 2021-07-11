package app

import (
	"fmt"
	"github.com/cecilozaur/tnemesum/domain"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

var (
	invalidCity = gin.H{
		"error": "invalid cityId requested",
	}

	invalidDateFormat = gin.H{
		"error": "invalid date format, should be YYYY-mm-dd",
	}
)

// GetAllCities
// GET /api/cities
// returns a list of all the cities we have inside the repository
func (a *App) GetAllCities(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"items": a.repo.GetItems(),
	})
}

// GetCity
// GET /api/cities/{cityId}
// returns the city object or an err in case the city is not found
func (a *App) GetCity(c *gin.Context) {
	cityId, err := strconv.Atoi(c.Param("cityId"))
	if err != nil || cityId == 0 {
		c.JSON(http.StatusBadRequest, invalidCity)
		return
	}

	city, err := a.repo.Get(uint64(cityId))
	if err != nil {
		c.JSON(http.StatusNotFound, "")
		return
	}

	c.JSON(http.StatusOK, city)
}

// GetCityForecast
// GET /api/cities/{cityId}/forecast/{day}
// returns the forecast for the specified city and day
// if day is not defined the function returns forecast for all days available
func (a *App) GetCityForecast(c *gin.Context) {
	cityId, err := strconv.Atoi(c.Param("cityId"))
	if err != nil || cityId == 0 {
		c.JSON(http.StatusBadRequest, invalidCity)
		return
	}

	forecast, err := a.repo.GetForecast(uint64(cityId))
	if err != nil {
		c.JSON(http.StatusNotFound, "")
		return
	}

	day := c.Param("day")
	// get all forecast by default
	if day == "" {
		c.JSON(http.StatusOK, gin.H{
			"forecast": forecast,
		})
		return
	}

	_, err = time.Parse("2006-01-02", day)
	if err != nil {
		c.JSON(http.StatusBadRequest, invalidDateFormat)
		return
	}

	for _, f := range forecast {
		if day == f.Day {
			c.JSON(http.StatusOK, gin.H{
				"forecast": f,
			})

			return
		}
	}

	// we didn't find the forecast for the requested day
	c.JSON(http.StatusNotFound, "")
}

// UpdateCityForecast
// POST /api/cities/{cityId}/forecast
// updates the forecast for specified city for whatever dates are sent in the request
func (a *App) UpdateCityForecast(c *gin.Context) {
	req := domain.Forecast{}
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})

		return
	}

	cityId, err := strconv.Atoi(c.Param("cityId"))
	if err != nil || cityId == 0 {
		c.JSON(http.StatusBadRequest, invalidCity)
		return
	}

	_, err = a.repo.Get(uint64(cityId))
	if err != nil {
		c.JSON(http.StatusNotFound, "")
		return
	}

	if !a.repo.UpdateForecast(uint64(cityId), req) {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": "unable to update city: " + c.Param("cityId"),
		})

		return
	}

	c.Header("Location", fmt.Sprintf("/api/cities/%d/forecast", cityId))
	c.JSON(http.StatusNoContent, "")
}
