package storage

import (
	"context"
	"fmt"
	"github.com/anchore/grype/grype/presenter/models"
	"github.com/docker/docker/api/types"
	"github.com/launchboxio/cript/internal/config"
	"github.com/redis/go-redis/v9"
	"k8s.io/apimachinery/pkg/util/json"
)

type RedisAdapter struct {
	Client *redis.Client
}

func NewRedisAdapter(conf config.Redis) (Storage, error) {
	opts := &redis.Options{}
	if conf.Address != "" {
		opts.Addr = conf.Address
	}
	if conf.Password != "" {
		opts.Password = conf.Password
	}
	opts.DB = conf.Database
	rdb := redis.NewClient(opts)
	adapter := RedisAdapter{Client: rdb}
	return adapter, nil
}

func (r RedisAdapter) StoreInspection(image string, inspection types.ImageInspect) error {
	out, err := json.Marshal(inspection)
	if err != nil {
		return err
	}
	return r.Client.Set(context.TODO(), fmt.Sprintf("inspection:%s", image), out, 0).Err()
}

func (r RedisAdapter) GetInspection(image string) (types.ImageInspect, error) {
	var res types.ImageInspect
	val, err := r.Client.Get(context.TODO(), fmt.Sprintf("inspection:%s", image)).Result()
	if err != nil {
		return res, err
	}
	err = json.Unmarshal([]byte(val), &res)
	return res, err
}

func (r RedisAdapter) StoreVulnerabilityReport(image string, report models.Document) error {
	out, err := json.Marshal(report)
	if err != nil {
		return err
	}
	return r.Client.Set(context.TODO(), fmt.Sprintf("vuln:%s", image), out, 0).Err()
}

func (r RedisAdapter) GetVulnerabilityReport(image string) (models.Document, error) {
	var res models.Document
	val, err := r.Client.Get(context.TODO(), fmt.Sprintf("vuln:%s", image)).Result()
	if err != nil {
		return res, err
	}
	err = json.Unmarshal([]byte(val), &res)
	return res, err
}
