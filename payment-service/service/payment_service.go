package service

import (
	"github.com/hasib-003/newsLetterMicroservice/payment-service/repository"
	"log"
)

type PaymentService struct {
	repo *repository.StripePaymentRepository
}

func NewPaymentService(repo *repository.StripePaymentRepository) *PaymentService {
	return &PaymentService{
		repo: repo,
	}
}
func (s *PaymentService) ProcessPayment(userId int, amount int) (string, error) {
	paymentID, err := s.repo.ProcessPayment(userId, amount)
	if err != nil {
		log.Println("payment process err:", err)
		return "", err
	}
	return paymentID, nil
}
