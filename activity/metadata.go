package amqp

import (
	"github.com/project-flogo/core/data/coerce"
)

type Settings struct {
	AmqpURI string `md:"amqpURI"`
	CaPem string `md:"caPem"`
	CertPem string `md:"certPem"`
	KeyPem string `md:"keyPem"`
}

type Input struct {
	LinkTargetAddress string `md:"linkTargetAddress"`
	MessageSubject string `md:"messageSubject"`
	MessageContentType string `md:"messageContentType"`
	Payload interface{} `md:"payload"`
}

func (i *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"linkTargetAddress":  i.LinkTargetAddress,
		"messageSubject": i.MessageSubject,
		"messageContentType":     i.MessageContentType,
		"payload":     i.Payload,
	}
}

func (i *Input) FromMap(values map[string]interface{}) error {

	var err error
	i.LinkTargetAddress, err = coerce.ToString(values["linkTargetAddress"])
	if err != nil {
		return err
	}
	i.MessageSubject, err = coerce.ToString(values["messageSubject"])
	if err != nil {
		return err
	}
	i.MessageContentType, err = coerce.ToString(values["messageContentType"])
	if err != nil {
		return err
	}
	i.Payload = values["payload"]

	return nil
}