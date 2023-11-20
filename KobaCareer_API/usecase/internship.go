package usecase

import (
	"errors"
	"semi_systems/KobaCareer_API/domain"
	"semi_systems/KobaCareer_API/repository"
	"semi_systems/KobaCareer_API/validator"
)

type IInternshipUsecase interface {
	GetAllInternships() ([]domain.InternshipsResponse, error)
	GetInternshipById(internshipId uint) (domain.InternshipsResponse, error)
	CreateInternship(internship domain.Internships) (domain.InternshipsResponse, error)
	UpdateInternship(internship domain.Internships, internshipId uint) (domain.InternshipsResponse, error)
	DeleteInternship(internshipId uint) error
}

type internshipUsecase struct {
	ir repository.IInternshipRepository
	iv validator.IInternshipValidator
}

func NewInternshipUsecase(ir repository.IInternshipRepository, iv validator.IInternshipValidator) IInternshipUsecase {
	return &internshipUsecase{ir, iv}
}

func (iu *internshipUsecase) GetAllInternships() ([]domain.InternshipsResponse, error) {
	internships := []domain.Internships{}
	if err := iu.ir.GetAllInternships(&internships); err != nil {
		return nil, err
	}
	resInternships := []domain.InternshipsResponse{}
	for _, v := range internships {
		internship := domain.InternshipsResponse{
			ID:          v.ID,
			Company:     v.Company,
			Title:       v.Title,
			Salary:      v.Salary,
			Period:      v.Period,
			Selection:   v.Selection,
			Deadline:    v.Deadline,
			Contributor: v.Contributor,
			Detail:      v.Detail,
			Occupation:  v.Occupation,
			Flow:        v.Flow,
			Method:      v.Method,
			URL:         v.URL,
			Duration:    v.Duration,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
		}
		resInternships = append(resInternships, internship)
	}
	return resInternships, nil
}

func (iu *internshipUsecase) GetInternshipById(internshipId uint) (domain.InternshipsResponse, error) {
	internship := domain.Internships{}
	if err := iu.ir.GetInternshipById(&internship, internshipId); err != nil {
		return domain.InternshipsResponse{}, err
	}
	resInternship := domain.InternshipsResponse{
		ID:          internship.ID,
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
		CreatedAt:   internship.CreatedAt,
		UpdatedAt:   internship.UpdatedAt,
	}
	return resInternship, nil
}

func (iu *internshipUsecase) CreateInternship(internship domain.Internships) (domain.InternshipsResponse, error) {
	if err := iu.iv.InternshipValidate(internship); err != nil {
		return domain.InternshipsResponse{}, err
	}
	if err := iu.ir.CreateInternship(&internship); err != nil {
		return domain.InternshipsResponse{}, err
	}
	resInternship := domain.InternshipsResponse{
		ID:          internship.ID,
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
		CreatedAt:   internship.CreatedAt,
		UpdatedAt:   internship.UpdatedAt,
	}
	return resInternship, nil
}

func (iu *internshipUsecase) UpdateInternship(internship domain.Internships, internshipId uint) (domain.InternshipsResponse, error) {
	if err := iu.iv.InternshipValidate(internship); err != nil {
		return domain.InternshipsResponse{}, err
	}
	if err := iu.ir.UpdateInternship(&internship, internshipId); err != nil {
		return domain.InternshipsResponse{}, err
	}
	resInternship := domain.InternshipsResponse{
		ID:          internship.ID,
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
		CreatedAt:   internship.CreatedAt,
		UpdatedAt:   internship.UpdatedAt,
	}
	return resInternship, nil
}

func (iu *internshipUsecase) DeleteInternship(internshipId uint) error {
	id, err := iu.ir.InternshipExist(internshipId)
	if err != nil {
		return err
	}
	if !id {
		return errors.New("id is not exist")
	}

	if err := iu.ir.DeleteInternship(internshipId); err != nil {
		return err
	}
	return nil
}
