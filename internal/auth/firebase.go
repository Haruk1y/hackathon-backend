// internal/auth/firebase.go

package auth

import (
    "context"
    "fmt"
    "os"

    firebase "firebase.google.com/go/v4"
    "firebase.google.com/go/v4/auth"
    "google.golang.org/api/option"
)

var firebaseAuth *auth.Client

func InitFirebase() error {
    credPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
    if credPath == "" {
        return fmt.Errorf("GOOGLE_APPLICATION_CREDENTIALS environment variable is not set")
    }

    log.Printf("Initializing Firebase with credentials from: %s", credPath)

    // 認証情報ファイルの存在確認
    if _, err := os.Stat(credPath); os.IsNotExist(err) {
        return fmt.Errorf("credentials file not found at: %s", credPath)
    }

    opt := option.WithCredentialsFile(credPath)
    app, err := firebase.NewApp(context.Background(), nil, opt)
    if err != nil {
        return fmt.Errorf("error initializing Firebase app: %v", err)
    }

    client, err := app.Auth(context.Background())
    if err != nil {
        return fmt.Errorf("error getting Firebase Auth client: %v", err)
    }

    firebaseAuth = client
    log.Printf("Firebase initialized successfully")
    return nil
}