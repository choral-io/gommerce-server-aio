package oss_v1beta

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awscfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/choral-io/gommerce-server-core/config"
	"github.com/choral-io/gommerce-server-core/logging"
	"github.com/choral-io/gommerce-server-core/server"
)

type ObjectStoreService struct {
	urltpl string
	client *s3.Client
	logger logging.Logger
}

func NewObjectStoreService(cfg config.ServerStorageConfig, logger logging.Logger) (*ObjectStoreService, error) {
	s3cfg, err := awscfg.LoadDefaultConfig(
		context.Background(),
		awscfg.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.GetAccessKey(), cfg.GetSecretKey(), "")),
		awscfg.WithRegion("auto"),
	)
	if err != nil {
		return nil, err
	}
	client := s3.NewFromConfig(s3cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(cfg.GetEndpoint())
	})
	urltpl := strings.TrimRight(cfg.GetEndpoint(), "/") + "/%s/%s"
	return &ObjectStoreService{
		urltpl: urltpl,
		client: client,
		logger: logger,
	}, nil
}

func (s *ObjectStoreService) ServerMuxRoutes() []server.ServerMuxRoute {
	return []server.ServerMuxRoute{
		{
			Methods: []string{"PUT"},
			Pattern: "/v1beta/objects/{bucket}/{object=**}",
			Handler: s.PutObject,
		},
	}
}

func (s *ObjectStoreService) PutObject(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	file, header, err := r.FormFile("attachment")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	data := make([]byte, 512)
	_, _ = file.Read(data)
	_, _ = file.Seek(0, 0)
	name := pathEscapePath(header.Filename)
	path := pathEscapePath(pathParams["object"])
	size := header.Size
	if _, err := s.client.PutObject(r.Context(), &s3.PutObjectInput{
		Bucket:             aws.String(pathParams["bucket"]),
		Key:                aws.String(path),
		Body:               file,
		ContentLength:      aws.Int64(size),
		ContentDisposition: aws.String("attachment; filename=" + name),
		ContentType:        aws.String(http.DetectContentType(data)),
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		s.logger.Error(r.Context(), "failed to put object", "error", err, "bucket", pathParams["bucket"], "object", path, "size", size)
	} else {
		w.Header().Add("Location", fmt.Sprintf(s.urltpl, pathParams["bucket"], path))
		w.WriteHeader(http.StatusCreated)
		s.logger.Info(r.Context(), "object uploaded successfully", "bucket", pathParams["bucket"], "object", path, "size", size)
	}
}

func pathEscapePath(path string) string {
	segs := strings.Split(path, "/")
	var j int
	for _, s := range segs {
		if trimmed := strings.TrimSpace(s); len(trimmed) > 0 {
			segs[j] = url.PathEscape(trimmed)
			j++
		}
	}
	return strings.Join(segs[:j], "/")
}
