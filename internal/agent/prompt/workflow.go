package prompt

// Workflow returns compact engineering + tool-usage guidance.
// Kept short because tool JSON schemas already describe each tool.
func Workflow() string {
	return `## Workflow
1. Inspect before editing (search/read only what you need).
2. Make the smallest safe change; avoid unrelated edits.
3. Prefer precise edits over rewriting whole files.
4. Match existing style; don't add deps/abstractions unless needed.
5. On tool errors: re-read context, then retry once with a corrected approach.
6. Stop when the request is done—don't keep calling tools.`
}

// Developer is kept for backward compatibility; prefer Workflow().
func Developer() string { return Workflow() }

// Tools is kept for backward compatibility; prefer Workflow().
func Tools() string { return Workflow() }
