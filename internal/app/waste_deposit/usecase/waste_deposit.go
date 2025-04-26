package usecase

import (
	"errors"
	"hackfest-uc/internal/app/waste_deposit/repository"
	"hackfest-uc/internal/domain/dto"
	"hackfest-uc/internal/domain/entity"
	"time"

	"github.com/google/uuid"
)

type WasteDepositUsecaseItf interface {
	CreateDeposit(userId uuid.UUID, req dto.DepositRequest) (*dto.DepositResponse, error)
	GetUserDeposits(userId uuid.UUID) ([]dto.DepositHistory, error)
	GetUserReward(userId uuid.UUID) ([]dto.DepositReward, error)
}

type WasteDepositUsecase struct {
	wasteDepositRepo repository.WasteDepositMySQLItf
}

func NewWasteDepositUsecase(wasteDepositRepo repository.WasteDepositMySQLItf) WasteDepositUsecaseItf {
	return &WasteDepositUsecase{
		wasteDepositRepo: wasteDepositRepo,
	}
}

func (u WasteDepositUsecase) CreateDeposit(userId uuid.UUID, req dto.DepositRequest) (*dto.DepositResponse, error) {
	reward := calculateReward(req.WasteType, req.WasteWeight)

	if req.WasteWeight <= 0 {
		return nil, errors.New("waste weight must be more than 0")
	}

	if req.WasteWeight > 1000 {
		return nil, errors.New("maximum waste weight 1000 kg")
	}

	deposit := entity.WasteDeposit{
		DepositId:    uuid.New(),
		UserId:       userId,
		Name:         req.Name,
		WasteType:    req.WasteType,
		WasteWeight:  req.WasteWeight,
		Reward:       reward,
		PickupMethod: req.PickupMethod,
		Status:       "Completed",
		PickupDate:   time.Now().UTC(),
	}

	if err := u.wasteDepositRepo.Create(deposit); err != nil {
		return nil, err
	}
	return &dto.DepositResponse{
		DepositId:    deposit.DepositId,
		Name:         deposit.Name,
		WasteType:    deposit.WasteType,
		WasteWeight:  deposit.WasteWeight,
		Reward:       deposit.Reward,
		PickupMethod: deposit.PickupMethod,
		Status:       deposit.Status,
		PickupDate:   deposit.PickupDate,
	}, nil
}

func (u WasteDepositUsecase) GetUserDeposits(userId uuid.UUID) ([]dto.DepositHistory, error) {
	deposits, err := u.wasteDepositRepo.GetByUserId(userId)
	if err != nil {
		return nil, err
	}

	var responses []dto.DepositHistory
	for _, d := range deposits {
		responses = append(responses, dto.DepositHistory{
			WasteType:   d.WasteType,
			WasteWeight: d.WasteWeight,
			Status:      d.Status,
			PickupDate:  d.PickupDate,
		})
	}

	return responses, nil
}

func (u WasteDepositUsecase) GetUserReward(userId uuid.UUID) ([]dto.DepositReward, error) {
	deposits, err := u.wasteDepositRepo.GetByUserId(userId)
	if err != nil {
		return nil, err
	}

	var responses []dto.DepositReward
	for _, d := range deposits {
		responses = append(responses, dto.DepositReward{
			Reward:     d.Reward,
			PickupDate: d.PickupDate,
		})
	}

	return responses, nil
}

func calculateReward(wasteType string, weight float64) float64 {
	rates := map[string]float64{
		"Limbah Organik Basah":  2000,
		"Limbah Organik Kering": 3000,
		"Limbah Campuran":       1500,
	}

	return rates[wasteType] * weight
}
