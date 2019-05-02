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
	  "amqpURI": "amqps://messaging@tcef56e88b16548f9a4a49cd5b92150af:QIktjmIzfJHfsvOlT5OF@messaging.bosch-iot-hub.com"
	},
	"handlers": [
	  {
			"action":{
				"id":"dummy"
			},
			"settings": {
				"sourceAddress":"event/tcef56e88b16548f9a4a49cd5b92150af"
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
	assert.Nil(t, err)

}
