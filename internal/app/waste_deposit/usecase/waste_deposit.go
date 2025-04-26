package usecase

type WasteDepositUsecaseItf interface {}

type WasteDepositUsecase struct {}

func NewWasteDepositUsecase() WasteDepositUsecaseItf {
    return &WasteDepositUsecase{}
}
