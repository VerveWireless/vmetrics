package main


import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/pubnative/vmetrics/vmetrics"
	"log"
	"os"
	"time"
)

var (
	personMetric = vmetrics.NewMessage("yield", "bo", "persons")
)

func init() {
	vmetrics.DefaultRegistry.Config.KafkaConfig.BrokerList = []string{"localhost:9092"}

	config := sarama.NewConfig()
	config.Net.DialTimeout = 10 * time.Second
	config.Version = sarama.V1_0_0_0
	config.Producer.Return.Successes = true

	vmetrics.DefaultRegistry.Config.KafkaConfig.Topic = "v-metrics2"
	vmetrics.DefaultRegistry.Config.KafkaConfig.Config = config

	vmetrics.DefaultRegistry.Config.Cycle = time.Second

	producer, err := sarama.NewSyncProducer(
		vmetrics.DefaultRegistry.Config.KafkaConfig.BrokerList,
		config)
	if err != nil {
		fmt.Println("Failed to start Sarama producer:", err)
		os.Exit(1)
	}
	vmetrics.DefaultRegistry.Producer = producer



	vrc := vmetrics.RegistryConfig{
		KafkaConfig: vmetrics.KafkaConfig{},
		Cycle:       time.Minute,
	}
	vmetrics.DefaultRegistry = vmetrics.NewRegistry(&vrc, log.New(os.Stdout,"", log.LstdFlags))
	vmetrics.Register(personMetric)
}

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	SomeOperation()

	time.Sleep(time.Second * 5)
}

func SomeOperation() {
	fmt.Println("Doing the operation")
	for i := 0; i < 3000; i++ {
		p := Person{Name: "test", Age: i}
		personMetric.Record(p)
		time.Sleep(time.Nanosecond)
	}
}

