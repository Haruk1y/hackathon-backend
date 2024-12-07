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
        return fmt.Errorf("GOOGLE_APPLICATION_CREDENTIALS is not set")
    }

    fmt.Printf("Loading credentials from: %s\n", credPath)
    
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

func VerifyIDToken(idToken string) (*auth.Token, error) {
    token, err := firebaseAuth.VerifyIDToken(context.Background(), idToken)
    if err != nil {
        return nil, err
    }
    return token, nil
}