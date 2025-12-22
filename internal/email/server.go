package email

import ( 
	"context"
	"log"
	
	pb "user-age-api/internal/email/proto"
)

type Server struct {
	pb.UnimplementedEmailServiceServer
}

func (s *Server) SendActivationEmail(ctx context.Context, req *pb.SendActivationEmailRequest)(*pb.SendActivationEmailResponse, error){
	log.Printf("Sending Email To: %s | Token: %s", req.Email,req.ActivationToken)

	return &pb.SendActivationEmailResponse{
		Success: true,
	},nil
}