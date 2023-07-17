package storage

import (
	"github.com/anchore/grype/grype/presenter/models"
	"github.com/docker/docker/api/types"
)

type Storage interface {
	StoreInspection(image string, inspection types.ImageInspect) error
	GetInspection(image string) (types.ImageInspect, error)
	StoreVulnerabilityReport(image string, report models.Document) error
	GetVulnerabilityReport(image string) (models.Document, error)
}
