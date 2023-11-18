package main

import (
	"fmt"
	"os"

	"github.com/kelseyhightower/envconfig"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/victor-felix/chat-bot/app/brokers"
	"github.com/victor-felix/chat-bot/app/clients"
	config "github.com/victor-felix/chat-bot/app/config"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	runLogFile, _ := os.OpenFile("broker.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
	multi := zerolog.MultiLevelWriter(runLogFile)
	log.Logger = zerolog.New(multi).With().Timestamp().Logger()

	var cfg config.Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Error().Msg(fmt.Sprintf("cannot load env variables %s", err))
	}

	
	conn, err := amqp.Dial(cfg.RabbitMQ.DSN)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	log.Info().Msg("Connected to RabbitMQ")
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	log.Info().Msg("Connected to RabbitMQ channel")
	defer channel.Close()

	stooqClient := clients.NewStooqClient(cfg.StooqBaseUrl, log.Logger)

	broker := brokers.NewStooqBroker(channel, stooqClient, cfg.RabbitMQ.ReceiverQueueName, cfg.RabbitMQ.PublisherQueueName, log.Logger)

	go broker.ReadMessages()
	select {}
}
