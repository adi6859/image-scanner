package client

import (
	"time"

	"github.com/caarlos0/env"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type PubSubClient struct {
	logger        *zap.SugaredLogger
	JetStrContext nats.JetStreamContext
	Conn          *nats.Conn
}

type PubSubConfig struct {
	NatsServerHost string `env:"NATS_SERVER_HOST" envDefault:"nats://localhost:4222"`
}

func NewPubSubClient(logger *zap.SugaredLogger) (*PubSubClient, error) {

	cfg := &PubSubConfig{}
	err := env.Parse(cfg)
	if err != nil {
		logger.Error("err", err)
		return &PubSubClient{}, err
	}

	nc, err := nats.Connect(cfg.NatsServerHost,
		nats.ReconnectWait(10*time.Second), nats.MaxReconnects(100),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			logger.Errorw("Nats Connection got disconnected!", "Reason", err)
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			logger.Infow("Nats Connection got reconnected", "url", nc.ConnectedUrl())
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			logger.Errorw("Nats Client Connection closed!", "Reason", nc.LastError())
		}))
	if err != nil {
		logger.Error("err", err)
		return &PubSubClient{}, err
	}
	//Create a jetstream context
	js, _ := nc.JetStream()

	natsClient := &PubSubClient{
		logger:        logger,
		JetStrContext: js,
	}
	return natsClient, nil
}
