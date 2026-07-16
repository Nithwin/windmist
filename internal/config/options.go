package config

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Nithwin/WindMist/internal/ui/selector"
)

//go:embed models.json
var embeddedModelsJSON []byte

const remoteModelsManifestURL = "https://raw.githubusercontent.com/Nithwin/WindMist/main/internal/config/models.json"

type modelEntry struct {
	Label       string `json:"label"`
	Description string `json:"description"`
	Value       string `json:"value"`
}

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

// GetModelOptions returns model options for the specified provider.
// Cloud providers fetch dynamically from remote/embedded models.json manifest.
// Ollama fetches locally installed models from http://localhost:11434/api/tags.
func GetModelOptions(providerName, ollamaBaseURL string) []selector.Option {
	var options []selector.Option

	if providerName == "ollama" {
		if ollamaBaseURL == "" {
			ollamaBaseURL = "http://localhost:11434"
		}
		localModels, err := fetchOllamaModels(ollamaBaseURL)
		if err != nil {
			options = []selector.Option{
				{
					Label:       "⚠️ Ollama daemon offline",
					Description: fmt.Sprintf("Ensure 'ollama serve' is running at %s", ollamaBaseURL),
					Value:       "__CUSTOM__",
				},
			}
		} else if len(localModels) == 0 {
			options = []selector.Option{
				{
					Label:       "📥 No local models downloaded",
					Description: "Run 'ollama pull <model-name>' in terminal first",
					Value:       "__CUSTOM__",
				},
			}
		} else {
			options = localModels
		}
	} else {
		// Cloud providers: fetch from remote models.json or embedded fallback
		manifest := loadModelsManifest()
		if entries, ok := manifest[providerName]; ok {
			for _, e := range entries {
				options = append(options, selector.Option{
					Label:       e.Label,
					Description: e.Description,
					Value:       e.Value,
				})
			}
		}
	}

	// Always append custom model escape hatch
	options = append(options, selector.Option{
		Label:       "Custom model ID...",
		Description: "Enter any model name or identifier manually",
		Value:       "__CUSTOM__",
	})

	return options
}

func loadModelsManifest() map[string][]modelEntry {
	manifest := make(map[string][]modelEntry)

	// Attempt remote HTTP fetch with 1.5s timeout
	client := &http.Client{Timeout: 1500 * time.Millisecond}
	resp, err := client.Get(remoteModelsManifestURL)
	if err == nil && resp.StatusCode == http.StatusOK {
		if err := json.NewDecoder(resp.Body).Decode(&manifest); err == nil {
			resp.Body.Close()
			return manifest
		}
		resp.Body.Close()
	}

	// Fallback to embedded models.json
	_ = json.Unmarshal(embeddedModelsJSON, &manifest)
	return manifest
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
