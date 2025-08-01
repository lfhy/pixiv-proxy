package server

import (
	"bytes"
	"context"
	"fmt"
	"go-pixiv-proxy/conf"
	"net/url"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/lfhy/log"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var hasRemote bool

var minioClient *minio.Client

var uploadLock sync.Map

func InitRemote() {
	info, err := url.Parse(conf.RemoteEndpoint)
	if err != nil {
		return
	}
	minioClient, err = minio.New(info.Host, &minio.Options{
		Creds:  credentials.NewStaticV4(conf.RemoteAk, conf.RemoteSk, ""),
		Secure: info.Scheme == "https", // 如果使用 HTTPS 设置为 true
	})
	if err != nil {
		log.Infoln("存储初始化失败:", err)
		return
	}
	hasRemote = true
}

func HeadRemote(objectName string) (string, bool) {
	remoteFilePath := filepath.Join(conf.RemoteDir, objectName)
	_, err := minioClient.StatObject(context.Background(), conf.RemoteBucket, remoteFilePath, minio.GetObjectOptions{})
	if err != nil {
		return "", false
	}
	return GetObjecRemoteUrl(remoteFilePath), true
}

func PutDataToRemote(data []byte, objectName string) (string, error) {
	defer uploadLock.Delete(objectName)
	remoteFilePath := filepath.Join(conf.RemoteDir, objectName)
	_, err := minioClient.PutObject(context.Background(), conf.RemoteBucket, remoteFilePath, bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		log.Warnln("上传失败:", err)
		return "", fmt.Errorf("上传至S3失败: %v", err)
	}
	log.Infoln("上传完成:", objectName)
	return GetObjecRemoteUrl(remoteFilePath), nil
}

func GetObjecRemoteUrl(object string) string {
	u, err := minioClient.PresignedGetObject(context.Background(), conf.RemoteBucket, strings.TrimPrefix(object, "/"), 24*time.Hour, nil)
	if err != nil {
		return ""
	}
	return u.String()
}
