package app

import (
	"Road_services/internal/app/controllers"
	"Road_services/internal/app/ds"
	"Road_services/internal/app/repository"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq" // Для PostgreSQL
)

func (a *Application) StartServer() {
	log.Println("Server start up")
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		searchQuery := c.DefaultQuery("fsearch", "")

		var result []ds.Road
		var err error

		if searchQuery == "" {
			result, err = a.repository.GetAllRoads()
		} else {
			result, err = a.repository.SearchRoads(searchQuery)
		}

		if err != nil {
			log.Printf("Error while fetching data: %v", err)
			c.HTML(http.StatusInternalServerError, "error.tmpl", nil) // Обработка ошибки
			return
		}

		c.HTML(http.StatusOK, "main_page.tmpl", gin.H{
			"roads": result,
			"Query": searchQuery,
		})
	})

	r.POST("/delete/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			//обработка ошибка
			log.Printf("cant get product by id %v", err)
			c.Redirect(http.StatusMovedPermanently, "/")
			return
		}
		a.repository.DeleteRoad(id)
		c.Redirect(http.StatusMovedPermanently, "/")
	})

	r.LoadHTMLGlob("templates/*")

	r.GET("/road/:id", func(c *gin.Context) {
		var roads []ds.Road
		roads, err := a.repository.GetAllRoads()
		if err != nil { // если не получилось
			log.Printf("cant get product by id %v", err)
			return
		}
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			// Обработка ошибки
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
			return
		}

		road := roads[id-1]
		c.HTML(http.StatusOK, "info.tmpl", road)
	})

	r.Static("/image", "./resources/image")
	r.Static("/css", "./resources/css")

	repo, err := repository.New("user=bmstu_user password=bmstu_password dbname=bmstu host=localhost port=5432 sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to create repository: %v", err)
	}

	roadController := controllers.NewRoadController(repo)
	r.GET("/roads", roadController.ListRoads)
	r.GET("/roads/:id", roadController.GetRoad)
	r.POST("/roads", roadController.CreateRoad)
	r.PUT("/roads/:id", roadController.UpdateRoad)
	r.DELETE("/roads/:id", roadController.DeleteRoad)
	r.POST("/add_road_to_last_request/:roadID", roadController.AddRoadToLastTravelRequest)
	r.POST("/road_travel_request/:roadID", roadController.AddRoadToTravelRequest)
	r.PUT("/road_add_image/:id", roadController.AddConsultationImage)

	travelRequestController := controllers.NewTravelRequestController(repo)
	r.GET("/travelrequests", travelRequestController.ListTravelRequests)
	r.GET("/travelrequests/:id", travelRequestController.GetTravelRequestByID)
	r.PUT("/travelrequests/:id", travelRequestController.UpdateTravelRequest)
	// r.DELETE("/travelrequests/:id", travelRequestController.DeleteTravelRequest)
	r.PUT("/travelrequests/change-status-user/:id", travelRequestController.ChangeRequestStatusByUser)
	r.PUT("/travelrequests/change-status-moderator/:id", travelRequestController.ChangeRequestStatusByModerator)
	r.DELETE("/travelrequests/:id", travelRequestController.SoftDeleteTravelRequest)

	travelRequestRoadController := controllers.NewTravelRequestRoadController(repo)
	r.DELETE("/travelrequestroads/:requestID/:roadID", travelRequestRoadController.PhysicalDeleteRoadFromTravelRequest)

	r.Run(":8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

	log.Println("Server down")
}
