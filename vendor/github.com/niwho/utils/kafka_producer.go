package utils

import (
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/Shopify/sarama"
	"github.com/niwho/logs"
)

const (
	DefaultWorkerNum = 2
)

var (
	kfkProducer    *KafkaProducer
	kfkLogProducer *KafkaProducer
	clients        map[string]*KafkaProducer
)

func init() {
	clients = map[string]*KafkaProducer{}
}

func SendTask(key, data []byte) error {
	return kfkProducer.AddMsg(key, data)
}

func SendMsgA(clientName string, key, data []byte) error {
	if cli, found := clients[clientName]; found {
		return cli.AddMsg(key, data)
	}
	return errors.New("not found")
}

func SendMsgAWithTopic(topic, clientName string, key, data []byte) error {
	if cli, found := clients[clientName]; found {
		return cli.AddMsgWithTopic(topic, key, data)
	}
	return errors.New("not found")
}

func SendLog(key, data []byte) error {
	return kfkLogProducer.AddMsg(key, data)
}

func KafkaClose() {
	kfkProducer.Close()
	kfkLogProducer.Close()
	for _, cli := range clients {
		cli.Close()
	}
}

func InitKafkaProducer(addrs []string, topic string, bufSize int, timeout time.Duration, workerNum int) {
	logs.Log(logs.F{"addrs": addrs, "topic": topic}).Info("InitKafkaProducer")
	kfkProducer = NewKafkaProducer(
		addrs,
		topic,
		bufSize,
		timeout,
		workerNum,
	)
	kfkProducer.Start()
}

func InitKafkaProducerA(clientName string, addrs []string, topic string, bufSize int, timeout time.Duration, workerNum int) {
	logs.Log(logs.F{"addrs": addrs, "topic": topic}).Info("InitKafkaProducer")
	client := NewKafkaProducer(
		addrs,
		topic,
		bufSize,
		timeout,
		workerNum,
	)
	client.Start()
	clients[clientName] = client
}

func InitKafkaLogProducer(addrs []string, topic string, bufSize int, timeout time.Duration, workerNum int) {
	kfkLogProducer = NewKafkaProducer(
		addrs,
		topic,
		bufSize,
		timeout,
		workerNum,
	)
	kfkLogProducer.Start()
}

type KafkaProducer struct {
	addrs     []string
	topic     string
	config    *sarama.Config
	msgCh     chan *sarama.ProducerMessage
	timeout   time.Duration
	workerNum int
	isRunning int32
	stop      chan struct{}
	stopFlag  bool
}

func NewKafkaProducer(addrs []string, topic string, bufSize int, timeout time.Duration, workerNum int) *KafkaProducer {
	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewRoundRobinPartitioner
	config.Producer.MaxMessageBytes = 10000000 //10m
	kw := &KafkaProducer{
		addrs:     addrs,
		topic:     topic,
		msgCh:     make(chan *sarama.ProducerMessage, bufSize),
		stop:      make(chan struct{}),
		config:    config,
		timeout:   timeout,
		workerNum: workerNum,
	}
	if workerNum <= 0 {
		kw.workerNum = DefaultWorkerNum
	}
	return kw
}

func (kw *KafkaProducer) AddMsg(key, data []byte) error {
	return kw.AddMsgWithTopic(kw.topic, key, data)
}
func (kw *KafkaProducer) AddMsgWithTopic(topic string, key, data []byte) error {
	pm := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(data),
		Key:   sarama.ByteEncoder(key),
	}
	select {
	case kw.msgCh <- pm:
	case <-time.After(kw.timeout):
		// config.EmitCounter("send_kafka_timeout", 1, map[string]string{"topic": topic})
		logs.Log(nil).Errorf("msg='KafkaProducer add msg timeout'")
		return fmt.Errorf("KafkaProducer add msg timeout")
	}
	return nil
}

func (kw *KafkaProducer) Start() {
	if !atomic.CompareAndSwapInt32(&kw.isRunning, 0, 1) {
		return
	}
	for i := 0; i < kw.workerNum; i += 1 {
		if kw.stopFlag {
			break
		}
		go kw.startWorker()
	}
}

func (kw *KafkaProducer) startWorker() {
	for {
		select {

		case <-kw.stop:
			kw.stop <- struct{}{}
			return
		default:
			worker, err := sarama.NewAsyncProducer(kw.addrs, kw.config)
			if err != nil {
				logs.Log(nil).Errorf("error=%v, msg='KafkaWorker new producer failed'", err)
				time.Sleep(500 * time.Millisecond)
				continue
			}
			kw.loop(worker)
		}
	}
}

func (kw *KafkaProducer) loop(worker sarama.AsyncProducer) {
	defer func() {
		if err := worker.Close(); err != nil {
			logs.Log(nil).Errorf("error=%v msg='KafkaWorker close producer failed'", err)
		}
	}()
	for {
		select {
		case msg := <-kw.msgCh:
			select {
			case worker.Input() <- msg:
				logs.Log(nil).Info("[KafkaWorker]send msg to chain")
			case err := <-worker.Errors():
				logs.Log(nil).Errorf("error=%v msg='KafkaWorker send msg failed'", err)
				return

			}

		case <-kw.stop:
			//clean
			logs.Log(nil).Debug("loop stop")
			kw.clean(worker)
			logs.Log(nil).Debug("loop stop after clean")
			kw.stop <- struct{}{}
			return
		}
	}

}
func (kw *KafkaProducer) clean(worker sarama.AsyncProducer) {
	// 等待消息全部发送到kafka
	for {
		select {
		case msg := <-kw.msgCh:
			worker.Input() <- msg
		default:
			logs.Log(nil).Debug("clean")
			// worker.Close()
			return

		}

	}
}

func (kw *KafkaProducer) Close() {
	if !atomic.CompareAndSwapInt32(&kw.isRunning, 1, 0) {
		return
	}
	kw.stopFlag = true

	logs.Log(nil).Debug("common close000\n")
	// 触发
	kw.stop <- struct{}{}
	logs.Log(nil).Debug("common close111\n")
	<-kw.stop
	logs.Log(nil).Debug("common close222\n")
}
