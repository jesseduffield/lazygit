package chroma

import (
	"fmt"
	"strings"
)

// A Mutator modifies the behaviour of the lexer.
type Mutator interface {
	// Mutate the lexer state machine as it is processing.
	Mutate(state *LexerState) error
}

// A LexerMutator is an additional interface that a Mutator can implement
// to modify the lexer when it is compiled.
type LexerMutator interface {
	// Rules are the lexer rules, state is the state key for the rule the mutator is associated with.
	MutateLexer(rules CompiledRules, state string, rule int) error
}

// A MutatorFunc is a Mutator that mutates the lexer state machine as it is processing.
type MutatorFunc func(state *LexerState) error

func (m MutatorFunc) Mutate(state *LexerState) error { return m(state) } // nolint

// Mutators applies a set of Mutators in order.
func Mutators(modifiers ...Mutator) MutatorFunc {
	return func(state *LexerState) error {
		for _, modifier := range modifiers {
			if err := modifier.Mutate(state); err != nil {
				return err
			}
		}
		return nil
	}
}

type includeMutator struct {
	state string
}

// Include the given state.
func Include(state string) Rule {
	return Rule{Mutator: &includeMutator{state}}
}

func (i *includeMutator) Mutate(s *LexerState) error {
	return fmt.Errorf("should never reach here Include(%q)", i.state)
}

func (i *includeMutator) MutateLexer(rules CompiledRules, state string, rule int) error {
	includedRules, ok := rules[i.state]
	if !ok {
		return fmt.Errorf("invalid include state %q", i.state)
	}
	rules[state] = append(rules[state][:rule], append(includedRules, rules[state][rule+1:]...)...)
	return nil
}

type combinedMutator struct {
	states []string
}

// Combined creates a new anonymous state from the given states, and pushes that state.
func Combined(states ...string) Mutator {
	return &combinedMutator{states}
}

func (c *combinedMutator) Mutate(s *LexerState) error {
	return fmt.Errorf("should never reach here Combined(%v)", c.states)
}

func (c *combinedMutator) MutateLexer(rules CompiledRules, state string, rule int) error {
	name := "__combined_" + strings.Join(c.states, "__")
	if _, ok := rules[name]; !ok {
		combined := []*CompiledRule{}
		for _, state := range c.states {
			rules, ok := rules[state]
			if !ok {
				return fmt.Errorf("invalid combine state %q", state)
			}
			combined = append(combined, rules...)
		}
		rules[name] = combined
	}
	rules[state][rule].Mutator = Push(name)
	return nil
}

// Push states onto the stack.
func Push(states ...string) MutatorFunc {
	return func(s *LexerState) error {
		if len(states) == 0 {
			s.Stack = append(s.Stack, s.State)
		} else {
			for _, state := range states {
				if state == "#pop" {
					s.Stack = s.Stack[:len(s.Stack)-1]
				} else {
					s.Stack = append(s.Stack, state)
				}
			}
		}
		return nil
	}
}

// Pop state from the stack when rule matches.
func Pop(n int) MutatorFunc {
	return func(state *LexerState) error {
		if len(state.Stack) == 0 {
			return fmt.Errorf("nothing to pop")
		}
		state.Stack = state.Stack[:len(state.Stack)-n]
		return nil
	}
}

// Default returns a Rule that applies a set of Mutators.
func Default(mutators ...Mutator) Rule {
	return Rule{Mutator: Mutators(mutators...)}
}

// Stringify returns the raw string for a set of tokens.
func Stringify(tokens ...Token) string {
	out := []string{}
	for _, t := range tokens {
		out = append(out, t.Value)
	}
	return strings.Join(out, "")
}
