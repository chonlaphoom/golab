# http-proto-buf

A simple HTTP messaging system using Protocol Buffers for message serialization.

## Overview

HTTP server and client for exchanging protobuf-encoded messages.

- **Receiver**: HTTP server listening on port 7777, accepts POST requests at `/message`
- **Sender**: Reads from stdin, serializes to protobuf, sends to receiver

## Message Format

```protobuf
message Message {
  string id = 1;
  string content = 2;
  int64 timestamp = 3;
}
```

## Usage

Run as receiver (server):

```bash
go run . receiver
```

Run as sender (client):

```bash
go run . sender
```

Then type messages in the sender terminal - they will appear in the receiver logs.

## Dependencies

- `google.golang.org/protobuf v1.36.11`
