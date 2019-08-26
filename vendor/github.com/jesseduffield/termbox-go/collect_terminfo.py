#!/usr/bin/env python

import sys, os, subprocess

def escaped(s):
	return repr(s)[1:-1]

def tput(term, name):
	try:
		return subprocess.check_output(['tput', '-T%s' % term, name]).decode()
	except subprocess.CalledProcessError as e:
		return e.output.decode()


def w(s):
	if s == None:
		return
	sys.stdout.write(s)

terminals = {
	'xterm' : 'xterm',
	'rxvt-256color' : 'rxvt_256color',
	'rxvt-unicode' : 'rxvt_unicode',
	'linux' : 'linux',
	'Eterm' : 'eterm',
	'screen' : 'screen'
}

keys = [
	"F1",		"kf1",
	"F2",		"kf2",
	"F3",		"kf3",
	"F4",		"kf4",
	"F5",		"kf5",
	"F6",		"kf6",
	"F7",		"kf7",
	"F8",		"kf8",
	"F9",		"kf9",
	"F10",		"kf10",
	"F11",		"kf11",
	"F12",		"kf12",
	"INSERT",	"kich1",
	"DELETE",	"kdch1",
	"HOME",		"khome",
	"END",		"kend",
	"PGUP",		"kpp",
	"PGDN",		"knp",
	"KEY_UP",	"kcuu1",
	"KEY_DOWN",	"kcud1",
	"KEY_LEFT",	"kcub1",
	"KEY_RIGHT",	"kcuf1"
]

funcs = [
	"T_ENTER_CA",		"smcup",
	"T_EXIT_CA",		"rmcup",
	"T_SHOW_CURSOR",	"cnorm",
	"T_HIDE_CURSOR",	"civis",
	"T_CLEAR_SCREEN",	"clear",
	"T_SGR0",		"sgr0",
	"T_UNDERLINE",		"smul",
	"T_BOLD",		"bold",
	"T_BLINK",		"blink",
	"T_REVERSE",            "rev",
	"T_ENTER_KEYPAD",	"smkx",
	"T_EXIT_KEYPAD",	"rmkx"
]

def iter_pairs(iterable):
	iterable = iter(iterable)
	while True:
		yield (next(iterable), next(iterable))

def do_term(term, nick):
	w("// %s\n" % term)
	w("var %s_keys = []string{\n\t" % nick)
	for k, v in iter_pairs(keys):
		w('"')
		w(escaped(tput(term, v)))
		w('",')
	w("\n}\n")
	w("var %s_funcs = []string{\n\t" % nick)
	for k,v in iter_pairs(funcs):
		w('"')
		if v == "sgr":
			w("\\033[3%d;4%dm")
		elif v == "cup":
			w("\\033[%d;%dH")
		else:
			w(escaped(tput(term, v)))
		w('", ')
	w("\n}\n\n")

def do_terms(d):
	w("var terms = []struct {\n")
	w("\tname  string\n")
	w("\tkeys  []string\n")
	w("\tfuncs []string\n")
	w("}{\n")
	for k, v in d.items():
		w('\t{"%s", %s_keys, %s_funcs},\n' % (k, v, v))
	w("}\n\n")

w("// +build !windows\n\npackage termbox\n\n")

for k,v in terminals.items():
	do_term(k, v)

do_terms(terminals)

