package repository

import (
	"fmt"
	"log"

	"strings"

	"Road_services/internal/app/ds"
	minioclient "Road_services/internal/minioClient"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	db          *gorm.DB
	minioClient *minioclient.MinioClient
}

func New(dsn string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	minioClient, err := minioclient.NewMinioClient()
	if err != nil {
		log.Println("error here start!")
		return nil, err
	}

	return &Repository{
		db:          db,
		minioClient: minioClient,
	}, nil
}

// func (r *Repository) GetConsultationByID(id int) (*ds.Consultation, error) {
// 	consultation := &ds.Consultation{}

// 	err := r.db.First(consultation, "id = ?", id).Error
// 	if err != nil {
// 		return nil, err
// 	}

// 	return consultation, nil
// }

func (r *Repository) GetUserID() int {
	// Вместо жестко заданного значения 1, здесь можете добавить логику получения `userID` из аутентификации или другого источника.
	return 1
}

func (r *Repository) SearchRoads(searchQuery string) ([]ds.Road, error) {
	var result []ds.Road

	err := r.db.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(searchQuery)+"%").Find(&result, "Statusroad::text='active'").Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *Repository) CreateRoad(road ds.Road) error {
	return r.db.Create(&road).Error
}

func (r *Repository) GetAllRoads() ([]ds.Road, error) {
	var roads []ds.Road

	err := r.db.Order("RoadID ASC").Find(&roads, "Statusroad::text='active'").Error
	if err != nil {
		return nil, err
	}

	return roads, nil
}

func (r *Repository) DeleteRoad(id int) error {
	return r.db.Exec("Update roads SET Statusroad = 'deleted' WHERE roadid=?", id).Error
}

// GetRoadByID - получение информации о дороге по ID
func (r *Repository) GetRoadByID(id int) (ds.Road, error) {
	var road ds.Road
	if err := r.db.First(&road, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ds.Road{}, nil // Дорога с таким ID не найдена
		}
		return ds.Road{}, err
	}
	return road, nil
}

// UpdateRoad - обновление информации о дороге по ID
func (r *Repository) UpdateRoad(id int, updatedRoad ds.Road) error {
	if err := r.db.Model(ds.Road{}).Where("roadid = ?", id).Updates(updatedRoad).Error; err != nil {
		return err
	}
	return nil
}

// AddRoadToLastTravelRequest добавляет дорогу к последней заявке.
func (r *Repository) AddRoadToLastTravelRequest(roadID uint) error {
	// Определите последнюю заявку
	var lastRequest ds.Travelrequest
	if err := r.db.Last(&lastRequest).Error; err != nil {
		return err
	}

	// Проверьте, существует ли дорога с указанным ID
	var road ds.Road
	if err := r.db.First(&road, roadID).Error; err != nil {
		return err
	}

	// Создайте новую запись для связи заявки и дороги
	travelRequestRoad := ds.Travelrequestroad{
		Travelrequestid: lastRequest.Travelrequestid,
		Roadid:          roadID,
	}

	// Начните транзакцию для выполнения операции
	tx := r.db.Begin()
	if err := tx.Create(&travelRequestRoad).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Завершите транзакцию
	tx.Commit()

	return nil
}

// Получение всех заявок
func (r *Repository) GetAllTravelRequests() ([]ds.Travelrequest, error) {
	var requests []ds.Travelrequest
	err := r.db.Find(&requests).Error
	if err != nil {
		return nil, err
	}
	return requests, nil
}

func (r *Repository) GetAllUserRequests(userID uint) ([]ds.Travelrequest, error) {
	var requests []ds.Travelrequest
	err := r.db.Find(&requests, "userid = ?", userID).Error
	if err != nil {
		return nil, err
	}
	return requests, nil
}

func (r *Repository) GetTravelRequestByID(userID uint) (ds.Travelrequest, error) {
	var travelRequest ds.Travelrequest

	err := r.db.Where("userid = ? AND requeststatus = ?", userID, "introduced").First(&travelRequest).Error
	if err != nil {
		return ds.Travelrequest{}, err
	}

	return travelRequest, nil
}

func (r *Repository) GetRoadsByMinLength(minLength int) ([]ds.Road, error) {
	var roads []ds.Road
	err := r.db.Where("Endofsection - Startofsection >= ?", minLength).Order("RoadID ASC").Find(&roads, "Statusroad::text='active'").Error
	if err != nil {
		return nil, err
	}
	return roads, nil
}

// Редактирование заявки
func (r *Repository) UpdateTravelRequest(id uint, updatedRequest ds.Travelrequest) error {
	if err := r.db.Model(ds.Travelrequest{}).Where("travelrequestid = ?", id).Updates(updatedRequest).Error; err != nil {
		return err
	}
	return nil
}

// Удаление заявки
func (r *Repository) DeleteTravelRequest(id uint) error {
	return r.db.Delete(&ds.Travelrequest{}, id).Error
}

// GetTravelRequestByUserAndStatus получает заявку пользователя с указанным статусом.
func (r *Repository) GetTravelRequestByUserAndStatus(userID int, status string) (ds.Travelrequest, error) {
	var request ds.Travelrequest
	if err := r.db.Where("userid = ? AND requeststatus = ?", userID, status).First(&request).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ds.Travelrequest{}, nil
		}
		return ds.Travelrequest{}, err
	}
	return request, nil
}

// GetLastTravelRequestByUser получает последнюю заявку пользователя.
func (r *Repository) GetLastTravelRequestByUser(userID uint) (ds.Travelrequest, error) {
	var request ds.Travelrequest
	if err := r.db.Where("userid = ?", userID).Last(&request).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ds.Travelrequest{}, nil
		}
		return ds.Travelrequest{}, err
	}
	return request, nil
}

