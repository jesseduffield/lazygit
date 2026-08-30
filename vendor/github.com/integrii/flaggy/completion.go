package flaggy

import (
	"fmt"
	"strings"
)

// EnableCompletion enables shell autocomplete outputs to be generated.
func EnableCompletion() {
	DefaultParser.ShowCompletion = true
}

// DisableCompletion disallows shell autocomplete outputs to be generated.
func DisableCompletion() {
	DefaultParser.ShowCompletion = false
}

// GenerateBashCompletion returns a bash completion script for the parser.
func GenerateBashCompletion(p *Parser) string {
	var b strings.Builder
	funcName := "_" + sanitizeName(p.Name) + "_complete"
	b.WriteString("# bash completion for " + p.Name + "\n")
	b.WriteString(funcName + "() {\n")
	b.WriteString("    local cur prev\n")
	b.WriteString("    COMPREPLY=()\n")
	b.WriteString("    cur=\"${COMP_WORDS[COMP_CWORD]}\"\n")
	b.WriteString("    prev=\"${COMP_WORDS[COMP_CWORD-1]}\"\n")
	b.WriteString("    case \"$prev\" in\n")
	bashCaseEntries(&p.Subcommand, &b)
	rootOpts := collectOptions(&p.Subcommand)
	b.WriteString("        *)\n            COMPREPLY=( $(compgen -W \"" + rootOpts + "\" -- \"$cur\") )\n            return 0\n            ;;\n    esac\n}\n")
	b.WriteString("complete -F " + funcName + " " + p.Name + "\n")
	return b.String()
}

// GenerateZshCompletion returns a zsh completion script for the parser.
func GenerateZshCompletion(p *Parser) string {
	var b strings.Builder
	funcName := "_" + sanitizeName(p.Name)
	b.WriteString("#compdef " + p.Name + "\n\n")
	b.WriteString(funcName + "() {\n")
	b.WriteString("    local cur prev\n")
	b.WriteString("    cur=${words[CURRENT]}\n")
	b.WriteString("    prev=${words[CURRENT-1]}\n")
	b.WriteString("    case \"$prev\" in\n")
	zshCaseEntries(&p.Subcommand, &b)
	rootOpts := collectOptions(&p.Subcommand)
	b.WriteString("        *)\n            compadd -- " + rootOpts + "\n            ;;\n    esac\n}\n")
	b.WriteString("compdef " + funcName + " " + p.Name + "\n")
	return b.String()
}

// GenerateFishCompletion returns a fish completion script for the parser.
func GenerateFishCompletion(p *Parser) string {
	var b strings.Builder
	b.WriteString("# fish completion for " + p.Name + "\n")
	writeFishEntries(&p.Subcommand, &b, p.Name, nil)
	return b.String()
}

// GeneratePowerShellCompletion returns a PowerShell completion script for the parser.
func GeneratePowerShellCompletion(p *Parser) string {
	var b strings.Builder
	b.WriteString("# PowerShell completion for " + p.Name + "\n")
	b.WriteString("Register-ArgumentCompleter -CommandName '" + p.Name + "' -ScriptBlock {\n")
	b.WriteString("    param($commandName, $parameterName, $wordToComplete, $commandAst, $fakeBoundParameters)\n")
	b.WriteString("    $completions = @(\n")
	writePowerShellEntries(&p.Subcommand, &b)
	b.WriteString("    )\n")
	b.WriteString("    $completions | Where-Object { $_.CompletionText -like \"$wordToComplete*\" }\n")
	b.WriteString("}\n")
	return b.String()
}

// GenerateNushellCompletion returns a Nushell completion script for the parser.
func GenerateNushellCompletion(p *Parser) string {
	var b strings.Builder
	command := p.Name
	funcName := "nu-complete " + command
	b.WriteString("# nushell completion for " + command + "\n")
	b.WriteString("def \"" + funcName + "\" [] {\n")
	b.WriteString("    [\n")
	writeNushellEntries(&p.Subcommand, &b)
	b.WriteString("    ]\n")
	b.WriteString("}\n\n")
	b.WriteString("extern \"" + command + "\" [\n")
	writeNushellFlagSignature(&p.Subcommand, &b)
	b.WriteString("    command?: string@\"" + funcName + "\"\n")
	b.WriteString("]\n")
	return b.String()
}

