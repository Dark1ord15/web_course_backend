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

	AuthGroup := r.Group("/auth")
	{
		AuthGroup.POST("/registration", a.Register)
		AuthGroup.POST("/login", a.Login)
		AuthGroup.GET("/logout", a.WithAuthCheck(role.Buyer, role.Manager, role.Admin), a.Logout)
	}

	RoadGroup := r.Group("/roads")
	{
		RoadGroup.GET("/", a.WithAuthCheck(role.Buyer, role.Manager, role.Admin), roadController.ListRoads)
		RoadGroup.GET("/:id", roadController.GetRoad)
		RoadGroup.POST("/", a.WithAuthCheck(role.Manager, role.Admin), roadController.CreateRoad)
		RoadGroup.PUT("/:id", a.WithAuthCheck(role.Manager, role.Admin), roadController.UpdateRoad)
		RoadGroup.DELETE("/:id", a.WithAuthCheck(role.Manager, role.Admin), roadController.DeleteRoad)
		RoadGroup.POST("/road_travel_request/:id", a.WithAuthCheck(role.Buyer, role.Manager, role.Admin), roadController.AddRoadToTravelRequest)
		RoadGroup.PUT("/road_add_image/:id", a.WithAuthCheck(role.Manager, role.Admin), roadController.AddRoadImage)
	}

	TravelRequestGroup := r.Group("/travelrequests")
	{
		TravelRequestGroup.GET("/", a.WithAuthCheck(role.Buyer, role.Manager, role.Admin), travelRequestController.ListTravelRequests)
		TravelRequestGroup.PUT("/change-status-moderator/:id", a.WithAuthCheck(role.Manager, role.Admin), travelRequestController.ChangeRequestStatusByModerator)
		TravelRequestGroup.GET("/:id", a.WithAuthCheck(role.Buyer, role.Manager, role.Admin), travelRequestController.GetTravelRequestByID)
		TravelRequestGroup.PUT("/:id", a.WithAuthCheck(role.Buyer, role.Manager, role.Admin), travelRequestController.UpdateTravelRequest)
		TravelRequestGroup.PUT("/change-status-user/:id", a.WithAuthCheck(role.Buyer, role.Manager, role.Admin), travelRequestController.ChangeRequestStatusByUser)
		TravelRequestGroup.DELETE("/:id", a.WithAuthCheck(role.Buyer, role.Manager, role.Admin), travelRequestController.SoftDeleteTravelRequest)
		TravelRequestGroup.POST("/pay/:id", a.WithAuthCheck(role.Buyer, role.Manager, role.Admin), travelRequestController.PayTravelRequest)
		TravelRequestGroup.PUT("/change-paidstatus", travelRequestController.ChangePaidStatus)

	}

	TravelRequestRoadGroup := r.Group("/travelrequestroads")
	{
		TravelRequestRoadGroup.DELETE("/:requestID/:roadID", a.WithAuthCheck(role.Buyer, role.Manager, role.Admin), travelRequestRoadController.PhysicalDeleteRoadFromTravelRequest)
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

	log.Println("Server down")

}
