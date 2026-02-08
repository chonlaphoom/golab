package main

import (
	"context"
	"encoding/json"
	"log"
)

func readAndPushMsgs(ctx context.Context, decoder *json.Decoder, msgChan chan<- json.RawMessage, errChan chan<- error) {
	defer func() {
		close(msgChan)
		log.Println("Message channel closed.")
	}()
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		var msg json.RawMessage
		if err := decoder.Decode(&msg); err != nil {
			errChan <- err
			return
		}
		msgChan <- msg
	}
}