// collectOptions builds a space-delimited list of flags, subcommands, and positional values
// for the provided subcommand.
func collectOptions(sc *Subcommand) string {
	var opts []string
	for _, f := range sc.Flags {
		if len(f.ShortName) > 0 {
			opts = append(opts, "-"+f.ShortName)
		}
		if len(f.LongName) > 0 {
			opts = append(opts, "--"+f.LongName)
		}
	}
	for _, p := range sc.PositionalFlags {
		if p.Name != "" {
			opts = append(opts, p.Name)
		}
	}
	for _, s := range sc.Subcommands {
		if s.Hidden {
			continue
		}
		if s.Name != "" {
			opts = append(opts, s.Name)
		}
		if s.ShortName != "" {
			opts = append(opts, s.ShortName)
		}
	}
	return strings.Join(opts, " ")
}

func bashCaseEntries(sc *Subcommand, b *strings.Builder) {
	for _, s := range sc.Subcommands {
		if s.Hidden {
			continue
		}
		opts := collectOptions(s)
		b.WriteString("        " + s.Name + ")\n            COMPREPLY=( $(compgen -W \"" + opts + "\" -- \"$cur\") )\n            return 0\n            ;;\n")
		if s.ShortName != "" {
			b.WriteString("        " + s.ShortName + ")\n            COMPREPLY=( $(compgen -W \"" + opts + "\" -- \"$cur\") )\n            return 0\n            ;;\n")
		}
		bashCaseEntries(s, b)
	}
}

func zshCaseEntries(sc *Subcommand, b *strings.Builder) {
	for _, s := range sc.Subcommands {
		if s.Hidden {
			continue
		}
		opts := collectOptions(s)
		b.WriteString("        " + s.Name + ")\n            compadd -- " + opts + "\n            return\n            ;;\n")
		if s.ShortName != "" {
			b.WriteString("        " + s.ShortName + ")\n            compadd -- " + opts + "\n            return\n            ;;\n")
		}
		zshCaseEntries(s, b)
	}
}

func sanitizeName(n string) string {
	return strings.ReplaceAll(n, "-", "_")
}

// writeFishEntries builds the fish completion statements for the provided subcommand path so
// the generated script mirrors Flaggy's flag and subcommand hierarchy for interactive use.
func writeFishEntries(sc *Subcommand, b *strings.Builder, command string, path []string) {
	condition := fishConditionForFlags(path)
	for _, f := range sc.Flags {
		if f.Hidden {
			continue
		}
		line := "complete -c " + command
		if condition != "" {
			line += " -n '" + condition + "'"
		}
		if f.ShortName != "" {
			line += " -s " + f.ShortName
		}
		if f.LongName != "" {
			line += " -l " + f.LongName
		}
		if f.Description != "" {
			line += " -d '" + escapeSingleQuotes(f.Description) + "'"
		}
		line += "\n"
		b.WriteString(line)
	}
	for _, p := range sc.PositionalFlags {
		if p.Hidden {
			continue
		}
		if p.Name == "" {
			continue
		}
		line := "complete -c " + command
		if condition != "" {
			line += " -n '" + condition + "'"
		}
		line += " -a '" + escapeSingleQuotes(p.Name) + "'"
		if p.Description != "" {
			line += " -d '" + escapeSingleQuotes(p.Description) + "'"
		}
		line += "\n"
		b.WriteString(line)
	}
	subCondition := fishConditionForSubcommands(path)
	for _, sub := range sc.Subcommands {
		if sub.Hidden {
			continue
		}
		line := "complete -c " + command
		if subCondition != "" {
			line += " -n '" + subCondition + "'"
		}
		line += " -a '" + escapeSingleQuotes(sub.Name) + "'"
		if sub.Description != "" {
			line += " -d '" + escapeSingleQuotes(sub.Description) + "'"
		}
		line += "\n"
		b.WriteString(line)
		if sub.ShortName != "" {
			aliasLine := "complete -c " + command
			if subCondition != "" {
				aliasLine += " -n '" + subCondition + "'"
			}
			aliasLine += " -a '" + escapeSingleQuotes(sub.ShortName) + "'"
			if sub.Description != "" {
				aliasLine += " -d '" + escapeSingleQuotes(sub.Description) + "'"
			}
			aliasLine += "\n"
			b.WriteString(aliasLine)
		}
		nextPath := appendPath(path, sub.Name)
		writeFishEntries(sub, b, command, nextPath)
	}
}

// fishConditionForFlags returns the fish condition needed to scope flag suggestions to the
// current subcommand path while leaving root flags globally available.
func fishConditionForFlags(path []string) string {
	if len(path) == 0 {
		return ""
	}
	return "__fish_seen_subcommand_from " + path[len(path)-1]
}

// fishConditionForSubcommands returns the fish condition that ensures subcommand suggestions
// appear only after their parent command token has been entered.
func fishConditionForSubcommands(path []string) string {
	if len(path) == 0 {
		return "__fish_use_subcommand"
	}
	return "__fish_seen_subcommand_from " + path[len(path)-1]
}

