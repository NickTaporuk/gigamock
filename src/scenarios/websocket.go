package scenarios

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
)

// WebSocketScenario describes scripted bidirectional WebSocket behavior.
type WebSocketScenario struct {
	Name          string
	DryRun        bool
	SendOnConnect []WebSocketMessage
	Steps         []WebSocketStep
}

func (ws WebSocketScenario) Validate() error {
	if !ws.DryRun && len(ws.SendOnConnect) == 0 && len(ws.Steps) == 0 {
		return fmt.Errorf("sendOnConnect or steps are required unless dryRun is true")
	}

	for index, message := range ws.SendOnConnect {
		if err := message.Validate(); err != nil {
			return fmt.Errorf("sendOnConnect message %d is invalid: %w", index, err)
		}
	}
	for index, step := range ws.Steps {
		if err := step.Validate(); err != nil {
			return fmt.Errorf("step %d is invalid: %w", index, err)
		}
	}

	return nil
}

// WebSocketMessage is a text or binary WebSocket message.
type WebSocketMessage struct {
	Type string
	Text string
}

func (wsm WebSocketMessage) Validate() error {
	return validation.ValidateStruct(
		&wsm,
		validation.Field(&wsm.Text, validation.Required),
		validation.Field(&wsm.Type, validation.In("", "text", "binary")),
	)
}

// WebSocketStep is one scripted receive/send/close action.
type WebSocketStep struct {
	Receive *WebSocketMessage
	Send    *WebSocketMessage
	Close   *WebSocketClose
	Delay   string
}

func (wss WebSocketStep) Validate() error {
	if wss.Receive == nil && wss.Send == nil && wss.Close == nil && wss.Delay == "" {
		return fmt.Errorf("receive, send, close, or delay is required")
	}
	if wss.Receive != nil {
		if err := wss.Receive.Validate(); err != nil {
			return fmt.Errorf("receive is invalid: %w", err)
		}
	}
	if wss.Send != nil {
		if err := wss.Send.Validate(); err != nil {
			return fmt.Errorf("send is invalid: %w", err)
		}
	}
	if wss.Close != nil {
		if err := wss.Close.Validate(); err != nil {
			return fmt.Errorf("close is invalid: %w", err)
		}
	}
	return nil
}

// WebSocketClose describes an intentional close frame.
type WebSocketClose struct {
	Code   int
	Reason string
}

func (wsc WebSocketClose) Validate() error {
	return validation.ValidateStruct(
		&wsc,
		validation.Field(&wsc.Code, validation.Min(1000), validation.Max(4999)),
	)
}

// WebSocketScenarios is a list of WebSocket scenarios.
type WebSocketScenarios []WebSocketScenario
