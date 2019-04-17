// Copyright © 2010 Fazlul Shahriar <fshahriar@gmail.com> and
// Copyright © 2014 Blake Gentry <blakesgentry@gmail.com>.
// See LICENSE file for license details.

package netrc

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

var expectedMachines = []*Machine{
	&Machine{Name: "mail.google.com", Login: "joe@gmail.com", Password: "somethingSecret", Account: "justagmail"},
	&Machine{Name: "ray", Login: "demo", Password: "mypassword", Account: ""},
	&Machine{Name: "weirdlogin", Login: "uname", Password: "pass#pass", Account: ""},
	&Machine{Name: "", Login: "anonymous", Password: "joe@example.com", Account: ""},
}
var expectedMacros = Macros{
	"allput":  "put src/*",
	"allput2": "  put src/*\nput src2/*",
}

func eqMachine(a *Machine, b *Machine) bool {
	return a.Name == b.Name &&
		a.Login == b.Login &&
		a.Password == b.Password &&
		a.Account == b.Account
}

func testExpected(n *Netrc, t *testing.T) {
	if len(expectedMachines) != len(n.machines) {
		t.Errorf("expected %d machines, got %d", len(expectedMachines), len(n.machines))
	} else {
		for i, e := range expectedMachines {
			if !eqMachine(e, n.machines[i]) {
				t.Errorf("bad machine; expected %v, got %v\n", e, n.machines[i])
			}
		}
	}

	if len(expectedMacros) != len(n.macros) {
		t.Errorf("expected %d macros, got %d", len(expectedMacros), len(n.macros))
	} else {
		for k, v := range expectedMacros {
			if v != n.macros[k] {
				t.Errorf("bad macro for %s; expected %q, got %q\n", k, v, n.macros[k])
			}
		}
	}
}

var newTokenTests = []struct {
	rawkind string
	tkind   tkType
}{
	{"machine", tkMachine},
	{"\n\n\tmachine", tkMachine},
	{"\n   machine", tkMachine},
	{"default", tkDefault},
	{"login", tkLogin},
	{"password", tkPassword},
	{"account", tkAccount},
	{"macdef", tkMacdef},
	{"\n # comment stuff ", tkComment},
	{"\n # I am another comment", tkComment},
	{"\n\t\n ", tkWhitespace},
}

var newTokenInvalidTests = []string{
	" junk",
	"sdfdsf",
	"account#unspaced comment",
}

func TestNewToken(t *testing.T) {
	for _, tktest := range newTokenTests {
		tok, err := newToken([]byte(tktest.rawkind))
		if err != nil {
			t.Fatal(err)
		}
		if tok.kind != tktest.tkind {
			t.Errorf("expected tok.kind %d, got %d", tktest.tkind, tok.kind)
		}
		if string(tok.rawkind) != tktest.rawkind {
			t.Errorf("expected tok.rawkind %q, got %q", tktest.rawkind, string(tok.rawkind))
		}
	}

	for _, tktest := range newTokenInvalidTests {
		_, err := newToken([]byte(tktest))
		if err == nil {
			t.Errorf("expected error with %q, got none", tktest)
		}
	}
}

func TestParse(t *testing.T) {
	r := netrcReader("examples/good.netrc", t)
	n, err := Parse(r)
	if err != nil {
		t.Fatal(err)
	}
	testExpected(n, t)
}

func TestParseFile(t *testing.T) {
	n, err := ParseFile("examples/good.netrc")
	if err != nil {
		t.Fatal(err)
	}
	testExpected(n, t)

	_, err = ParseFile("examples/bad_default_order.netrc")
	if err == nil {
		t.Error("expected an error parsing bad_default_order.netrc, got none")
	} else if !err.(*Error).BadDefaultOrder() {
		t.Error("expected BadDefaultOrder() to be true, got false")
	}

	_, err = ParseFile("examples/this_file_doesnt_exist.netrc")
	if err == nil {
		t.Error("expected an error loading this_file_doesnt_exist.netrc, got none")
	} else if _, ok := err.(*os.PathError); !ok {
		t.Errorf("expected *os.Error, got %v", err)
	}
}