// appendPath creates a new slice with the next subcommand name appended so recursive
// completion builders can keep the traversal stack immutable.
func appendPath(path []string, value string) []string {
	next := make([]string, len(path)+1)
	copy(next, path)
	next[len(path)] = value
	return next
}

// escapeSingleQuotes prepares text for inclusion in single-quoted shell strings so flag
// descriptions render safely in the generated scripts.
func escapeSingleQuotes(s string) string {
	return strings.ReplaceAll(s, "'", "\\'")
}

// escapeDoubleQuotes prepares text for inclusion in double-quoted shell strings which is
// required for PowerShell and Nushell emission.
func escapeDoubleQuotes(s string) string {
	return strings.ReplaceAll(s, "\"", "\\\"")
}

// writePowerShellEntries walks the parser tree and emits CompletionResult entries so the
// PowerShell script can surface flags, positionals, and subcommands interactively.
func writePowerShellEntries(sc *Subcommand, b *strings.Builder) {
	for _, f := range sc.Flags {
		if f.Hidden {
			continue
		}
		if f.LongName != "" {
			writePowerShellLine("--"+f.LongName, f.Description, "ParameterName", b)
		}
		if f.ShortName != "" {
			writePowerShellLine("-"+f.ShortName, f.Description, "ParameterName", b)
		}
	}
	for _, p := range sc.PositionalFlags {
		if p.Hidden {
			continue
		}
		if p.Name == "" {
			continue
		}
		writePowerShellLine(p.Name, p.Description, "ParameterValue", b)
	}
	for _, sub := range sc.Subcommands {
		if sub.Hidden {
			continue
		}
		writePowerShellLine(sub.Name, sub.Description, "Command", b)
		if sub.ShortName != "" {
			writePowerShellLine(sub.ShortName, sub.Description, "Command", b)
		}
		writePowerShellEntries(sub, b)
	}
}

// writePowerShellLine emits a single CompletionResult definition with the supplied tooltip and
// completion type for consumption by Register-ArgumentCompleter.
func writePowerShellLine(value, description, kind string, b *strings.Builder) {
	tooltip := description
	if tooltip == "" {
		tooltip = value
	}
	line := fmt.Sprintf("        [System.Management.Automation.CompletionResult]::new(\"%s\", \"%s\", \"%s\", \"%s\")\n", escapeDoubleQuotes(value), escapeDoubleQuotes(value), kind, escapeDoubleQuotes(tooltip))
	b.WriteString(line)
}

// writeNushellEntries collects all completion values into Nushell's structured format so
// external commands can expose their interactive help inside the shell.
func writeNushellEntries(sc *Subcommand, b *strings.Builder) {
	for _, f := range sc.Flags {
		if f.Hidden {
			continue
		}
		if f.LongName != "" {
			writeNushellLine("--"+f.LongName, f.Description, b)
		}
		if f.ShortName != "" {
			writeNushellLine("-"+f.ShortName, f.Description, b)
		}
	}
	for _, p := range sc.PositionalFlags {
		if p.Hidden {
			continue
		}
		if p.Name == "" {
			continue
		}
		writeNushellLine(p.Name, p.Description, b)
	}
	for _, sub := range sc.Subcommands {
		if sub.Hidden {
			continue
		}
		writeNushellLine(sub.Name, sub.Description, b)
		if sub.ShortName != "" {
			writeNushellLine(sub.ShortName, sub.Description, b)
		}
		writeNushellEntries(sub, b)
	}
}

// writeNushellLine emits a single structured completion item for Nushell with a value and
// friendly description.
func writeNushellLine(value, description string, b *strings.Builder) {
	tooltip := description
	if tooltip == "" {
		tooltip = value
	}
	line := fmt.Sprintf("        { value: \"%s\", description: \"%s\" }\n", escapeDoubleQuotes(value), escapeDoubleQuotes(tooltip))
	b.WriteString(line)
}

// writeNushellFlagSignature appends flag signature stubs so Nushell understands which
// switches are available when invoking the external command.
func writeNushellFlagSignature(sc *Subcommand, b *strings.Builder) {
	for _, f := range sc.Flags {
		if f.Hidden {
			continue
		}
		if f.LongName != "" || f.ShortName != "" {
			line := "    "
			if f.LongName != "" {
				line += "--" + f.LongName
			}
			if f.ShortName != "" {
				if f.LongName != "" {
					line += "(-" + f.ShortName + ")"
				} else {
					line += "-" + f.ShortName
				}
			}
			line += "\n"
			b.WriteString(line)
		}
	}
	for _, sub := range sc.Subcommands {
		if sub.Hidden {
			continue
		}
		writeNushellFlagSignature(sub, b)
	}
}