// Создать новую заявку
func (r *Repository) CreateTravelRequest(request ds.Travelrequest) error {
	if err := r.db.Create(&request).Error; err != nil {
		return err
	}

	// Проверьте значение request.Travelrequestid здесь.
	fmt.Printf("Created request with ID: %d\n", request.Travelrequestid)

	return nil
}

// Добавить дорогу к заявке
func (r *Repository) AddRoadToTravelRequest(requestID, roadID uint) error {
	travelRequestRoad := ds.Travelrequestroad{
		Travelrequestid: requestID,
		Roadid:          roadID,
	}

	if err := r.db.Create(&travelRequestRoad).Error; err != nil {
		return err
	}

	return nil
}

// GetRoadsByTravelRequest возвращает список дорог, связанных с заявкой пользователя.
// GetRoadsByTravelRequest возвращает список дорог, связанных с заявкой.
func (r *Repository) GetRoadsByTravelRequest(travelRequestID uint, requestStatus string) ([]ds.Road, error) {
	var roads []ds.Road

	// Здесь вы должны выполнить запрос, чтобы получить все дороги, связанные с указанной заявкой.
	// Пример:
	err := r.db.Raw("SELECT r.* FROM roads r JOIN travelrequestroads tr ON r.roadid = tr.roadid WHERE tr.travelrequestid = ? AND r.Statusroad = ?", travelRequestID, "active").Scan(&roads).Error
	if err != nil {
		return nil, err
	}

	return roads, nil
}

func (r *Repository) IsRoadAlreadyAddedToRequest(requestID, roadID uint) (bool, error) {
	var count int64
	if err := r.db.Model(&ds.Travelrequestroad{}).
		Where("roadid = ? AND travelrequestid = ?", roadID, requestID).
		Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}
func (r *Repository) FindOrCreateRequest(userID int) (uint, error) {
	// Попробуем сначала найти заявку с статусом "introduced" для данного пользователя.
	var existingRequest ds.Travelrequest
	if err := r.db.Where("userid = ? AND requeststatus = ?", userID, "introduced").First(&existingRequest).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Если заявка не найдена, создаем новую заявку для пользователя.
			newRequest := ds.Travelrequest{
				Userid:        uint(userID),
				Requeststatus: "introduced",
			}
			if err := r.CreateTravelRequest(newRequest); err != nil {
				// Ошибка при создании новой заявки
				fmt.Printf("Ошибка при создании новой заявки: %v", err)
				return 0, err
			}

			// Возвращаем функцию саму себя снова, чтобы повторно попытаться найти или создать заявку
			return r.FindOrCreateRequest(userID)
		}
		return 0, err
	}

	// Существующая заявка успешно найдена, возвращаем её ID
	fmt.Printf("Найдена существующая заявка с ID: %d", existingRequest.Travelrequestid)
	return existingRequest.Travelrequestid, nil
}

// DeleteRoadFromTravelRequest удаляет связь между дорогой и заявкой.
func (r *Repository) DeleteRoadFromTravelRequest(requestID uint, roadID uint) error {
	query := "DELETE FROM travelrequestroads WHERE travelrequestid = ? AND roadid = ?"
	result := r.db.Exec(query, requestID, roadID)

	if result.Error != nil {
		return result.Error
	}
	return nil
}

// IsRoadConnectedToRequest проверяет, существует ли связь между дорогой и заявкой.
func (r *Repository) IsRoadConnectedToRequest(requestID, roadID uint) bool {
	var count int64
	if err := r.db.Model(&ds.Travelrequestroad{}).
		Where("travelrequestid = ? AND roadid = ?", requestID, roadID).
		Count(&count).Error; err != nil {
		return false
	}

	return count > 0
}

func (r *Repository) AddConsultationImage(id int, imageBytes []byte, contentType string) error {
	// Удаление существующего изображения (если есть)
	err := r.minioClient.RemoveServiceImage(id)
	if err != nil {
		return err
	}

	// Загрузка нового изображения в MinIO
	imageURL, err := r.minioClient.UploadServiceImage(id, imageBytes, contentType)
	if err != nil {
		return err
	}

	// Обновление информации об изображении в БД (например, ссылки на MinIO)
	err = r.db.Model(&ds.Road{}).Where("roadid = ?", id).Update("image", imageURL).Error
	if err != nil {
		return err
	}

	return nil
}

// Метод для получения id заявки со статусом "introduced".
func (r *Repository) GetRequestIdWithStatusAndUser(status string, userID uint) uint {
	var requestID uint

	// Ваш запрос к базе данных для получения id заявки с указанным статусом для конкретного пользователя.
	if err := r.db.Model(&ds.Travelrequest{}).Where("Requeststatus = ? AND Userid = ?", status, userID).Select("Travelrequestid").First(&requestID).
		Error; err != nil {
		// Если запись не найдена, возвращаем 0.
		return 0
	}

	return requestID
}

func (r *Repository) Login() (*ds.User, error) {
	return nil, nil
}

func (r *Repository) Register(user *ds.User) error {
	return r.db.Create(user).Error
}

func (r *Repository) Logout() (*ds.User, error) {
	return nil, nil
}

func (r *Repository) GetUserByLogin(login string) (*ds.User, error) {
	user := &ds.User{
		Login: login,
	}

	err := r.db.Where(user).First(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *Repository) GetTravelRequestByID2(id uint) (ds.Travelrequest, error) {
	var request ds.Travelrequest
	err := r.db.First(&request, id).Error
	if err != nil {
		return ds.Travelrequest{}, err
	}
	return request, nil
}
