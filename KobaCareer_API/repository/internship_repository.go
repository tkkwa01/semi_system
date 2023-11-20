package repository

import (
	"gorm.io/gorm"
	"semi_systems/KobaCareer_API/domain"
)

type IInternshipRepository interface {
	GetAllInternships(internships *[]domain.Internships) error
	GetInternshipById(internship *domain.Internships, internshipId uint) error
	CreateInternship(internship *domain.Internships) error
	UpdateInternship(internship *domain.Internships, internshipId uint) error
	DeleteInternship(internshipId uint) error
	InternshipExist(internshipId uint) (bool, error)
}

type internshipRepository struct {
	db *gorm.DB
}

func NewInternshipRepository(db *gorm.DB) IInternshipRepository {
	return &internshipRepository{db}
}

func (i internshipRepository) GetAllInternships(internships *[]domain.Internships) error {
	if err := i.db.Find(internships).Error; err != nil {
		return err
	}
	return nil
}

func (i internshipRepository) GetInternshipById(internship *domain.Internships, internshipId uint) error {
	if err := i.db.First(internship, internshipId).Error; err != nil {
		return err
	}
	return nil
}

func (i internshipRepository) CreateInternship(internship *domain.Internships) error {
	if err := i.db.Create(internship).Error; err != nil {
		return err
	}
	return nil
}

func (i internshipRepository) UpdateInternship(internship *domain.Internships, internshipId uint) error {
	result := i.db.Model(internship).Where("id=?", internshipId).Updates(domain.Internships{
		Company:     internship.Company,
		Title:       internship.Title,
		Salary:      internship.Salary,
		Period:      internship.Period,
		Selection:   internship.Selection,
		Deadline:    internship.Deadline,
		Contributor: internship.Contributor,
		Detail:      internship.Detail,
		Occupation:  internship.Occupation,
		Flow:        internship.Flow,
		Method:      internship.Method,
		URL:         internship.URL,
		Duration:    internship.Duration,
	})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (i internshipRepository) DeleteInternship(internshipId uint) error {
	result := i.db.Where("id=?", internshipId).Delete(&domain.Internships{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (i internshipRepository) InternshipExist(internshipId uint) (bool, error) {
	var count int64
	if err := i.db.Model(&domain.Internships{}).Where(&domain.Internships{ID: internshipId}).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
