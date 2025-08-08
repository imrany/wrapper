package v1

import (
    "context"
    "log"

    "google.golang.org/genai"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
    v1pb "github.com/imrany/wrapper/proto/gen/api/v1"
)

func (s *APIV1Service) GenAi(ctx context.Context, req *v1pb.GenAiRequest) (*v1pb.GenAiResponse, error) {
    if req.Prompt == "" {
        return nil, status.Error(codes.InvalidArgument, "prompt cannot be empty") 
    }

    client, err := genai.NewClient(ctx, &genai.ClientConfig{
        APIKey:  s.APIKey,
        Backend: genai.BackendGeminiAPI,
    })
    if err != nil {
        return nil, status.Error(codes.Canceled,"Gemini generation failed")
    }

    result, err := client.Models.GenerateContent(
        ctx,
        "gemini-2.0-flash",
        genai.Text(req.Prompt),
        nil,
    )
    if err != nil {
        log.Printf("Gemini generation failed: %v", err)
        return nil, status.Error(codes.Canceled,"Gemini generation failed")
    }

    return &v1pb.GenAiResponse{
        Prompt:   req.Prompt,
        Response: result.Text(),
    }, nil
}
