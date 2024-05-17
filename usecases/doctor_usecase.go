package usecases

import (
	"context"

	"github.com/tsanaativa/sehatin-backend-v0.1/entities"
	"github.com/tsanaativa/sehatin-backend-v0.1/repositories"
)

type DoctorUsecaseOpts struct {
	DoctorRepo repositories.DoctorRepository
}

type DoctorUsecase interface {
	GetDoctorById(ctx context.Context, doctorId int64) (*entities.Doctor, error)
	UpdateDoctor(ctx context.Context, doctor entities.Doctor) error
	ToggleDoctorIsOnline(ctx context.Context, doctorId int64) error
	DeleteDoctor(ctx context.Context, doctorId int64) error
	GetAllDoctor(ctx context.Context, params entities.DoctorParams, isPublic bool) ([]entities.Doctor, *entities.PaginationInfo, error)
}

type DoctorUsecaseImpl struct {
	DoctorRepository repositories.DoctorRepository
}

func NewDoctorUsecaseImpl(doctorOpts *DoctorUsecaseOpts) DoctorUsecase {
	return &DoctorUsecaseImpl{
		DoctorRepository: doctorOpts.DoctorRepo,
	}
}

func (u *DoctorUsecaseImpl) GetDoctorById(ctx context.Context, doctorId int64) (*entities.Doctor, error) {
	doctor, err := u.DoctorRepository.FindOneById(ctx, doctorId)
	if err != nil {
		return nil, err
	}

	return doctor, nil
}

func (u *DoctorUsecaseImpl) UpdateDoctor(ctx context.Context, doctor entities.Doctor) error {
	_, err := u.DoctorRepository.FindOneById(ctx, doctor.Id)
	if err != nil {
		return err
	}

	err = u.DoctorRepository.UpdateOne(ctx, doctor)
	if err != nil {
		return err
	}

	return nil
}

func (u *DoctorUsecaseImpl) ToggleDoctorIsOnline(ctx context.Context, doctorId int64) error {
	doctor, err := u.DoctorRepository.FindOneById(ctx, doctorId)
	if err != nil {
		return err
	}

	doctor.IsOnline = !doctor.IsOnline

	err = u.DoctorRepository.UpdateIsOnline(ctx, *doctor)
	if err != nil {
		return err
	}

	return nil
}

func (u *DoctorUsecaseImpl) DeleteDoctor(ctx context.Context, doctorId int64) error {
	err := u.DoctorRepository.Delete(ctx, doctorId)
	if err != nil {
		return err
	}

	return nil
}

func (u *DoctorUsecaseImpl) GetAllDoctor(ctx context.Context, params entities.DoctorParams, isPublic bool) ([]entities.Doctor, *entities.PaginationInfo, error) {
	var doctors []entities.Doctor
	var totalData int
	var err error

	if isPublic {
		doctors, totalData, err = u.DoctorRepository.FindAll(ctx, params, true)
		if err != nil {
			return nil, nil, err
		}
	} else {
		doctors, totalData, err = u.DoctorRepository.FindAll(ctx, params, false)
		if err != nil {
			return nil, nil, err
		}
	}

	totalPage := totalData / params.Limit
	if totalData%params.Limit > 0 {
		totalPage++
	}

	pagination := entities.PaginationInfo{
		Page:      params.Page,
		Limit:     params.Limit,
		TotalData: totalData,
		TotalPage: totalPage,
	}

	return doctors, &pagination, nil
}
