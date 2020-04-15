package amqp


type Settings struct {
	AmqpURI string `md:"amqpURI"`
	CaPem   string `md:"caPem'`
	CertPem string `md:"certPem"`
	KeyPem  string `md:"keyPem"`
	Username string `md:"username"`
	Password string `md:"password"`
}

type HandlerSettings struct {
	SourceAddress string `md:"sourceAddress"`
}

type Output struct {
	Data interface{} `md:"data"`
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"data": o.Data,
	}
}

func (o *Output) FromMap(values map[string]interface{}) error {
	
	o.Data = values["data"]
	
	return nil
}
