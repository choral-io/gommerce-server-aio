package v1beta

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/choral-io/gommerce-server-core/config"
	"github.com/choral-io/gommerce-server-core/logging"
	"github.com/choral-io/gommerce-server-core/server"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type ObjectStoreService struct {
	client *minio.Client
	logger logging.Logger
}

func NewObjectStoreService(config config.ServerMinIOConfig, logger logging.Logger) (*ObjectStoreService, error) {
	client, err := minio.New(config.GetEndpoint(), &minio.Options{
		Creds:  credentials.NewStaticV4(config.GetAccessKey(), config.GetSecretKey(), ""),
		Secure: config.GetUseSSL(),
	})
	if err != nil {
		return nil, err
	}
	return &ObjectStoreService{
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
	file, header, _ := r.FormFile("attachment")
	data := make([]byte, 512)
	file.Read(data)
	file.Seek(0, 0)
	name := queryEscapePath(header.Filename)
	path := queryEscapePath(pathParams["object"])
	size := header.Size
	if info, err := s.client.PutObject(context.Background(), pathParams["bucket"], path, file, size, minio.PutObjectOptions{
		ContentDisposition: "attachment; filename=" + name,
		ContentType:        http.DetectContentType(data),
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		s.logger.Info(context.Background(), "PutObject", "info", info, "error", err)
		w.Header().Add("Location", fmt.Sprintf("%s/%s/%s", s.client.EndpointURL(), info.Bucket, queryEscapePath(info.Key)))
		w.WriteHeader(http.StatusCreated)
	}
}

func queryEscapePath(path string) string {
	segs := strings.Split(path, "/")
	var j int
	for _, s := range segs {
		u := strings.TrimSpace(s)
		if len(u) > 0 {
			segs[j] = url.QueryEscape(u)
			j++
		}
	}
	return strings.Join(segs[:j], "/")
}
