package app

import (
	"Road_services/internal/app/controllers"
	"Road_services/internal/app/repository"
	"Road_services/internal/app/role"
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "Road_services/docs"

	_ "github.com/lib/pq" // Для PostgreSQL
)

func (a *Application) StartServer() {
	log.Println("Server start up")
	r := gin.Default()
	// c := controllers.NewController(a.repository)
	repo, err := repository.New("user=bmstu_user password=bmstu_password dbname=bmstu host=localhost port=5432 sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to create repository: %v", err)
	}

	roadController := controllers.NewRoadController(repo)
	travelRequestController := controllers.NewTravelRequestController(repo)
	travelRequestRoadController := controllers.NewTravelRequestRoadController(repo)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	AuthGroup := r.Group("/auth")
	{
		AuthGroup.POST("/registration", a.Register)
		AuthGroup.POST("/login", a.Login)
		AuthGroup.Use(a.WithAuthCheck(role.Buyer, role.Manager, role.Admin)).GET("/logout", a.Logout)

	}

	RoadGroup := r.Group("/roads")
	{
		RoadGroup.Use(a.WithAuthCheck(role.Buyer, role.Manager, role.Admin)).GET("/", roadController.ListRoads)
		RoadGroup.GET("/:id", roadController.GetRoad)
		RoadGroup.Use(a.WithAuthCheck(role.Manager, role.Admin)).POST("/", roadController.CreateRoad)
		RoadGroup.Use(a.WithAuthCheck(role.Manager, role.Admin)).PUT("/:id", roadController.UpdateRoad)
		RoadGroup.Use(a.WithAuthCheck(role.Manager, role.Admin)).DELETE("/:id", roadController.DeleteRoad)
		// RoadGroup.POST("/add_road_to_last_request/:roadID", roadController.AddRoadToLastTravelRequest)
		RoadGroup.Use(a.WithAuthCheck(role.Buyer)).POST("/road_travel_request/:id", roadController.AddRoadToTravelRequest)
		RoadGroup.Use(a.WithAuthCheck(role.Manager, role.Admin)).PUT("/road_add_image/:id", roadController.AddRoadImage)
	}

	TravelRequestGroup := r.Group("/travelrequests")
	{
		TravelRequestGroup.Use(a.WithAuthCheck(role.Manager, role.Admin)).PUT("/change-status-moderator/:id", travelRequestController.ChangeRequestStatusByModerator)
		TravelRequestGroup.Use(a.WithAuthCheck(role.Buyer, role.Manager, role.Admin)).GET("/", travelRequestController.ListTravelRequests)
		TravelRequestGroup.Use(a.WithAuthCheck(role.Buyer)).GET("/introduced", travelRequestController.GetTravelRequestByID)
		TravelRequestGroup.Use(a.WithAuthCheck(role.Buyer)).PUT("/:id", travelRequestController.UpdateTravelRequest)
		// r.DELETE("/travelrequests/:id", travelRequestController.DeleteTravelRequest)
		TravelRequestGroup.Use(a.WithAuthCheck(role.Buyer)).PUT("/change-status-user/:id", travelRequestController.ChangeRequestStatusByUser)
		TravelRequestGroup.Use(a.WithAuthCheck(role.Buyer)).DELETE("/:id", travelRequestController.SoftDeleteTravelRequest)
	}

	TravelRequestRoadGroup := r.Group("/travelrequestroads")
	{
		TravelRequestRoadGroup.Use(a.WithAuthCheck(role.Buyer)).DELETE("/:requestID/:roadID", travelRequestRoadController.PhysicalDeleteRoadFromTravelRequest)
	}

	r.Run(":8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

	log.Println("Server down")
}

// repo, err := repository.New("user=bmstu_user password=bmstu_password dbname=bmstu host=localhost port=5432 sslmode=disable")
// if err != nil {
// 	log.Fatalf("Failed to create repository: %v", err)
// }

// roadController := controllers.NewRoadController(repo)
// r.GET("/roads", roadController.ListRoads)
// r.GET("/roads/:id", roadController.GetRoad)
// r.POST("/roads", roadController.CreateRoad)
// r.PUT("/roads/:id", roadController.UpdateRoad)
// r.DELETE("/roads/:id", roadController.DeleteRoad)
// r.POST("/add_road_to_last_request/:roadID", roadController.AddRoadToLastTravelRequest)
// r.POST("/road_travel_request/:roadID", roadController.AddRoadToTravelRequest)
// r.PUT("/road_add_image/:id", roadController.AddConsultationImage)

// travelRequestController := controllers.NewTravelRequestController(repo)
// r.GET("/travelrequests", travelRequestController.ListTravelRequests)
// r.GET("/travelrequests/:id", travelRequestController.GetTravelRequestByID)
// r.PUT("/travelrequests/:id", travelRequestController.UpdateTravelRequest)
// // r.DELETE("/travelrequests/:id", travelRequestController.DeleteTravelRequest)
// r.PUT("/travelrequests/change-status-user/:id", travelRequestController.ChangeRequestStatusByUser)
// r.PUT("/travelrequests/change-status-moderator/:id", travelRequestController.ChangeRequestStatusByModerator)
// r.DELETE("/travelrequests/:id", travelRequestController.SoftDeleteTravelRequest)

// travelRequestRoadController := controllers.NewTravelRequestRoadController(repo)
// r.DELETE("/travelrequestroads/:requestID/:roadID", travelRequestRoadController.PhysicalDeleteRoadFromTravelRequest)

// r.POST("/login", func(c *gin.Context) {
// 	a.Login(c)
// })

// r.POST("/registration", func(c *gin.Context) {
// 	a.Register(a.repository, c)
// })

// r.GET("/logout", func(c *gin.Context) {
// 	a.Logout(a.repository, c)
// })

// // r.Use(a.WithAuthCheck()).GET("/ping", func(c *gin.Context) {
// // 	c.JSON(http.StatusOK, gin.H{
// // 		"Status":  "OkNo",
// // 		"Message": "GG",
// // 	})
// // })
// // или ниженаписанное значит что доступ имеют менеджер и админ
// r.Use(a.WithAuthCheck(role.Manager, role.Admin)).GET("/ping", func(c *gin.Context) {
// 	c.JSON(http.StatusOK, gin.H{
// 		"Status":  "Ok",
// 		"Message": "GG",
// 	})
// })
