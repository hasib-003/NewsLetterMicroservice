package handler

import (
	"context"
	payment "github.com/hasib-003/newsLetterMicroservice/payment-service/proto"
	"github.com/hasib-003/newsLetterMicroservice/payment-service/service"
)

type PaymentServiceHandler struct {
	payment.UnimplementedPaymentServiceServer
	PaymentService *service.PaymentService
}

func NewPaymentServiceHandler(paymentService *service.PaymentService) *PaymentServiceHandler {
	return &PaymentServiceHandler{
		PaymentService: paymentService,
	}
}
func (h *PaymentServiceHandler) ProcessPayment(ctx context.Context, request *payment.PaymentRequest) (*payment.PaymentResponse, error) {
	paymentId, err := h.PaymentService.ProcessPayment(int(request.UserId), int(request.Amount))
	if err != nil {
		return &payment.PaymentResponse{
			PaymentId: "",
			Message:   "payment failed" + err.Error(),
			Success:   false,
		}, nil
	}
	return &payment.PaymentResponse{
		PaymentId: paymentId,
		Message:   "payment Successful",
		Success:   true,
	}, nil
}
