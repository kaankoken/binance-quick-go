package gcloudbucket

import (
	"context"
	"io"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"github.com/kaankoken/binance-quick-go/helper"
	"google.golang.org/api/option"
)

type ClientUploader struct {
	cl         *storage.Client
	projectId  string
	bucketName string
	uploadPath string
}

var Uploader *ClientUploader

const (
	fileName      string = "gcloud_storage"
	fileNameJSON  string = "credentials"
	extension     string = "yaml"
	extensionJSON string = "json"
)

func init() {
	config := helper.ReadGCloudConfig(fileName, extension)
	configJSON := helper.ReadGCloudConfigJson(fileNameJSON, extensionJSON)

	client, err := storage.NewClient(context.Background(), option.WithCredentialsFile(fileNameJSON+"."+extensionJSON))
	helper.CheckError(err)

	Uploader = &ClientUploader{
		cl:         client,
		bucketName: config["name"],
		projectId:  configJSON["installe"],
		uploadPath: config["path"],
	}
}

func (c *ClientUploader) UploadFile(file os.File, object string) error {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	wc := c.cl.Bucket(c.bucketName).Object(c.uploadPath + object).NewWriter(ctx)

	_, err := io.Copy(wc, &file)
	helper.CheckError(err)

	err = wc.Close()
	helper.CheckError(err)

	return nil
}
