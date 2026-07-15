#!/usr/bin/env bash
# Syncs GitHub repository labels from .github/labels.json to the upstream repository.
# Usage: ./.github/scripts/sync_labels.sh [repository_owner/repository_name]
# Example: ./.github/scripts/sync_labels.sh Nithwin/windmist

set -euo pipefail

REPO="${1:-}"

if ! command -v gh &> /dev/null; then
  echo "Error: GitHub CLI (gh) is not installed or not in PATH."
  echo "Please install gh from https://cli.github.com/ and run 'gh auth login' before executing this script."
  exit 1
fi

# Check if repo argument was supplied, or allow gh to detect from git remote
GH_REPO_FLAG=""
if [ -n "$REPO" ]; then
  GH_REPO_FLAG="--repo $REPO"
fi

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
LABELS_JSON="$SCRIPT_DIR/../labels.json"

if [ ! -f "$LABELS_JSON" ]; then
  echo "Error: Cannot find labels definition file at $LABELS_JSON"
  exit 1
fi

echo "Syncing 13 labels to GitHub repository ($REPO)..."

# Array of labels with (Name Color Description)
declare -a LABELS=(
  "bug|d73a4a|Something isn't working as expected or crashes"
  "enhancement|a2eeef|New feature, capability, or CLI request"
  "documentation|0075ca|Improvements or additions to documentation and help text"
  "good first issue|7057ff|Good for newcomers diving into the WindMist codebase"
  "help wanted|008672|Extra attention or community assistance is requested"
  "question|d876e3|Further information is requested or general inquiry"
  "performance|ff9f1c|Optimizing execution speed, memory footprint, or token latency"
  "security|e11d48|Security vulnerability, secret handling, or safety check"
  "provider|3b82f6|Relates to AI model providers (Gemini, OpenAI, Anthropic, Ollama)"
  "tooling|10b981|Relates to CLI flags, Glamour/Lip Gloss TUI, or build tools"
  "testing|fbbf24|Adding or improving unit tests, golden tests, or CI checks"
  "refactor|8b5cf6|Code cleanup or architectural restructuring without changing behavior"
  "dependencies|0366d6|Pull requests that update go.mod packages or external dependencies"
)

for label_entry in "${LABELS[@]}"; do
  IFS="|" read -r name color desc <<< "$label_entry"
  echo -n "Syncing label: '$name'... "
  # Try to edit existing label first; if not found, create it
  if gh label edit "$name" --color "$color" --description "$desc" $GH_REPO_FLAG &>/dev/null; then
    echo "Updated existing label."
  else
    if gh label create "$name" --color "$color" --description "$desc" --force $GH_REPO_FLAG &>/dev/null; then
      echo "Created new label."
    else
      echo "Failed to sync '$name'."
    fi
  fi
done

echo "Label sync complete!"
