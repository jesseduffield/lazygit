package plural

import "golang.org/x/text/language"

// Rules is a set of plural rules by language tag.
type Rules map[language.Tag]*Rule

// Rule returns the closest matching plural rule for the language tag
// or nil if no rule could be found.
func (r Rules) Rule(tag language.Tag) *Rule {
	t := tag
	for {
		if rule := r[t]; rule != nil {
			return rule
		}
		t = t.Parent()
		if t.IsRoot() {
			break
		}
	}
	base, _ := tag.Base()
	baseTag, _ := language.Parse(base.String())
	return r[baseTag]
}
