package kafka

import (
	"context"
	"errors"
	"fmt"
	"github.com/Shopify/sarama"
	"library/dao/config"
	"log"
)

type KafkaConsumer struct {
	consumerConn sarama.ConsumerGroup
	isClose      bool
}

func InitConsumer(cfg *config.KafkaCfg) *KafkaConsumer {
	consumerCfg := sarama.NewConfig()
	consumerCfg.Version = sarama.V2_2_0_0
	consumerCfg.Producer.Return.Errors = true
	consumerCfg.Net.SASL.Enable = false
	consumerCfg.Producer.Return.Successes = true
	consumerCfg.Producer.RequiredAcks = sarama.WaitForAll
	consumerCfg.Producer.Partitioner = sarama.NewManualPartitioner
	consumerCfg.Consumer.Return.Errors = true
	consumerCfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	consumerCfg.ClientID = cfg.ClientName

	conn, err := sarama.NewClient(cfg.Address, consumerCfg)
	if err != nil {
		log.Fatal(err.Error())
		return &KafkaConsumer{}
	}

	gConn, err := sarama.NewConsumerGroupFromClient(cfg.GroupName, conn)
	if err != nil {
		log.Fatal(err.Error())
		return &KafkaConsumer{}
	} else {
		return &KafkaConsumer{
			consumerConn: gConn,
			isClose:      false,
		}
	}
}

func (consumer KafkaConsumer) topicConsumer(topic string, hand sarama.ConsumerGroupHandler) {
	go func() {
		for err := range consumer.consumerConn.Errors() {
			fmt.Println(err)
		}
	}()

	ctx, _ := context.WithCancel(context.Background())

	for {
		err := consumer.consumerConn.Consume(ctx, []string{topic}, hand)
		if err != nil {
			fmt.Println(err)
			break
		}
		if ctx.Err() != nil {
			break
		}
	}
}

type MainHandler struct {
	recvFunc func(recv []byte) bool
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

	mo := make(map[int32]int64, 0)

	for message := range claim.Messages() {

		if claim.InitialOffset()+int64(mo[message.Partition]) != message.Offset && claim.InitialOffset()+int64(mo[message.Partition]) >= 0 {
			fmt.Println("-----", claim.InitialOffset()+mo[message.Partition], message.Partition)
			return errors.New(fmt.Sprintf("Topic:%v, Partition:%v, Offset:%v,  Value:%v", message.Topic, message.Partition, message.Offset, message.Value))
		} else {
			if ok := m.recvFunc(message.Value); ok {
				sess.MarkMessage(message, "") //确认消费， 偏移量+1
			}
		}
		mo[message.Partition]++
	}
	return nil
}

func TestConsumer() {

	//var (
	//	Topic      = "demo_service998"
	//	GroupName  = "chat_Group_LLH91"
	//	ClientName = "client3"
	//)
	//
	//consumer, err := InitConsumer([]string{"192.168.28.25:9092", "192.168.28.26:9092", "192.168.28.27:9092"}, GroupName, ClientName)
	//if err != nil {
	//	fmt.Println(err.Error())
	//	return
	//}
	//
	//hand := &MainHandler{}
	//consumer.topicConsumer(Topic, hand)
}
