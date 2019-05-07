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



func TestEval(t *testing.T) {

	settings := &Settings{AmqpURI: "amqps://messaging@tcef56e88b16548f9a4a49cd5b92150af:QIktjmIzfJHfsvOlT5OF@messaging.bosch-iot-hub.com"}

	iCtx := test.NewActivityInitContext(settings, nil)
	act, err := New(iCtx)
	assert.Nil(t, err)
	
	tc := test.NewActivityContext(act.Metadata())

	input := &Input{LinkTargetAddress:"control/tb02afe3819814694a9e1383ca19eeae6_hub/4711",Payload : "Hello", MessageContentType :"string", MessageSubject: "Hello"}
	tc.SetInputObject(input)

	act.Eval(tc)
}
