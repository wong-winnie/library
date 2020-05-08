package kafka

import (
	"context"
	"errors"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/wong-winnie/library/dao/config"
	"log"
)

type KafkaConsumer struct {
	ConsumerConn    sarama.ConsumerGroup
	ConsumerHandler MainHandler
}

func InitConsumer(cfg *config.KafkaCfg) *KafkaConsumer {
	consumerCfg := sarama.NewConfig()
	consumerCfg.Version = sarama.V2_2_0_0
	consumerCfg.Producer.Return.Errors = true
	consumerCfg.Producer.Return.Successes = true
	consumerCfg.Producer.RequiredAcks = sarama.WaitForAll
	consumerCfg.Producer.Partitioner = sarama.NewManualPartitioner
	consumerCfg.Consumer.Return.Errors = true
	consumerCfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	consumerCfg.ClientID = cfg.ClientName

	conn, err := sarama.NewClient(cfg.Address, consumerCfg)
	if err != nil {
		log.Fatal("InitConsumer Failed", err.Error())
	}

	gConn, err := sarama.NewConsumerGroupFromClient(cfg.GroupName, conn)
	if err != nil {
		log.Fatal("InitConsumer Failed", err.Error())
	}

	return &KafkaConsumer{ConsumerConn: gConn}
}

func (consumer KafkaConsumer) TopicConsumer(topic string, hand sarama.ConsumerGroupHandler) error {
	ctx, _ := context.WithCancel(context.Background())

	for {
		if err := consumer.ConsumerConn.Consume(ctx, []string{topic}, hand); err != nil || ctx.Err() != nil {
			return err
		}
	}
}

type MainHandler struct {
	RecvFunc func(recv []byte) bool
}

func (m *MainHandler) Setup(sess sarama.ConsumerGroupSession) error {
	// 如果极端情况下markOffset失败，需要手动同步offset
	return nil
}

func (m *MainHandler) Cleanup(sess sarama.ConsumerGroupSession) error {
	// 如果极端情况下markOffset失败，需要手动同步offset
	return nil
}

//此方法会自动控制偏移值，当分组里的主题消息被接收到时，则偏移值会进行加1 他是跟着主题走的
func (m *MainHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	offsetList := make(map[int32]int64, 0)

	for message := range claim.Messages() {
		//claim.InitialOffset() 初始偏移量是从-2开始
		if claim.InitialOffset()+int64(offsetList[message.Partition]) != message.Offset && claim.InitialOffset()+int64(offsetList[message.Partition]) >= 0 {
			return errors.New(fmt.Sprintf("Topic:%v, Partition:%v, Offset:%v,  Value:%v", message.Topic, message.Partition, message.Offset, message.Value))
		} else {
			if ok := m.RecvFunc(message.Value); ok {
				sess.MarkMessage(message, "") //确认消费， 偏移量+1
			}
		}
		offsetList[message.Partition]++
	}
	return nil
}
