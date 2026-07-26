package web

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/Nithwin/WindMist/internal/tools"
)

type WebSearchTool struct{}

func NewWebSearchTool() *WebSearchTool {
	return &WebSearchTool{}
}

func (t *WebSearchTool) Definition() tools.Definition {
	return tools.Definition{
		Name:        "web_search",
		Description: "Searches the internet for a given query and returns a summary of the results with URLs. Useful for looking up documentation, error codes, and tutorials.",
		Category:    tools.CategoryWeb,
		Permission:  tools.PermReadOnly,
		Parameters: []tools.Parameter{
			{
				Name:        "query",
				Type:        "string",
				Description: "The search query.",
				Required:    true,
			},
		},
	}
}

type SearchResult struct {
	Title   string `json:"title"`
	Snippet string `json:"snippet"`
	URL     string `json:"url"`
}

func (t *WebSearchTool) Run(ctx context.Context, call tools.Call) tools.Result {
	query, ok := call.Args["query"].(string)
	if !ok || query == "" {
		return tools.Result{Error: fmt.Errorf("query is required")}
	}

	searchURL := fmt.Sprintf("https://html.duckduckgo.com/html/?q=%s", url.QueryEscape(query))
	req, err := http.NewRequestWithContext(ctx, "GET", searchURL, nil)
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return tools.Result{Error: err}
	}
	content := string(body)

	// Naive HTML parsing for DuckDuckGo results
	var results []SearchResult

	// Extract results using regex to avoid external HTML parser dependencies
	titleRe := regexp.MustCompile(`(?s)<a class="result__url" href="([^"]+)".*?>(.*?)</a>`)
	snippetRe := regexp.MustCompile(`(?s)<a class="result__snippet[^>]+>(.*?)</a>`)

	titles := titleRe.FindAllStringSubmatch(content, 5)
	snippets := snippetRe.FindAllStringSubmatch(content, 5)

	count := len(titles)
	if len(snippets) < count {
		count = len(snippets)
	}

	for i := 0; i < count; i++ {
		urlStr := titles[i][1]
		if strings.HasPrefix(urlStr, "//duckduckgo.com/l/?uddg=") {
			// Extract actual URL
			urlStr = strings.TrimPrefix(urlStr, "//duckduckgo.com/l/?uddg=")
			urlStr = strings.Split(urlStr, "&rut=")[0]
			decoded, _ := url.QueryUnescape(urlStr)
			if decoded != "" {
				urlStr = decoded
			}
		}

		title := stripHTML(titles[i][2])
		snippet := stripHTML(snippets[i][1])

		results = append(results, SearchResult{
			Title:   title,
			Snippet: snippet,
			URL:     urlStr,
		})
	}

	if len(results) == 0 {
		return tools.Result{Output: "No results found or rate limited."}
	}

	return tools.Result{Output: results}
}

func stripHTML(str string) string {
	re := regexp.MustCompile(`<[^>]*>`)
	return strings.TrimSpace(re.ReplaceAllString(str, ""))
}
