package amqp

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/support/log"
	"github.com/project-flogo/core/trigger"
	"pack.ag/amqp"
)

var triggerMd = trigger.NewMetadata(&Settings{}, &HandlerSettings{}, &Output{})

func init() {
	_ = trigger.Register(&AMQPTrigger{}, &Factory{})
}

type AMQPTrigger struct {
	settings        *Settings
	logger          log.Logger
	client          *amqp.Client
	receiverHandler []RecieverHandler
}
type RecieverHandler struct {
	receiver *amqp.Receiver
	handler  trigger.Handler
}

type Factory struct {
}

func (f *Factory) Metadata() *trigger.Metadata {
	return triggerMd
}

// New implements trigger.Factory.New
func (f *Factory) New(config *trigger.Config) (trigger.Trigger, error) {
	s := &Settings{}

	err := metadata.MapToStruct(config.Settings, s, true)
	if err != nil {
		return nil, err
	}
	var client *amqp.Client

	cfg := new(tls.Config)

	if s.CaPem != "" || s.CertPem != "" {
		cfg.RootCAs = x509.NewCertPool()

		if ca, err := ioutil.ReadFile(s.CaPem); err == nil {
			cfg.RootCAs.AppendCertsFromPEM(ca)
		}

		if s.CertPem != "" && s.KeyPem != "" {

			if cert, err := tls.LoadX509KeyPair(s.CertPem, s.KeyPem); err == nil {
				cfg.Certificates = append(cfg.Certificates, cert)
			}

		}
		client, err = amqp.Dial(s.AmqpURI,
			amqp.ConnTLS(true),
			amqp.ConnTLSConfig(cfg))

		if err != nil {
			return nil, fmt.Errorf("Dialing AMQP server: %v", err)
		}
	} else if s.Password != "" || s.Username != "" {

		client, err = amqp.Dial(s.AmqpURI, amqp.ConnSASLPlain(s.Username, s.Password))

		if err != nil {
			return nil, fmt.Errorf("Dialing AMQP server: %v", err)
		}

	} else {

		client, err = amqp.Dial(s.AmqpURI)
		if err != nil {
			return nil, fmt.Errorf("Dialing AMQP server: %v", err)
		}
	}

	return &AMQPTrigger{settings: s, client: client}, nil
}

func (t *AMQPTrigger) Initialize(ctx trigger.InitContext) error {

	t.logger = ctx.Logger()

	for _, handler := range ctx.GetHandlers() {

		s := &HandlerSettings{}
		err := metadata.MapToStruct(handler.Settings(), s, true)
		if err != nil {
			return err
		}

		session, err := t.client.NewSession()
		if err != nil {
			t.logger.Infof("Creating AMQP session:", err)
		}
		receiver, err := session.NewReceiver(
			amqp.LinkSourceAddress(s.SourceAddress),
			amqp.LinkCredit(10),
		)
		if err != nil {
			return err
		}
		t.logger.Info("Handler ..",handler)
		t.receiverHandler = append(t.receiverHandler, RecieverHandler{receiver: receiver, handler: handler})
	}

	return nil
}

func (t *AMQPTrigger) Start() error {
	
	for _, handler := range t.receiverHandler {
		t.logger.Infof("Starting Handler")
		handler.Start()
	}

	return nil
}

func (t *AMQPTrigger) Stop() error {
	for _, handler := range t.receiverHandler {
		handler.Stop()
	}
	t.client.Close()
	return nil
}

func (recieverHandler *RecieverHandler) Start() error {
	ctx := context.Background()

	for {
		
		msg, err := recieverHandler.receiver.Receive(ctx)
		if err != nil {
			return fmt.Errorf("Reading message from AMQP: %v", err)
		}
		
		// Accept message
		msg.Accept()
		out := &Output{}
		out.Data = msg.GetData()

		_, err = recieverHandler.handler.Handle(context.Background(), out)
		if err != nil {
			fmt.Println("Error in Handler..", err)
		}

	}
}

func (recieverHandler *RecieverHandler) Stop() error {

	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)

	recieverHandler.receiver.Close(ctx)

	cancel()

	return nil

}
