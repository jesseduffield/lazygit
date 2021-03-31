package gui

import "github.com/go-errors/errors"

// SentinelErrors are the errors that have special meaning and need to be checked
// by calling functions. The less of these, the better
type SentinelErrors struct {
	ErrSubProcess error
	ErrNoFiles    error
	ErrSwitchRepo error
	ErrRestart    error
}

const UNKNOWN_VIEW_ERROR_MSG = "unknown view"

// GenerateSentinelErrors makes the sentinel errors for the gui. We're defining it here
// because we can't do package-scoped errors with localization, and also because
// it seems like package-scoped variables are bad in general
// https://dave.cheney.net/2017/06/11/go-without-package-scoped-variables
// In the future it would be good to implement some of the recommendations of
// that article. For now, if we don't need an error to be a sentinel, we will just
// define it inline. This has implications for error messages that pop up everywhere
// in that we'll be duplicating the default values. We may need to look at
// having a default localisation bundle defined, and just using keys-only when
// localising things in the code.
func (gui *Gui) GenerateSentinelErrors() {
	gui.Errors = SentinelErrors{
		ErrSubProcess: errors.New(gui.Tr.RunningSubprocess),
		ErrNoFiles:    errors.New(gui.Tr.NoChangedFiles),
		ErrSwitchRepo: errors.New("switching repo"),
		ErrRestart:    errors.New("restarting"),
	}
}

func (gui *Gui) sentinelErrorsArr() []error {
	return []error{
		gui.Errors.ErrSubProcess,
		gui.Errors.ErrNoFiles,
		gui.Errors.ErrSwitchRepo,
		gui.Errors.ErrRestart,
	}
}
