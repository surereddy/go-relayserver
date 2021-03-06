package main

import (
	"flag"
	"fmt"
	"strings"
	"time"
)

const recordToFile = false

const defaultPortStream = ":8081"
const defaultPortWS = ":8082"
const bufferSize = 8 * 1000 * 1024 // 8MB

var recordName = fmt.Sprintf("%s_record.mpeg", time.Now().Format("20060102_1504"))

type config struct {
	portStream   string
	portWS       string
	secretStream string
	printHelp    bool
}

func main() {
	config := readCmdArguments()
	if config.printHelp {
		fmt.Println("Usage: ")
		fmt.Println("go-relayserver.exe optional: -port-stream <port> -port-ws <port> -s <secret>")
		return
	}

	RunRelayServer(config.portStream, config.portWS, config.secretStream)

	fmt.Println("Relay started, hit Enter-key to close")
	fmt.Scanln()
	fmt.Println("Shuting down...")
}

// RunRelayServer starts the relayserver.
// Listen for single incoming stream on: localhost:portStream/streamPassword.
// Listen for websockets on: localhost:portWS
func RunRelayServer(portStream string, portWS string, streamPassword string) {
	stream := waitForStream(portStream, streamPassword)
	if recordToFile {
		fmt.Println("Recording stream to " + recordName)
		fmt.Println("Warning: Recording stream may decrease performance and should be used for testing only")
		stream = recordStream(stream, recordName)
	}

	clients := waitForWSClients(portWS)
	relayStreamToWSClients(stream, clients)
}

func readCmdArguments() config {
	help := flag.Bool("h", false, "print help")

	portStream := flag.String("port-stream", defaultPortStream, "Port to listen for stream")
	portWS := flag.String("port-ws", defaultPortWS, "Port to listen for websockets")
	secretStream := flag.String("s", "", "Secure stream with this password")

	flag.Parse()

	return config{
		portStream:   normalizePort(*portStream),
		portWS:       normalizePort(*portWS),
		secretStream: *secretStream,
		printHelp:    *help,
	}
}

func normalizePort(port string) string {
	if !strings.HasPrefix(port, ":") {
		return ":" + port
	}
	return port
}