func TestFindMachine(t *testing.T) {
	m, err := FindMachine("examples/good.netrc", "ray")
	if err != nil {
		t.Fatal(err)
	}
	if !eqMachine(m, expectedMachines[1]) {
		t.Errorf("bad machine; expected %v, got %v\n", expectedMachines[1], m)
	}
	if m.IsDefault() {
		t.Errorf("expected m.IsDefault() to be false")
	}

	m, err = FindMachine("examples/good.netrc", "non.existent")
	if err != nil {
		t.Fatal(err)
	}
	if !eqMachine(m, expectedMachines[3]) {
		t.Errorf("bad machine; expected %v, got %v\n", expectedMachines[3], m)
	}
	if !m.IsDefault() {
		t.Errorf("expected m.IsDefault() to be true")
	}
}

func TestNetrcFindMachine(t *testing.T) {
	n, err := ParseFile("examples/good.netrc")
	if err != nil {
		t.Fatal(err)
	}

	m := n.FindMachine("ray")
	if !eqMachine(m, expectedMachines[1]) {
		t.Errorf("bad machine; expected %v, got %v\n", expectedMachines[1], m)
	}
	if m.IsDefault() {
		t.Errorf("expected def to be false")
	}

	n = &Netrc{}
	m = n.FindMachine("nonexistent")
	if m != nil {
		t.Errorf("expected nil, got %v", m)
	}
}

func TestMarshalText(t *testing.T) {
	// load up expected netrc Marshal output
	expected, err := ioutil.ReadAll(netrcReader("examples/good.netrc", t))
	if err != nil {
		t.Fatal(err)
	}

	n, err := ParseFile("examples/good.netrc")
	if err != nil {
		t.Fatal(err)
	}

	result, err := n.MarshalText()
	if err != nil {
		t.Fatal(err)
	}
	if string(result) != string(expected) {
		t.Errorf("expected:\n%q\ngot:\n%q", string(expected), string(result))
	}

	// make sure tokens w/ no value are not serialized
	m := n.FindMachine("mail.google.com")
	m.UpdatePassword("")
	result, err = n.MarshalText()
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(result), "\tpassword \n") {
		fmt.Println(string(result))
		t.Errorf("expected zero-value password token to not be serialzed")
	}
}

var newMachineTests = []struct {
	name     string
	login    string
	password string
	account  string
}{
	{"heroku.com", "dodging-samurai-42@heroku.com", "octocatdodgeballchampions", "2011+2013"},
	{"bgentry.io", "special@test.com", "noacct", ""},
	{"github.io", "2@test.com", "", "acctwithnopass"},
	{"someotherapi.com", "", "passonly", ""},
}

func TestNewMachine(t *testing.T) {
	n, err := ParseFile("examples/good.netrc")
	if err != nil {
		t.Fatal(err)
	}
	testNewMachine(t, n)
	n = &Netrc{}
	testNewMachine(t, n)

	// make sure that tokens without a value are not serialized at all
	for _, test := range newMachineTests {
		n = &Netrc{}
		_ = n.NewMachine(test.name, test.login, test.password, test.account)

		bodyb, _ := n.MarshalText()
		body := string(bodyb)

		// ensure desired values are present when they should be
		if !strings.Contains(body, "machine") {
			t.Errorf("NewMachine() %s missing keyword 'machine'", test.name)
		}
		if !strings.Contains(body, test.name) {
			t.Errorf("NewMachine() %s missing value %q", test.name, test.name)
		}
		if test.login != "" && !strings.Contains(body, "login "+test.login) {
			t.Errorf("NewMachine() %s missing value %q", test.name, "login "+test.login)
		}
		if test.password != "" && !strings.Contains(body, "password "+test.password) {
			t.Errorf("NewMachine() %s missing value %q", test.name, "password "+test.password)
		}
		if test.account != "" && !strings.Contains(body, "account "+test.account) {
			t.Errorf("NewMachine() %s missing value %q", test.name, "account "+test.account)
		}

		// ensure undesired values are not present when they shouldn't be
		if test.login == "" && strings.Contains(body, "login") {
			t.Errorf("NewMachine() %s contains unexpected value %q", test.name, "login")
		}
		if test.password == "" && strings.Contains(body, "password") {
			t.Errorf("NewMachine() %s contains unexpected value %q", test.name, "password")
		}
		if test.account == "" && strings.Contains(body, "account") {
			t.Errorf("NewMachine() %s contains unexpected value %q", test.name, "account")
		}
	}
}

