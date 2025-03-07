package main

import (
	"certiblock/base"
	"certiblock/configurations"
	"certiblock/controllers"
	"certiblock/services/database"

	_ "certiblock/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	var err error

	context := base.ApplicationContext{}
	context.Config, err = configurations.Load()
	if err != nil {
		panic(err)
	}

	err = database.InitConnection(&context)
	if err != nil {
		panic(err)
	}
	defer context.DB.Close()

	router := gin.Default()
	apiRouter := router.Group("/api")
	controllers.StudentsAPI(&context, apiRouter.Group("/students"))
	controllers.CountriesAPI(&context, apiRouter.Group("/countries"))
	controllers.UniversitiesAPI(&context, apiRouter.Group("/universities"))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run("0.0.0.0:3000")
}
