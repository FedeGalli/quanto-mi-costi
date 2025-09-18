package utils

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

// Initialize Firebase (call this once in your main function or init)
func InitFirebase() (*firestore.Client, error) {
	ctx := context.Background()

	// Option 1: Using service account key file
	opt := option.WithCredentialsFile("secrets/quanto-mi-costi.json")
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, err
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func updateUserProStatus(client *firestore.Client, uid string, isPro bool) error {
	ctx := context.Background()

	// Get reference to the user document
	userDocRef := client.Collection("users").Doc(uid)

	// Update the document
	_, err := userDocRef.Update(ctx, []firestore.Update{
		{
			Path:  "is_pro",
			Value: isPro,
		},
	})

	if err != nil {
		log.Printf("Error updating user pro status: %v", err)
		return err
	}

	log.Println("User pro status updated successfully")
	return nil
}

// Usage example in your payment processing function
func HandleSuccessfulPayment(client *firestore.Client, userID string) error {
	err := updateUserProStatus(client, userID, true)
	if err != nil {
		log.Printf("Failed to update user pro status for user %s: %v", userID, err)
		// Handle error appropriately
		return err
	}

	log.Printf("Successfully activated pro status for user: %s", userID)
	return nil
}
