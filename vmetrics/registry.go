package vmetrics

import (
	"fmt"
	"github.com/Shopify/sarama"
	"log"
	"os"
	"strconv"
	"time"
)

var (
	DefaultRegistry = NewRegistry(&RegistryConfig{}, log.New(os.Stdout, "v-metric ", log.LstdFlags))
)

func init() {
	DefaultRegistry.Config.KafkaConfig.BrokerList = []string{"localhost:9092"}

	config := sarama.NewConfig()
	config.Net.DialTimeout = 10 * time.Second
	config.Version = sarama.V1_0_0_0
	config.Producer.Return.Successes = true

	DefaultRegistry.Config.KafkaConfig.Topic = "v-metrics"
	DefaultRegistry.Config.KafkaConfig.Config = config

	DefaultRegistry.Config.Cycle = time.Second

	producer, err := sarama.NewSyncProducer(
		DefaultRegistry.Config.KafkaConfig.BrokerList,
		config)
	if err != nil {
		fmt.Println("Failed to start Sarama producer:", err)
		os.Exit(1)
	}
	DefaultRegistry.Producer = producer

	DefaultRegistry.Start()
}

type KafkaConfig struct {
	BrokerList []string
	Config     *sarama.Config
	Topic      string
}

type RegistryConfig struct {
	KafkaConfig
	Cycle    time.Duration
}

type Registry struct {
	Metrics        []Metric
	Config *RegistryConfig
	Producer sarama.SyncProducer
	Logger *log.Logger
}

func NewRegistry(config *RegistryConfig, logger *log.Logger) *Registry {
	return &Registry{
		Metrics:        make([]Metric, 0),
		Config: config,
		Logger: logger,
	}
}

func Register(metric Metric) {
	DefaultRegistry.Register(metric)
}

func (r *Registry) Register(metric Metric) {
	r.Metrics = append(r.Metrics, metric)
}

func (r *Registry) Start() {
	go func() {
		for {
			for _, metric := range r.Metrics {
				messages := metric.Consume()
				go r.writeToKafka(messages)
				metric.Clear()
			}
			time.Sleep(r.Config.Cycle)
		}
	}()
}

func (r *Registry) writeToKafka(messages []string) {
	var pmsg []*sarama.ProducerMessage
	for _, message := range messages {
		msg := &sarama.ProducerMessage{
			Topic: r.Config.Topic,
			Key: sarama.StringEncoder(strconv.FormatInt(time.Now().UnixNano(), 10)),
			Value: sarama.StringEncoder(message),
		}
		pmsg = append(pmsg, msg)
	}
	if len(pmsg) > 0 {
		err := r.Producer.SendMessages(pmsg)
		if err != nil {
			r.Logger.Println(err)
		} else {
			r.Logger.Printf("(%d) messages processed\n", len(pmsg))
		}
	}
}