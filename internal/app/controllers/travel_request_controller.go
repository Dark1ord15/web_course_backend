package controllers

import (
	"Road_services/internal/app/ds"
	"Road_services/internal/app/repository"
	"Road_services/internal/app/role"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type TravelRequestController struct {
	repo *repository.Repository
}

func NewTravelRequestController(repo *repository.Repository) *TravelRequestController {
	return &TravelRequestController{
		repo: repo,
	}
}

type GetResponse struct {
	ID             uint      `json:"Travelrequestid"`
	User           string    `json:"User"`
	RequestStatus  string    `json:"Requeststatus"`
	CreationDate   time.Time `json:"Creationdate"`
	FormationDate  time.Time `json:"Formationdate"`
	CompletionDate time.Time `json:"Completiondate"`
	Moderator      string    `json:"Moderator"`
	PaidStatus     string    `json:"Paidstatus"`
}

// @Summary Get Requests
// @Security ApiKeyAuth
// @Description Get all travelrequests
// @Tags TravelRequests
// @ID get-travelrequests
// @Produce json
// @Success 200 {object} ds.Travelrequest
// @Failure 400 {object} ds.Travelrequest "Некорректный запрос"
// @Failure 404 {object} ds.Travelrequest "Некорректный запрос"
// @Failure 500 {object} ds.Travelrequest "Ошибка сервера"
// @Router /travelrequests [get]
// ListTravelRequests обрабатывает запрос на получение списка заявок с дорогами.
// ListTravelRequests обрабатывает запрос на получение списка заявок с дорогами.
func (tc *TravelRequestController) ListTravelRequests(c *gin.Context) {
	userID, contextError := c.Value("userID").(uint)
	if !contextError {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status":  "Failed",
			"Message": "ошибка при авторизации",
		})
		return
	}

	var userRole role.Role
	userRole, contextError = c.Value("userRole").(role.Role)
	if !contextError {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status":  "Failed",
			"Message": "ошибка при авторизации",
		})
		return
	}

	// Получите параметры статуса и диапазона дат из запроса.
	status := c.DefaultQuery("status", "")          // Получаем статус из параметра status
	startDateStr := c.DefaultQuery("startDate", "") // Получаем начальную дату из параметра startDate
	endDateStr := c.DefaultQuery("endDate", "")     // Получаем конечную дату из параметра endDate

	var startDate, endDate time.Time
	var err error

	if userRole == role.Buyer {
		requests, err := tc.repo.GetAllUserRequests(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Status":  "Failed",
				"Message": "Заявки не обнаружены",
			})
			return
		}

		c.JSON(http.StatusOK, requests)
		return
	}

	// Преобразуйте строки с датами в объекты time.Time, если они заданы.
	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат начальной даты"})
			return
		}
		// Добавьте начало дня (00:00:00.000000) к startDate
		startDate = startDate.Add(0 * time.Hour).Add(0 * time.Minute).Add(0 * time.Second).Add(0 * time.Nanosecond)
	}

	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат конечной даты"})
			return
		}
		// Добавьте конец дня (23:59:59.999999) к endDate
		endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second + 999999*time.Nanosecond)
	}

	// Срез для хранения соответствующих заявок
	var matchingRequests []GetResponse

	// Получите список всех заявок из репозитория.
	requests, err := tc.repo.GetAllTravelRequests()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, request := range requests {
		if request.Requeststatus != "deleted" &&
			(status == "" || request.Requeststatus == status) &&
			(startDate.IsZero() || request.Formationdate.After(startDate)) &&
			(endDate.IsZero() || request.Formationdate.Before(endDate)) {

			// Получить информацию о пользователе
			user, err := tc.repo.GetUserByID(request.Userid)
			if err != nil {
				// Пропустить заявку в случае ошибки, продолжить с следующей
				continue
			}

			// Получить информацию о модераторе
			moderator := ""
			if request.Moderatorid != 0 {
				moderatorUser, err := tc.repo.GetUserByID(request.Moderatorid)
				if err != nil {
					// Пропустить заявку в случае ошибки, продолжить с следующей
					continue
				}
				moderator = moderatorUser.Name
			}

			response := GetResponse{
				ID:             request.Travelrequestid,
				User:           user.Name,
				RequestStatus:  request.Requeststatus,
				CreationDate:   request.Creationdate,
				FormationDate:  request.Formationdate,
				CompletionDate: request.Completiondate,
				Moderator:      moderator,
				PaidStatus:     request.Paidstatus,
			}

			// Добавить соответствующую заявку в срез
			matchingRequests = append(matchingRequests, response)
		}
	}

	// Отправить срез в JSON после завершения цикла
	c.JSON(http.StatusOK, matchingRequests)
}

