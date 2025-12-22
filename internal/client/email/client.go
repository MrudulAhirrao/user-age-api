package email

import (
	"context"
	"log"
	"time"
	breaker "user-age-api/internal/client"
	pb "user-age-api/internal/email/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type EmailClinet struct{
	service pb.EmailServiceClient
	conn *grpc.ClientConn
	cb *breaker.CircuitBreaker
}

func NewEmailClient(address string) (*EmailClinet, error){
	conn, err:= grpc.NewClient(address,grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil{
		return nil,err
	}

	client := pb.NewEmailServiceClient(conn)
	return &EmailClinet{
		service: client,
		conn: conn,
		cb : breaker.NewCircuitBreaker(3,10*time.Second),
	},nil
}

func (c *EmailClinet) SendActivationEmail(email string, token string) error {
	// Define the action we want to take
	action := func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		req := &pb.SendActivationEmailRequest{
			Email:           email,
			ActivationToken: token,
		}

		// Real Network Call
		resp, err := c.service.SendActivationEmail(ctx, req)
		if err != nil {
			return err // Failure triggers the breaker logic
		}
		
		if resp.Success {
			log.Printf("✅ Email sent to %s", email)
		}
		return nil
	}

	// EXECUTE WITH PROTECTION
	// The breaker runs the logic above. 
	// If breaker is OPEN, 'action' is never called, and it returns error instantly.
	return c.cb.Execute(action)
}
func (c *EmailClinet) Close(){
	c.conn.Close()
}