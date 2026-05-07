package common

import (
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	is "github.com/go-ozzo/ozzo-validation/v4/is"
)

const (
	// HTTPScenarioType represents type "http" for a scenario
	HTTPScenarioType = "http"
	// GraphQLScenarioType represents type "graphql" for a scenario
	GraphQLScenarioType = "graphql"
	// KafkaScenarioType represents type "kafka" for a scenario
	KafkaScenarioType = "kafka"
	// GRPCScenarioType represents type "grpc" for a scenario
	GRPCScenarioType = "grpc"
	// NATSScenarioType represents type "nats" for a scenario
	NATSScenarioType = "nats"
	// RabbitMQScenarioType represents type "rabbitmq" for a scenario
	RabbitMQScenarioType = "rabbitmq"
	// MQTTScenarioType represents type "mqtt" for a scenario
	MQTTScenarioType = "mqtt"
	// WebSocketScenarioType represents type "websocket" for a scenario
	WebSocketScenarioType = "websocket"
	// S3ScenarioType represents type "s3" for a scenario
	S3ScenarioType = "s3"
	// SQSScenarioType represents type "sqs" for a scenario
	SQSScenarioType = "sqs"
	// SNSScenarioType represents type "sns" for a scenario
	SNSScenarioType = "sns"
	// PubSubScenarioType represents type "pubsub" for a scenario
	PubSubScenarioType = "pubsub"
	// ServiceBusScenarioType represents type "servicebus" for a scenario
	ServiceBusScenarioType = "servicebus"
	// SOAPScenarioType represents type "soap" for a scenario
	SOAPScenarioType = "soap"
)

var (
	// ScenarioTypeValidator is a validator rule for the type of a scenario
	// can be http, graphql, kafka, grpc, nats, rabbitmq, mqtt, websocket, s3, sqs, sns, pubsub, servicebus, or soap
	ScenarioTypeValidator = []validation.Rule{
		validation.Required,
		validation.In(
			HTTPScenarioType,
			GraphQLScenarioType,
			KafkaScenarioType,
			GRPCScenarioType,
			NATSScenarioType,
			RabbitMQScenarioType,
			MQTTScenarioType,
			WebSocketScenarioType,
			S3ScenarioType,
			SQSScenarioType,
			SNSScenarioType,
			PubSubScenarioType,
			ServiceBusScenarioType,
			SOAPScenarioType,
		),
	}
	// ScenarioMethodValidator is a validator rule for the method type of a scenario
	// must be any type of HTTP methods
	ScenarioMethodValidator = []validation.Rule{
		validation.Required,
		validation.In(http.MethodPost, http.MethodGet, http.MethodPut,
			http.MethodConnect, http.MethodDelete, http.MethodHead,
			http.MethodOptions, http.MethodPatch, http.MethodTrace),
	}
	// BaseScenariosValidator is a validator rule for a base validation of the field scenarios
	// must be required
	BaseScenariosValidator = []validation.Rule{
		validation.Required,
	}
	// CodeStatus is a validator rule for the field codeStatus
	// must be required
	CodeStatusValidator = []validation.Rule{
		validation.Required,
		validation.Min(http.StatusOK),
		validation.Max(http.StatusNetworkAuthenticationRequired),
	}
	// URLPathValidator is a validator rule for the URL path
	// must be required
	// a string must be a valid URL
	URLPathValidator = []validation.Rule{
		validation.Required,
		is.URL,
	}
)
