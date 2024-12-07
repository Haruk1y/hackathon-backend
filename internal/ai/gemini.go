// internal/ai/gemini.go

package ai

import (
    "context"
    "fmt"
    "os"

    "github.com/google/generative-ai-go/genai"
    "google.golang.org/api/option"
)

var client *genai.Client
var genModel *genai.GenerativeModel

func InitGemini() error {
    ctx := context.Background()
    var err error
    client, err = genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
    if err != nil {
        return fmt.Errorf("failed to create Gemini client: %v", err)
    }

    genModel = client.GenerativeModel("gemini-pro")
    return nil
}

func SummarizeText(content string) (string, error) {
    ctx := context.Background()
    prompt := fmt.Sprintf("以下の文章を3行程度で要約してください:\n\n%s", content)

    resp, err := genModel.GenerateContent(ctx, genai.Text(prompt))
    if err != nil {
        return "", fmt.Errorf("failed to generate summary: %v", err)
    }

    if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
        return "", fmt.Errorf("no content generated")
    }

    text, ok := resp.Candidates[0].Content.Parts[0].(genai.Text)
    if !ok {
        return "", fmt.Errorf("unexpected response format")
    }

    return string(text), nil
}

func Close() {
    if client != nil {
        client.Close()
    }
}