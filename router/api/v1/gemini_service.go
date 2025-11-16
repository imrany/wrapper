package v1

import (
	"context"
	"fmt"
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

	switch strings.Split(s.Model, "-")[0] {
	case "gemini":
		config := geminiwrapper.GeminiClientConfig{
			APIKey: s.APIKey,
			Model:  s.Model,
		}
		result, err := config.GenerateGeminiContent(ctx, req.Prompt)

		if err != nil {
			log.Printf("Gemini generation failed: %v", err)
			return nil, status.Error(codes.Canceled, "Gemini generation failed")
		}

		return &v1pb.GenAiResponse{
			Prompt:   req.Prompt,
			Response: result,
		}, nil
	case "openai":
		config := openaiwrapper.OpenAIClient{
			APIKey: s.APIKey,
			Model:  s.Model,
		}
		result, err := config.GenerateOpenAIContent(ctx, req.Prompt)

		if err != nil {
			log.Printf("OpenAI generation failed: %v", err)
			return nil, status.Error(codes.Canceled, "OpenAI generation failed")
		}

		return &v1pb.GenAiResponse{
			Prompt:   req.Prompt,
			Response: result,
		}, nil
	default:
		log.Printf("Unsupported model: %s", s.Model)
		return &v1pb.GenAiResponse{
			Prompt:   req.Prompt,
			Response: fmt.Sprintf("Unsupported model: %s", s.Model),
		}, nil
	}
}
