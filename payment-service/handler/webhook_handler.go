package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/webhook"
	"log"
	"net/http"
)

type WebhookHandler struct {
	StripeWebhookSecret string
}

func NewWebhookHandler(secret string) *WebhookHandler {
	return &WebhookHandler{StripeWebhookSecret: secret}
}

func (h *WebhookHandler) HandleWebhook(c *gin.Context) {
	body, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read request body"})
		return
	}

	// Verify the webhook signature
	event, err := webhook.ConstructEvent(body, c.GetHeader("Stripe-Signature"), h.StripeWebhookSecret)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid webhook signature"})
		return
	}

	// Process the event
	switch event.Type {
	case "payment_intent.succeeded":
		// Unmarshal the event data into the expected PaymentIntent type
		var paymentIntent stripe.PaymentIntent
		if err := json.Unmarshal(event.Data.Raw, &paymentIntent); err != nil {
			log.Printf("Error unmarshaling PaymentIntent: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process PaymentIntent"})
			return
		}
		log.Printf("PaymentIntent succeeded: %s", paymentIntent.ID)

	case "payment_intent.payment_failed":
		// Unmarshal the event data into the expected PaymentIntent type
		var paymentIntent stripe.PaymentIntent
		if err := json.Unmarshal(event.Data.Raw, &paymentIntent); err != nil {
			log.Printf("Error unmarshaling PaymentIntent: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process PaymentIntent"})
			return
		}
		log.Printf("PaymentIntent failed: %s", paymentIntent.ID)

	default:
		log.Printf("Unhandled event type: %s", event.Type)
	}

	c.JSON(http.StatusOK, gin.H{"received": true})
}
