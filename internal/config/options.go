package config

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Nithwin/WindMist/internal/ui/selector"
)

// GetProviderOptions returns the list of all supported AI providers with descriptions.
func GetProviderOptions() []selector.Option {
	return []selector.Option{
		{
			Label:       "gemini",
			Description: "Google Gemini — Fast, highly capable multimodal AI (Default)",
			Value:       "gemini",
		},
		{
			Label:       "openai",
			Description: "OpenAI — Flagship models like GPT-4o, o1, o3-mini",
			Value:       "openai",
		},
		{
			Label:       "anthropic",
			Description: "Anthropic — Claude 3.5 Sonnet, Haiku, Opus models",
			Value:       "anthropic",
		},
		{
			Label:       "groq",
			Description: "Groq — Ultra-fast Llama 3 and Mixtral inference",
			Value:       "groq",
		},
		{
			Label:       "ollama",
			Description: "Ollama — Run open-source models locally on your system",
			Value:       "ollama",
		},
	}
}

// GetModelOptions returns the curated model options for a provider.
// For Ollama, it attempts to fetch locally installed models via HTTP.
func GetModelOptions(providerName, ollamaBaseURL string) []selector.Option {
	var options []selector.Option

	switch providerName {
	case "gemini":
		options = []selector.Option{
			{Label: "gemini-2.5-flash", Description: "Fast, highly capable multimodal model (Default)", Value: "gemini-2.5-flash"},
			{Label: "gemini-2.5-pro", Description: "Advanced reasoning and complex task model", Value: "gemini-2.5-pro"},
			{Label: "gemini-1.5-flash", Description: "Previous generation high-speed model", Value: "gemini-1.5-flash"},
			{Label: "gemini-1.5-pro", Description: "Previous generation professional model", Value: "gemini-1.5-pro"},
		}

	case "openai":
		options = []selector.Option{
			{Label: "gpt-4o", Description: "Flagship multimodal intelligence (Default)", Value: "gpt-4o"},
			{Label: "gpt-4o-mini", Description: "Fast, cost-effective small model", Value: "gpt-4o-mini"},
			{Label: "o3-mini", Description: "Latest fast reasoning model", Value: "o3-mini"},
			{Label: "o1", Description: "Advanced reasoning model", Value: "o1"},
			{Label: "o1-mini", Description: "Fast reasoning model", Value: "o1-mini"},
		}

	case "anthropic":
		options = []selector.Option{
			{Label: "claude-3-5-sonnet-latest", Description: "Most intelligent Claude model for coding and reasoning (Default)", Value: "claude-3-5-sonnet-latest"},
			{Label: "claude-3-5-haiku-latest", Description: "Fastest Claude model for quick tasks", Value: "claude-3-5-haiku-latest"},
			{Label: "claude-3-opus-latest", Description: "Powerful deep reasoning model", Value: "claude-3-opus-latest"},
		}

	case "groq":
		options = []selector.Option{
			{Label: "llama-3.3-70b-versatile", Description: "Fast, versatile Llama 3.3 70B model (Default)", Value: "llama-3.3-70b-versatile"},
			{Label: "llama-3.1-8b-instant", Description: "Ultra-fast low latency 8B model", Value: "llama-3.1-8b-instant"},
			{Label: "mixtral-8x7b-32768", Description: "Mixtral MoE fast model", Value: "mixtral-8x7b-32768"},
			{Label: "gemma2-9b-it", Description: "Google Gemma 2 9B model on Groq", Value: "gemma2-9b-it"},
		}

	case "ollama":
		if ollamaBaseURL == "" {
			ollamaBaseURL = "http://localhost:11434"
		}
		localModels, err := fetchOllamaModels(ollamaBaseURL)
		if err == nil && len(localModels) > 0 {
			options = append(options, localModels...)
		} else {
			options = []selector.Option{
				{Label: "qwen2.5:8b", Description: "Recommended local model (ensure pulled via 'ollama pull qwen2.5:8b')", Value: "qwen2.5:8b"},
				{Label: "llama3.2:3b", Description: "Lightweight local model (ensure pulled via 'ollama pull llama3.2:3b')", Value: "llama3.2:3b"},
			}
		}
	}

	options = append(options, selector.Option{
		Label:       "Custom model ID...",
		Description: "Enter any model name or identifier manually",
		Value:       "__CUSTOM__",
	})

	return options
}

type ollamaTagsResponse struct {
	Models []struct {
		Name    string `json:"name"`
		Details struct {
			ParameterSize     string `json:"parameter_size"`
			QuantizationLevel string `json:"quantization_level"`
		} `json:"details"`
	} `json:"models"`
}

func fetchOllamaModels(baseURL string) ([]selector.Option, error) {
	client := &http.Client{Timeout: 2 * time.Second}
	url := strings.TrimRight(baseURL, "/") + "/api/tags"
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var data ollamaTagsResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	var options []selector.Option
	for _, m := range data.Models {
		desc := "Installed local model"
		if m.Details.ParameterSize != "" {
			desc = fmt.Sprintf("Installed local (%s, %s)", m.Details.ParameterSize, m.Details.QuantizationLevel)
		}
		options = append(options, selector.Option{
			Label:       m.Name,
			Description: desc,
			Value:       m.Name,
		})
	}

	return options, nil
}