func testNewMachine(t *testing.T, n *Netrc) {
	for _, test := range newMachineTests {
		mcount := len(n.machines)
		// sanity check
		bodyb, _ := n.MarshalText()
		body := string(bodyb)
		for _, value := range []string{test.name, test.login, test.password, test.account} {
			if value != "" && strings.Contains(body, value) {
				t.Errorf("MarshalText() before NewMachine() contained unexpected %q", value)
			}
		}

		// test prefix for machine token
		prefix := "\n"
		if len(n.tokens) == 0 {
			prefix = ""
		}

		m := n.NewMachine(test.name, test.login, test.password, test.account)
		if m == nil {
			t.Fatalf("NewMachine() returned nil")
		}

		if len(n.machines) != mcount+1 {
			t.Errorf("n.machines count expected %d, got %d", mcount+1, len(n.machines))
		}
		// check values
		if m.Name != test.name {
			t.Errorf("m.Name expected %q, got %q", test.name, m.Name)
		}
		if m.Login != test.login {
			t.Errorf("m.Login expected %q, got %q", test.login, m.Login)
		}
		if m.Password != test.password {
			t.Errorf("m.Password expected %q, got %q", test.password, m.Password)
		}
		if m.Account != test.account {
			t.Errorf("m.Account expected %q, got %q", test.account, m.Account)
		}
		// check tokens
		checkToken(t, "nametoken", m.nametoken, tkMachine, prefix+"machine", test.name)
		checkToken(t, "logintoken", m.logintoken, tkLogin, "\n\tlogin", test.login)
		checkToken(t, "passtoken", m.passtoken, tkPassword, "\n\tpassword", test.password)
		checkToken(t, "accounttoken", m.accounttoken, tkAccount, "\n\taccount", test.account)
		// check marshal output
		bodyb, _ = n.MarshalText()
		body = string(bodyb)
		for _, value := range []string{test.name, test.login, test.password, test.account} {
			if !strings.Contains(body, value) {
				t.Errorf("MarshalText() after NewMachine() did not include %q as expected", value)
			}
		}
	}
}

func checkToken(t *testing.T, name string, tok *token, kind tkType, rawkind, value string) {
	if tok == nil {
		t.Errorf("%s not defined", name)
		return
	}
	if tok.kind != kind {
		t.Errorf("%s expected kind %d, got %d", name, kind, tok.kind)
	}
	if string(tok.rawkind) != rawkind {
		t.Errorf("%s expected rawkind %q, got %q", name, rawkind, string(tok.rawkind))
	}
	if tok.value != value {
		t.Errorf("%s expected value %q, got %q", name, value, tok.value)
	}
	if tok.value != value {
		t.Errorf("%s expected value %q, got %q", name, value, tok.value)
	}
}

func TestNewMachineGoesBeforeDefault(t *testing.T) {
	n, err := ParseFile("examples/good.netrc")
	if err != nil {
		t.Fatal(err)
	}
	m := n.NewMachine("mymachine", "mylogin", "mypassword", "myaccount")
	if m2 := n.machines[len(n.machines)-2]; m2 != m {
		t.Errorf("expected machine %v, got %v", m, m2)
	}
}

func TestRemoveMachine(t *testing.T) {
	n, err := ParseFile("examples/good.netrc")
	if err != nil {
		t.Fatal(err)
	}

	tests := []string{"mail.google.com", "weirdlogin"}

	for _, name := range tests {
		mcount := len(n.machines)
		// sanity check
		m := n.FindMachine(name)
		if m == nil {
			t.Fatalf("machine %q not found", name)
		}
		if m.IsDefault() {
			t.Fatalf("expected machine %q, got default instead", name)
		}
		n.RemoveMachine(name)

		if len(n.machines) != mcount-1 {
			t.Errorf("n.machines count expected %d, got %d", mcount-1, len(n.machines))
		}

		// make sure Machine is no longer returned by FindMachine()
		if m2 := n.FindMachine(name); m2 != nil && !m2.IsDefault() {
			t.Errorf("Machine %q not removed from Machines list", name)
		}

		// make sure tokens are not present in tokens list
		for _, token := range []*token{m.nametoken, m.logintoken, m.passtoken, m.accounttoken} {
			if token != nil {
				for _, tok2 := range n.tokens {
					if tok2 == token {
						t.Errorf("token not removed from tokens list: %v", token)
						break
					}
				}
			}
		}

		bodyb, _ := n.MarshalText()
		body := string(bodyb)
		for _, value := range []string{m.Name, m.Login, m.Password, m.Account} {
			if value != "" && strings.Contains(body, value) {
				t.Errorf("MarshalText() after RemoveMachine() contained unexpected %q", value)
			}
		}
	}
}

