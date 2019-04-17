package translation

import (
	"bytes"
	"fmt"
	//"launchpad.net/goyaml"
	"testing"
	gotemplate "text/template"
)

func TestNilTemplate(t *testing.T) {
	expected := "hello"
	tmpl := &template{
		tmpl: nil,
		src:  expected,
	}
	if actual := tmpl.Execute(nil); actual != expected {
		t.Errorf("Execute(nil) returned %s; expected %s", actual, expected)
	}
}

func TestMarshalText(t *testing.T) {
	tmpl := &template{
		tmpl: gotemplate.Must(gotemplate.New("id").Parse("this is a {{.foo}} template")),
		src:  "boom",
	}
	expectedBuf := []byte(tmpl.src)
	if buf, err := tmpl.MarshalText(); !bytes.Equal(buf, expectedBuf) || err != nil {
		t.Errorf("MarshalText() returned %#v, %#v; expected %#v, nil", buf, err, expectedBuf)
	}
}

func TestUnmarshalText(t *testing.T) {
	tmpl := &template{}
	tmpl.UnmarshalText([]byte("hello {{.World}}"))
	result := tmpl.Execute(map[string]string{
		"World": "world!",
	})
	expected := "hello world!"
	if result != expected {
		t.Errorf("expected %#v; got %#v", expected, result)
	}
}

/*
func TestYAMLMarshal(t *testing.T) {
	src := "hello {{.World}}"
	tmpl, err := newTemplate(src)
	if err != nil {
		t.Fatal(err)
	}
	buf, err := goyaml.Marshal(tmpl)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(buf, []byte(src)) {
		t.Fatalf(`expected "%s"; got "%s"`, src, buf)
	}
}

func TestYAMLUnmarshal(t *testing.T) {
	buf := []byte(`Tmpl: "hello"`)

	var out struct {
		Tmpl *template
	}
	var foo map[string]string
	if err := goyaml.Unmarshal(buf, &foo); err != nil {
		t.Fatal(err)
	}
	if out.Tmpl == nil {
		t.Fatalf("out.Tmpl was nil")
	}
	if out.Tmpl.tmpl == nil {
		t.Fatalf("out.Tmpl.tmpl was nil")
	}
	if expected := "hello {{.World}}"; out.Tmpl.src != expected {
		t.Fatalf("expected %s; got %s", expected, out.Tmpl.src)
	}
}

func TestGetYAML(t *testing.T) {
	src := "hello"
	tmpl := &template{
		tmpl: nil,
		src:  src,
	}
	if tag, value := tmpl.GetYAML(); tag != "" || value != src {
		t.Errorf("GetYAML() returned (%#v, %#v); expected (%#v, %#v)", tag, value, "", src)
	}
}

func TestSetYAML(t *testing.T) {
	tmpl := &template{}
	tmpl.SetYAML("tagDoesntMatter", "hello {{.World}}")
	result := tmpl.Execute(map[string]string{
		"World": "world!",
	})
	expected := "hello world!"
	if result != expected {
		t.Errorf("expected %#v; got %#v", expected, result)
	}
}
*/

func BenchmarkExecuteNilTemplate(b *testing.B) {
	template := &template{src: "hello world"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		template.Execute(nil)
	}
}

func BenchmarkExecuteHelloWorldTemplate(b *testing.B) {
	template, err := newTemplate("hello world")
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		template.Execute(nil)
	}
}

// Executing a simple template like this is ~6x slower than Sprintf
// but it is still only a few microseconds which should be sufficiently fast.
// The benefit is that we have nice semantic tags in the translation.
func BenchmarkExecuteHelloNameTemplate(b *testing.B) {
	template, err := newTemplate("hello {{.Name}}")
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		template.Execute(map[string]string{
			"Name": "Nick",
		})
	}
}

var sprintfResult string

func BenchmarkSprintf(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sprintfResult = fmt.Sprintf("hello %s", "nick")
	}
}
