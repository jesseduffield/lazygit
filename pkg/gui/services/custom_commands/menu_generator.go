package custom_commands

import (
	"bytes"
	"errors"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
)

type MenuGenerator struct {
	c *common.Common
}

// takes the output of a command and returns a list of menu entries based on a filter
// and value/label format templates provided by the user
func NewMenuGenerator(c *common.Common) *MenuGenerator {
	return &MenuGenerator{c: c}
}

type commandMenuEntry struct {
	label string
	value string
}

func (self *MenuGenerator) call(commandOutput, filter, valueFormat, labelFormat string) ([]*commandMenuEntry, error) {
	regex, err := regexp.Compile(filter)
	if err != nil {
		return nil, errors.New("unable to parse filter regex, error: " + err.Error())
	}

	valueTemplateAux, err := template.New("format").Parse(valueFormat)
	if err != nil {
		return nil, errors.New("unable to parse value format, error: " + err.Error())
	}
	valueTemplate := NewTrimmerTemplate(valueTemplateAux)

	var labelTemplate *TrimmerTemplate
	if labelFormat != "" {
		colorFuncMap := style.TemplateFuncMapAddColors(template.FuncMap{})
		labelTemplateAux, err := template.New("format").Funcs(colorFuncMap).Parse(labelFormat)
		if err != nil {
			return nil, errors.New("unable to parse label format, error: " + err.Error())
		}
		labelTemplate = NewTrimmerTemplate(labelTemplateAux)
	} else {
		labelTemplate = valueTemplate
	}

	candidates := []*commandMenuEntry{}
	for _, line := range strings.Split(commandOutput, "\n") {
		if line == "" {
			continue
		}

		candidate, err := self.generateMenuCandidate(
			line,
			regex,
			valueTemplate,
			labelTemplate,
		)
		if err != nil {
			return nil, err
		}

		candidates = append(candidates, candidate)
	}

	return candidates, err
}

func (self *MenuGenerator) generateMenuCandidate(
	line string,
	regex *regexp.Regexp,
	valueTemplate *TrimmerTemplate,
	labelTemplate *TrimmerTemplate,
) (*commandMenuEntry, error) {
	tmplData := self.parseLine(line, regex)

	entry := &commandMenuEntry{}

	var err error
	entry.value, err = valueTemplate.execute(tmplData)
	if err != nil {
		return nil, err
	}

	entry.label, err = labelTemplate.execute(tmplData)
	if err != nil {
		return nil, err
	}

	return entry, nil
}

func (self *MenuGenerator) parseLine(line string, regex *regexp.Regexp) map[string]string {
	tmplData := map[string]string{}
	out := regex.FindAllStringSubmatch(line, -1)
	if len(out) > 0 {
		for groupIdx, group := range regex.SubexpNames() {
			// Record matched group with group ids
			matchName := "group_" + strconv.Itoa(groupIdx)
			tmplData[matchName] = out[0][groupIdx]
			// Record last named group non-empty matches as group matches
			if group != "" {
				tmplData[group] = out[0][groupIdx]
			}
		}
	}

	return tmplData
}

// wrapper around a template which trims the output
type TrimmerTemplate struct {
	template *template.Template
	buffer   *bytes.Buffer
}

func NewTrimmerTemplate(template *template.Template) *TrimmerTemplate {
	return &TrimmerTemplate{
		template: template,
		buffer:   bytes.NewBuffer(nil),
	}
}

func (self *TrimmerTemplate) execute(tmplData map[string]string) (string, error) {
	self.buffer.Reset()
	err := self.template.Execute(self.buffer, tmplData)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(self.buffer.String()), nil
}
