package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/wong-winnie/library/dao/config"
	"log"
)

type KafkaProducer struct {
	ProductConn sarama.SyncProducer
}

func InitProducer(cfg *config.KafkaCfg) *KafkaProducer {
	producerCfg := sarama.NewConfig()
	producerCfg.Producer.RequiredAcks = sarama.WaitForAll
	producerCfg.Producer.Partitioner = sarama.NewRandomPartitioner
	producerCfg.Producer.Return.Successes = true

	conn, err := sarama.NewSyncProducer(cfg.Address, producerCfg)
	if err != nil {
		log.Fatal("InitProducer Failed", err.Error())
	}

	return &KafkaProducer{ProductConn: conn}
}

func (product *KafkaProducer) SendMessage(topic string, data []byte) (partition int32, offset int64, err error) {
	var msg sarama.ProducerMessage
	msg.Topic = topic
	msg.Value = sarama.StringEncoder(data)

	return product.ProductConn.SendMessage(&msg)
}
