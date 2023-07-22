package models

import "fmt"

// A commit author
type Author struct {
	Name  string
	Email string
}

func (self *Author) Combined() string {
	return fmt.Sprintf("%s <%s>", self.Name, self.Email)
}
