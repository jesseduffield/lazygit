package ai

import (
    "fmt"
    "strings"
)

func systemPrompt(style string, wrap int) string {
    style = strings.ToLower(strings.TrimSpace(style))
    var rules []string
    rules = append(rules, "Write a high-quality git commit message from the provided diff.")
    switch style {
    case "conventional", "conv":
        rules = append(rules,
            "Use Conventional Commits format: <type>(<scope>): <subject>",
            "Types: feat, fix, docs, style, refactor, perf, test, chore, build, ci",
        )
    default:
        // plain: no extra rule
    }
    rules = append(rules,
        "Subject should be concise, ideally <= 72 chars.",
        fmt.Sprintf("Wrap body at ~%d chars when useful.", wrap),
        "Focus on what and why; avoid restating diff line-by-line.",
        "Do not include code fences or markdown headings.",
    )
    return strings.Join(rules, "\n")
}

func buildUserPrompt(diff string) string {
    return "Diff of changes (unified):\n" + diff + "\n\nReturn subject on first line, optional body after."
}

