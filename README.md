# Adapter project

## Project Introduction

### Purpose

This project provides a flexible adapter that enables communication between various infrastructure components and protocols. It acts as a bridge between devices and external systems (e.g., cloud services or other devices), facilitating protocol translation through “processors”.

The adapter can be deployed on devices or as a cloud service, allowing for dynamic integration in both testing and production environments.

### Example Adapter Use Cases

- Testing Adapter
	- Purpose: Enables testing interactions and validation of internal events
	- Implementation:
      - HTTP for sending internal commands, events, and queries
      - WebSocket for receiving internal events
- Factory End-of-Line (EOL) Adapter
    - Purpose: Facilitates device operations like calibration, configuration, and tests during manufacturing
    - Implementation:
      - HTTP requests for sending commands, events, and queries
      - WebSocket for receiving events and tracking operation status
- Cloud Service Adapter
    - Purpose: Connects a device to the cloud, facilitating monitoring, remote operations and get device status
    - Implementation:
      - MQTT for device commands
      - MQTT for receiving internal events


## Shared language

- Transporter:
	- Infrastructure component that is responsible to connect one side of the system to another side
	- Examples:
		- HTTP request to Redis
		- Redis to Websocket
		- MQTT to Redis
		- Redis to MQTT
		- ...
- HTTP handler:
	- It is an infrastructure component that let read from a HTTP request and write to a HTTP response
- TopicMessage:
	- It is a concept that represents a message that is coming or going to a topic
	- It is wrapping a message and a topic
- TopicMessageProcessor:
	- It is a concept that is used to process TopicMessages
	- It can be used to:
		- Translate messages from one protocol to another
		- Filter messages


## Project Architecture

This project includes:
- Transporters:
	- HTTP to Redis channel (`HttpToRedisChannelTransporter`)
	- Redis channel to Websocket (`RedisChannelToWebsocketTransporter`)
- Processors:
	- Command processor (`commands.Processor`)
	- Event processor (`events.Processor`)

Adapter Functionalities:
- Using the provided transporters and processors, the adapter can:
  - Process “Command” HTTP POST requests and send them to a Redis channel
  - Receive “Event” messages from a Redis channel and forward them to a WebSocket


## Running the Project

1. Setup
	- `make setup`
2. Run development environment
	- `make start-dev`
3. Run tests
	- `make test`


## Building and Running

1. Build the project
	- `make build`
2. Run the project
	- `./build/http-ws-server`


## Additional Information

- This project is following [golang-standards/project-layout](https://github.com/golang-standards/project-layout)
- This project is using a dependency injection pattern (found in `cmd/di`) based on my previous projects
- Tests:
	- Unit tests are in the `internal/app/http-ws-server/domain` package
      - These tests are testing the business logic without any external infrastructure
	- Integration tests are in the `internal/app/http-ws-server/infrastructure/transporter` package
      - These tests are testing the infrastructure components with specific implementations (e.g., Redis, Websocket) 
    - Acceptance tests are in the `test` folder
      - These tests are end-to-end. They are testing the whole system with real implementations (e.g., Redis, Websocket)
