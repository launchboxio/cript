package storage

import (
	"github.com/anchore/grype/grype/presenter/models"
	"github.com/docker/docker/api/types"
	"gorm.io/gorm"
)

type DatabaseAdapter struct {
	DB *gorm.DB
}

func (db DatabaseAdapter) StoreInspection(image string, inspection types.ImageInspect) error {
	return nil
}

func (db DatabaseAdapter) GetInspection(image string) (types.ImageInspect, error) {
	return types.ImageInspect{}, nil
}

func (db DatabaseAdapter) StoreVulnerabilityReport(image string, report models.Document) error {
	return nil
}

func (db DatabaseAdapter) GetVulnerabilityReport(image string) (models.Document, error) {
	return models.Document{}, nil
}
