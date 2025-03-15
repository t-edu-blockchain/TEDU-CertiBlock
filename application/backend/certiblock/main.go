package main

import (
	"CertiBlock/application/backend/certiblock/base"
	"CertiBlock/application/backend/certiblock/configurations"
	"CertiBlock/application/backend/certiblock/controllers"
	"CertiBlock/application/backend/certiblock/services/database"
	educert "CertiBlock/application/backend/gateway"
	"fmt"

	_ "CertiBlock/application/backend/certiblock/docs"

	"github.com/gin-contrib/cors"
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

	//////////////////////////////////////////////
	////////// BEGIN BLOCKCHAIN STUFF ////////////
	//////////////////////////////////////////////

	// Khởi tạo kết nối Fabric Gateway
	gw, conn, err := educert.NewGateway()
	if err != nil {
		panic(fmt.Errorf("failed to initialize gateway: %w", err))
	}
	defer gw.Close()
	defer conn.Close()

	// Kết nối tới channel và chaincode
	network := gw.GetNetwork("mychannel")
	context.Contract = network.GetContract("certicontract")

	// Khởi tạo ledger
	// result, err := educert.InitLedger(contract)
	// if err != nil {
	// 	fmt.Printf("Failed to initialize ledger: %v\n", err)
	// 	return
	// }
	// fmt.Println("Ledger initialized:", result)

	//////////////////////////////////////////////
	/////////// END BLOCKCHAIN STUFF /////////////
	//////////////////////////////////////////////

	router := gin.Default()

	// CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = false
	config.AllowOrigins = []string{"http://localhost:3001"}
	config.AllowMethods = []string{"POST", "GET", "PUT", "OPTIONS", "DELETE"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept", "User-Agent", "Cache-Control", "Pragma"}
	config.ExposeHeaders = []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"}
	config.AllowCredentials = true
	router.Use(cors.New(config))

	apiRouter := router.Group("/api")
	controllers.CountriesAPI(&context, apiRouter.Group("/countries"))
	controllers.StudentsAPI(&context, apiRouter.Group("/students"))
	controllers.UniversitiesAPI(&context, apiRouter.Group("/universities"))
	controllers.QRsAPI(&context, apiRouter.Group("/qrs"))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Println("=====================================")
	fmt.Println("View the API documentation at: http://localhost:3000/swagger/index.html")
	fmt.Println("Of course, change the domain name and port number appropriately if this is production environment.")
	fmt.Println("=====================================")
	fmt.Println()
	fmt.Println()
	fmt.Println()

	router.Run("0.0.0.0:3000")
}
