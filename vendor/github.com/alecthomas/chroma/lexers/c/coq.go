package c

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Coq lexer.
var Coq = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "Coq",
		Aliases:   []string{"coq"},
		Filenames: []string{"*.v"},
		MimeTypes: []string{"text/x-coq"},
	},
	coqRules,
))

func coqRules() Rules {
	return Rules{
		"root": {
			{`\s+`, Text, nil},
			{`false|true|\(\)|\[\]`, NameBuiltinPseudo, nil},
			{`\(\*`, Comment, Push("comment")},
			{Words(`\b`, `\b`, `Section`, `Module`, `End`, `Require`, `Import`, `Export`, `Variable`, `Variables`, `Parameter`, `Parameters`, `Axiom`, `Hypothesis`, `Hypotheses`, `Notation`, `Local`, `Tactic`, `Reserved`, `Scope`, `Open`, `Close`, `Bind`, `Delimit`, `Definition`, `Let`, `Ltac`, `Fixpoint`, `CoFixpoint`, `Morphism`, `Relation`, `Implicit`, `Arguments`, `Set`, `Unset`, `Contextual`, `Strict`, `Prenex`, `Implicits`, `Inductive`, `CoInductive`, `Record`, `Structure`, `Canonical`, `Coercion`, `Theorem`, `Lemma`, `Corollary`, `Proposition`, `Fact`, `Remark`, `Example`, `Proof`, `Goal`, `Save`, `Qed`, `Defined`, `Hint`, `Resolve`, `Rewrite`, `View`, `Search`, `Show`, `Print`, `Printing`, `All`, `Graph`, `Projections`, `inside`, `outside`, `Check`, `Global`, `Instance`, `Class`, `Existing`, `Universe`, `Polymorphic`, `Monomorphic`, `Context`), KeywordNamespace, nil},
			{Words(`\b`, `\b`, `forall`, `exists`, `exists2`, `fun`, `fix`, `cofix`, `struct`, `match`, `end`, `in`, `return`, `let`, `if`, `is`, `then`, `else`, `for`, `of`, `nosimpl`, `with`, `as`), Keyword, nil},
			{Words(`\b`, `\b`, `Type`, `Prop`), KeywordType, nil},
			{Words(`\b`, `\b`, `pose`, `set`, `move`, `case`, `elim`, `apply`, `clear`, `hnf`, `intro`, `intros`, `generalize`, `rename`, `pattern`, `after`, `destruct`, `induction`, `using`, `refine`, `inversion`, `injection`, `rewrite`, `congr`, `unlock`, `compute`, `ring`, `field`, `replace`, `fold`, `unfold`, `change`, `cutrewrite`, `simpl`, `have`, `suff`, `wlog`, `suffices`, `without`, `loss`, `nat_norm`, `assert`, `cut`, `trivial`, `revert`, `bool_congr`, `nat_congr`, `symmetry`, `transitivity`, `auto`, `split`, `left`, `right`, `autorewrite`, `tauto`, `setoid_rewrite`, `intuition`, `eauto`, `eapply`, `econstructor`, `etransitivity`, `constructor`, `erewrite`, `red`, `cbv`, `lazy`, `vm_compute`, `native_compute`, `subst`), Keyword, nil},
			{Words(`\b`, `\b`, `by`, `done`, `exact`, `reflexivity`, `tauto`, `romega`, `omega`, `assumption`, `solve`, `contradiction`, `discriminate`, `congruence`), KeywordPseudo, nil},
			{Words(`\b`, `\b`, `do`, `last`, `first`, `try`, `idtac`, `repeat`), KeywordReserved, nil},
			{`\b([A-Z][\w\']*)`, Name, nil},
			{"(\u03bb|\u03a0|\\|\\}|\\{\\||\\\\/|/\\\\|=>|~|\\}|\\|]|\\||\\{<|\\{|`|_|]|\\[\\||\\[>|\\[<|\\[|\\?\\?|\\?|>\\}|>]|>|=|<->|<-|<|;;|;|:>|:=|::|:|\\.\\.|\\.|->|-\\.|-|,|\\+|\\*|\\)|\\(|&&|&|#|!=)", Operator, nil},
			{`([=<>@^|&+\*/$%-]|[!?~])?[!$%&*+\./:<=>?@^|~-]`, Operator, nil},
			{`\b(unit|nat|bool|string|ascii|list)\b`, KeywordType, nil},
			{`[^\W\d][\w']*`, Name, nil},
			{`\d[\d_]*`, LiteralNumberInteger, nil},
			{`0[xX][\da-fA-F][\da-fA-F_]*`, LiteralNumberHex, nil},
			{`0[oO][0-7][0-7_]*`, LiteralNumberOct, nil},
			{`0[bB][01][01_]*`, LiteralNumberBin, nil},
			{`-?\d[\d_]*(.[\d_]*)?([eE][+\-]?\d[\d_]*)`, LiteralNumberFloat, nil},
			{`'(?:(\\[\\\"'ntbr ])|(\\[0-9]{3})|(\\x[0-9a-fA-F]{2}))'`, LiteralStringChar, nil},
			{`'.'`, LiteralStringChar, nil},
			{`'`, Keyword, nil},
			{`"`, LiteralStringDouble, Push("string")},
			{`[~?][a-z][\w\']*:`, Name, nil},
		},
		"comment": {
			{`[^(*)]+`, Comment, nil},
			{`\(\*`, Comment, Push()},
			{`\*\)`, Comment, Pop(1)},
			{`[(*)]`, Comment, nil},
		},
		"string": {
			{`[^"]+`, LiteralStringDouble, nil},
			{`""`, LiteralStringDouble, nil},
			{`"`, LiteralStringDouble, Pop(1)},
		},
		"dotted": {
			{`\s+`, Text, nil},
			{`\.`, Punctuation, nil},
			{`[A-Z][\w\']*(?=\s*\.)`, NameNamespace, nil},
			{`[A-Z][\w\']*`, NameClass, Pop(1)},
			{`[a-z][a-z0-9_\']*`, Name, Pop(1)},
			Default(Pop(1)),
		},
	}
}
