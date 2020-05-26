package elastice

import (
	elastic "github.com/olivere/elastic/v7"
	"github.com/wong-winnie/library/dao/config"
	"log"
)

type ElasticMgr struct {
	Conn *elastic.Client
}

func InitElastic(cfg *config.ElasticCfg) *ElasticMgr {
	client, err := elastic.NewClient(elastic.SetURL(cfg.Url))
	if err != nil {
		log.Fatal("InitElastic Failed", err.Error())
	}
	return &ElasticMgr{Conn: client}
}
