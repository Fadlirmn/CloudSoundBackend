package repository

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

func InitDB() (*firestore.Client, error) {
	ctx := context.Background()
	
	var app *firebase.App
	var err error

	// Prioritaskan FIREBASE_SERVICE_ACCOUNT_JSON (string JSON) jika ada
	// Jika tidak ada, baru cek FIREBASE_SERVICE_ACCOUNT_PATH (path file)
	serviceAccountJSON := os.Getenv("FIREBASE_SERVICE_ACCOUNT_JSON")
	serviceAccountPath := os.Getenv("FIREBASE_SERVICE_ACCOUNT_PATH")

	if serviceAccountJSON != "" {
		opt := option.WithCredentialsJSON([]byte(serviceAccountJSON))
		app, err = firebase.NewApp(ctx, nil, opt)
	} else if serviceAccountPath != "" {
		opt := option.WithCredentialsFile(serviceAccountPath)
		app, err = firebase.NewApp(ctx, nil, opt)
	} else {
		// Gunakan default credentials
		app, err = firebase.NewApp(ctx, nil)
	}

	if err != nil {
		return nil, err
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, err
	}

	log.Println("Firebase Firestore connected successfully")
	return client, nil
}