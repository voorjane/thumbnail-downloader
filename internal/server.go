package internal

import (
	"context"
	"errors"
	"fmt"
	pg "getThumb/database"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"io"
	"net/http"
	"net/url"
	"strings"

	"getThumb/api/proto"
)

type Server struct {
	proto.UnimplementedThumbnailServer
	db *gorm.DB
}

func NewServer(db *gorm.DB) *Server {
	return &Server{db: db}
}

func (s *Server) GetThumbnail(_ context.Context, req *proto.Request) (*proto.Response, error) {
	rawUrl := req.GetUrl()
	videoId, err := ParseURL(rawUrl)
	if err != nil {
		return nil, err
	}

	thumbnail := pg.Thumbnail{}

	s.db.Table("thumbnails").Where("url = ?", videoId).Scan(&thumbnail)

	if thumbnail.Image != nil {
		logrus.Printf("URL %s found in cache", rawUrl)
		return &proto.Response{Image: thumbnail.Image}, nil
	}

	resp, err := http.Get(fmt.Sprintf("https://img.youtube.com/vi/%s/maxresdefault.jpg", videoId))
	if err != nil {
		return nil, err
	}

	if resp.Status == "404 Not Found" {
		return nil, errors.New(fmt.Sprintf("video %s not found", rawUrl))
	}

	defer resp.Body.Close()

	image, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	thumbnail = pg.Thumbnail{
		Url:   videoId,
		Image: image,
	}

	s.db.Create(&thumbnail)
	logrus.Printf("new url %s added to database", rawUrl)
	return &proto.Response{Image: thumbnail.Image}, nil
}

func ParseURL(rawUrl string) (string, error) {
	parsedURL, err := url.Parse(rawUrl)
	if err != nil {
		return "", err
	}

	if parsedURL.Host == "www.youtube.com" || parsedURL.Host == "youtube.com" {
		queryParams := parsedURL.Query()
		return queryParams.Get("v"), nil
	}

	if parsedURL.Host == "youtu.be" {
		return strings.TrimPrefix(parsedURL.Path, "/"), nil
	}

	return "", fmt.Errorf("invalid youtube URL")
}
