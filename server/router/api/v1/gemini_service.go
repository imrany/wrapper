package v1

import (
    "context"
    "fmt"
    "log"

    "google.golang.org/genai"
    "google.golang.org/genai/google"
    pb "proto/gen/api/v1"
)

type AiService struct {
    pb.UnimplementedAiServiceServer
    APIKey string
}

func (s *AiService) GenAi(ctx context.Context, req *pb.GenAiRequest) (*pb.GenAiResponse, error) {
    if req.Prompt == "" {
        return nil, fmt.Errorf("prompt is required")
    }

    client, err := genai.NewClient(ctx, &genai.ClientOptions{
        APIKey: s.APIKey,
    })
    if err != nil {
        log.Printf("Failed to create Gemini client: %v", err)
        return nil, err
    }
    defer client.Close()

    result, err := client.Models.GenerateContent(
        ctx,
        "gemini-1.5-flash", // or "gemini-2.5-flash" if supported
        genai.Text(req.Prompt),
        &genai.GenerateOptions{
            Model: &google.Model{
                Name: "gemini-1.5-flash",
            },
        },
    )
    if err != nil {
        log.Printf("Gemini generation failed: %v", err)
        return nil, err
    }

    return &pb.GenAiResponse{
        Prompt:   req.Prompt,
        Response: result.Text(),
    }, nil
}
