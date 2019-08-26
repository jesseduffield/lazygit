[![Build Status](https://travis-ci.org/heroku/rollrus.svg?branch=master)](https://travis-ci.org/heroku/rollrus)&nbsp;[![GoDoc](https://godoc.org/github.com/heroku/rollrus?status.svg)](https://godoc.org/github.com/heroku/rollrus)

# What

Rollrus is what happens when [Logrus](https://github.com/sirupsen/logrus) meets [Roll](https://github.com/stvp/roll).

When a .Error, .Fatal or .Panic logging function is called, report the details to rollbar via a Logrus hook.

Delivery is synchronous to help ensure that logs are delivered.

If the error includes a [`StackTrace`](https://godoc.org/github.com/pkg/errors#StackTrace), that `StackTrace` is reported to rollbar.

# Usage

Examples available in the [tests](https://github.com/heroku/rollrus/blob/master/rollrus_test.go) or on [GoDoc](https://godoc.org/github.com/heroku/rollrus).
