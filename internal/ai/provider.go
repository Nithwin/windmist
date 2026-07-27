// Package ai defines the provider-agnostic interfaces and types for interacting
// with AI language models. It includes the Provider interface, request/response types,
// tool definitions, and a registry for dynamically loading provider implementations.
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
	) (*GenerateResponse, error)
}
