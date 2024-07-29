package internal

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/credentials/insecure"
	"os"
	"sync"
	"time"

	"getThumb/api/proto"
	"google.golang.org/grpc"
)

func downloadThumbnail(client proto.ThumbnailClient, url string, wg *sync.WaitGroup) {
	defer wg.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.GetThumbnail(ctx, &proto.Request{Url: url})
	if err != nil {
		logrus.Errorf("Failed to download thumbnail: %v", err)
		return
	}
	filename, err := ParseURL(url)
	err = os.WriteFile(fmt.Sprintf("./downloads/%s.jpg", filename), resp.Image, 0777)
	if err != nil {
		logrus.Errorf("Failed to save thumbnail: %v", err)
		return
	}
}

func RunClient(urls []string, async bool) error {
	conn, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to connect: %v", err))
	}
	defer conn.Close()

	client := proto.NewThumbnailClient(conn)

	os.Mkdir("./downloads", 0777)

	var wg sync.WaitGroup
	for _, url := range urls {
		wg.Add(1)
		if async {
			logrus.Println("downloading asynchronously")
			go downloadThumbnail(client, url, &wg)
		} else {
			downloadThumbnail(client, url, &wg)
		}
	}
	wg.Wait()
	return nil
}
