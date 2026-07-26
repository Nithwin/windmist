# Build WindMist Step by Step

This guide explains how to build WindMist in a simple way.

WindMist is a terminal app that talks to an AI assistant. The assistant can answer questions, read files, and help with code.

## 1. Start With a Small Go CLI

First, make the program run from the terminal.

Example:

```bash
go run .
```

At first, the app can just print a short message:

```text
WindMist is starting...
```

The goal is only to confirm that the app starts.

## 2. Add One Command

Next, add a chat command.

Example:

```bash
windmist chat "Explain this project"
```

This command should take the text after `chat` and send it to the assistant.

Simple flow:

```text
user types prompt -> program reads prompt -> assistant answers
```

## 3. Add Simple Settings

The app needs a few settings:

- provider name
- model name
- API key

A small config file is enough.

Example:

```yaml
provider: gemini
model: gemini-2.0-flash
```

The app should load these settings before it starts chatting.

## 4. Build the Assistant Loop

This is the main idea.

The loop works like this:

1. The user sends a prompt.
2. The app sends the prompt to the model.
3. The model answers or asks for a tool.
4. The app runs the tool.
5. The tool result goes back to the model.
6. The model gives the final answer.

Example:

```text
User: Summarize the files in this folder.
AI: I need to check the folder first.
Tool: list files
AI: Here is the summary.
```

Keep this loop small and easy to read.

## 5. Add Basic Tools

Start with only a few tools.

Good beginner tools:

- `read` to open a file
- `write` to save a file
- `list` to show files in a folder

Example:

```text
AI asks to read README.md
The app reads the file
AI explains what it found
```

Later, you can add more tools if you want, but the basic version should stay small.

## 6. Connect the Model

The model is the part that thinks and responds.

Send these things to the model:

- the user prompt
- earlier messages
- tool results

Then wait for:

- a normal reply, or
- a request to use a tool

Example:

```text
Prompt: What does this function do?
Model: Please read internal/agent/loop.go
Tool result: file contents
Model: This function repeats the conversation until it finishes.
```

## 7. Add a Simple Terminal Screen

WindMist also has an interactive terminal screen.

The screen should only do a few things:

- show a place to type
- show the answer
- show streaming text while the model is replying

Example flow:

```text
Type a question
Press Enter
See the answer appear
```

You do not need a complex UI for the first version.

## 8. Keep the Code Split by Job

Put each part in its own place:

- CLI commands
- chat logic
- tools
- config
- terminal UI

That makes the project easier to understand.

## 9. Test the Main Parts

Start with a few simple tests.

Test:

- config loading
- reading and writing files
- prompt handling
- basic assistant responses

Example test idea:

```text
given a prompt
when the assistant runs
then it returns an answer
```

## 10. Write the Guide Like a Tutorial

If you submit this to a build-your-own-x style repo, keep the guide beginner friendly.

Use:

- short sections
- simple words
- clear steps
- small examples

Avoid advanced architecture language unless you explain it first.

## Summary

WindMist is a good tutorial topic because it shows how to build a real terminal AI assistant from simple parts.