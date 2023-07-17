package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/anchore/grype/grype/presenter/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/docker/docker/api/types"
	"github.com/launchboxio/cript/internal/config"
	"io"
	"strings"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
)

type S3Adapter struct {
	Svc    *s3.Client
	Bucket string
}

func NewS3Adapter(conf config.S3) (Storage, error) {
	sdkConfig, err := awsconfig.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}
	s3Client := s3.NewFromConfig(sdkConfig)
	return S3Adapter{
		Svc:    s3Client,
		Bucket: conf.Bucket,
	}, nil
}

func (adapter S3Adapter) StoreInspection(image string, inspection types.ImageInspect) error {
	out, err := json.Marshal(inspection)
	if err != nil {
		return err
	}
	reader := strings.NewReader(string(out))
	_, err = adapter.Svc.PutObject(context.TODO(), &s3.PutObjectInput{
		Key:    aws.String(fmt.Sprintf("inspections/%s", image)),
		Bucket: aws.String(adapter.Bucket),
		Body:   reader,
	})
	return err
}

func (adapter S3Adapter) GetInspection(image string) (types.ImageInspect, error) {
	var res types.ImageInspect
	result, err := adapter.Svc.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(adapter.Bucket),
		Key:    aws.String(fmt.Sprintf("inspections/%s", image)),
	})
	body, err := io.ReadAll(result.Body)
	if err != nil {
		return res, err
	}
	err = json.Unmarshal(body, &res)
	return res, err
}

func (adapter S3Adapter) StoreVulnerabilityReport(image string, report models.Document) error {
	out, err := json.Marshal(report)
	if err != nil {
		return err
	}
	reader := strings.NewReader(string(out))
	_, err = adapter.Svc.PutObject(context.TODO(), &s3.PutObjectInput{
		Key:    aws.String(fmt.Sprintf("reports/%s", image)),
		Bucket: aws.String(adapter.Bucket),
		Body:   reader,
	})
	return err
}

func (adapter S3Adapter) GetVulnerabilityReport(image string) (models.Document, error) {
	var res models.Document
	result, err := adapter.Svc.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(adapter.Bucket),
		Key:    aws.String(fmt.Sprintf("reports/%s", image)),
	})
	body, err := io.ReadAll(result.Body)
	if err != nil {
		return res, err
	}
	err = json.Unmarshal(body, &res)
	return res, err
}
