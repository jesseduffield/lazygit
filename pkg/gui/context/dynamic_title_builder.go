package context

import "fmt"

type DynamicTitleBuilder struct {
	formatStr string // e.g. 'remote branches for %s'

	titleRef string // e.g. 'origin'
}

func NewDynamicTitleBuilder(formatStr string) *DynamicTitleBuilder {
	return &DynamicTitleBuilder{
		formatStr: formatStr,
	}
}

func (self *DynamicTitleBuilder) SetTitleRef(titleRef string) {
	self.titleRef = titleRef
}

func (self *DynamicTitleBuilder) Title() string {
	return fmt.Sprintf(self.formatStr, self.titleRef)
}
