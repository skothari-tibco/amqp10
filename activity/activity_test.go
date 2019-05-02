package amqp

import (

	"github.com/project-flogo/core/activity"
	
	"github.com/stretchr/testify/assert"
	"github.com/project-flogo/core/support/test"
	
	"testing"
)


func TestRegister(t *testing.T) {

	ref := activity.GetRef(&Activity{})
	act := activity.Get(ref)

	assert.NotNil(t, act)
}

func TestSettings(t *testing.T) {
	settings := &Settings{AmqpURI: "amqps://messaging@tcef56e88b16548f9a4a49cd5b92150af:QIktjmIzfJHfsvOlT5OF@messaging.bosch-iot-hub.com"}

	iCtx := test.NewActivityInitContext(settings, nil)
	_, err := New(iCtx)
	assert.Nil(t, err)

}