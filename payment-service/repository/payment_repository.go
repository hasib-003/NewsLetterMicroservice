package repository

import (
	"errors"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/paymentintent"
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
	if intent.Status == "requires_confirmation" {
		intent, err = paymentintent.Confirm(intent.ID, nil)
		if err != nil {
			log.Printf("Payment intent failed with %v", err)
		}
	}
	if intent.Status != "succeeded" {
		log.Printf("Payment intent failed with %v", intent)
		return "", errors.New("Payment intent failed ")
	}
	log.Printf("Payment intent succeeded with %v", intent)
	return string(intent.Status), nil
}
