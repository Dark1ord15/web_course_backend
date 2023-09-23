package api

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"database/sql"

	_ "github.com/lib/pq" // Для PostgreSQL
)

// var roads = []Road{
// 	{
// 		Id:              1,
// 		Name:            "М-1 'Беларусь' 33км-66км",
// 		TrustManagement: 516,
// 		Lenght:          456,
// 		PaidLenght:      33,
// 		Category:        "1Б,1В",
// 		NumberOfStripes: "4-8",
// 		Speed:           110,
// 		Price:           150,
// 		Image:           "/image/М1.png",
// 	},
// 	{
// 		Id:              2,
// 		Name:            "М-3 'Украина' 124км-194км",
// 		TrustManagement: 454,
// 		Lenght:          510,
// 		PaidLenght:      70,
// 		Category:        "1Б,II,III",
// 		NumberOfStripes: "2-8",
// 		Speed:           110,
// 		Price:           250,
// 		Image:           "/image/М3.png",
// 	},
// 	{
// 		Id:              3,
// 		Name:            "М-4 'Дон' 21км-1319км",
// 		TrustManagement: 1930,
// 		Lenght:          1542,
// 		PaidLenght:      1042,
// 		Category:        "1А,1Б,1В,II",
// 		NumberOfStripes: "2-8",
// 		Speed:           130,
// 		Price:           3500,
// 		Image:           "/image/М4.png",
// 	},
// 	{
// 		Id:              4,
// 		Name:            "М-11 'Нева' 77км-684км",
// 		TrustManagement: 619,
// 		Lenght:          669,
// 		PaidLenght:      607,
// 		Category:        "1А",
// 		NumberOfStripes: "4-10",
// 		Speed:           130,
// 		Price:           5000,
// 		Image:           "/image/М11.png",
// 	},
// 	{
// 		Id:              5,
// 		Name:            "М-12 'Восток' 56км-471км",
// 		TrustManagement: 415,
// 		Lenght:          811,
// 		PaidLenght:      415,
// 		Category:        "1Б",
// 		NumberOfStripes: "4-6",
// 		Speed:           110,
// 		Price:           2800,
// 		Image:           "/image/М12.png",
// 	},
// 	{
// 		Id:              6,
// 		Name:            "ЦКАД",
// 		TrustManagement: 336,
// 		Lenght:          336,
// 		PaidLenght:      260,
// 		Category:        "1A, II",
// 		NumberOfStripes: "4",
// 		Speed:           110,
// 		Price:           2000,
// 		Image:           "/image/ЦКАД.png",
// 	},
// }

func setupDB() (*sql.DB, error) {
	// Здесь используйте значения из ваших переменных окружения или файла конфигурации
	connectionString := "user=bmstu_user password=bmstu_password dbname=bmstu sslmode=disable"

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func StartServer() {
	log.Println("Server start up")

	db, err := setupDB()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	rows, err := db.Query(`SELECT "RoadID", "Name", "TrustManagment", "Lenght", "PaidLenght", "Category", "NumberOfStripes", "Speed", "Price", "Image" FROM public."Road"`)
	if err != nil {
		log.Fatalf("Failed to query the database: %v", err)
	}
	defer rows.Close()

	var roads []Road

	for rows.Next() {
		var road Road
		if err := rows.Scan(&road.Id, &road.Name, &road.TrustManagement, &road.Lenght, &road.PaidLenght, &road.Category, &road.NumberOfStripes, &road.Speed, &road.Price, &road.Image); err != nil {
			log.Fatalf("Failed to scan row: %v", err)
		}
		roads = append(roads, road)
	}

	r := gin.Default()

	r.LoadHTMLGlob("templates/*")

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

	r.GET("/", func(c *gin.Context) {

		searchQuery := c.DefaultQuery("fsearch", "")

		if searchQuery == "" {
			c.HTML(http.StatusOK, "main_page.tmpl", gin.H{
				"roads": roads,
			})
			return
		}

		var result []Road
		for _, road := range roads {
			if strings.Contains(strings.ToLower(road.Name), strings.ToLower(searchQuery)) {
				result = append(result, road)
			}
		}

		c.HTML(http.StatusOK, "main_page.tmpl", gin.H{
			"roads": result,
			"Query": searchQuery,
		})
	})
	r.Static("/image", "./resources/image")
	r.Static("/css", "./resources/css")

	r.Run(":8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

	log.Println("Server down")
}
