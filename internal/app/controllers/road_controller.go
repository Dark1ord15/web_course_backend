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

type response struct {
	RequestID string    `json:"requestID"`
	Roads     []ds.Road `json:"roads"`
}

// @Summary Get Roads
// @Description Get all roads
// @Tags Roads
// @ID get-roads
// @Produce json
// @Success 200 {object} response
// @Failure 400 {object} ds.Road "Некорректный запрос"
// @Failure 404 {object} ds.Road "Некорректный запрос"
// @Failure 500 {object} ds.Road "Ошибка сервера"
// @Router /roads [get]
func (rc *RoadController) ListRoads(c *gin.Context) {
	// Получите параметр minLength из запроса, если предоставлен.
	minLengthStr := c.DefaultQuery("minLength", "")

	var roads []ds.Road
	var err error
	userID, _ := c.Value("userID").(uint)
	// Ваш запрос для получения id заявки со статусом introduced
	requestID := rc.repo.GetRequestIdWithStatusAndUser("introduced", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

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

	// Добавьте id заявки в ответ
	response := gin.H{
		"requestID": requestID,
		"roads":     roads,
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Get Road by ID
// @Description Show road by ID
// @Tags Roads
// @ID id
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID дороги"
// @Success 200 {object} ds.Road
// @Failure 400 {object} ds.Road "Некорректный запрос"
// @Failure 404 {object} ds.Road "Некорректный запрос"
// @Failure 500 {object} ds.Road "Ошибка сервера"
// @Router /roads/{id} [get]
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

// @Summary create road
// @Security ApiKeyAuth
// @Description create road
// @Tags Roads
// @ID create-road
// @Accept json
// @Produce json
// @Param input body ds.Road true "road info"
// @Success 200 {string} string
// @Failure 400 {object} ds.Road "Некорректный запрос"
// @Failure 404 {object} ds.Road "Некорректный запрос"
// @Failure 500 {object} ds.Road "Ошибка сервера"
// @Router /roads [post]
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

// @Summary update croad
// @Security ApiKeyAuth
// @Description update road
// @Tags Roads
// @ID update-road
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID дороги"
// @Param input body ds.Road true "consultation info"
// @Success 200 {string} string
// @Failure 400 {object} ds.Road "Некорректный запрос"
// @Failure 404 {object} ds.Road "Некорректный запрос"
// @Failure 500 {object} ds.Road "Ошибка сервера"
// @Router /roads/{id} [put]
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

// @Summary Delete road by ID
// @Security ApiKeyAuth
// @Description Delete road by ID
// @Tags Roads
// @ID delete-road-by-id
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID дороги"
// @Success 200 {string} string
// @Failure 400 {object} ds.Road "Некорректный запрос"
// @Failure 404 {object} ds.Road "Некорректный запрос"
// @Failure 500 {object} ds.Road "Ошибка сервера"
// @Router /roads/{id} [delete]
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

// @Summary add road to travelrequest
// @Security ApiKeyAuth
// @Description add road to travelrequest
// @Tags Roads
// @ID add-road-to-request
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID дороги"
// @Success 200 {string} string
// @Failure 400 {object} ds.Road "Некорректный запрос"
// @Failure 404 {object} ds.Road "Некорректный запрос"
// @Failure 500 {object} ds.Road "Ошибка сервера"
// @Router /roads/road_travel_request/{id} [post]
func (rc *RoadController) AddRoadToTravelRequest(c *gin.Context) {
	userID, contextError := c.Value("userID").(uint)
	if !contextError {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status":  "Failed",
			"Message": "ошибка при авторизации",
		})
		return
	}
	roadID, err := strconv.Atoi(c.Param("id"))
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

// @Summary Add road image
// @Security ApiKeyAuth
// @Description Add an road to a specific consultation by ID.
// @Tags Roads
// @ID add-road-image
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID дороги"
// @Param image formData file true "Image file to be uploaded"
// @Success 200 {string} string
// @Failure 400 {object} ds.Road "Некорректный запрос"
// @Failure 404 {object} ds.Road "Некорректный запрос"
// @Failure 500 {object} ds.Road "Ошибка сервера"
// @Router /roads/road_add_image/{id} [post]
func (rc *RoadController) AddRoadImage(c *gin.Context) {
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
