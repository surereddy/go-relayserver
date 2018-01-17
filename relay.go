package main

import (
	"fmt"
	"log"
	"time"

	"github.com/grafov/bcast"
)

// relayStreamToWSClients waits for clients to connect and relays the given stream to
// connected websocket clients. If a client disconnects, it does no longer receives the stream.
func relayStreamToWSClients(stream <-chan *[]byte, clients <-chan *wsClient) {
	wsGroup := bcast.NewGroup()
	wsGroup.Broadcast(10 * time.Millisecond)

	// handle incoming clients
	go func() {
		for {
			// wait for clients to connect
			client := <-clients
			log.Println("Client connected: " + client.remoteAddress)
			subscription := wsGroup.Join()

			// start goroutine to monitor the connection
			go func() {
				<-client.isClosed
				wsGroup.Leave(subscription)
				log.Println("Client disconnected: " + client.remoteAddress)
			}()

			// start goroutine to consume broadcast
			go func() {
				for {
					fmt.Println("wait for recv")
					data := subscription.Recv().(*[]byte)
					fmt.Println("recv from bc")
					client.writeStream <- data
				}
			}()
		}
	}()

	// write stream into broadcast group
	go func() {
		for {
			data := <-stream
			fmt.Println("write to bc")
			wsGroup.Send(data)
		}
	}()
}
