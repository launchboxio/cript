package storage

import (
	"errors"
	"fmt"
	"github.com/anchore/grype/grype/presenter/models"
	"github.com/docker/docker/api/types"
	"github.com/launchboxio/cript/internal/config"
)

type Storage interface {
	StoreInspection(image string, inspection types.ImageInspect) error
	GetInspection(image string) (types.ImageInspect, error)
	StoreVulnerabilityReport(image string, report models.Document) error
	GetVulnerabilityReport(image string) (models.Document, error)
}

func NewForConfig(config *config.Config) (Storage, error) {
	provider := config.Store
	switch provider {
	case "redis":
		return NewRedisAdapter(config.Redis)
	default:
		return nil, errors.New(fmt.Sprintf("Invalid storage provider: %s", provider))
	}
}
