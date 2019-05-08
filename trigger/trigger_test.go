package amqp

import (
	"encoding/json"
	"testing"

	"github.com/project-flogo/core/action"
	"github.com/project-flogo/core/support"
	"github.com/project-flogo/core/support/test"
	"github.com/project-flogo/core/trigger"
	"github.com/stretchr/testify/assert"
)

const testConfig string = `{
	"id": "flogo-amqp",
	"ref": "github.com/skothari-tibco/amqp10trigger",
	"settings": {
	  "amqpURI": "amqps://messaging@t90ab69dfe0e54fb9bb38c9083ebb2936_hub:FIX51VQah7NhJZYlDXEx@messaging.bosch-iot-hub.com"
	},
	"handlers": [
	  {
			"action":{
				"id":"dummy"
			},
			"settings": {
				"sourceAddress":"telemetry/t90ab69dfe0e54fb9bb38c9083ebb2936_hub"
			}
	  }
	]
	
  }`

func TestTrigger_Register(t *testing.T) {

	ref := support.GetRef(&AMQPTrigger{})
	f := trigger.GetFactory(ref)
	assert.NotNil(t, f)
}

func TestAMQPTrigger_Initialize(t *testing.T) {
	f := &Factory{}

	config := &trigger.Config{}
	err := json.Unmarshal([]byte(testConfig), config)
	assert.Nil(t, err)

	actions := map[string]action.Action{"dummy": test.NewDummyAction(func() {
		//do nothing
	})}

	trg, err := test.InitTrigger(f, config, actions)
	assert.Nil(t, err)
	assert.NotNil(t, trg)

	err = trg.Start()
	for {

	}

}
