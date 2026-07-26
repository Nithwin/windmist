package web

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/Nithwin/WindMist/internal/tools"
)

type FetchTool struct{}

func NewFetchTool() *FetchTool {
	return &FetchTool{}
}

func (t *FetchTool) Definition() tools.Definition {
	return tools.Definition{
		Name:        "fetch",
		Description: "Fetches the text content of a given URL. Useful for reading documentation pages or articles.",
		Parameters: []tools.Parameter{
			{
				Name:        "url",
				Type:        "string",
				Description: "The URL to fetch.",
				Required:    true,
			},
		},
	}
}

func (t *FetchTool) Run(ctx context.Context, call tools.Call) tools.Result {
	targetURL, ok := call.Args["url"].(string)
	if !ok || targetURL == "" {
		return tools.Result{Error: fmt.Errorf("url is required")}
	}

	req, err := http.NewRequestWithContext(ctx, "GET", targetURL, nil)
	if err != nil {
		return tools.Result{Error: err}
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return tools.Result{Error: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return tools.Result{Error: fmt.Errorf("HTTP %d", resp.StatusCode)}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return tools.Result{Error: err}
	}
	content := string(body)

	// Very naive HTML to text conversion to save tokens
	// Remove script and style tags
	scriptRe := regexp.MustCompile(`(?is)<script.*?>.*?</script>`)
	styleRe := regexp.MustCompile(`(?is)<style.*?>.*?</style>`)
	content = scriptRe.ReplaceAllString(content, "")
	content = styleRe.ReplaceAllString(content, "")

	// Remove all HTML tags
	tagRe := regexp.MustCompile(`(?is)<[^>]*>`)
	content = tagRe.ReplaceAllString(content, " ")

	// Condense whitespace
	wsRe := regexp.MustCompile(`\s+`)
	content = wsRe.ReplaceAllString(content, " ")
	content = strings.TrimSpace(content)

	// Truncate to reasonable length (e.g. 15000 chars) to not blow up context
	if len(content) > 15000 {
		content = content[:15000] + "\n... (truncated)"
	}

	return tools.Result{Output: content}
}
