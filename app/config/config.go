package models

type Config struct {
	DevMode bool `split_words:"true" required:"true" default:"true"`
	MongoURL string `split_words:"true" required:"true" default:"mongodb://admin:Qwe12345@localhost:27017"`
	DatabaseName string `split_words:"true" required:"true" default:"chat"`
	RabbitMQ struct {
		DSN string `split_words:"true" required:"true" default:"amqp://guest:guest@localhost:5672/"`
		ReceiverQueueName string `split_words:"true" required:"true" default:"stock-process-request"`
		PublisherQueueName string `split_words:"true" required:"true" default:"stock-result"`
	}
	StooqBaseUrl string `split_words:"true" required:"true" default:"https://stooq.com"`
}
