package v1

import(
	v1pb "github.com/imrany/wrapper/proto/gen/api/v1"
)

type APIV1Service struct {
    v1pb.UnimplementedAiServiceServer
    APIKey string
}