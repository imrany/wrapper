package v1

import (
	"log/slog"

	v1pb "github.com/imrany/wrapper/proto/gen/api/v1"
)

type APIV1Service struct {
	v1pb.UnimplementedAiServiceServer
	APIKey string
	Model  string
	Logger *slog.Logger
}
