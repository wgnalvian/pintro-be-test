package service

import (
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"go.mongodb.org/mongo-driver/mongo"
	"wgnalvian.com/payment-server/exception"
)

type TransactionService struct {
	Db *mongo.Database
}

func (t *TransactionService) CreateMidtransTransaction(orderID string, amount int) (string, error) {
	// Buat objek transaksi
	snapClient := snap.Client{}
	snapClient.New("YOUR_SERVER_KEY", midtrans.Sandbox)

	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: int64(amount),
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
	}

	// Buat transaksi Midtrans
	resp, err := snapClient.CreateTransaction(req)
	if err != nil {
		exception.LogError(err)
		return "", err
	}

	return resp.RedirectURL, nil
}
