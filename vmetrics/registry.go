package vmetrics

import (
	"github.com/Shopify/sarama"
	"log"
	"os"
	"strconv"
	"time"
)

var (
	DefaultRegistry *Registry
)

type BrokerList []string

func SetupDefaultRegistry(brokers []string, topic string, logger *log.Logger) {
	if logger == nil {
		logger = log.New(os.Stdout, "v-metrics: ", log.LstdFlags)
	}
	DefaultRegistry = NewRegistry(&RegistryConfig{}, logger)

	if len(brokers) <= 0 {
		logger.Fatal("Kafka brokers should be specified")
	}
	DefaultRegistry.Config.BrokerList = brokers

	saramaConfig := sarama.NewConfig()
	saramaConfig.Net.DialTimeout = 10 * time.Second
	saramaConfig.Version = sarama.V1_0_0_0
	saramaConfig.Producer.Return.Successes = true

	DefaultRegistry.Config.Topic = topic

	DefaultRegistry.Config.Cycle = time.Second

	producer, err := sarama.NewSyncProducer(
		DefaultRegistry.Config.BrokerList,
		saramaConfig)
	if err != nil {
		DefaultRegistry.Logger.Println(err)
	}
	DefaultRegistry.Producer = producer
	if err != nil {
		logger.Println(err)
	} else {
		DefaultRegistry.Start()
	}
}

type RegistryConfig struct {
	BrokerList []string
	Topic      string
	Cycle      time.Duration
}

type Registry struct {
	Metrics  []Metric
	Config   *RegistryConfig
	Producer sarama.SyncProducer
	Logger   *log.Logger
}

func NewRegistry(config *RegistryConfig, logger *log.Logger) *Registry {
	return &Registry{
		Metrics: make([]Metric, 0),
		Config:  config,
		Logger:  logger,
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
			Key:   sarama.StringEncoder(strconv.FormatInt(time.Now().UnixNano(), 10)),
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
