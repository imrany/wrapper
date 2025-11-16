package v1

import (
	"context"
	"log"
	"strings"

	geminiwrapper "github.com/imrany/wrapper/pkg/gemini"
	openaiwrapper "github.com/imrany/wrapper/pkg/openai"
	v1pb "github.com/imrany/wrapper/proto/gen/api/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *APIV1Service) GenAi(ctx context.Context, req *v1pb.GenAiRequest) (*v1pb.GenAiResponse, error) {
	if req.Prompt == "" {
		return nil, status.Error(codes.InvalidArgument, "prompt cannot be empty")
	}

	if ctx.Err() != nil {
		return nil, status.FromContextError(ctx.Err()).Err()
	}

	// Extract provider from model name (case-insensitive)
	modelParts := strings.Split(s.Model, "-")
	if len(modelParts) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid model format: %s", s.Model)
	}
	provider := strings.ToLower(modelParts[0])

	switch provider {
	case "gemini":
		config := geminiwrapper.GeminiClientConfig{
			APIKey: s.APIKey,
			Model:  s.Model,
		}
		result, err := config.GenerateGeminiContent(ctx, req.Prompt)
		if err != nil {
			log.Printf("Gemini generation failed: %v", err)
			// Check if it's a context error
			if ctx.Err() != nil {
				return nil, status.FromContextError(ctx.Err()).Err()
			}
			// Return the actual error to the client for debugging
			return nil, status.Errorf(codes.Internal, "Gemini generation failed: %v", err)
		}
		return &v1pb.GenAiResponse{
			Prompt:   req.Prompt,
			Response: result,
		}, nil

	case "gpt", "o1":
		config := openaiwrapper.OpenAIClient{
			APIKey: s.APIKey,
			Model:  s.Model,
		}
		result, err := config.GenerateOpenAIContent(ctx, req.Prompt)
		if err != nil {
			log.Printf("OpenAI generation failed: %v", err)
			// Check if it's a context error
			if ctx.Err() != nil {
				return nil, status.FromContextError(ctx.Err()).Err()
			}
			// Return the actual error to the client for debugging
			return nil, status.Errorf(codes.Internal, "OpenAI generation failed: %v", err)
		}
		return &v1pb.GenAiResponse{
			Prompt:   req.Prompt,
			Response: result,
		}, nil

	default:
		log.Printf("Unsupported model: %s", s.Model)
		return nil, status.Errorf(codes.InvalidArgument, "unsupported model: %s", s.Model)
	}
}
