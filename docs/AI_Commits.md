# AI Commit Messages

Lazygit can generate commit messages using an OpenAI‑compatible API (OpenAI, OpenRouter, or a self‑hosted endpoint).

## Setup

Add an `ai` section to your `config.yml` (open via `e` in the Status panel):

```yaml
ai:
  # One of: openai | openrouter | custom (used only for sensible defaults)
  provider: openai
  # If empty, defaults to provider base (openai: https://api.openai.com/v1, openrouter: https://openrouter.ai/api/v1)
  baseURL: ""
  # Required: model id (e.g. "gpt-4o-mini" or an OpenRouter model)
  model: ""
  # If empty, defaults to OPENAI_API_KEY (openai) or OPENROUTER_API_KEY (openrouter)
  apiKeyEnv: ""
  temperature: 0.2
  maxTokens: 300
  stagedOnly: true
  # One of: conventional | plain
  commitStyle: conventional
```

Then export your API key (shell example):

```bash
# For OpenAI
export OPENAI_API_KEY=your_key

# For OpenRouter
# export OPENROUTER_API_KEY=your_key
```

## Usage

1. Stage changes you want included in the message.
2. Press `c` to open the commit message panel.
3. Press `Ctrl+O` to open the commit menu.
4. Choose “Generate AI commit message” (shortcut `g`).

The summary and description fields will be populated with the generated message. You can edit them before confirming the commit as usual.

Notes:
- By default only the staged diff is sent. Set `ai.stagedOnly: false` to allow a fallback to consider tracked changes when nothing is staged.
- The client uses the OpenAI Chat Completions API schema and should work with any compatible endpoint.

