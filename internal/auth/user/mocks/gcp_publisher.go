package mocks

import (
	"context"
	"fmt"
)

type userPublisher struct {
	topic string
}

func NewPublisher(ctx context.Context, topic string) (*userPublisher, error) {
	if topic == "" {
		return nil, fmt.Errorf("topic cannot be empty")
	}

	return &userPublisher{topic: topic}, nil
}

func (p *userPublisher) Publish(ctx context.Context, message string) error {
	return nil
}
