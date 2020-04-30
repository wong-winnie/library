package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"library/dao/config"
	"log"
)

type KafkaProducer struct {
	productConn sarama.SyncProducer
	writeChan   chan sarama.ProducerMessage
	isClose     bool
}

func InitProducer(cfg *config.KafkaCfg) *KafkaProducer {
	producerCfg := sarama.NewConfig()
	producerCfg.Producer.RequiredAcks = sarama.WaitForAll
	producerCfg.Producer.Partitioner = sarama.NewRandomPartitioner
	producerCfg.Producer.Return.Successes = true

	if conn, err := sarama.NewSyncProducer(cfg.Address, producerCfg); err != nil {
		log.Fatal(err.Error())
		return &KafkaProducer{}
	} else {
		pd := &KafkaProducer{
			productConn: conn,
			writeChan:   make(chan sarama.ProducerMessage, 1024),
			isClose:     false,
		}

		go pd.writeMessage()

		return pd
	}
}

func (product *KafkaProducer) SendMessage(topic string, data []byte) {
	if product.isClose {
		return
	}

	var msg sarama.ProducerMessage
	msg.Topic = topic
	msg.Value = sarama.StringEncoder(data)
	product.writeChan <- msg
}

func (product *KafkaProducer) writeChannel(data sarama.ProducerMessage) {
	if product.isClose {
		return
	}
	product.writeChan <- data
}

func (product *KafkaProducer) writeMessage() {
	for msg := range product.writeChan {
		_, _, err := product.productConn.SendMessage(&msg)
		if err != nil {
			fmt.Println("send message failed ", err)
			product.writeChannel(msg)
			continue
		}
	}
}

func TestProduct() {
	////初始化
	//producer, err := InitProducer([]string{"192.168.28.25:9092", "192.168.28.26:9092", "192.168.28.27:9092"})
	//if err != nil {
	//	fmt.Println(err.Error())
	//	return
	//}
	//
	////处理主题的消息
	//go producer.writeMessage()
	//
	//var msg sarama.ProducerMessage
	//msg.Topic = "demo_service"
	//
	////给主题发送消息
	//for i := 0; i < 100; i++ {
	//	msg.Value = sarama.StringEncoder(i)
	//	producer.WriteChan <- msg
	//}
	//
	//time.Sleep(time.Second * 60)
}
