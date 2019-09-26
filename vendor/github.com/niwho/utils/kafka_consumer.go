package utils

import (
	"fmt"
	//	"os"
	"runtime"
	"sync/atomic"

	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"github.com/niwho/logs"
)

type ConsumerCallback interface {
	Consume(message *sarama.ConsumerMessage) error
}

type KafkaConsumer struct {
	stopCh    chan struct{}
	isRunning int32
	workerNum int
	brokers   []string
	topics    []string
	group     string
	config    *cluster.Config

	consumecb ConsumerCallback
}

func NewKafkaConsumer(brokers []string, topic, group string, workerNum int, cb ConsumerCallback) *KafkaConsumer {
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Offsets.Retention = 0

	return &KafkaConsumer{
		stopCh:    make(chan struct{}),
		workerNum: workerNum,
		brokers:   brokers,
		topics:    []string{topic},
		config:    config,
		group:     group,
		consumecb: cb,
	}
}

func NewKafkaConsumerV2(brokers []string, topics []string, group string, workerNum int, cb ConsumerCallback) *KafkaConsumer {
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Offsets.Retention = 0

	return &KafkaConsumer{
		stopCh:    make(chan struct{}),
		workerNum: workerNum,
		brokers:   brokers,
		topics:    topics,
		config:    config,
		group:     group,
		consumecb: cb,
	}
}

func (kc *KafkaConsumer) Run() {

	if !atomic.CompareAndSwapInt32(&kc.isRunning, 0, 1) &&
		kc.workerNum <= 0 {
		return
	}
	fmt.Println("runrun")
	for i := 0; i < kc.workerNum; i++ {
		go func() {
			for {
				select {
				case <-kc.stopCh:
					kc.stopCh <- struct{}{}
					return
				default:
					kc.runSandbox()
				}
			}
		}()
	}
}

func (kc *KafkaConsumer) runSandbox() {

	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 20
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			fmt.Printf("AsyncFrame panic=%v\n%s\n", err, buf)
		}
	}()
	logs.Log(logs.F{"brokes": kc.brokers, "group": kc.group, "topic": kc.topics}).Info("KafkaConsumer")
	consumer, err := cluster.NewConsumer(kc.brokers, kc.group, kc.topics, kc.config)
	if err != nil {
		panic(err)
	}
	defer consumer.Close()
	for {
		select {
		case msg, more := <-consumer.Messages():
			if more {
				//fmt.Fprintf(os.Stdout, "%s/%d/%d\t%s\t%s\n", msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)
				// 不管处理状态?
				if kc.consumecb != nil {
					kc.consumecb.Consume(msg)
				}
				consumer.MarkOffset(msg, "") // mark message as processed
			}
		case err, more := <-consumer.Errors():
			if more {
				logs.Log(nil).Printf("Error: %s\n", err.Error())
			}
		case ntf, more := <-consumer.Notifications():
			if more {
				_ = ntf
				//logs.Log(nil).Printf("Rebalanced: %+v\n", ntf)
			}
		case <-kc.stopCh:
			// doing some worker at here
			kc.stopCh <- struct{}{}
			return
		}
	}

}

func (kc *KafkaConsumer) Close() {
	if !atomic.CompareAndSwapInt32(&kc.isRunning, 1, 0) {
		return
	}
	kc.stopCh <- struct{}{}
	<-kc.stopCh
}

func test() {
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true

}
