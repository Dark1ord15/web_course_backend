package controllers

import (
	"Road_services/internal/app/ds"
	"Road_services/internal/app/repository"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RoadController struct {
	repo *repository.Repository
}

func NewRoadController(repo *repository.Repository) *RoadController {
	return &RoadController{repo: repo}
}

func (rc *RoadController) ListRoads(c *gin.Context) {
	// Получите параметр minLength из запроса, если предоставлен.
	minLengthStr := c.DefaultQuery("minLength", "")

	var roads []ds.Road
	var err error

	if minLengthStr != "" {
		// Фильтруйте дороги по минимальной длине.
		minLength, conversionErr := strconv.Atoi(minLengthStr)
		if conversionErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный minLength"})
			return
		}
		roads, err = rc.repo.GetRoadsByMinLength(minLength)
	} else {
		// Верните все дороги, так как minLength не задан.
		roads, err = rc.repo.GetAllRoads()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, roads)
}

func (rc *RoadController) GetRoad(c *gin.Context) {
	roadID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID дороги"})
		return
	}

	road, err := rc.repo.GetRoadByID(roadID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Проверяем, найдена ли дорога. Если она не найдена, отправляем статус 404.
	if road == (ds.Road{}) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Дорога не найдена"})
		return
	}

	c.JSON(http.StatusOK, road)
}

func (rc *RoadController) CreateRoad(c *gin.Context) {
	var newRoad ds.Road
	if err := c.ShouldBindJSON(&newRoad); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
		return
	}

	if err := rc.repo.CreateRoad(newRoad); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newRoad)
}

func (rc *RoadController) UpdateRoad(c *gin.Context) {
	roadID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID дороги"})
		return
	}

	var updatedRoad ds.Road
	if err := c.ShouldBindJSON(&updatedRoad); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
		return
	}

	if err := rc.repo.UpdateRoad(roadID, updatedRoad); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedRoad)
}

func (rc *RoadController) DeleteRoad(c *gin.Context) {
	roadID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID дороги"})
		return
	}

	if err := rc.repo.DeleteRoad(roadID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// AddRoadToLastTravelRequest добавляет дорогу к последней заявке.
func (rc *RoadController) AddRoadToLastTravelRequest(c *gin.Context) {
	// Прочитайте roadID из запроса
	roadID, err := strconv.Atoi(c.Param("roadID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
		return
	}

	// Вызовите метод репозитория для добавления дороги к последней заявке
	if err := rc.repo.AddRoadToLastTravelRequest(uint(roadID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Дорога успешно добавлена к последней заявке"})
}
func (rc *RoadController) AddRoadToTravelRequest(c *gin.Context) {
	userID := rc.repo.GetUserID()
	roadID, err := strconv.Atoi(c.Param("roadID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный roadID"})
		return
	}

	// Получите информацию о дороге по roadID
	road, err := rc.repo.GetRoadByID(roadID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Проверьте существование и статус дороги
	if road.Roadid == 0 || road.Statusroad != "active" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный roadID или услуга 'неактивна'"})
		return
	}

	// Поиск или создание заявки пользователя
	requestID, err := rc.repo.FindOrCreateRequest(int(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if requestID == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create or find the travel request"})
		return
	}
	// Проверьте, существует ли уже такая дорога в текущей заявке пользователя.
	exists, err := rc.repo.IsRoadAlreadyAddedToRequest(requestID, uint(roadID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Такая дорога уже существует в заявке"})
		return
	}

	// Добавьте дорогу к заявке пользователя.
	if err := rc.repo.AddRoadToTravelRequest(requestID, uint(roadID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Дорога успешно добавлена к заявке пользователя"})
}
func (rc *RoadController) AddConsultationImage(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if id < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status":  "Failed",
			"Message": "неверное значение id",
		})
		return
	}
	// Чтение изображения из запроса
	image, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image"})
		return
	}

	// Чтение содержимого изображения в байтах
	file, err := image.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при открытии"})
		return
	}
	defer file.Close()

	imageBytes, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка чтения"})
		return
	}
	// Получение Content-Type из заголовков запроса
	contentType := image.Header.Get("Content-Type")

	// Вызов функции репозитория для добавления изображения
	err = rc.repo.AddConsultationImage(id, imageBytes, contentType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сервера"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Картинка обнавлена"})

}
