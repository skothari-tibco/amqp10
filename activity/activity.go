package amqp

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"encoding/gob"
    "bytes"
	"fmt"
	"pack.ag/amqp"
	"github.com/nu7hatch/gouuid"
	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/metadata"
)
func init() {
	_ = activity.Register(&Activity{}, New)
}

type Activity struct {
	settings *Settings
	client  *amqp.Client

}
var activityMd = activity.ToMetadata(&Settings{}, &Input{})

func New(ctx activity.InitContext) (activity.Activity, error) {

	s := &Settings{}
	
	err := metadata.MapToStruct(ctx.Settings(), s, true)

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
	} else {

		client, err = amqp.Dial(s.AmqpURI)
		if err != nil {
			return nil, fmt.Errorf("Dialing AMQP server: %v", err)
		}
	}

	return &Activity{settings: s, client: client}, nil
}
func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}

func (a *Activity) Eval(ctx activity.Context) (done bool, err error) {

	logger := ctx.Logger()

	input := &Input{}

	err = ctx.GetInputObject(input)
	
	if err != nil {
		return true, err
	}
	session, err := a.client.NewSession()
	if err != nil {
		logger.Infof("Creating AMQP session:", err)
	}
	
	data , err := GetBytes(input.Payload)
	if err != nil {
		return true, err
	}
	
	amqpMessage := amqp.NewMessage(data)

	u, err := uuid.NewV4()

	if err != nil {
		return true, err
	}
	
	properties := &amqp.MessageProperties{
		ContentType : input.MessageContentType,
		Subject : input.MessageSubject,
		MessageID : u,
	}

	amqpMessage.Properties = properties

	bCtx := context.Background()

    sender, err := session.NewSender(
        amqp.LinkTargetAddress(input.LinkTargetAddress),
    )
    if err != nil {
		return true, err
       
	}
	logger.Debug("Creating sender link:", sender)

	// Send message
	logger.Debug("Sending Message ...", amqpMessage)
    err = sender.Send(bCtx, amqpMessage)
    if err != nil {
		logger.Debug("Error in sending message:", err)
		return true, nil
    }

	sender.Close(bCtx)
	
	a.client.Close()

	return true, nil
}
func GetBytes(key interface{}) ([]byte, error) {
    var buf bytes.Buffer
    enc := gob.NewEncoder(&buf)
    err := enc.Encode(key)
    if err != nil {
        return nil, err
    }
    return buf.Bytes(), nil
}
