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

	// Получите список всех заявок из репозитория.
	requests, err := tc.repo.GetAllTravelRequests()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Создайте слайс для хранения заявок, которые соответствуют фильтру.
	filteredRequests := []ds.Travelrequest{}

	// Фильтруйте заявки в соответствии с заданным статусом и диапазоном дат.
	for _, request := range requests {
		if request.Requeststatus != "deleted" { // Исключаем удаленные заявки
			if (status == "" || request.Requeststatus == status) &&
				(startDate.IsZero() || request.Formationdate.After(startDate)) &&
				(endDate.IsZero() || request.Formationdate.Before(endDate)) {
				filteredRequests = append(filteredRequests, request)
			}
		}
	}

	if len(filteredRequests) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Нет результатов, соответствующих заданным параметрам"})
		return
	}
	// Верните список заявок, соответствующих фильтру.
	c.JSON(http.StatusOK, filteredRequests)
}

// RoadsResponse описывает структуру ответа о дорогах.
type RoadsResponse struct {
	RoadNames  []string `json:"RoadNames"`
	RoadImages []string `json:"RoadImages"`
}

// Response описывает структуру ответа о заявке с дорогами.
type Response struct {
	ID            uint          `json:"ID"`
	CreationDate  time.Time     `json:"CreationDate"`
	FormationDate time.Time     `json:"FormationDate"`
	RequestStatus string        `json:"RequestStatus"`
	RoadsResponse RoadsResponse `json:"RoadsResponse"`
}

// @Summary Get Roads by request ID
// @Security ApiKeyAuth
// @Description Show roads by ID of request
// @Tags TravelRequests
// @ID get-roads-by-id-of-request
// @Accept       json
// @Produce      json
// @Success 200 {object} Response
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

	// Проверьте, что заявка не удалена.
	// if request.Requeststatus == "deleted" {
	// 	c.JSON(http.StatusNotFound, gin.H{"error": "Заявка не найдена"})
	// 	return
	// }

	// Получите связанные с заявкой дороги.
	roads, err := tc.repo.GetRoadsByTravelRequest(request.Travelrequestid, request.Requeststatus)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Создайте структуру для хранения информации о дорогах и заявке.
	type RoadsResponse struct {
		RoadNames  []string `json:"RoadNames"`
		RoadImages []string `json:"RoadImages"`
	}

	roadsResponse := RoadsResponse{
		RoadNames:  []string{},
		RoadImages: []string{},
	}

	for _, road := range roads {
		roadsResponse.RoadNames = append(roadsResponse.RoadNames, road.Name)
		roadsResponse.RoadImages = append(roadsResponse.RoadImages, road.Image)
	}

	// Включите информацию о заявке в ответ.
	type Response struct {
		ID            uint          `json:"ID"`
		CreationDate  time.Time     `json:"CreationDate"`
		FormationDate time.Time     `json:"FormationDate"`
		RequestStatus string        `json:"RequestStatus"`
		RoadsResponse RoadsResponse `json:"RoadsResponse"`
	}

	response := Response{
		ID:            request.Travelrequestid,
		CreationDate:  request.Creationdate,
		FormationDate: request.Formationdate,
		RequestStatus: request.Requeststatus,
		RoadsResponse: roadsResponse,
	}

	c.JSON(http.StatusOK, response)
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
	// Получите ID заявки из параметров запроса.
	fmt.Println("me tut 1")
	requestIDStr := c.Param("id")
	requestID, err := strconv.Atoi(requestIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный request ID"})
		return
	}
	fmt.Println("me tut 2")
	// Получите статус заявки из параметров запроса.
	status := c.Query("status")

	// Проверьте, что статус валиден (допустимы 'completed' и 'rejected').
	if status != "completed" && status != "rejected" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный статус. Он должен быть 'completed' или 'rejected'"})
		return
	}
	fmt.Println("me tut 3")
	// Здесь проверьте, что заявка с указанным ID существует и имеет статус 'formed'.
	request, err := tc.repo.GetTravelRequestByID2(uint(requestID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Запрос не найден"})
		return
	}
	fmt.Println("me tut4")
	if request.Requeststatus != "formed" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Для изменений статус запрос должен быть'formed'"})
		return
	}
	fmt.Println("me tut 5")
	// Здесь обновите статус заявки на 'completed' или 'rejected' в репозитории.
	request.Requeststatus = status
	if err := tc.repo.UpdateTravelRequest(request.Travelrequestid, request); err != nil {
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
