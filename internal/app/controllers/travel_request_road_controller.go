// travel_request_road_controller.go

package controllers

import (
	"Road_services/internal/app/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TravelRequestRoadController struct {
	repo *repository.Repository
}

func NewTravelRequestRoadController(repo *repository.Repository) *TravelRequestRoadController {
	return &TravelRequestRoadController{
		repo: repo,
	}
}

// @Summary Delete Road From TravelRequest
// @Security ApiKeyAuth
// @Description delete road from travelrequest
// @Tags Road-TravelRequest
// @ID delete-road-from-travelrequest
// @Accept       json
// @Produce      json
// @Param        requestID   path      int  true  "ID заявки"
// @Param        roadID   path      int  true  "ID консультации"
// @Success 200 {string} string "Консультация была удалена из заявки"
// @Failure 400 {string} string "Некорректный запрос"
// @Failure 404 {string} string "Некорректный запрос"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /travelrequestroads/{requestID}/{roadID} [delete]
func (tc *TravelRequestRoadController) PhysicalDeleteRoadFromTravelRequest(c *gin.Context) {
	// Получите ID заявки и ID дороги из параметров запроса.
	requestIDStr := c.Param("requestID")
	roadIDStr := c.Param("roadID")

	requestID, err := strconv.Atoi(requestIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный request ID"})
		return
	}

	roadID, err := strconv.Atoi(roadIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный road ID"})
		return
	}

	// Проверьте, существует ли связь между дорогой и заявкой.
	if !tc.repo.IsRoadConnectedToRequest(uint(requestID), uint(roadID)) {
		c.JSON(http.StatusNotFound, gin.H{"error": "В заявке нет такой дороги"})
		return
	}

	// Выполните физическое удаление связи между дорогой и заявкой.
	if err := tc.repo.DeleteRoadFromTravelRequest(uint(requestID), uint(roadID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Запись успешно удалена"})
}
