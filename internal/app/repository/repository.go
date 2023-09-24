package repository

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"Road_services/internal/app/ds"
)

type Repository struct {
	db *gorm.DB
}

func New(dsn string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	return &Repository{
		db: db,
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

func (r *Repository) CreateRoad(road ds.Road) error {
	return r.db.Create(road).Error
}

func (r *Repository) GetAllRoads() ([]ds.Road, error) {
	var roads []ds.Road

	err := r.db.Find(&roads).Error
	if err != nil {
		return nil, err
	}

	return roads, nil
}
