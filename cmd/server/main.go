package main

import (
	"getThumb/api/proto"
	pg "getThumb/database"
	"getThumb/internal"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := godotenv.Load(); err != nil {
		logrus.Fatal("Error loading .env file")
	}

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		logrus.Errorf("Failed to listen: %v", err)
	}

	db, err := pg.ConnectToDB()

	dbInstance, _ := db.DB()
	defer dbInstance.Close()

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

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGQUIT, syscall.SIGTERM)
	<-quit
	logrus.Printf("Shutting down server...")
	s.GracefulStop()
}
