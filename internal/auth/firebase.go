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
    credPath := os.Getenv("FIREBASE_CREDENTIALS")
    if credPath == "" {
        return fmt.Errorf("FIREBASE_CREDENTIALS environment variable is not set")
    }

    opt := option.WithCredentialsFile(credPath)
    app, err := firebase.NewApp(context.Background(), nil, opt)
    if err != nil {
        return fmt.Errorf("error initializing app: %v", err)
    }

    client, err := app.Auth(context.Background())
    if err != nil {
        return fmt.Errorf("error getting Auth client: %v", err)
    }

    firebaseAuth = client
    return nil
}

func VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
    if firebaseAuth == nil {
        return nil, fmt.Errorf("firebase auth client not initialized")
    }
    
    token, err := firebaseAuth.VerifyIDToken(ctx, idToken)
    if err != nil {
        return nil, fmt.Errorf("error verifying ID token: %v", err)
    }
    
    return token, nil
}