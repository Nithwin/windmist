package ai

import "context"

// Provider defines the behavior every AI provider must implement.
type Provider interface {
	Generate(
		ctx context.Context,
		req *GenerateRequest,
	) (*GenerateResponse, error)

	Stream(
		ctx context.Context,
		req *GenerateRequest,
		onChunk func(string),
	) error
}
