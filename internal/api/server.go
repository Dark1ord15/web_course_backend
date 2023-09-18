package api

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

var roads = []Road{
	{
		Id:              1,
		Name:            "М-1 'Беларусь'",
		TrustManagement: 516,
		Lenght:          456,
		PaidLenght:      33,
		Category:        "1Б,1В",
		NumberOfStripes: "4-8",
		Speed:           110,
		Image:           "/image/М1.png",
	},
	{
		Id:              2,
		Name:            "М-3 'Украина'",
		TrustManagement: 454,
		Lenght:          510,
		PaidLenght:      70,
		Category:        "1Б,II,III",
		NumberOfStripes: "2-8",
		Speed:           110,
		Image:           "/image/М3.png",
	},
	{
		Id:              3,
		Name:            "М-4 'Дон'",
		TrustManagement: 1930,
		Lenght:          1542,
		PaidLenght:      1042,
		Category:        "1А,1Б,1В,II",
		NumberOfStripes: "2-8",
		Speed:           130,
		Image:           "/image/М4.png",
	},
	{
		Id:              4,
		Name:            "М-11 'Нева'",
		TrustManagement: 619,
		Lenght:          669,
		PaidLenght:      607,
		Category:        "1А",
		NumberOfStripes: "4-10",
		Speed:           130,
		Image:           "/image/М11.png",
	},
	{
		Id:              5,
		Name:            "М-12 'Восток'",
		TrustManagement: 415,
		Lenght:          811,
		PaidLenght:      415,
		Category:        "1Б",
		NumberOfStripes: "4-6",
		Speed:           110,
		Image:           "/image/М12.png",
	},
	{
		Id:              6,
		Name:            "ЦКАД",
		TrustManagement: 336,
		Lenght:          336,
		PaidLenght:      260,
		Category:        "1A, II",
		NumberOfStripes: "4",
		Speed:           110,
		Image:           "/image/ЦКАД.png",
	},
}

func StartServer() {
	log.Println("Server start up")

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "main_page.tmpl", gin.H{
			"roads": roads,
		})
	})

	r.GET("/road/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			// Обработка ошибки
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
			return
		}

		road := roads[id-1]
		c.HTML(http.StatusOK, "info.tmpl", road)
	})

	r.GET("/search", func(c *gin.Context) {

		searchQuery := c.DefaultQuery("fsearch", "")

		var result []Road

		for _, road := range roads {
			if strings.Contains(strings.ToLower(road.Name), strings.ToLower(searchQuery)) {
				result = append(result, road)
			}
		}

		c.HTML(http.StatusOK, "main_page.tmpl", gin.H{
			"roads": result,
		})
	})
	r.Static("/image", "./resources/image")
	r.Static("/css", "./resources/css")

	r.Run(":8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

	log.Println("Server down")
}
