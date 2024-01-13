package css

import "fmt"

// Declaration represents a parsed style property
type Declaration struct {
	Property  string
	Value     string
	Important bool
}

// NewDeclaration instanciates a new Declaration
func NewDeclaration() *Declaration {
	return &Declaration{}
}

// Returns string representation of the Declaration
func (decl *Declaration) String() string {
	return decl.StringWithImportant(true)
}

// StringWithImportant returns string representation with optional !important part
func (decl *Declaration) StringWithImportant(option bool) string {
	result := fmt.Sprintf("%s: %s", decl.Property, decl.Value)

	if option && decl.Important {
		result += " !important"
	}

	result += ";"

	return result
}

// Equal returns true if both Declarations are equals
func (decl *Declaration) Equal(other *Declaration) bool {
	return (decl.Property == other.Property) && (decl.Value == other.Value) && (decl.Important == other.Important)
}

//
// DeclarationsByProperty
//

// DeclarationsByProperty represents sortable style declarations
type DeclarationsByProperty []*Declaration

// Implements sort.Interface
func (declarations DeclarationsByProperty) Len() int {
	return len(declarations)
}

// Implements sort.Interface
func (declarations DeclarationsByProperty) Swap(i, j int) {
	declarations[i], declarations[j] = declarations[j], declarations[i]
}

// Implements sort.Interface
func (declarations DeclarationsByProperty) Less(i, j int) bool {
	return declarations[i].Property < declarations[j].Property
}