// RoadInfo представляет информацию о дороге.
type RoadInfo struct {
	ID    uint   `json:"Id"`
	Name  string `json:"Name"`
	Price int    `json:"Price"`
}

// @Summary Get Roads by request ID
// @Security ApiKeyAuth
// @Description Show roads by ID of request
// @Tags TravelRequests
// @ID get-roads-by-id-of-request
// @Accept       json
// @Produce      json
// @Success 200 {array} RoadInfo
// @Failure 400 {object} ds.Road "Некорректный запрос"
// @Failure 404 {object} ds.Road "Некорректный запрос"
// @Failure 500 {object} ds.Road "Ошибка сервера"
// @Router /travelrequests/introduced [get]
func (tc *TravelRequestController) GetTravelRequestByID(c *gin.Context) {
	userID, contextError := c.Value("userID").(uint)
	if !contextError {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status":  "Failed",
			"Message": "ошибка при авторизации",
		})
		return
	}

	// Получите заявку с указанным ID из репозитория.
	request, err := tc.repo.GetTravelRequestByID(uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Заявка не найдена"})
		return
	}

	// Получите связанные с заявкой дороги.
	roads, err := tc.repo.GetRoadsByTravelRequest(request.Travelrequestid, request.Requeststatus)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Создайте структуру для хранения информации о дорогах.
	var roadInfoList []RoadInfo

	for _, road := range roads {
		roadInfo := RoadInfo{
			ID:    road.Roadid,
			Name:  road.Name,
			Price: road.Price,
		}
		roadInfoList = append(roadInfoList, roadInfo)
	}

	c.JSON(http.StatusOK, roadInfoList)
}

// @Summary Update TravelRequest by ID
// @Security ApiKeyAuth
// @Description Update travelrequest by ID
// @Tags TravelRequests
// @ID update-travelrequest-by-id
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID заявки"
// @Param input body ds.Travelrequest true "request info"
// @Success 200 {string} string
// @Failure 400 {object} ds.Travelrequest "Некорректный запрос"
// @Failure 404 {object} ds.Travelrequest "Некорректный запрос"
// @Failure 500 {object} ds.Travelrequest "Ошибка сервера"
// @Router /travelrequests/{id} [put]
func (tc *TravelRequestController) UpdateTravelRequest(c *gin.Context) {
	// Получите ID заявки из параметров запроса.
	requestIDStr := c.Param("id")
	requestID, err := strconv.Atoi(requestIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный request ID"})
		return
	}

	// Прочитайте данные заявки из запроса.
	var updatedRequest ds.Travelrequest
	if err := c.ShouldBindJSON(&updatedRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный request data"})
		return
	}

	// Здесь обновите заявку с указанным ID в репозитории.
	err = tc.repo.UpdateTravelRequest(uint(requestID), updatedRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Заявка успешно обновлена"})
}

// func (tc *TravelRequestController) DeleteTravelRequest(c *gin.Context) {
// 	// Получите ID заявки из параметров запроса.
// 	requestIDStr := c.Param("id")
// 	requestID, err := strconv.Atoi(requestIDStr)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный request ID"})
// 		return
// 	}

// 	// Здесь удалите заявку с указанным ID из репозитория.
// 	err = tc.repo.DeleteTravelRequest(uint(requestID))
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Заявка успешно удалена"})
// }

// @Summary Update TravelRequest Status By User
// @Security ApiKeyAuth
// @Description Update travelrequest status by user
// @Tags TravelRequests
// @ID update-travelrequest-status-by-user
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID заявки"
// @Success 200 {string} string
// @Failure 400 {object} ds.Travelrequest "Некорректный запрос"
// @Failure 404 {object} ds.Travelrequest "Некорректный запрос"
// @Failure 500 {object} ds.Travelrequest "Ошибка сервера"
// @Router /travelrequests/change-status-user/{id} [put]
func (tc *TravelRequestController) ChangeRequestStatusByUser(c *gin.Context) {

	userID, contextError := c.Value("userID").(uint)
	if !contextError {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status":  "Failed",
			"Message": "ошибка при авторизации",
		})
		return
	}
	// Получите ID заявки из параметров запроса.
	requestIDStr := c.Param("id")
	requestID, err := strconv.Atoi(requestIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный request ID"})
		return
	}

	request, err := tc.repo.GetTravelRequestByID(uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Заявка не найдена"})
		return
	}

	// Проверьте, что статус заявки "introduced".
	if request.Requeststatus != "introduced" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Статус заявки не 'introduced'"})
		return
	}

	// Измените статус заявки на "formed".
	request.Requeststatus = "formed"

	// Сохраните обновленную заявку в репозитории.
	err = tc.repo.UpdateTravelRequest(uint(requestID), request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Статус заявки обновлен на 'formed'"})
}

