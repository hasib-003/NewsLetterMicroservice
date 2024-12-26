package repository

import (
	"github.com/stripe/stripe-go/v72/paymentintent"

	"github.com/stripe/stripe-go/v72"
	"log"
)

type StripePaymentRepository struct {
	apiKey string
}

func NewStripePaymentRepository(apiKey string) *StripePaymentRepository {
	stripe.Key = apiKey
	return &StripePaymentRepository{apiKey: apiKey}
}
func (r *StripePaymentRepository) ProcessPayment(userId int, amount int) (string, error) {
	fakePaymentMethod := "pm_card_visa"
	params := &stripe.PaymentIntentParams{
		Amount:        stripe.Int64(int64(amount)),
		Currency:      stripe.String(string(stripe.CurrencyUSD)),
		PaymentMethod: &fakePaymentMethod,
	}
	intent, err := paymentintent.New(params)
	if err != nil {
		log.Printf("Payment intent failed with %v", err)
		return "", err
	}
	log.Printf("Payment intent is %v", intent)
	return intent.ID, nil

}
