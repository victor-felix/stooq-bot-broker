package brokers

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"github.com/victor-felix/chat-bot/app/models"
)

type stooqBroker struct {
	Receiver amqp.Queue
	Publisher amqp.Queue
	Channel *amqp.Channel
	stooqClient models.StooqClient
	log zerolog.Logger
}

func NewStooqBroker(channel *amqp.Channel, stooqClient models.StooqClient, receiverQueue string, publisherQueue string, log zerolog.Logger) *stooqBroker {
	queueReceiver, err := channel.QueueDeclare(
		receiverQueue,
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	queuePublisher, err := channel.QueueDeclare(
		publisherQueue,
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	return &stooqBroker{
		Receiver: queueReceiver,
		Publisher: queuePublisher,
		Channel: channel,
		stooqClient: stooqClient,
		log: log,
	}
}

func (sb *stooqBroker) Publish(response models.BotResponse) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := json.Marshal(response)
	if err != nil {
		sb.log.Error().Msg(err.Error())
		return
	}

	err = sb.Channel.PublishWithContext(ctx, "", sb.Publisher.Name, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body: body,
	})

	if err != nil {
		sb.log.Error().Msg(err.Error())
		return
	}
}

func (sb *stooqBroker) ReadMessages() {
	messages, err := sb.Channel.Consume(
		sb.Receiver.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		sb.log.Fatal().Msg(err.Error())
	}

	receivedMessages := make(chan models.BotRequest)
	go sb.messageTransformer(messages, receivedMessages)
	go sb.processResponse(receivedMessages)
}


func (sb *stooqBroker) messageTransformer(entries <-chan amqp.Delivery, receivedMessages chan models.BotRequest) {
	var botRequest models.BotRequest
	for message := range entries {
		sb.log.Info().Msg(fmt.Sprintf("Received a message: %s", message.Body))
		err := json.Unmarshal([]byte(message.Body), &botRequest)
		if err != nil {
			sb.log.Error().Msg(err.Error())
			continue
		}
		receivedMessages <- botRequest
	}
}

func (sb *stooqBroker) processResponse(requests <-chan models.BotRequest) {
	for request := range requests {
		sb.log.Info().Msg(fmt.Sprintf("Processing request for room: %s", request.RoomID))

		messageContent := request.Content
		messageContent = strings.Replace(messageContent, "/stock=", "", 1)

		response := models.BotResponse{
			RoomID:  request.RoomID,
			RoomName: request.RoomName,
			UserID:  request.UserID,
			CreatedAt: request.CreatedAt,
		}

		messageResponse, err := sb.stooqClient.GetStockPrice(messageContent)
		if err != nil {
			sb.log.Error().Msg(err.Error())
			continue
		}

		response.Content = messageResponse

		sb.log.Info().Msg(fmt.Sprintf("Sending response for room: %s, content: %s", request.RoomID, response.Content))
		go sb.Publish(response)
	}
}