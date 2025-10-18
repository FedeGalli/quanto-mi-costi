package utils

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/paymentintent"
)

// Request/Response structs
type CreatePaymentIntentRequest struct {
	Amount   int64  `json:"amount" binding:"required"`
	Currency string `json:"currency" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Address  struct {
		Line1      string `json:"line1" binding:"required"`
		City       string `json:"city" binding:"required"`
		PostalCode string `json:"postal_code" binding:"required"`
		Country    string `json:"country" binding:"required"`
	} `json:"address" binding:"required"`
	Uid                string `json:"uid" binding:"required"`
	PaymentMethodTypes []*string
}

type ProcessPaymentResponse struct {
	Success         bool   `json:"success"`
	PaymentIntentID string `json:"payment_intent_id,omitempty"`
	Status          string `json:"status"`
	Message         string `json:"message"`
	RequiresAction  bool   `json:"requires_action,omitempty"`
	ClientSecret    string `json:"client_secret,omitempty"`
}

type CreatePaymentIntentResponse struct {
	ClientSecret   string `json:"client_secret"`
	DpmCheckerLink string `json:"dpm_checker_link,omitempty"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// Environment variables
var (
	stripeSecretKey string
)

var FirebaseClient *firestore.Client

func Init() {
	// Load environment variables
	stripeSecretKey := os.Getenv("STRIPE_SECRET_KEY")
	if stripeSecretKey == "" {
		log.Println("Warning: .env file not found")
	}

	var err error
	// Set Stripe API key
	stripe.Key = stripeSecretKey

	FirebaseClient, err = InitFirebase()
	if err != nil {
		log.Fatalf("Failed to initialize Firebase: %v", err)
	}
}

func getEnvVar(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func CreatePaymentIntent(c *gin.Context) {
	var req CreatePaymentIntentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	if req.Amount != 399 {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_amount",
			Message: fmt.Sprintf("Invalid amount: %d", req.Amount),
		})
		return
	}

	// Create the PaymentIntent WITHOUT a hardcoded payment method
	// This allows the frontend to collect the payment method
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(req.Amount),
		Currency: stripe.String(req.Currency),
		Shipping: &stripe.ShippingDetailsParams{
			Address: &stripe.AddressParams{
				Line1:      stripe.String(req.Address.Line1),
				City:       stripe.String(req.Address.City),
				PostalCode: stripe.String(req.Address.PostalCode),
				Country:    stripe.String(req.Address.Country),
			},
			Name: stripe.String(req.Name),
		},
		PaymentMethodTypes: []*string{stripe.String("card")},
		ReceiptEmail:       stripe.String(req.Email),
		// Remove the hardcoded PaymentMethod to allow frontend card collection
		// PaymentMethod:      stripe.String("pm_card_visa"), // REMOVED THIS LINE
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		log.Printf("Error creating payment intent: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "payment_intent_creation_failed",
			Message: "Unable to create payment intent",
		})
		return
	}

	log.Printf("Payment Intent created with status: %s", pi.Status)

	// For a newly created PaymentIntent without a payment method,
	// it should have status "requires_payment_method"
	// Return the client secret so frontend can collect payment method and confirm
	c.JSON(http.StatusOK, CreatePaymentIntentResponse{
		ClientSecret: pi.ClientSecret,
	})
}

// Add a separate endpoint to handle payment confirmation after 3D Secure
func ConfirmPayment(c *gin.Context) {
	type ConfirmPaymentRequest struct {
		PaymentIntentID string `json:"payment_intent_id" binding:"required"`
		Uid             string `json:"uid" binding:"required"`
	}

	var req ConfirmPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	// Retrieve the payment intent to check its current status
	pi, err := paymentintent.Get(req.PaymentIntentID, nil)
	if err != nil {
		log.Printf("Error retrieving payment intent: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "payment_retrieval_failed",
			Message: "Unable to retrieve payment intent",
		})
		return
	}

	log.Printf("Payment Intent status: %s", pi.Status)

	// Handle different payment statuses
	switch pi.Status {
	case stripe.PaymentIntentStatusSucceeded:
		// Payment succeeded, update user's pro status
		err := HandleSuccessfulPayment(FirebaseClient, req.Uid)
		if err != nil {
			log.Printf("Error updating user pro status: %v", err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "user_update_failed",
				Message: "Payment succeeded but failed to update user status",
			})
			return
		}

		c.JSON(http.StatusOK, ProcessPaymentResponse{
			Success:         true,
			PaymentIntentID: pi.ID,
			Status:          string(pi.Status),
			Message:         "Payment completed successfully! Pro access activated.",
		})

	case stripe.PaymentIntentStatusRequiresAction:
		c.JSON(http.StatusOK, ProcessPaymentResponse{
			Success:         false,
			PaymentIntentID: pi.ID,
			Status:          string(pi.Status),
			Message:         "Payment requires additional authentication",
			RequiresAction:  true,
			ClientSecret:    pi.ClientSecret,
		})

	default:
		c.JSON(http.StatusOK, ProcessPaymentResponse{
			Success:         false,
			PaymentIntentID: pi.ID,
			Status:          string(pi.Status),
			Message:         fmt.Sprintf("Payment status: %s", pi.Status),
		})
	}
}