// @Summary Update TravelRequest Status By Moderator
// @Security ApiKeyAuth
// @Description Update travelrequest by moderator
// @Tags TravelRequests
// @ID update-travelrequest-status-by-moderator
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID заявки"
// @Success 200 {string} string
// @Failure 400 {object} ds.Travelrequest "Некорректный запрос"
// @Failure 404 {object} ds.Travelrequest "Некорректный запрос"
// @Failure 500 {object} ds.Travelrequest "Ошибка сервера"
// @Router /travelrequests/change-status-moderator/{id} [put]
func (tc *TravelRequestController) ChangeRequestStatusByModerator(c *gin.Context) {
	adminID, contextError := c.Value("userID").(uint)
	if !contextError {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status":  "Failed",
			"Message": "ошибка при авторизации",
		})
		return
	}
	// Получите ID заявки из параметров запроса.
	requestIDStr := c.Param("id")
	requestID, err := strconv.Atoi(requestIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный request ID"})
		return
	}
	// Получите статус заявки из параметров запроса.
	status := c.Query("status")

	// Проверьте, что статус валиден (допустимы 'completed' и 'rejected').
	if status != "completed" && status != "rejected" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный статус. Он должен быть 'completed' или 'rejected'"})
		return
	}
	// Здесь проверьте, что заявка с указанным ID существует и имеет статус 'formed'.
	request, err := tc.repo.GetTravelRequestByID2(uint(requestID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Запрос не найден"})
		return
	}
	if request.Requeststatus != "formed" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Для изменений статус запрос должен быть'formed'"})
		return
	}
	// Здесь обновите статус заявки на 'completed' или 'rejected' в репозитории.
	request.Requeststatus = status
	if err := tc.repo.UpdateTravelRequest2(request.Travelrequestid, adminID, request); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("me tut 6")
	c.JSON(http.StatusOK, gin.H{"message": "Статус успешно обнавлен"})
}

// @Summary Delete TravelRequest by ID
// @Security ApiKeyAuth
// @Description Delete travelrequest by ID
// @Tags TravelRequests
// @ID delete-request-by-id
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID заявки"
// @Success 200 {string} string
// @Failure 400 {object} ds.Travelrequest "Некорректный запрос"
// @Failure 404 {object} ds.Travelrequest "Некорректный запрос"
// @Failure 500 {object} ds.Travelrequest "Ошибка сервера"
// @Router /travelrequests/{id} [delete]
func (tc *TravelRequestController) SoftDeleteTravelRequest(c *gin.Context) {
	userID, contextError := c.Value("userID").(uint)
	if !contextError {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status":  "Failed",
			"Message": "ошибка при авторизации",
		})
		return
	}

	// Получите ID заявки из параметров запроса.
	requestIDStr := c.Param("id")
	requestID, err := strconv.Atoi(requestIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный request ID"})
		return
	}

	request, err := tc.repo.GetTravelRequestByID(uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Заявка не найдена"})
		return
	}

	// Проверьте, можно ли удалить заявку в текущем статусе.
	if request.Requeststatus != "introduced" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Заявка с текущем статусом не может быть удалена."})
		return
	}

	// Измените статус заявки на "deleted" в репозитории.
	request.Requeststatus = "deleted"
	err = tc.repo.UpdateTravelRequest(uint(requestID), request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Заявка успешно удалена"})
}

// PayTravelRequest обрабатывает оплату заявки
// PayTravelRequest обрабатывает оплату заявки
func (t *TravelRequestController) PayTravelRequest(c *gin.Context) {
	id := c.Param("id")

	// Вызываем метод оплаты из репозитория
	paidStatus, err := t.repo.PayTravelRequest(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process payment", "details": err.Error()})
		return
	}

	// Обновляем статус оплаты в базе данных
	if err := t.repo.UpdatePaidStatus(id, paidStatus); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update paid status", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": paidStatus})
}

// ChangePaidStatus обновляет статус оплаты в заявке
func (t *TravelRequestController) ChangePaidStatus(c *gin.Context) {
	// Парсим JSON из тела запроса
	var requestBody map[string]interface{}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	// Проверяем наличие ключа "12345" в JSON
	if apiKey, ok := requestBody["key"]; !ok || apiKey != "12345" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
		return
	}

	// Извлекаем значение id из JSON
	id, ok := requestBody["id"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'id' in JSON"})
		return
	}

	// Извлекаем значение paidStatus из JSON
	paidStatus, ok := requestBody["paidstatus"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'paidstatus' in JSON"})
		return
	}

	// Вызываем метод обновления статуса оплаты из репозитория
	if err := t.repo.UpdatePaidStatus(id, paidStatus); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update paid status" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "Статус оплаты изменен"})
}
