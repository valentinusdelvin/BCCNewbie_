package usecase

import (
	"hackfest-uc/internal/app/waste_deposit/repository"
	"hackfest-uc/internal/domain/dto"
	"hackfest-uc/internal/domain/entity"
	"time"

	"github.com/google/uuid"
)

type WasteDepositUsecaseItf interface {
	CreateDeposit(userId uuid.UUID, req dto.DepositRequest) (*dto.DepositResponse, error)
	GetUserDeposits(userId uuid.UUID) ([]dto.DepositResponse, error)
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

	deposit := entity.WasteDeposit{
		DepositId:    uuid.New(),
		UserId:       userId,
		Name:         req.Name,
		WasteType:    req.WasteType,
		WasteWeight:  req.WasteWeight,
		Reward:       reward,
		PickupMethod: req.PickupMethod,
		Status:       "Pending",
		PickupDate:   time.Now(),
	}

	if err := u.wasteDepositRepo.Create(deposit); err != nil {
		return nil, err
	}
	return &dto.DepositResponse{
		DepositId:   deposit.DepositId,
		WasteType:   deposit.WasteType,
		WasteWeight: deposit.WasteWeight,
		Reward:      deposit.Reward,
		Status:      deposit.Status,
	}, nil
}

func (u WasteDepositUsecase) GetUserDeposits(userId uuid.UUID) ([]dto.DepositResponse, error) {
	deposits, err := u.wasteDepositRepo.GetByUserId(userId)
	if err != nil {
		return nil, err
	}

	var responses []dto.DepositResponse
	for _, d := range deposits {
		responses = append(responses, dto.DepositResponse{
			DepositId:   d.DepositId,
			WasteType:   d.WasteType,
			WasteWeight: d.WasteWeight,
			Reward:      d.Reward,
			Status:      d.Status,
			PickupDate:  d.PickupDate,
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
