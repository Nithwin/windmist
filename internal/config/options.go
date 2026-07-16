package config

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
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
// Ollama intelligently checks daemon state, auto-starting or auto-pulling models upon user confirmation.
func GetModelOptions(providerName, ollamaBaseURL string) []selector.Option {
	var options []selector.Option

	if providerName == "ollama" {
		if ollamaBaseURL == "" {
			ollamaBaseURL = "http://localhost:11434"
		}
		options = ensureOllamaReadyAndGetModels(ollamaBaseURL)
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

func ensureOllamaReadyAndGetModels(baseURL string) []selector.Option {
	// 1. Check if Ollama binary is installed on the system
	if _, err := exec.LookPath("ollama"); err != nil {
		return []selector.Option{
			{
				Label:       "❌ Ollama CLI not installed",
				Description: "Please install Ollama from https://ollama.com first",
				Value:       "__CUSTOM__",
			},
		}
	}

	// 2. Attempt to fetch local models from the running daemon
	localModels, err := fetchOllamaModels(baseURL)
	if err != nil {
		// Daemon offline: Ask user if they want WindMist to start 'ollama serve'
		opt, runErr := selector.Run(
			"Ollama Daemon Not Running",
			fmt.Sprintf("Ollama server is offline at %s.\nWould you like WindMist to automatically start 'ollama serve' in the background?", baseURL),
			[]selector.Option{
				{Label: "Yes (Start 'ollama serve' right now and retry)", Description: "Launch Ollama background service automatically", Value: "yes"},
				{Label: "No (Skip auto-start)", Description: "Enter model ID manually or start Ollama yourself later", Value: "no"},
			},
		)
		if runErr == nil && opt.Value == "yes" {
			fmt.Println("\n⏳ Starting 'ollama serve' in background (waiting for initialization)...")
			cmd := exec.Command("ollama", "serve")
			if startErr := cmd.Start(); startErr == nil {
				// Poll the server up to 12 times (every 500ms, max 6 seconds) until ready
				for i := 0; i < 12; i++ {
					time.Sleep(500 * time.Millisecond)
					localModels, err = fetchOllamaModels(baseURL)
					if err == nil {
						break
					}
				}
			} else {
				fmt.Printf("⚠️ Could not start 'ollama serve': %v\n", startErr)
			}
		}
	}

	// 3. If daemon is online but no models are downloaded: Ask user if they want to auto-pull qwen2.5:8b
	if err == nil && len(localModels) == 0 {
		opt, runErr := selector.Run(
			"No Local Models Downloaded",
			"Ollama is running, but you have 0 models pulled to your system.\nWould you like WindMist to automatically pull 'qwen2.5:8b' right now?",
			[]selector.Option{
				{Label: "Yes (Run 'ollama pull qwen2.5:8b' right now)", Description: "Download recommended 8B local model (shows live progress)", Value: "yes"},
				{Label: "No (Skip and pull later)", Description: "Enter model ID manually or run 'ollama pull' yourself", Value: "no"},
			},
		)
		if runErr == nil && opt.Value == "yes" {
			fmt.Println("\n📥 Running 'ollama pull qwen2.5:8b' (please wait for download to complete)...")
			pullCmd := exec.Command("ollama", "pull", "qwen2.5:8b")
			pullCmd.Stdout = os.Stdout
			pullCmd.Stderr = os.Stderr
			if pullErr := pullCmd.Run(); pullErr == nil {
				fmt.Println("✔ Successfully pulled qwen2.5:8b!")
				localModels, _ = fetchOllamaModels(baseURL)
			} else {
				fmt.Printf("⚠️ Error pulling model: %v\n", pullErr)
			}
		}
	}

	// 4. Return whatever local models are available (if still empty after prompts, return explicit status item)
	if len(localModels) > 0 {
		return localModels
	}

	return []selector.Option{
		{
			Label:       "⚠️ Ollama offline or empty",
			Description: fmt.Sprintf("Run 'ollama serve' and 'ollama pull <model>' at %s", baseURL),
			Value:       "__CUSTOM__",
		},
	}
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
