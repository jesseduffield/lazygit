package packet

// This file implements the pushdown automata (PDA) from PGPainless (Paul Schaub)
// to verify pgp packet sequences. See Paul's blogpost for more details:
// https://blog.jabberhead.tk/2022/10/26/implementing-packet-sequence-validation-using-pushdown-automata/
import (
	"fmt"

	"github.com/ProtonMail/go-crypto/openpgp/errors"
)

func NewErrMalformedMessage(from State, input InputSymbol, stackSymbol StackSymbol) errors.ErrMalformedMessage {
	return errors.ErrMalformedMessage(fmt.Sprintf("state %d, input symbol %d, stack symbol %d ", from, input, stackSymbol))
}

// InputSymbol defines the input alphabet of the PDA
type InputSymbol uint8

const (
	LDSymbol InputSymbol = iota
	SigSymbol
	OPSSymbol
	CompSymbol
	ESKSymbol
	EncSymbol
	EOSSymbol
	UnknownSymbol
)

// StackSymbol defines the stack alphabet of the PDA
type StackSymbol int8

const (
	MsgStackSymbol StackSymbol = iota
	OpsStackSymbol
	KeyStackSymbol
	EndStackSymbol
	EmptyStackSymbol
)

// State defines the states of the PDA
type State int8

const (
	OpenPGPMessage State = iota
	ESKMessage
	LiteralMessage
	CompressedMessage
	EncryptedMessage
	ValidMessage
)

// transition represents a state transition in the PDA
type transition func(input InputSymbol, stackSymbol StackSymbol) (State, []StackSymbol, bool, error)

// SequenceVerifier is a pushdown automata to verify
// PGP messages packet sequences according to rfc4880.
type SequenceVerifier struct {
	stack []StackSymbol
	state State
}

// Next performs a state transition with the given input symbol.
// If the transition fails a ErrMalformedMessage is returned.
func (sv *SequenceVerifier) Next(input InputSymbol) error {
	for {
		stackSymbol := sv.popStack()
		transitionFunc := getTransition(sv.state)
		nextState, newStackSymbols, redo, err := transitionFunc(input, stackSymbol)
		if err != nil {
			return err
		}
		if redo {
			sv.pushStack(stackSymbol)
		}
		for _, newStackSymbol := range newStackSymbols {
			sv.pushStack(newStackSymbol)
		}
		sv.state = nextState
		if !redo {
			break
		}
	}
	return nil
}

// Valid returns true if RDA is in a valid state.
func (sv *SequenceVerifier) Valid() bool {
	return sv.state == ValidMessage && len(sv.stack) == 0
}

func (sv *SequenceVerifier) AssertValid() error {
	if !sv.Valid() {
		return errors.ErrMalformedMessage("invalid message")
	}
	return nil
}

func NewSequenceVerifier() *SequenceVerifier {
	return &SequenceVerifier{
		stack: []StackSymbol{EndStackSymbol, MsgStackSymbol},
		state: OpenPGPMessage,
	}
}

func (sv *SequenceVerifier) popStack() StackSymbol {
	if len(sv.stack) == 0 {
		return EmptyStackSymbol
	}
	elemIndex := len(sv.stack) - 1
	stackSymbol := sv.stack[elemIndex]
	sv.stack = sv.stack[:elemIndex]
	return stackSymbol
}

func (sv *SequenceVerifier) pushStack(stackSymbol StackSymbol) {
	sv.stack = append(sv.stack, stackSymbol)
}

func getTransition(from State) transition {
	switch from {
	case OpenPGPMessage:
		return fromOpenPGPMessage
	case LiteralMessage:
		return fromLiteralMessage
	case CompressedMessage:
		return fromCompressedMessage
	case EncryptedMessage:
		return fromEncryptedMessage
	case ESKMessage:
		return fromESKMessage
	case ValidMessage:
		return fromValidMessage
	}
	return nil
}

// fromOpenPGPMessage is the transition for the state OpenPGPMessage.
func fromOpenPGPMessage(input InputSymbol, stackSymbol StackSymbol) (State, []StackSymbol, bool, error) {
	if stackSymbol != MsgStackSymbol {
		return 0, nil, false, NewErrMalformedMessage(OpenPGPMessage, input, stackSymbol)
	}
	switch input {
	case LDSymbol:
		return LiteralMessage, nil, false, nil
	case SigSymbol:
		return OpenPGPMessage, []StackSymbol{MsgStackSymbol}, false, nil
	case OPSSymbol:
		return OpenPGPMessage, []StackSymbol{OpsStackSymbol, MsgStackSymbol}, false, nil
	case CompSymbol:
		return CompressedMessage, nil, false, nil
	case ESKSymbol:
		return ESKMessage, []StackSymbol{KeyStackSymbol}, false, nil
	case EncSymbol:
		return EncryptedMessage, nil, false, nil
	}
	return 0, nil, false, NewErrMalformedMessage(OpenPGPMessage, input, stackSymbol)
}

// fromESKMessage is the transition for the state ESKMessage.
func fromESKMessage(input InputSymbol, stackSymbol StackSymbol) (State, []StackSymbol, bool, error) {
	if stackSymbol != KeyStackSymbol {
		return 0, nil, false, NewErrMalformedMessage(ESKMessage, input, stackSymbol)
	}
	switch input {
	case ESKSymbol:
		return ESKMessage, []StackSymbol{KeyStackSymbol}, false, nil
	case EncSymbol:
		return EncryptedMessage, nil, false, nil
	}
	return 0, nil, false, NewErrMalformedMessage(ESKMessage, input, stackSymbol)
}

// fromLiteralMessage is the transition for the state LiteralMessage.
func fromLiteralMessage(input InputSymbol, stackSymbol StackSymbol) (State, []StackSymbol, bool, error) {
	switch input {
	case SigSymbol:
		if stackSymbol == OpsStackSymbol {
			return LiteralMessage, nil, false, nil
		}
	case EOSSymbol:
		if stackSymbol == EndStackSymbol {
			return ValidMessage, nil, false, nil
		}
	}
	return 0, nil, false, NewErrMalformedMessage(LiteralMessage, input, stackSymbol)
}

// fromLiteralMessage is the transition for the state CompressedMessage.
func fromCompressedMessage(input InputSymbol, stackSymbol StackSymbol) (State, []StackSymbol, bool, error) {
	switch input {
	case SigSymbol:
		if stackSymbol == OpsStackSymbol {
			return CompressedMessage, nil, false, nil
		}
	case EOSSymbol:
		if stackSymbol == EndStackSymbol {
			return ValidMessage, nil, false, nil
		}
	}
	return OpenPGPMessage, []StackSymbol{MsgStackSymbol}, true, nil
}

// fromEncryptedMessage is the transition for the state EncryptedMessage.
func fromEncryptedMessage(input InputSymbol, stackSymbol StackSymbol) (State, []StackSymbol, bool, error) {
	switch input {
	case SigSymbol:
		if stackSymbol == OpsStackSymbol {
			return EncryptedMessage, nil, false, nil
		}
	case EOSSymbol:
		if stackSymbol == EndStackSymbol {
			return ValidMessage, nil, false, nil
		}
	}
	return OpenPGPMessage, []StackSymbol{MsgStackSymbol}, true, nil
}

// fromValidMessage is the transition for the state ValidMessage.
func fromValidMessage(input InputSymbol, stackSymbol StackSymbol) (State, []StackSymbol, bool, error) {
	return 0, nil, false, NewErrMalformedMessage(ValidMessage, input, stackSymbol)
}
