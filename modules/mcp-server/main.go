package main

import (
	"log/slog"
	j "mcp-server/jsonrpc2"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	var buf [1024]byte
	for {
		select {
		case <-sigChan:
			slog.Info("Received shutdown signal, gracefully closing...")
			return
		default:
		}
		n, err := os.Stdin.Read(buf[:])

		if err != nil {
			slog.Error("Error reading from stdin:", err)
			return
		}

		if n == 0 {
			slog.Error("No data read from stdin")
			return
		}

		if n == len(buf) {
			slog.Error("Input too large to fit in buffer")
			return
		}

		msg := j.NewMessage()
		err = msg.ParseRequest(buf[:n])
		if err != nil {
			slog.Error("Error parsing JSON-RPC request:", err)
			return
		}

		msg.NewSuccessResponseWithResult("success calling method: " + msg.Request.Method)
		res, err := msg.SuccessResponse.UnmarshalJSON()
		if err != nil {
			slog.Error("Error generating success response JSON:", err)
			return
		}
		os.Stdout.Write(res)
	}
}
