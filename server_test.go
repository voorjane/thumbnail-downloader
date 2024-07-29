package test

import (
	"context"
	"getThumb/api/proto"
	pg "getThumb/database"
	"getThumb/internal"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"testing"
)

func init() {
	if err := godotenv.Load(); err != nil {
		logrus.Fatal("Error loading .env file")
	}

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		logrus.Errorf("Failed to listen: %v", err)
	}

	db, err := pg.ConnectToDB()

	if err != nil {
		logrus.Errorf("error connecting to database: %v", err)
	}

	s := grpc.NewServer()
	proto.RegisterThumbnailServer(s, internal.NewServer(db))

	logrus.Printf("Server is running at %s", lis.Addr())
	go func() {
		if err = s.Serve(lis); err != nil {
			logrus.Errorf("Failed to serve: %v", err)
		}
	}()
}

func TestThumbnail(t *testing.T) {
	conn, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := proto.NewThumbnailClient(conn)
	request := "https://www.youtube.com/watch?v=dQw4w9WgXcQ"
	_, err = client.GetThumbnail(context.Background(), &proto.Request{Url: request})
	if err != nil {
		t.Fatalf("failed to get thumbnail: %v", err)
	}
}

func TestWrongUrl(t *testing.T) {
	conn, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := proto.NewThumbnailClient(conn)
	request := "https://www.you.com/123"
	_, err = client.GetThumbnail(context.Background(), &proto.Request{Url: request})
	if err == nil {
		t.Fatalf("expected error")
	}
	logrus.Printf("error found as expected: %v", err)
}

func TestWrongId(t *testing.T) {
	conn, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := proto.NewThumbnailClient(conn)
	request := "https://www.youtube.com/watch?v=dQw4w9"
	_, err = client.GetThumbnail(context.Background(), &proto.Request{Url: request})
	if err == nil {
		t.Fatalf("expected error")
	}
	logrus.Printf("error found as expected: %v", err)
}
