package agent

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

var (
	greetingRE = regexp.MustCompile(`(?i)^(hi|hello|hey|yo|sup|hola|howdy|good\s+(morning|afternoon|evening)|thanks|thank\s*you|thx|ty|ok|okay|cool|great|nice|bye|goodbye|see\s+ya)[\s!.?]*$`)
	buildRE    = regexp.MustCompile(`(?i)\b(implement|fix|create|add|write|edit|modify|update|refactor|rename|delete|remove|install|migrate|generate|scaffold|build|make|patch|apply|replace|insert|debug|resolve)\b`)
	planRE     = regexp.MustCompile(`(?i)\b(plan|analyze|analyse|review|explain|search|find|locate|how\s+(do|does|to)|what\s+(is|are|does)|why\s+(is|are|does)|where\s+(is|are)|compare|outline|architect|design\s+a|walk\s*me\s*through|summarize|summarise)\b`)
	codeHintRE = regexp.MustCompile(`(?i)(\x60{3}|@[\w./\\-]+|\.(go|ts|tsx|js|jsx|py|rs|java|kt|c|cpp|h|cs|rb|php|swift)\b|func\s+\w+|class\s+\w+|def\s+\w+)`)
)

// ClassifyIntent resolves auto-mode without an LLM call when confidence is high.
// Returns (mode, ok). When ok is false, the caller should use the LLM router.
func ClassifyIntent(userPrompt string) (Mode, bool) {
	p := strings.TrimSpace(userPrompt)
	if p == "" {
		return ModeChat, true
	}

	lower := strings.ToLower(p)
	runes := utf8.RuneCountInString(p)

	// Pure greetings / acknowledgements — never need plan or build.
	if greetingRE.MatchString(strings.TrimSpace(lower)) {
		return ModeChat, true
	}

	// Very short small-talk without code signals.
	if runes <= 40 && !codeHintRE.MatchString(p) && !buildRE.MatchString(lower) {
		// "hi there", "how's it going", "who are you", etc.
		if !planRE.MatchString(lower) || runes <= 20 {
			return ModeChat, true
		}
	}

	hasBuild := buildRE.MatchString(lower)
	hasPlan := planRE.MatchString(lower)
	hasCode := codeHintRE.MatchString(p)

	// Explicit implementation intent.
	if hasBuild && (hasCode || runes > 25) {
		return ModeBuild, true
	}
	if hasBuild && !hasPlan {
		return ModeBuild, true
	}

	// Analysis / explanation without edit verbs.
	if hasPlan && !hasBuild {
		// Short conceptual questions stay in chat (no tools/repo map tax).
		if runes < 100 && !hasCode && !strings.Contains(lower, "codebase") &&
			!strings.Contains(lower, "this project") && !strings.Contains(lower, "this repo") {
			return ModeChat, true
		}
		return ModePlan, true
	}

	// Short questions with no engineering verbs → chat.
	if runes < 60 && !hasBuild && !hasCode {
		return ModeChat, true
	}

	return "", false
}

// IsTrivialPrompt reports whether a message is too trivial to spend an
// auto-title API call on (greetings, thanks, etc.).
func IsTrivialPrompt(userPrompt string) bool {
	mode, ok := ClassifyIntent(userPrompt)
	return ok && mode == ModeChat && utf8.RuneCountInString(strings.TrimSpace(userPrompt)) <= 40
}
