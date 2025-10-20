package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

// Initialize Firebase (call this once in your main function or init)
func InitFirebase() (*firestore.Client, error) {
	ctx := context.Background()

	var app *firebase.App
	var err error

	// Try loading credentials from environment variable first
	credsJSON := os.Getenv("PROJECT_CREDENTIALS")
	if credsJSON != "" {
		log.Println("Initializing Firebase using credentials from environment variable...")
		opt := option.WithCredentialsJSON([]byte(credsJSON))
		app, err = firebase.NewApp(ctx, nil, opt)
	} else {
		// Fall back to local key file for local development
		log.Println("Initializing Firebase using local credentials file...")
		opt := option.WithCredentialsFile("secrets/quanto-mi-costi.json")
		app, err = firebase.NewApp(ctx, nil, opt)
	}

	if err != nil {
		return nil, err
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func IsUserStillPro(client *firestore.Client, uid string) error {
	ctx := context.Background()

	// Get reference to the user document
	userDocRef := client.Collection("users").Doc(uid)
	doc, err := userDocRef.Get(ctx)
	if err != nil {
		return fmt.Errorf("error getting document: %v", err)
	}

	isProRaw, err := doc.DataAt("is_pro")
	if err != nil {
		return fmt.Errorf("error getting pro: %v", err)
	}

	isPro := isProRaw.(bool)

	if isPro {
		proDueDateRaw, err := doc.DataAt("pro_due_date")
		if err != nil {
			return fmt.Errorf("error getting pro_due_date: %v", err)
		}

		proDueDate := proDueDateRaw.(int64)

		now := time.Now()

		today := int64(now.Year()*10000 + int(now.Month())*100 + now.Day())

		if today > proDueDate {
			UpdateUserProStatus(client, uid, false)
		}
	}

	log.Println("User pro status updated successfully")
	return nil
}

func UpdateUserProStatus(client *firestore.Client, uid string, isPro bool) error {
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

	// Adding 4 month for the PRO plan
	now := time.Now().AddDate(0, 4, 0)
	pro_due_date := int64(now.Year()*10000 + int(now.Month())*100 + now.Day())

	_, err = userDocRef.Update(ctx, []firestore.Update{
		{
			Path:  "pro_due_date",
			Value: pro_due_date,
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
	err := UpdateUserProStatus(client, userID, true)
	if err != nil {
		log.Printf("Failed to update user pro status for user %s: %v", userID, err)
		// Handle error appropriately
		return err
	}

	log.Printf("Successfully activated pro status for user: %s", userID)
	return nil
}
