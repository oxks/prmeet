package natsio

import (
	"log"
	"prmeet/config"
	"prmeet/internal/er"
	"time"

	"github.com/nats-io/nats.go"
)

// connect encoded nats
func ConnectNatsEncoded() *nats.EncodedConn {

	nc := config.ConnectNats()
	enc, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		log.Fatal(err)
	}
	return enc
}

func AskPipe(subj string, request map[string]any) (response map[string]any) {

	// Connect to NATS encoded
	enc := ConnectNatsEncoded()
	defer enc.Close()

	err := enc.Request(subj, request, &response, 2*time.Second)
	er.ErrorPrint(err)

	return response
}
