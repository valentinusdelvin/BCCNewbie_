package usecase

import (
	MarketRepo "hackfest-uc/internal/app/market/repository"
	PaymentRepo "hackfest-uc/internal/app/payment/repository"
	"hackfest-uc/internal/domain/entity"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"gorm.io/gorm"
)

type InterPaymentUsecase interface {
	Purchase(payment entity.Payment) (string, error)
	Validate(MidtransNotifications map[string]interface{}) error
}

type PaymentUsecase struct {
	rp PaymentRepo.InterPaymentRepository
	rt MarketRepo.MarketMySQLItf
	db *gorm.DB
}

func NewPaymentUsecase(paymentRepo PaymentRepo.InterPaymentRepository, marketRepo MarketRepo.MarketMySQLItf, db *gorm.DB) InterPaymentUsecase {
	return &PaymentUsecase{
		rp: paymentRepo,
		rt: marketRepo,
		db: db,
	}
}

func (p *PaymentUsecase) Purchase(payment entity.Payment) (string, error) {
	product, err := p.rt.GetProductByID(payment.ProductID.String())
	if err != nil {
		return "", err
	}

	snapReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  payment.OrderID.String(),
			GrossAmt: int64(product.ProductPrice),
		},
		Items: &[]midtrans.ItemDetails{
			{
				ID:           webinar.ID.String(),
				Name:         webinar.WebinarName,
				Price:        int64(webinar.Price),
				Qty:          1,
				Category:     "Webinar VeloVent",
				MerchantName: "Velo Mom",
			},
		},
	}

	ServiceFee := midtrans.ItemDetails{
		ID:    "biaya_layanan",
		Name:  "Biaya Layanan",
		Price: 2000,
		Qty:   1,
	}

	*snapReq.Items = append(*snapReq.Items, ServiceFee)
	payment.Price = uint64(snapReq.TransactionDetails.GrossAmt)

	snapReq.TransactionDetails.GrossAmt += ServiceFee.Price
	payment.FinalPrice = uint64(snapReq.TransactionDetails.GrossAmt)

	paymentLink, paymentErr := snap.CreateTransactionUrl(snapReq)
	if paymentErr != nil {
		return "", paymentErr
	}
	payment.PaymentLink = paymentLink
	payment.ProductName = (*snapReq.Items)[0].Name

	err = p.prsc.CreatePayment(payment)
	if err != nil {
		return "", err
	}
	return paymentLink, nil
}

func (p *PaymentUsecase) Validate(MidtransNotifications map[string]interface{}) error {
	transactionStatus := MidtransNotifications["transaction_status"]
	orderID := MidtransNotifications["order_id"].(string)
	fraudStatus := MidtransNotifications["fraud_status"]

	err := p.db.Transaction(func(tx *gorm.DB) error {
		switch transactionStatus {
		case "capture":
			switch fraudStatus {
			case "challenge":
				return p.prsc.UpdatePaymentStatus(tx, "challenge", orderID)
			case "accept":
				if err := p.prsc.UpdatePaymentStatus(tx, "success", orderID); err != nil {
					return err
				}
				invoice, err := p.prsc.GetInvoice(orderID)
				if err != nil {
					return err
				}
				attendee := entity.WebinarAttendee{
					UserID:    invoice.UserID,
					WebinarID: invoice.ProductID,
				}

				if err := p.wrsc.CreateWebinarAttendee(tx, attendee); err != nil {
					return err
				}

				if err := p.wrsc.UpdateWebinarInfo(tx, invoice.ProductID); err != nil {
					return err
				}
				return nil
			}

		case "settlement":
			if err := p.prsc.UpdatePaymentStatus(tx, "success", orderID); err != nil {
				return err
			}
			invoice, err := p.prsc.GetInvoice(orderID)
			if err != nil {
				return err
			}
			attendee := entity.WebinarAttendee{
				UserID:    invoice.UserID,
				WebinarID: invoice.ProductID,
			}

			if err := p.wrsc.CreateWebinarAttendee(tx, attendee); err != nil {
				return err
			}

			if err := p.wrsc.UpdateWebinarInfo(tx, invoice.ProductID); err != nil {
				return err
			}
			return nil

		case "cancel", "expire":
			return p.prsc.UpdatePaymentStatus(tx, "failure", orderID)

		case "pending":
			return p.prsc.UpdatePaymentStatus(tx, "pending", orderID)

		case "deny":
			return p.prsc.UpdatePaymentStatus(tx, "denied", orderID)
		}

		return nil
	})

	if err != nil {
		return err
	}
	return nil
}
