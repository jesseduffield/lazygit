package utils

type StringStack struct {
	stack []string
}

func (self *StringStack) Push(s string) {
	self.stack = append(self.stack, s)
}

func (self *StringStack) Pop() string {
	if len(self.stack) == 0 {
		return ""
	}
	n := len(self.stack) - 1
	last := self.stack[n]
	self.stack = self.stack[:n]
	return last
}

func (self *StringStack) IsEmpty() bool {
	return len(self.stack) == 0
}

func (self *StringStack) Clear() {
	self.stack = []string{}
}