func TestUpdateLogin(t *testing.T) {
	n, err := ParseFile("examples/good.netrc")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		exists   bool
		name     string
		oldlogin string
		newlogin string
	}{
		{true, "mail.google.com", "joe@gmail.com", "joe2@gmail.com"},
		{false, "heroku.com", "", "dodging-samurai-42@heroku.com"},
	}

	bodyb, _ := n.MarshalText()
	body := string(bodyb)
	for _, test := range tests {
		if strings.Contains(body, test.newlogin) {
			t.Errorf("MarshalText() before UpdateLogin() contained unexpected %q", test.newlogin)
		}
	}

	for _, test := range tests {
		m := n.FindMachine(test.name)
		if m.IsDefault() == test.exists {
			t.Errorf("expected machine %s to not exist, but it did", test.name)
		} else {
			if !test.exists {
				m = n.NewMachine(test.name, test.newlogin, "", "")
			}
			if m == nil {
				t.Errorf("machine %s was nil", test.name)
				continue
			}
			m.UpdateLogin(test.newlogin)
			m := n.FindMachine(test.name)
			if m.Login != test.newlogin {
				t.Errorf("expected new login %q, got %q", test.newlogin, m.Login)
			}
			if m.logintoken.value != test.newlogin {
				t.Errorf("expected m.logintoken %q, got %q", test.newlogin, m.logintoken.value)
			}
		}
	}

	bodyb, _ = n.MarshalText()
	body = string(bodyb)
	for _, test := range tests {
		if test.exists && strings.Contains(body, test.oldlogin) {
			t.Errorf("MarshalText() after UpdateLogin() contained unexpected %q", test.oldlogin)
		}
		if !strings.Contains(body, test.newlogin) {
			t.Errorf("MarshalText after UpdatePassword did not contain %q as expected", test.newlogin)
		}
	}
}

func TestUpdatePassword(t *testing.T) {
	n, err := ParseFile("examples/good.netrc")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		exists      bool
		name        string
		oldpassword string
		newpassword string
	}{
		{true, "ray", "mypassword", "supernewpass"},
		{false, "heroku.com", "", "octocatdodgeballchampions"},
	}

	bodyb, _ := n.MarshalText()
	body := string(bodyb)
	for _, test := range tests {
		if test.exists && !strings.Contains(body, test.oldpassword) {
			t.Errorf("MarshalText() before UpdatePassword() did not include %q as expected", test.oldpassword)
		}
		if strings.Contains(body, test.newpassword) {
			t.Errorf("MarshalText() before UpdatePassword() contained unexpected %q", test.newpassword)
		}
	}

	for _, test := range tests {
		m := n.FindMachine(test.name)
		if m.IsDefault() == test.exists {
			t.Errorf("expected machine %s to not exist, but it did", test.name)
		} else {
			if !test.exists {
				m = n.NewMachine(test.name, "", test.newpassword, "")
			}
			if m == nil {
				t.Errorf("machine %s was nil", test.name)
				continue
			}
			m.UpdatePassword(test.newpassword)
			m = n.FindMachine(test.name)
			if m.Password != test.newpassword {
				t.Errorf("expected new password %q, got %q", test.newpassword, m.Password)
			}
			if m.passtoken.value != test.newpassword {
				t.Errorf("expected m.passtoken %q, got %q", test.newpassword, m.passtoken.value)
			}
		}
	}

	bodyb, _ = n.MarshalText()
	body = string(bodyb)
	for _, test := range tests {
		if test.exists && strings.Contains(body, test.oldpassword) {
			t.Errorf("MarshalText() after UpdatePassword() contained unexpected %q", test.oldpassword)
		}
		if !strings.Contains(body, test.newpassword) {
			t.Errorf("MarshalText() after UpdatePassword() did not contain %q as expected", test.newpassword)
		}
	}
}

func TestNewFile(t *testing.T) {
	var n Netrc

	result, err := n.MarshalText()
	if err != nil {
		t.Fatal(err)
	}
	if string(result) != "" {
		t.Errorf("expected empty result=\"\", got %q", string(result))
	}

	n.NewMachine("netrctest.heroku.com", "auser", "apassword", "")

	result, err = n.MarshalText()
	if err != nil {
		t.Fatal(err)
	}
	expected := `machine netrctest.heroku.com
	login auser
	password apassword`

	if string(result) != expected {
		t.Errorf("expected result:\n%q\ngot:\n%q", expected, string(result))
	}
}

func netrcReader(filename string, t *testing.T) io.Reader {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	return bytes.NewReader(b)
}
