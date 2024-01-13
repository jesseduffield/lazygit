package r

import (
	"regexp"
	"strings"
	"unicode/utf8"

	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
	"github.com/dlclark/regexp2"
)

// Raku lexer.
var Raku Lexer = internal.Register(MustNewLazyLexer(
	&Config{
		Name:    "Raku",
		Aliases: []string{"perl6", "pl6", "raku"},
		Filenames: []string{
			"*.pl", "*.pm", "*.nqp", "*.p6", "*.6pl", "*.p6l", "*.pl6", "*.6pm",
			"*.p6m", "*.pm6", "*.t", "*.raku", "*.rakumod", "*.rakutest", "*.rakudoc",
		},
		MimeTypes: []string{
			"text/x-perl6", "application/x-perl6",
			"text/x-raku", "application/x-raku",
		},
		DotAll: true,
	},
	rakuRules,
))

func rakuRules() Rules {
	type RakuToken int

	const (
		rakuQuote RakuToken = iota
		rakuNameAttribute
		rakuPod
		rakuPodFormatter
		rakuPodDeclaration
		rakuMultilineComment
		rakuMatchRegex
		rakuSubstitutionRegex
	)

	const (
		colonPairOpeningBrackets = `(?:<<|<|«|\(|\[|\{)`
		colonPairClosingBrackets = `(?:>>|>|»|\)|\]|\})`
		colonPairPattern         = `(?<!:)(?<colon>:)(?<key>\w[\w'-]*)(?<opening_delimiters>` + colonPairOpeningBrackets + `)`
		colonPairLookahead       = `(?=(:['\w-]+` +
			colonPairOpeningBrackets + `.+?` + colonPairClosingBrackets + `)?`
		namePattern           = `(?:(?!` + colonPairPattern + `)(?:::|[\w':-]))+`
		variablePattern       = `[$@%&]+[.^:?=!~]?` + namePattern
		globalVariablePattern = `[$@%&]+\*` + namePattern
	)

	keywords := []string{
		`BEGIN`, `CATCH`, `CHECK`, `CLOSE`, `CONTROL`, `DOC`, `END`, `ENTER`, `FIRST`, `INIT`,
		`KEEP`, `LAST`, `LEAVE`, `NEXT`, `POST`, `PRE`, `QUIT`, `UNDO`, `anon`, `augment`, `but`,
		`class`, `constant`, `default`, `does`, `else`, `elsif`, `enum`, `for`, `gather`, `given`,
		`grammar`, `has`, `if`, `import`, `is`, `of`, `let`, `loop`, `made`, `make`, `method`,
		`module`, `multi`, `my`, `need`, `orwith`, `our`, `proceed`, `proto`, `repeat`, `require`,
		`where`, `return`, `return-rw`, `returns`, `->`, `-->`, `role`, `state`, `sub`, `no`,
		`submethod`, `subset`, `succeed`, `supersede`, `try`, `unit`, `unless`, `until`,
		`use`, `when`, `while`, `with`, `without`, `export`, `native`, `repr`, `required`, `rw`,
		`symbol`, `default`, `cached`, `DEPRECATED`, `dynamic`, `hidden-from-backtrace`, `nodal`,
		`pure`, `raw`, `start`, `react`, `supply`, `whenever`, `also`, `rule`, `token`, `regex`,
		`dynamic-scope`, `built`, `temp`,
	}

	keywordsPattern := Words(`(?<!['\w:-])`, `(?!['\w:-])`, keywords...)

	wordOperators := []string{
		`X`, `Z`, `R`, `after`, `and`, `andthen`, `before`, `cmp`, `div`, `eq`, `eqv`, `extra`, `ge`,
		`gt`, `le`, `leg`, `lt`, `mod`, `ne`, `or`, `orelse`, `x`, `xor`, `xx`, `gcd`, `lcm`,
		`but`, `min`, `max`, `^fff`, `fff^`, `fff`, `^ff`, `ff^`, `ff`, `so`, `not`, `unicmp`,
		`TR`, `o`, `(&)`, `(.)`, `(|)`, `(+)`, `(-)`, `(^)`, `coll`, `(elem)`, `(==)`,
		`(cont)`, `(<)`, `(<=)`, `(>)`, `(>=)`, `minmax`, `notandthen`, `S`,
	}

	wordOperatorsPattern := Words(`(?<=^|\b|\s)`, `(?=$|\b|\s)`, wordOperators...)

	operators := []string{
		`++`, `--`, `-`, `**`, `!`, `+`, `~`, `?`, `+^`, `~^`, `?^`, `^`, `*`, `/`, `%`, `%%`, `+&`,
		`+<`, `+>`, `~&`, `~<`, `~>`, `?&`, `+|`, `+^`, `~|`, `~^`, `?`, `?|`, `?^`, `&`, `^`,
		`<=>`, `^…^`, `^…`, `…^`, `…`, `...`, `...^`, `^...`, `^...^`, `..`, `..^`, `^..`, `^..^`,
		`::=`, `:=`, `!=`, `==`, `<=`, `<`, `>=`, `>`, `~~`, `===`, `&&`, `||`, `|`, `^^`, `//`,
		`??`, `!!`, `^fff^`, `^ff^`, `<==`, `==>`, `<<==`, `==>>`, `=>`, `=`, `<<`, `«`, `>>`, `»`,
		`,`, `>>.`, `».`, `.&`, `.=`, `.^`, `.?`, `.+`, `.*`, `.`, `∘`, `∩`, `⊍`, `∪`, `⊎`, `∖`,
		`⊖`, `≠`, `≤`, `≥`, `=:=`, `=~=`, `≅`, `∈`, `∉`, `≡`, `≢`, `∋`, `∌`, `⊂`, `⊄`, `⊆`, `⊈`,
		`⊃`, `⊅`, `⊇`, `⊉`, `:`, `!!!`, `???`, `¯`, `×`, `÷`, `−`, `⁺`, `⁻`,
	}

	operatorsPattern := Words(``, ``, operators...)

	builtinTypes := []string{
		`False`, `True`, `Order`, `More`, `Less`, `Same`, `Any`, `Array`, `Associative`, `AST`,
		`atomicint`, `Attribute`, `Backtrace`, `Backtrace::Frame`, `Bag`, `Baggy`, `BagHash`,
		`Blob`, `Block`, `Bool`, `Buf`, `Callable`, `CallFrame`, `Cancellation`, `Capture`,
		`CArray`, `Channel`, `Code`, `compiler`, `Complex`, `ComplexStr`, `CompUnit`,
		`CompUnit::PrecompilationRepository`, `CompUnit::Repository`, `Empty`,
		`CompUnit::Repository::FileSystem`, `CompUnit::Repository::Installation`, `Cool`,
		`CurrentThreadScheduler`, `CX::Warn`, `CX::Take`, `CX::Succeed`, `CX::Return`, `CX::Redo`,
		`CX::Proceed`, `CX::Next`, `CX::Last`, `CX::Emit`, `CX::Done`, `Cursor`, `Date`, `Dateish`,
		`DateTime`, `Distribution`, `Distribution::Hash`, `Distribution::Locally`,
		`Distribution::Path`, `Distribution::Resource`, `Distro`, `Duration`, `Encoding`,
		`Encoding::Registry`, `Endian`, `Enumeration`, `Exception`, `Failure`, `FatRat`, `Grammar`,
		`Hash`, `HyperWhatever`, `Instant`, `Int`, `int`, `int16`, `int32`, `int64`, `int8`, `str`,
		`IntStr`, `IO`, `IO::ArgFiles`, `IO::CatHandle`, `IO::Handle`, `IO::Notification`,
		`IO::Notification::Change`, `IO::Path`, `IO::Path::Cygwin`, `IO::Path::Parts`,
		`IO::Path::QNX`, `IO::Path::Unix`, `IO::Path::Win32`, `IO::Pipe`, `IO::Socket`,
		`IO::Socket::Async`, `IO::Socket::Async::ListenSocket`, `IO::Socket::INET`, `IO::Spec`,
		`IO::Spec::Cygwin`, `IO::Spec::QNX`, `IO::Spec::Unix`, `IO::Spec::Win32`, `IO::Special`,
		`Iterable`, `Iterator`, `Junction`, `Kernel`, `Label`, `List`, `Lock`, `Lock::Async`,
		`Lock::ConditionVariable`, `long`, `longlong`, `Macro`, `Map`, `Match`,
		`Metamodel::AttributeContainer`, `Metamodel::C3MRO`, `Metamodel::ClassHOW`,
		`Metamodel::ConcreteRoleHOW`, `Metamodel::CurriedRoleHOW`, `Metamodel::DefiniteHOW`,
		`Metamodel::Documenting`, `Metamodel::EnumHOW`, `Metamodel::Finalization`,
		`Metamodel::MethodContainer`, `Metamodel::Mixins`, `Metamodel::MROBasedMethodDispatch`,
		`Metamodel::MultipleInheritance`, `Metamodel::Naming`, `Metamodel::Primitives`,
		`Metamodel::PrivateMethodContainer`, `Metamodel::RoleContainer`, `Metamodel::RolePunning`,
		`Metamodel::Stashing`, `Metamodel::Trusting`, `Metamodel::Versioning`, `Method`, `Mix`,
		`MixHash`, `Mixy`, `Mu`, `NFC`, `NFD`, `NFKC`, `NFKD`, `Nil`, `Num`, `num32`, `num64`,
		`Numeric`, `NumStr`, `ObjAt`, `Order`, `Pair`, `Parameter`, `Perl`, `Pod::Block`,
		`Pod::Block::Code`, `Pod::Block::Comment`, `Pod::Block::Declarator`, `Pod::Block::Named`,
		`Pod::Block::Para`, `Pod::Block::Table`, `Pod::Heading`, `Pod::Item`, `Pointer`,
		`Positional`, `PositionalBindFailover`, `Proc`, `Proc::Async`, `Promise`, `Proxy`,
		`PseudoStash`, `QuantHash`, `RaceSeq`, `Raku`, `Range`, `Rat`, `Rational`, `RatStr`,
		`Real`, `Regex`, `Routine`, `Routine::WrapHandle`, `Scalar`, `Scheduler`, `Semaphore`,
		`Seq`, `Sequence`, `Set`, `SetHash`, `Setty`, `Signature`, `size_t`, `Slip`, `Stash`,
		`Str`, `StrDistance`, `Stringy`, `Sub`, `Submethod`, `Supplier`, `Supplier::Preserving`,
		`Supply`, `Systemic`, `Tap`, `Telemetry`, `Telemetry::Instrument::Thread`,
		`Telemetry::Instrument::ThreadPool`, `Telemetry::Instrument::Usage`, `Telemetry::Period`,
		`Telemetry::Sampler`, `Thread`, `Test`, `ThreadPoolScheduler`, `UInt`, `uint16`, `uint32`,
		`uint64`, `uint8`, `Uni`, `utf8`, `ValueObjAt`, `Variable`, `Version`, `VM`, `Whatever`,
		`WhateverCode`, `WrapHandle`, `NativeCall`,
		// Pragmas
		`precompilation`, `experimental`, `worries`, `MONKEY-TYPING`, `MONKEY-SEE-NO-EVAL`,
		`MONKEY-GUTS`, `fatal`, `lib`, `isms`, `newline`, `nqp`, `soft`,
		`strict`, `trace`, `variables`,
	}

	builtinTypesPattern := Words(`(?<!['\w:-])`, `(?::[_UD])?(?!['\w:-])`, builtinTypes...)

	builtinRoutines := []string{
		`ACCEPTS`, `abs`, `abs2rel`, `absolute`, `accept`, `accepts_type`, `accessed`, `acos`,
		`acosec`, `acosech`, `acosh`, `acotan`, `acotanh`, `acquire`, `act`, `action`, `actions`,
		`add`, `add_attribute`, `add_enum_value`, `add_fallback`, `add_method`, `add_parent`,
		`add_private_method`, `add_role`, `add_stash`, `add_trustee`, `addendum`, `adverb`, `after`,
		`all`, `allocate`, `allof`, `allowed`, `alternative-names`, `annotations`, `antipair`,
		`antipairs`, `any`, `anyof`, `api`, `app_lifetime`, `append`, `arch`, `archetypes`,
		`archname`, `args`, `ARGS-TO-CAPTURE`, `arity`, `Array`, `asec`, `asech`, `asin`, `asinh`,
		`ASSIGN-KEY`, `ASSIGN-POS`, `assuming`, `ast`, `at`, `atan`, `atan2`, `atanh`, `AT-KEY`,
		`atomic-assign`, `atomic-dec-fetch`, `atomic-fetch`, `atomic-fetch-add`, `atomic-fetch-dec`,
		`atomic-fetch-inc`, `atomic-fetch-sub`, `atomic-inc-fetch`, `AT-POS`, `attributes`, `auth`,
		`await`, `backend`, `backtrace`, `Bag`, `bag`, `Baggy`, `BagHash`, `bail-out`, `base`,
		`basename`, `base-repeating`, `base_type`, `batch`, `BIND-KEY`, `BIND-POS`, `bind-stderr`,
		`bind-stdin`, `bind-stdout`, `bind-udp`, `bits`, `bless`, `block`, `Bool`, `bool-only`,
		`bounds`, `break`, `Bridge`, `broken`, `BUILD`, `TWEAK`, `build-date`, `bytes`, `cache`,
		`callframe`, `calling-package`, `CALL-ME`, `callsame`, `callwith`, `can`, `cancel`,
		`candidates`, `cando`, `can-ok`, `canonpath`, `caps`, `caption`, `Capture`, `capture`,
		`cas`, `catdir`, `categorize`, `categorize-list`, `catfile`, `catpath`, `cause`, `ceiling`,
		`cglobal`, `changed`, `Channel`, `channel`, `chars`, `chdir`, `child`, `child-name`,
		`child-typename`, `chmod`, `chomp`, `chop`, `chr`, `chrs`, `chunks`, `cis`, `classify`,
		`classify-list`, `cleanup`, `clone`, `close`, `closed`, `close-stdin`, `cmp-ok`, `code`,
		`codename`, `codes`, `coerce_type`, `coll`, `collate`, `column`, `comb`, `combinations`,
		`command`, `comment`, `compiler`, `Complex`, `compose`, `composalizer`, `compose_type`,
		`compose_values`, `composer`, `compute_mro`, `condition`, `config`, `configure_destroy`,
		`configure_type_checking`, `conj`, `connect`, `constraints`, `construct`, `contains`,
		`content`, `contents`, `copy`, `cos`, `cosec`, `cosech`, `cosh`, `cotan`, `cotanh`, `count`,
		`count-only`, `cpu-cores`, `cpu-usage`, `CREATE`, `create_type`, `cross`, `cue`, `curdir`,
		`curupdir`, `d`, `Date`, `DateTime`, `day`, `daycount`, `day-of-month`, `day-of-week`,
		`day-of-year`, `days-in-month`, `dd-mm-yyyy`, `declaration`, `decode`, `decoder`, `deepmap`,
		`default`, `defined`, `DEFINITE`, `definite`, `delayed`, `delete`, `delete-by-compiler`,
		`DELETE-KEY`, `DELETE-POS`, `denominator`, `desc`, `DESTROY`, `destroyers`, `devnull`,
		`diag`, `did-you-mean`, `die`, `dies-ok`, `dir`, `dirname`, `distribution`, `dir-sep`,
		`DISTROnames`, `do`, `does`, `does-ok`, `done`, `done-testing`, `duckmap`, `dynamic`, `e`,
		`eager`, `earlier`, `elems`, `emit`, `enclosing`, `encode`, `encoder`, `encoding`, `end`,
		`endian`, `ends-with`, `enum_from_value`, `enum_value_list`, `enum_values`, `enums`, `EOF`,
		`eof`, `EVAL`, `eval-dies-ok`, `EVALFILE`, `eval-lives-ok`, `event`, `exception`,
		`excludes-max`, `excludes-min`, `EXISTS-KEY`, `EXISTS-POS`, `exit`, `exitcode`, `exp`,
		`expected`, `explicitly-manage`, `expmod`, `export_callback`, `extension`, `f`, `fail`,
		`FALLBACK`, `fails-like`, `fc`, `feature`, `file`, `filename`, `files`, `find`,
		`find_method`, `find_method_qualified`, `finish`, `first`, `flat`, `first-date-in-month`,
		`flatmap`, `flip`, `floor`, `flunk`, `flush`, `flush_cache`, `fmt`, `format`, `formatter`,
		`free-memory`, `freeze`, `from`, `from-list`, `from-loop`, `from-posix`, `from-slurpy`,
		`full`, `full-barrier`, `GENERATE-USAGE`, `generate_mixin`, `get`, `get_value`, `getc`,
		`gist`, `got`, `grab`, `grabpairs`, `grep`, `handle`, `handled`, `handles`, `hardware`,
		`has_accessor`, `Hash`, `hash`, `head`, `headers`, `hh-mm-ss`, `hidden`, `hides`, `hostname`,
		`hour`, `how`, `hyper`, `id`, `illegal`, `im`, `in`, `in-timezone`, `indent`, `index`,
		`indices`, `indir`, `infinite`, `infix`, `postcirumfix`, `cicumfix`, `install`,
		`install_method_cache`, `Instant`, `instead`, `Int`, `int-bounds`, `interval`, `in-timezone`,
		`invalid-str`, `invert`, `invocant`, `IO`, `IO::Notification.watch-path`, `is_trusted`,
		`is_type`, `isa`, `is-absolute`, `isa-ok`, `is-approx`, `is-deeply`, `is-hidden`,
		`is-initial-thread`, `is-int`, `is-lazy`, `is-leap-year`, `isNaN`, `isnt`, `is-prime`,
		`is-relative`, `is-routine`, `is-setting`, `is-win`, `item`, `iterator`, `join`, `keep`,
		`kept`, `KERNELnames`, `key`, `keyof`, `keys`, `kill`, `kv`, `kxxv`, `l`, `lang`, `last`,
		`lastcall`, `later`, `lazy`, `lc`, `leading`, `level`, `like`, `line`, `lines`, `link`,
		`List`, `list`, `listen`, `live`, `lives-ok`, `load`, `load-repo-id`, `load-unit`, `loaded`,
		`loads`, `local`, `lock`, `log`, `log10`, `lookup`, `lsb`, `made`, `MAIN`, `make`, `Map`,
		`map`, `match`, `max`, `maxpairs`, `merge`, `message`, `method`, `meta`, `method_table`,
		`methods`, `migrate`, `min`, `minmax`, `minpairs`, `minute`, `misplaced`, `Mix`, `mix`,
		`MixHash`, `mixin`, `mixin_attribute`, `Mixy`, `mkdir`, `mode`, `modified`, `month`, `move`,
		`mro`, `msb`, `multi`, `multiness`, `name`, `named`, `named_names`, `narrow`,
		`nativecast`, `native-descriptor`, `nativesizeof`, `need`, `new`, `new_type`,
		`new-from-daycount`, `new-from-pairs`, `next`, `nextcallee`, `next-handle`, `nextsame`,
		`nextwith`, `next-interesting-index`, `NFC`, `NFD`, `NFKC`, `NFKD`, `nice`, `nl-in`,
		`nl-out`, `nodemap`, `nok`, `normalize`, `none`, `norm`, `not`, `note`, `now`, `nude`,
		`Num`, `numerator`, `Numeric`, `of`, `offset`, `offset-in-hours`, `offset-in-minutes`,
		`ok`, `old`, `on-close`, `one`, `on-switch`, `open`, `opened`, `operation`, `optional`,
		`ord`, `ords`, `orig`, `os-error`, `osname`, `out-buffer`, `pack`, `package`, `package-kind`,
		`package-name`, `packages`, `Pair`, `pair`, `pairs`, `pairup`, `parameter`, `params`,
		`parent`, `parent-name`, `parents`, `parse`, `parse-base`, `parsefile`, `parse-names`,
		`parts`, `pass`, `path`, `path-sep`, `payload`, `peer-host`, `peer-port`, `periods`, `perl`,
		`permutations`, `phaser`, `pick`, `pickpairs`, `pid`, `placeholder`, `plan`, `plus`,
		`polar`, `poll`, `polymod`, `pop`, `pos`, `positional`, `posix`, `postfix`, `postmatch`,
		`precomp-ext`, `precomp-target`, `precompiled`, `pred`, `prefix`, `prematch`, `prepend`,
		`primary`, `print`, `printf`, `print-nl`, `print-to`, `private`, `private_method_names`,
		`private_method_table`, `proc`, `produce`, `Promise`, `promise`, `prompt`, `protect`,
		`protect-or-queue-on-recursion`, `publish_method_cache`, `pull-one`, `push`, `push-all`,
		`push-at-least`, `push-exactly`, `push-until-lazy`, `put`, `qualifier-type`, `quaternary`,
		`quit`, `r`, `race`, `radix`, `raku`, `rand`, `Range`, `range`, `Rat`, `raw`, `re`, `read`,
		`read-bits`, `read-int128`, `read-int16`, `read-int32`, `read-int64`, `read-int8`,
		`read-num32`, `read-num64`, `read-ubits`, `read-uint128`, `read-uint16`, `read-uint32`,
		`read-uint64`, `read-uint8`, `readchars`, `readonly`, `ready`, `Real`, `reallocate`,
		`reals`, `reason`, `rebless`, `receive`, `recv`, `redispatcher`, `redo`, `reduce`,
		`rel2abs`, `relative`, `release`, `remove`, `rename`, `repeated`, `replacement`,
		`replace-with`, `repo`, `repo-id`, `report`, `required`, `reserved`, `resolve`, `restore`,
		`result`, `resume`, `rethrow`, `return`, `return-rw`, `returns`, `reverse`, `right`,
		`rindex`, `rmdir`, `role`, `roles_to_compose`, `rolish`, `roll`, `rootdir`, `roots`,
		`rotate`, `rotor`, `round`, `roundrobin`, `routine-type`, `run`, `RUN-MAIN`, `rw`, `rwx`,
		`samecase`, `samemark`, `samewith`, `say`, `schedule-on`, `scheduler`, `scope`, `sec`,
		`sech`, `second`, `secondary`, `seek`, `self`, `send`, `Seq`, `Set`, `set`, `serial`,
		`set_hidden`, `set_name`, `set_package`, `set_rw`, `set_value`, `set_api`, `set_auth`,
		`set_composalizer`, `set_export_callback`, `set_is_mixin`, `set_mixin_attribute`,
		`set_package`, `set_ver`, `set_why`, `SetHash`, `Setty`, `set-instruments`,
		`setup_finalization`, `setup_mixin_cache`, `shape`, `share`, `shell`, `short-id`,
		`short-name`, `shortname`, `shift`, `sibling`, `sigil`, `sign`, `signal`, `signals`,
		`signature`, `sin`, `sinh`, `sink`, `sink-all`, `skip`, `skip-at-least`,
		`skip-at-least-pull-one`, `skip-one`, `skip-rest`, `sleep`, `sleep-timer`, `sleep-until`,
		`Slip`, `slip`, `slurp`, `slurp-rest`, `slurpy`, `snap`, `snapper`, `so`, `socket-host`,
		`socket-port`, `sort`, `source`, `source-package`, `spawn`, `SPEC`, `splice`, `split`,
		`splitdir`, `splitpath`, `sprintf`, `spurt`, `sqrt`, `squish`, `srand`, `stable`, `start`,
		`started`, `starts-with`, `status`, `stderr`, `stdout`, `STORE`, `store-file`,
		`store-repo-id`, `store-unit`, `Str`, `Stringy`, `sub_signature`, `subbuf`, `subbuf-rw`,
		`subname`, `subparse`, `subst`, `subst-mutate`, `substr`, `substr-eq`, `substr-rw`,
		`subtest`, `succ`, `sum`, `suffix`, `summary`, `Supply`, `symlink`, `T`, `t`, `tail`,
		`take`, `take-rw`, `tan`, `tanh`, `tap`, `target`, `target-name`, `tc`, `tclc`, `tell`,
		`term`, `tertiary`, `then`, `throttle`, `throw`, `throws-like`, `time`, `timezone`,
		`tmpdir`, `to`, `today`, `todo`, `toggle`, `to-posix`, `total`, `total-memory`, `trailing`,
		`trans`, `tree`, `trim`, `trim-leading`, `trim-trailing`, `truncate`, `truncated-to`,
		`trusts`, `try_acquire`, `trying`, `twigil`, `type`, `type_captures`, `type_check`,
		`typename`, `uc`, `udp`, `uncaught_handler`, `undefine`, `unimatch`, `unicmp`, `uniname`,
		`uninames`, `uninstall`, `uniparse`, `uniprop`, `uniprops`, `unique`, `unival`, `univals`,
		`unlike`, `unlink`, `unlock`, `unpack`, `unpolar`, `unset`, `unshift`, `unwrap`, `updir`,
		`USAGE`, `usage-name`, `use-ok`, `utc`, `val`, `value`, `values`, `VAR`, `variable`, `ver`,
		`verbose-config`, `Version`, `version`, `VMnames`, `volume`, `vow`, `w`, `wait`, `warn`,
		`watch`, `watch-path`, `week`, `weekday-of-month`, `week-number`, `week-year`, `WHAT`,
		`what`, `when`, `WHERE`, `WHEREFORE`, `WHICH`, `WHO`, `whole-second`, `WHY`, `why`,
		`with-lock-hidden-from-recursion-check`, `wordcase`, `words`, `workaround`, `wrap`,
		`write`, `write-bits`, `write-int128`, `write-int16`, `write-int32`, `write-int64`,
		`write-int8`, `write-num32`, `write-num64`, `write-ubits`, `write-uint128`, `write-uint16`,
		`write-uint32`, `write-uint64`, `write-uint8`, `write-to`, `x`, `yada`, `year`, `yield`,
		`yyyy-mm-dd`, `z`, `zip`, `zip-latest`, `HOW`, `s`, `DEPRECATED`, `trait_mod`,
	}

	builtinRoutinesPattern := Words(`(?<!['\w:-])`, `(?!['\w-])`, builtinRoutines...)

	// A map of opening and closing brackets
	brackets := map[rune]rune{
		'\u0028': '\u0029', '\u003c': '\u003e', '\u005b': '\u005d',
		'\u007b': '\u007d', '\u00ab': '\u00bb', '\u0f3a': '\u0f3b',
		'\u0f3c': '\u0f3d', '\u169b': '\u169c', '\u2018': '\u2019',
		'\u201a': '\u2019', '\u201b': '\u2019', '\u201c': '\u201d',
		'\u201e': '\u201d', '\u201f': '\u201d', '\u2039': '\u203a',
		'\u2045': '\u2046', '\u207d': '\u207e', '\u208d': '\u208e',
		'\u2208': '\u220b', '\u2209': '\u220c', '\u220a': '\u220d',
		'\u2215': '\u29f5', '\u223c': '\u223d', '\u2243': '\u22cd',
		'\u2252': '\u2253', '\u2254': '\u2255', '\u2264': '\u2265',
		'\u2266': '\u2267', '\u2268': '\u2269', '\u226a': '\u226b',
		'\u226e': '\u226f', '\u2270': '\u2271', '\u2272': '\u2273',
		'\u2274': '\u2275', '\u2276': '\u2277', '\u2278': '\u2279',
		'\u227a': '\u227b', '\u227c': '\u227d', '\u227e': '\u227f',
		'\u2280': '\u2281', '\u2282': '\u2283', '\u2284': '\u2285',
		'\u2286': '\u2287', '\u2288': '\u2289', '\u228a': '\u228b',
		'\u228f': '\u2290', '\u2291': '\u2292', '\u2298': '\u29b8',
		'\u22a2': '\u22a3', '\u22a6': '\u2ade', '\u22a8': '\u2ae4',
		'\u22a9': '\u2ae3', '\u22ab': '\u2ae5', '\u22b0': '\u22b1',
		'\u22b2': '\u22b3', '\u22b4': '\u22b5', '\u22b6': '\u22b7',
		'\u22c9': '\u22ca', '\u22cb': '\u22cc', '\u22d0': '\u22d1',
		'\u22d6': '\u22d7', '\u22d8': '\u22d9', '\u22da': '\u22db',
		'\u22dc': '\u22dd', '\u22de': '\u22df', '\u22e0': '\u22e1',
		'\u22e2': '\u22e3', '\u22e4': '\u22e5', '\u22e6': '\u22e7',
		'\u22e8': '\u22e9', '\u22ea': '\u22eb', '\u22ec': '\u22ed',
		'\u22f0': '\u22f1', '\u22f2': '\u22fa', '\u22f3': '\u22fb',
		'\u22f4': '\u22fc', '\u22f6': '\u22fd', '\u22f7': '\u22fe',
		'\u2308': '\u2309', '\u230a': '\u230b', '\u2329': '\u232a',
		'\u23b4': '\u23b5', '\u2768': '\u2769', '\u276a': '\u276b',
		'\u276c': '\u276d', '\u276e': '\u276f', '\u2770': '\u2771',
		'\u2772': '\u2773', '\u2774': '\u2775', '\u27c3': '\u27c4',
		'\u27c5': '\u27c6', '\u27d5': '\u27d6', '\u27dd': '\u27de',
		'\u27e2': '\u27e3', '\u27e4': '\u27e5', '\u27e6': '\u27e7',
		'\u27e8': '\u27e9', '\u27ea': '\u27eb', '\u2983': '\u2984',
		'\u2985': '\u2986', '\u2987': '\u2988', '\u2989': '\u298a',
		'\u298b': '\u298c', '\u298d': '\u298e', '\u298f': '\u2990',
		'\u2991': '\u2992', '\u2993': '\u2994', '\u2995': '\u2996',
		'\u2997': '\u2998', '\u29c0': '\u29c1', '\u29c4': '\u29c5',
		'\u29cf': '\u29d0', '\u29d1': '\u29d2', '\u29d4': '\u29d5',
		'\u29d8': '\u29d9', '\u29da': '\u29db', '\u29f8': '\u29f9',
		'\u29fc': '\u29fd', '\u2a2b': '\u2a2c', '\u2a2d': '\u2a2e',
		'\u2a34': '\u2a35', '\u2a3c': '\u2a3d', '\u2a64': '\u2a65',
		'\u2a79': '\u2a7a', '\u2a7d': '\u2a7e', '\u2a7f': '\u2a80',
		'\u2a81': '\u2a82', '\u2a83': '\u2a84', '\u2a8b': '\u2a8c',
		'\u2a91': '\u2a92', '\u2a93': '\u2a94', '\u2a95': '\u2a96',
		'\u2a97': '\u2a98', '\u2a99': '\u2a9a', '\u2a9b': '\u2a9c',
		'\u2aa1': '\u2aa2', '\u2aa6': '\u2aa7', '\u2aa8': '\u2aa9',
		'\u2aaa': '\u2aab', '\u2aac': '\u2aad', '\u2aaf': '\u2ab0',
		'\u2ab3': '\u2ab4', '\u2abb': '\u2abc', '\u2abd': '\u2abe',
		'\u2abf': '\u2ac0', '\u2ac1': '\u2ac2', '\u2ac3': '\u2ac4',
		'\u2ac5': '\u2ac6', '\u2acd': '\u2ace', '\u2acf': '\u2ad0',
		'\u2ad1': '\u2ad2', '\u2ad3': '\u2ad4', '\u2ad5': '\u2ad6',
		'\u2aec': '\u2aed', '\u2af7': '\u2af8', '\u2af9': '\u2afa',
		'\u2e02': '\u2e03', '\u2e04': '\u2e05', '\u2e09': '\u2e0a',
		'\u2e0c': '\u2e0d', '\u2e1c': '\u2e1d', '\u2e20': '\u2e21',
		'\u3008': '\u3009', '\u300a': '\u300b', '\u300c': '\u300d',
		'\u300e': '\u300f', '\u3010': '\u3011', '\u3014': '\u3015',
		'\u3016': '\u3017', '\u3018': '\u3019', '\u301a': '\u301b',
		'\u301d': '\u301e', '\ufd3e': '\ufd3f', '\ufe17': '\ufe18',
		'\ufe35': '\ufe36', '\ufe37': '\ufe38', '\ufe39': '\ufe3a',
		'\ufe3b': '\ufe3c', '\ufe3d': '\ufe3e', '\ufe3f': '\ufe40',
		'\ufe41': '\ufe42', '\ufe43': '\ufe44', '\ufe47': '\ufe48',
		'\ufe59': '\ufe5a', '\ufe5b': '\ufe5c', '\ufe5d': '\ufe5e',
		'\uff08': '\uff09', '\uff1c': '\uff1e', '\uff3b': '\uff3d',
		'\uff5b': '\uff5d', '\uff5f': '\uff60', '\uff62': '\uff63',
	}

	bracketsPattern := `[` + regexp.QuoteMeta(joinRuneMap(brackets)) + `]`

	// Finds opening brackets and their closing counterparts (including pod and heredoc)
	// and modifies state groups and position accordingly
	findBrackets := func(tokenClass RakuToken) MutatorFunc {
		return func(state *LexerState) error {
			var openingChars []rune
			var adverbs []rune
			switch tokenClass {
			case rakuPod:
				openingChars = []rune(strings.Join(state.Groups[1:5], ``))
			default:
				adverbs = []rune(state.NamedGroups[`adverbs`])
				openingChars = []rune(state.NamedGroups[`opening_delimiters`])
			}

			openingChar := openingChars[0]

			nChars := len(openingChars)

			var closingChar rune
			var closingCharExists bool
			var closingChars []rune

			switch tokenClass {
			case rakuPod:
				closingCharExists = true
			default:
				closingChar, closingCharExists = brackets[openingChar]
			}

			switch tokenClass {
			case rakuPodFormatter:
				formatter := StringOther

				switch state.NamedGroups[`keyword`] {
				case "B":
					formatter = GenericStrong
				case "I":
					formatter = GenericEmph
				case "U":
					formatter = GenericUnderline
				}

				formatterRule := ruleReplacingConfig{
					pattern:      `.+?`,
					tokenType:    formatter,
					mutator:      nil,
					stateName:    `pod-formatter`,
					rulePosition: bottomRule,
				}

				err := replaceRule(formatterRule)(state)
				if err != nil {
					panic(err)
				}

				err = replaceRule(ruleReplacingConfig{
					delimiter:              []rune{closingChar},
					tokenType:              Punctuation,
					stateName:              `pod-formatter`,
					pushState:              true,
					numberOfDelimiterChars: nChars,
					appendMutator:          popRule(formatterRule),
				})(state)
				if err != nil {
					panic(err)
				}

				return nil
			case rakuMatchRegex:
				var delimiter []rune
				if closingCharExists {
					delimiter = []rune{closingChar}
				} else {
					delimiter = openingChars
				}

				err := replaceRule(ruleReplacingConfig{
					delimiter: delimiter,
					tokenType: Punctuation,
					stateName: `regex`,
					popState:  true,
					pushState: true,
				})(state)
				if err != nil {
					panic(err)
				}

				return nil
			case rakuSubstitutionRegex:
				delimiter := regexp2.Escape(string(openingChars))

				err := replaceRule(ruleReplacingConfig{
					pattern:      `(` + delimiter + `)` + `((?:\\\\|\\/|.)*?)` + `(` + delimiter + `)`,
					tokenType:    ByGroups(Punctuation, UsingSelf(`qq`), Punctuation),
					rulePosition: topRule,
					stateName:    `regex`,
					popState:     true,
					pushState:    true,
				})(state)
				if err != nil {
					panic(err)
				}

				return nil
			}

			text := state.Text

			var endPos int

			var nonMirroredOpeningCharPosition int

			if !closingCharExists {
				// it's not a mirrored character, which means we
				// just need to look for the next occurrence
				closingChars = openingChars
				nonMirroredOpeningCharPosition = indexAt(text, closingChars, state.Pos)
				endPos = nonMirroredOpeningCharPosition
			} else {
				var podRegex *regexp2.Regexp
				if tokenClass == rakuPod {
					podRegex = regexp2.MustCompile(
						state.NamedGroups[`ws`]+`=end`+`\s+`+regexp2.Escape(state.NamedGroups[`name`]),
						0,
					)
				} else {
					closingChars = []rune(strings.Repeat(string(closingChar), nChars))
				}

				// we need to look for the corresponding closing character,
				// keep nesting in mind
				nestingLevel := 1

				searchPos := state.Pos - nChars

				var nextClosePos int

				for nestingLevel > 0 {
					if tokenClass == rakuPod {
						match, err := podRegex.FindRunesMatchStartingAt(text, searchPos+nChars)
						if err == nil {
							closingChars = match.Runes()
							nextClosePos = match.Index
						} else {
							nextClosePos = -1
						}
					} else {
						nextClosePos = indexAt(text, closingChars, searchPos+nChars)
					}

					nextOpenPos := indexAt(text, openingChars, searchPos+nChars)

					switch {
					case nextClosePos == -1:
						nextClosePos = len(text)
						nestingLevel = 0
					case nextOpenPos != -1 && nextOpenPos < nextClosePos:
						nestingLevel++
						nChars = len(openingChars)
						searchPos = nextOpenPos
					default: // next_close_pos < next_open_pos
						nestingLevel--
						nChars = len(closingChars)
						searchPos = nextClosePos
					}
				}

				endPos = nextClosePos
			}

			if endPos < 0 {
				// if we didn't find a closer, just highlight the
				// rest of the text in this class
				endPos = len(text)
			}

			adverbre := regexp.MustCompile(`:to\b|:heredoc\b`)
			var heredocTerminator []rune
			var endHeredocPos int
			if adverbre.MatchString(string(adverbs)) {
				if endPos != len(text) {
					heredocTerminator = text[state.Pos:endPos]
					nChars = len(heredocTerminator)
				} else {
					endPos = state.Pos + 1
					heredocTerminator = []rune{}
					nChars = 0
				}

				if nChars > 0 {
					endHeredocPos = indexAt(text[endPos:], heredocTerminator, 0)
					if endHeredocPos > -1 {
						endPos += endHeredocPos
					} else {
						endPos = len(text)
					}
				}
			}

			textBetweenBrackets := string(text[state.Pos:endPos])
			switch tokenClass {
			case rakuPod, rakuPodDeclaration, rakuNameAttribute:
				state.NamedGroups[`value`] = textBetweenBrackets
				state.NamedGroups[`closing_delimiters`] = string(closingChars)
			case rakuQuote:
				if len(heredocTerminator) > 0 {
					// Length of heredoc terminator + closing chars + `;`
					heredocFristPunctuationLen := nChars + len(openingChars) + 1

					state.NamedGroups[`opening_delimiters`] = string(openingChars) +
						string(text[state.Pos:state.Pos+heredocFristPunctuationLen])

					state.NamedGroups[`value`] =
						string(text[state.Pos+heredocFristPunctuationLen : endPos])

					if endHeredocPos > -1 {
						state.NamedGroups[`closing_delimiters`] = string(heredocTerminator)
					}
				} else {
					state.NamedGroups[`value`] = textBetweenBrackets
					if nChars > 0 {
						state.NamedGroups[`closing_delimiters`] = string(closingChars)
					}
				}
			default:
				state.Groups = []string{state.Groups[0] + string(text[state.Pos:endPos+nChars])}
			}

			state.Pos = endPos + nChars

			return nil
		}
	}

	// Raku rules
	// Empty capture groups are placeholders and will be replaced by mutators
	// DO NOT REMOVE THEM!
	return Rules{
		"root": {
			// Placeholder, will be overwritten by mutators, DO NOT REMOVE!
			{`\A\z`, nil, nil},
			Include("common"),
			{`{`, Punctuation, Push(`root`)},
			{`\(`, Punctuation, Push(`root`)},
			{`[)}]`, Punctuation, Pop(1)},
			{`;`, Punctuation, nil},
			{`\[|\]`, Operator, nil},
			{`.+?`, Text, nil},
		},
		"common": {
			{`^#![^\n]*$`, CommentHashbang, nil},
			Include("pod"),
			// Multi-line, Embedded comment
			{
				"#`(?<opening_delimiters>(?<delimiter>" + bracketsPattern + `)\k<delimiter>*)`,
				CommentMultiline,
				findBrackets(rakuMultilineComment),
			},
			{`#[^\n]*$`, CommentSingle, nil},
			// /regex/
			{
				`(?<=(?:^|\(|=|:|~~|\[|{|,|=>)\s*)(/)(?!\]|\))((?:\\\\|\\/|.)*?)((?<!(?<!\\)\\)/(?!'|"))`,
				ByGroups(Punctuation, UsingSelf("regex"), Punctuation),
				nil,
			},
			Include("variable"),
			// ::?VARIABLE
			{`::\?\w+(?::[_UD])?`, NameVariableGlobal, nil},
			// Version
			{
				`\b(v)(\d+)((?:\.(?:\*|[\d\w]+))*)(\+)?`,
				ByGroups(Keyword, NumberInteger, NameEntity, Operator),
				nil,
			},
			Include("number"),
			// Hyperoperator | »*«
			{`(>>)(\S+?)(<<)`, ByGroups(Operator, UsingSelf("root"), Operator), nil},
			{`(»)(\S+?)(«)`, ByGroups(Operator, UsingSelf("root"), Operator), nil},
			// Hyperoperator | «*«
			{`(<<)(\S+?)(<<)`, ByGroups(Operator, UsingSelf("root"), Operator), nil},
			{`(«)(\S+?)(«)`, ByGroups(Operator, UsingSelf("root"), Operator), nil},
			// Hyperoperator | »*»
			{`(>>)(\S+?)(>>)`, ByGroups(Operator, UsingSelf("root"), Operator), nil},
			{`(»)(\S+?)(»)`, ByGroups(Operator, UsingSelf("root"), Operator), nil},
			// <<quoted words>>
			{`(?<!(?:\d+|\.(?:Int|Numeric)|[$@%]\*?[\w':-]+\s+|[\])}]\s+)\s*)(<<)(?!(?:(?!>>)[^\n])+?[},;] *\n)(?!(?:(?!>>).)+?>>\S+?>>)`, Punctuation, Push("<<")},
			// «quoted words»
			{`(?<!(?:\d+|\.(?:Int|Numeric)|[$@%]\*?[\w':-]+\s+|[\])}]\s+)\s*)(«)(?![^»]+?[},;] *\n)(?![^»]+?»\S+?»)`, Punctuation, Push("«")},
			// [<]
			{`(?<=\[\\?)<(?=\])`, Operator, nil},
			// < and > operators | something < onething > something
			{
				`(?<=[$@%&]?\w[\w':-]* +)(<=?)( *[^ ]+? *)(>=?)(?= *[$@%&]?\w[\w':-]*)`,
				ByGroups(Operator, UsingSelf("root"), Operator),
				nil,
			},
			// <quoted words>
			{
				`(?<!(?:\d+|\.(?:Int|Numeric)|[$@%]\*?[\w':-]+\s+|[\])}]\s+)\s*)(<)((?:(?![,;)}] *(?:#[^\n]+)?\n)[^<>])+?)(>)(?!\s*(?:\d+|\.(?:Int|Numeric)|[$@%]\*?\w[\w':-]*[^(]|\s+\[))`,
				ByGroups(Punctuation, String, Punctuation),
				nil,
			},
			{`C?X::['\w:-]+`, NameException, nil},
			Include("metaoperator"),
			// Pair | key => value
			{
				`(\w[\w'-]*)(\s*)(=>)`,
				ByGroups(String, Text, Operator),
				nil,
			},
			Include("colon-pair"),
			// Token
			{
				`(?<=(?:^|\s)(?:regex|token|rule)(\s+))` + namePattern + colonPairLookahead + `\s*[({])`,
				NameFunction,
				Push("token", "name-adverb"),
			},
			// Substitution
			{`(?<=^|\b|\s)(?<!\.)(ss|S|s|TR|tr)\b(\s*)`, ByGroups(Keyword, Text), Push("substitution")},
			{keywordsPattern, Keyword, nil},
			{builtinTypesPattern, NameBuiltin, nil},
			{builtinRoutinesPattern, NameBuiltin, nil},
			// Class name
			{
				`(?<=(?:^|\s)(?:class|grammar|role|does|but|is|subset|of)\s+)` + namePattern,
				NameClass,
				Push("name-adverb"),
			},
			//  Routine
			{
				`(?<=(?:^|\s)(?:sub|method|multi sub|multi)\s+)!?` + namePattern + colonPairLookahead + `\s*[({])`,
				NameFunction,
				Push("name-adverb"),
			},
			// Constant
			{`(?<=\bconstant\s+)` + namePattern, NameConstant, Push("name-adverb")},
			// Namespace
			{`(?<=\b(?:use|module|package)\s+)` + namePattern, NameNamespace, Push("name-adverb")},
			Include("operator"),
			Include("single-quote"),
			{`(?<!(?<!\\)\\)"`, Punctuation, Push("double-quotes")},
			// m,rx regex
			{`(?<=^|\b|\s)(ms|m|rx)\b(\s*)`, ByGroups(Keyword, Text), Push("rx")},
			// Quote constructs
			{
				`(?<=^|\b|\s)(?<keyword>(?:qq|q|Q))(?<adverbs>(?::?(?:heredoc|to|qq|ww|q|w|s|a|h|f|c|b|to|v|x))*)(?<ws>\s*)(?<opening_delimiters>(?<delimiter>[^0-9a-zA-Z:\s])\k<delimiter>*)`,
				EmitterFunc(quote),
				findBrackets(rakuQuote),
			},
			// Function
			{
				`\b` + namePattern + colonPairLookahead + `\()`,
				NameFunction,
				Push("name-adverb"),
			},
			// Method
			{
				`(?<!\.\.[?^*+]?)(?<=(?:\.[?^*+&]?)|self!)` + namePattern + colonPairLookahead + `\b)`,
				NameFunction,
				Push("name-adverb"),
			},
			// Indirect invocant
			{namePattern + `(?=\s+\W?['\w:-]+:\W)`, NameFunction, Push("name-adverb")},
			{`(?<=\W)(?:∅|i|e|𝑒|tau|τ|pi|π|Inf|∞)(?=\W)`, NameConstant, nil},
			{`(｢)([^｣]*)(｣)`, ByGroups(Punctuation, String, Punctuation), nil},
			{`(?<=^ *)\b` + namePattern + `(?=:\s*(?:for|while|loop))`, NameLabel, nil},
			// Sigilless variable
			{
				`(?<=\b(?:my|our|constant|let|temp)\s+)\\` + namePattern,
				NameVariable,
				Push("name-adverb"),
			},
			{namePattern, Name, Push("name-adverb")},
		},
		"rx": {
			Include("colon-pair-attribute"),
			{
				`(?<opening_delimiters>(?<delimiter>[^\w:\s])\k<delimiter>*)`,
				ByGroupNames(
					map[string]Emitter{
						`opening_delimiters`: Punctuation,
						`delimiter`:          nil,
					},
				),
				findBrackets(rakuMatchRegex),
			},
		},
		"substitution": {
			Include("colon-pair-attribute"),
			// Substitution | s{regex} = value
			{
				`(?<opening_delimiters>(?<delimiter>` + bracketsPattern + `)\k<delimiter>*)`,
				ByGroupNames(map[string]Emitter{
					`opening_delimiters`: Punctuation,
					`delimiter`:          nil,
				}),
				findBrackets(rakuMatchRegex),
			},
			// Substitution | s/regex/string/
			{
				`(?<opening_delimiters>[^\w:\s])`,
				Punctuation,
				findBrackets(rakuSubstitutionRegex),
			},
		},
		"number": {
			{`0_?[0-7]+(_[0-7]+)*`, LiteralNumberOct, nil},
			{`0x[0-9A-Fa-f]+(_[0-9A-Fa-f]+)*`, LiteralNumberHex, nil},
			{`0b[01]+(_[01]+)*`, LiteralNumberBin, nil},
			{
				`(?i)(\d*(_\d*)*\.\d+(_\d*)*|\d+(_\d*)*\.\d+(_\d*)*)(e[+-]?\d+)?`,
				LiteralNumberFloat,
				nil,
			},
			{`(?i)\d+(_\d*)*e[+-]?\d+(_\d*)*`, LiteralNumberFloat, nil},
			{`(?<=\d+)i`, NameConstant, nil},
			{`\d+(_\d+)*`, LiteralNumberInteger, nil},
		},
		"name-adverb": {
			Include("colon-pair-attribute-keyvalue"),
			Default(Pop(1)),
		},
		"colon-pair": {
			// :key(value)
			{colonPairPattern, colonPair(String), findBrackets(rakuNameAttribute)},
			// :123abc
			{
				`(:)(\d+)(\w[\w'-]*)`,
				ByGroups(Punctuation, UsingSelf("number"), String),
				nil,
			},
			// :key
			{`(:)(!?)(\w[\w'-]*)`, ByGroups(Punctuation, Operator, String), nil},
			{`\s+`, Text, nil},
		},
		"colon-pair-attribute": {
			// :key(value)
			{colonPairPattern, colonPair(NameAttribute), findBrackets(rakuNameAttribute)},
			// :123abc
			{
				`(:)(\d+)(\w[\w'-]*)`,
				ByGroups(Punctuation, UsingSelf("number"), NameAttribute),
				nil,
			},
			// :key
			{`(:)(!?)(\w[\w'-]*)`, ByGroups(Punctuation, Operator, NameAttribute), nil},
			{`\s+`, Text, nil},
		},
		"colon-pair-attribute-keyvalue": {
			// :key(value)
			{colonPairPattern, colonPair(NameAttribute), findBrackets(rakuNameAttribute)},
		},
		"escape-qq": {
			{
				`(?<!(?<!\\)\\)(\\qq)(\[)(.+?)(\])`,
				ByGroups(StringEscape, Punctuation, UsingSelf("qq"), Punctuation),
				nil,
			},
		},
		`escape-char`: {
			{`(?<!(?<!\\)\\)(\\[abfrnrt])`, StringEscape, nil},
		},
		`escape-single-quote`: {
			{`(?<!(?<!\\)\\)(\\)(['\\])`, ByGroups(StringEscape, StringSingle), nil},
		},
		"escape-c-name": {
			{
				`(?<!(?<!\\)\\)(\\[c|C])(\[)(.+?)(\])`,
				ByGroups(StringEscape, Punctuation, String, Punctuation),
				nil,
			},
		},
		"escape-hexadecimal": {
			{
				`(?<!(?<!\\)\\)(\\[x|X])(\[)([0-9a-fA-F]+)(\])`,
				ByGroups(StringEscape, Punctuation, NumberHex, Punctuation),
				nil,
			},
			{`(\\[x|X])([0-9a-fA-F]+)`, ByGroups(StringEscape, NumberHex), nil},
		},
		"regex": {
			// Placeholder, will be overwritten by mutators, DO NOT REMOVE!
			{`\A\z`, nil, nil},
			Include("regex-escape-class"),
			Include(`regex-character-escape`),
			// $(code)
			{
				`([$@])((?<!(?<!\\)\\)\()`,
				ByGroups(Keyword, Punctuation),
				replaceRule(ruleReplacingConfig{
					delimiter: []rune(`)`),
					tokenType: Punctuation,
					stateName: `root`,
					pushState: true,
				}),
			},
			// Exclude $/ from variables, because we can't get out of the end of the slash regex: $/;
			{`\$(?=/)`, NameEntity, nil},
			// Exclude $ from variables
			{`\$(?=\z|\s|[^<(\w*!.])`, NameEntity, nil},
			Include("variable"),
			Include("escape-c-name"),
			Include("escape-hexadecimal"),
			Include("number"),
			Include("single-quote"),
			// :my variable code ...
			{
				`(?<!(?<!\\)\\)(:)(my|our|state|constant|temp|let)`,
				ByGroups(Operator, KeywordDeclaration),
				replaceRule(ruleReplacingConfig{
					delimiter: []rune(`;`),
					tokenType: Punctuation,
					stateName: `root`,
					pushState: true,
				}),
			},
			// <{code}>
			{
				`(?<!(?<!\\)\\)(<)([?!.]*)((?<!(?<!\\)\\){)`,
				ByGroups(Punctuation, Operator, Punctuation),
				replaceRule(ruleReplacingConfig{
					delimiter: []rune(`}>`),
					tokenType: Punctuation,
					stateName: `root`,
					pushState: true,
				}),
			},
			// {code}
			Include(`closure`),
			// Properties
			{`(:)(\w+)`, ByGroups(Punctuation, NameAttribute), nil},
			// Operator
			{`\|\||\||&&|&|\.\.|\*\*|%%|%|:|!|<<|«|>>|»|\+|\*\*|\*|\?|=|~|<~~>`, Operator, nil},
			// Anchors
			{`\^\^|\^|\$\$|\$`, NameEntity, nil},
			{`\.`, NameEntity, nil},
			{`#[^\n]*\n`, CommentSingle, nil},
			// Lookaround
			{
				`(?<!(?<!\\)\\)(<)(\s*)([?!.]+)(\s*)(after|before)`,
				ByGroups(Punctuation, Text, Operator, Text, OperatorWord),
				replaceRule(ruleReplacingConfig{
					delimiter: []rune(`>`),
					tokenType: Punctuation,
					stateName: `regex`,
					pushState: true,
				}),
			},
			{
				`(?<!(?<!\\)\\)(<)([|!?.]*)(wb|ww|ws|w)(>)`,
				ByGroups(Punctuation, Operator, OperatorWord, Punctuation),
				nil,
			},
			// <$variable>
			{
				`(?<!(?<!\\)\\)(<)([?!.]*)([$@]\w[\w:-]*)(>)`,
				ByGroups(Punctuation, Operator, NameVariable, Punctuation),
				nil,
			},
			// Capture markers
			{`(?<!(?<!\\)\\)<\(|\)>`, Operator, nil},
			{
				`(?<!(?<!\\)\\)(<)(\w[\w:-]*)(=\.?)`,
				ByGroups(Punctuation, NameVariable, Operator),
				Push(`regex-variable`),
			},
			{
				`(?<!(?<!\\)\\)(<)([|!?.&]*)(\w(?:(?!:\s)[\w':-])*)`,
				ByGroups(Punctuation, Operator, NameFunction),
				Push(`regex-function`),
			},
			{`(?<!(?<!\\)\\)<`, Punctuation, Push("regex-property")},
			{`(?<!(?<!\\)\\)"`, Punctuation, Push("double-quotes")},
			{`(?<!(?<!\\)\\)(?:\]|\))`, Punctuation, Pop(1)},
			{`(?<!(?<!\\)\\)(?:\[|\()`, Punctuation, Push("regex")},
			{`.+?`, StringRegex, nil},
		},
		"regex-class-builtin": {
			{
				`\b(?:alnum|alpha|blank|cntrl|digit|graph|lower|print|punct|space|upper|xdigit|same|ident)\b`,
				NameBuiltin,
				nil,
			},
		},
		"regex-function": {
			// <function>
			{`(?<!(?<!\\)\\)>`, Punctuation, Pop(1)},
			// <function(parameter)>
			{
				`\(`,
				Punctuation,
				replaceRule(ruleReplacingConfig{
					delimiter: []rune(`)>`),
					tokenType: Punctuation,
					stateName: `root`,
					popState:  true,
					pushState: true,
				}),
			},
			// <function value>
			{
				`\s+`,
				StringRegex,
				replaceRule(ruleReplacingConfig{
					delimiter: []rune(`>`),
					tokenType: Punctuation,
					stateName: `regex`,
					popState:  true,
					pushState: true,
				}),
			},
			// <function: value>
			{
				`:`,
				Punctuation,
				replaceRule(ruleReplacingConfig{
					delimiter: []rune(`>`),
					tokenType: Punctuation,
					stateName: `root`,
					popState:  true,
					pushState: true,
				}),
			},
		},
		"regex-variable": {
			Include(`regex-starting-operators`),
			// <var=function(
			{
				`(&)?(\w(?:(?!:\s)[\w':-])*)(?=\()`,
				ByGroups(Operator, NameFunction),
				Mutators(Pop(1), Push(`regex-function`)),
			},
			// <var=function>
			{`(&)?(\w[\w':-]*)(>)`, ByGroups(Operator, NameFunction, Punctuation), Pop(1)},
			// <var=
			Default(Pop(1), Push(`regex-property`)),
		},
		"regex-property": {
			{`(?<!(?<!\\)\\)>`, Punctuation, Pop(1)},
			Include("regex-class-builtin"),
			Include("variable"),
			Include(`regex-starting-operators`),
			Include("colon-pair-attribute"),
			{`(?<!(?<!\\)\\)\[`, Punctuation, Push("regex-character-class")},
			{`\+|\-`, Operator, nil},
			{`@[\w':-]+`, NameVariable, nil},
			{`.+?`, StringRegex, nil},
		},
		`regex-starting-operators`: {
			{`(?<=<)[|!?.]+`, Operator, nil},
		},
		"regex-escape-class": {
			{`(?i)\\n|\\t|\\h|\\v|\\s|\\d|\\w`, StringEscape, nil},
		},
		`regex-character-escape`: {
			{`(?<!(?<!\\)\\)(\\)(.)`, ByGroups(StringEscape, StringRegex), nil},
		},
		"regex-character-class": {
			{`(?<!(?<!\\)\\)\]`, Punctuation, Pop(1)},
			Include("regex-escape-class"),
			Include("escape-c-name"),
			Include("escape-hexadecimal"),
			Include(`regex-character-escape`),
			Include("number"),
			{`\.\.`, Operator, nil},
			{`.+?`, StringRegex, nil},
		},
		"metaoperator": {
			// Z[=>]
			{
				`\b([RZX]+)\b(\[)([^\s\]]+?)(\])`,
				ByGroups(OperatorWord, Punctuation, UsingSelf("root"), Punctuation),
				nil,
			},
			// Z=>
			{`\b([RZX]+)\b([^\s\]]+)`, ByGroups(OperatorWord, UsingSelf("operator")), nil},
		},
		"operator": {
			// Word Operator
			{wordOperatorsPattern, OperatorWord, nil},
			// Operator
			{operatorsPattern, Operator, nil},
		},
		"pod": {
			// Single-line pod declaration
			{`(#[|=])\s`, Keyword, Push("pod-single")},
			// Multi-line pod declaration
			{
				"(?<keyword>#[|=])(?<opening_delimiters>(?<delimiter>" + bracketsPattern + `)\k<delimiter>*)(?<value>)(?<closing_delimiters>)`,
				ByGroupNames(
					map[string]Emitter{
						`keyword`:            Keyword,
						`opening_delimiters`: Punctuation,
						`delimiter`:          nil,
						`value`:              UsingSelf("pod-declaration"),
						`closing_delimiters`: Punctuation,
					}),
				findBrackets(rakuPodDeclaration),
			},
			Include("pod-blocks"),
		},
		"pod-blocks": {
			// =begin code
			{
				`(?<=^ *)(?<ws> *)(?<keyword>=begin)(?<ws2> +)(?<name>code)(?<config>[^\n]*)(?<value>.*?)(?<ws3>^\k<ws>)(?<end_keyword>=end)(?<ws4> +)\k<name>`,
				EmitterFunc(podCode),
				nil,
			},
			// =begin
			{
				`(?<=^ *)(?<ws> *)(?<keyword>=begin)(?<ws2> +)(?!code)(?<name>\w[\w'-]*)(?<config>[^\n]*)(?<value>)(?<closing_delimiters>)`,
				ByGroupNames(
					map[string]Emitter{
						`ws`:                 Comment,
						`keyword`:            Keyword,
						`ws2`:                StringDoc,
						`name`:               Keyword,
						`config`:             EmitterFunc(podConfig),
						`value`:              UsingSelf("pod-begin"),
						`closing_delimiters`: Keyword,
					}),
				findBrackets(rakuPod),
			},
			// =for ...
			{
				`(?<=^ *)(?<ws> *)(?<keyword>=(?:for|defn))(?<ws2> +)(?<name>\w[\w'-]*)(?<config>[^\n]*\n)`,
				ByGroups(Comment, Keyword, StringDoc, Keyword, EmitterFunc(podConfig)),
				Push("pod-paragraph"),
			},
			// =config
			{
				`(?<=^ *)(?<ws> *)(?<keyword>=config)(?<ws2> +)(?<name>\w[\w'-]*)(?<config>[^\n]*\n)`,
				ByGroups(Comment, Keyword, StringDoc, Keyword, EmitterFunc(podConfig)),
				nil,
			},
			// =alias
			{
				`(?<=^ *)(?<ws> *)(?<keyword>=alias)(?<ws2> +)(?<name>\w[\w'-]*)(?<value>[^\n]*\n)`,
				ByGroups(Comment, Keyword, StringDoc, Keyword, StringDoc),
				nil,
			},
			// =encoding
			{
				`(?<=^ *)(?<ws> *)(?<keyword>=encoding)(?<ws2> +)(?<name>[^\n]+)`,
				ByGroups(Comment, Keyword, StringDoc, Name),
				nil,
			},
			// =para ...
			{
				`(?<=^ *)(?<ws> *)(?<keyword>=(?:para|table|pod))(?<config>(?<!\n\s*)[^\n]*\n)`,
				ByGroups(Comment, Keyword, EmitterFunc(podConfig)),
				Push("pod-paragraph"),
			},
			// =head1 ...
			{
				`(?<=^ *)(?<ws> *)(?<keyword>=head\d+)(?<ws2> *)(?<config>#?)`,
				ByGroups(Comment, Keyword, GenericHeading, Keyword),
				Push("pod-heading"),
			},
			// =item ...
			{
				`(?<=^ *)(?<ws> *)(?<keyword>=(?:item\d*|comment|data|[A-Z]+))(?<ws2> *)(?<config>#?)`,
				ByGroups(Comment, Keyword, StringDoc, Keyword),
				Push("pod-paragraph"),
			},
			{
				`(?<=^ *)(?<ws> *)(?<keyword>=finish)(?<config>[^\n]*)`,
				ByGroups(Comment, Keyword, EmitterFunc(podConfig)),
				Push("pod-finish"),
			},
			// ={custom} ...
			{
				`(?<=^ *)(?<ws> *)(?<name>=\w[\w'-]*)(?<ws2> *)(?<config>#?)`,
				ByGroups(Comment, Name, StringDoc, Keyword),
				Push("pod-paragraph"),
			},
			// = podconfig
			{
				`(?<=^ *)(?<keyword> *=)(?<ws> *)(?<config>(?::\w[\w'-]*(?:` + colonPairOpeningBrackets + `.+?` +
					colonPairClosingBrackets + `) *)*\n)`,
				ByGroups(Keyword, StringDoc, EmitterFunc(podConfig)),
				nil,
			},
		},
		"pod-begin": {
			Include("pod-blocks"),
			Include("pre-pod-formatter"),
			{`.+?`, StringDoc, nil},
		},
		"pod-declaration": {
			Include("pre-pod-formatter"),
			{`.+?`, StringDoc, nil},
		},
		"pod-paragraph": {
			{`\n *\n|\n(?=^ *=)`, StringDoc, Pop(1)},
			Include("pre-pod-formatter"),
			{`.+?`, StringDoc, nil},
		},
		"pod-single": {
			{`\n`, StringDoc, Pop(1)},
			Include("pre-pod-formatter"),
			{`.+?`, StringDoc, nil},
		},
		"pod-heading": {
			{`\n *\n|\n(?=^ *=)`, GenericHeading, Pop(1)},
			Include("pre-pod-formatter"),
			{`.+?`, GenericHeading, nil},
		},
		"pod-finish": {
			{`\z`, nil, Pop(1)},
			Include("pre-pod-formatter"),
			{`.+?`, StringDoc, nil},
		},
		"pre-pod-formatter": {
			// C<code>, B<bold>, ...
			{
				`(?<keyword>[CBIUDTKRPAELZVMSXN])(?<opening_delimiters><+|«)`,
				ByGroups(Keyword, Punctuation),
				findBrackets(rakuPodFormatter),
			},
		},
		"pod-formatter": {
			// Placeholder rule, will be replaced by mutators. DO NOT REMOVE!
			{`>`, Punctuation, Pop(1)},
			Include("pre-pod-formatter"),
			// Placeholder rule, will be replaced by mutators. DO NOT REMOVE!
			{`.+?`, StringOther, nil},
		},
		"variable": {
			{variablePattern, NameVariable, Push("name-adverb")},
			{globalVariablePattern, NameVariableGlobal, Push("name-adverb")},
			{`[$@]<[^>]+>`, NameVariable, nil},
			{`\$[/!¢]`, NameVariable, nil},
			{`[$@%]`, NameVariable, nil},
		},
		"single-quote": {
			{`(?<!(?<!\\)\\)'`, Punctuation, Push("single-quote-inner")},
		},
		"single-quote-inner": {
			{`(?<!(?<!(?<!\\)\\)\\)'`, Punctuation, Pop(1)},
			Include("escape-single-quote"),
			Include("escape-qq"),
			{`(?:\\\\|\\[^\\]|[^'\\])+?`, StringSingle, nil},
		},
		"double-quotes": {
			{`(?<!(?<!\\)\\)"`, Punctuation, Pop(1)},
			Include("qq"),
		},
		"<<": {
			{`>>(?!\s*(?:\d+|\.(?:Int|Numeric)|[$@%]\*?[\w':-]+|\s+\[))`, Punctuation, Pop(1)},
			Include("ww"),
		},
		"«": {
			{`»(?!\s*(?:\d+|\.(?:Int|Numeric)|[$@%]\*?[\w':-]+|\s+\[))`, Punctuation, Pop(1)},
			Include("ww"),
		},
		"ww": {
			Include("single-quote"),
			Include("qq"),
		},
		"qq": {
			Include("qq-variable"),
			Include("closure"),
			Include(`escape-char`),
			Include("escape-hexadecimal"),
			Include("escape-c-name"),
			Include("escape-qq"),
			{`.+?`, StringDouble, nil},
		},
		"qq-variable": {
			{
				`(?<!(?<!\\)\\)(?:` + variablePattern + `|` + globalVariablePattern + `)` + colonPairLookahead + `)`,
				NameVariable,
				Push("qq-variable-extras", "name-adverb"),
			},
		},
		"qq-variable-extras": {
			// Method
			{
				`(?<operator>\.)(?<method_name>` + namePattern + `)` + colonPairLookahead + `\()`,
				ByGroupNames(map[string]Emitter{
					`operator`:    Operator,
					`method_name`: NameFunction,
				}),
				Push(`name-adverb`),
			},
			// Function/Signature
			{
				`\(`, Punctuation, replaceRule(
					ruleReplacingConfig{
						delimiter: []rune(`)`),
						tokenType: Punctuation,
						stateName: `root`,
						pushState: true,
					}),
			},
			Default(Pop(1)),
		},
		"Q": {
			Include("escape-qq"),
			{`.+?`, String, nil},
		},
		"Q-closure": {
			Include("escape-qq"),
			Include("closure"),
			{`.+?`, String, nil},
		},
		"Q-variable": {
			Include("escape-qq"),
			Include("qq-variable"),
			{`.+?`, String, nil},
		},
		"closure": {
			{`(?<!(?<!\\)\\){`, Punctuation, replaceRule(
				ruleReplacingConfig{
					delimiter: []rune(`}`),
					tokenType: Punctuation,
					stateName: `root`,
					pushState: true,
				}),
			},
		},
		"token": {
			// Token signature
			{`\(`, Punctuation, replaceRule(
				ruleReplacingConfig{
					delimiter: []rune(`)`),
					tokenType: Punctuation,
					stateName: `root`,
					pushState: true,
				}),
			},
			{`{`, Punctuation, replaceRule(
				ruleReplacingConfig{
					delimiter: []rune(`}`),
					tokenType: Punctuation,
					stateName: `regex`,
					popState:  true,
					pushState: true,
				}),
			},
			{`\s*`, Text, nil},
			Default(Pop(1)),
		},
	}
}

// Joins keys of rune map
func joinRuneMap(m map[rune]rune) string {
	runes := make([]rune, 0, len(m))
	for k := range m {
		runes = append(runes, k)
	}

	return string(runes)
}

// Finds the index of substring in the string starting at position n
func indexAt(str []rune, substr []rune, pos int) int {
	strFromPos := str[pos:]
	text := string(strFromPos)

	idx := strings.Index(text, string(substr))
	if idx > -1 {
		idx = utf8.RuneCountInString(text[:idx])

		// Search again if the substr is escaped with backslash
		if (idx > 1 && strFromPos[idx-1] == '\\' && strFromPos[idx-2] != '\\') ||
			(idx == 1 && strFromPos[idx-1] == '\\') {
			idx = indexAt(str[pos:], substr, idx+1)

			idx = utf8.RuneCountInString(text[:idx])

			if idx < 0 {
				return idx
			}
		}
		idx += pos
	}

	return idx
}

// Tells if an array of string contains a string
func contains(s []string, e string) bool {
	for _, value := range s {
		if value == e {
			return true
		}
	}
	return false
}

type rulePosition int

const (
	topRule    rulePosition = 0
	bottomRule              = -1
)

type ruleMakingConfig struct {
	delimiter              []rune
	pattern                string
	tokenType              Emitter
	mutator                Mutator
	numberOfDelimiterChars int
}

type ruleReplacingConfig struct {
	delimiter              []rune
	pattern                string
	tokenType              Emitter
	numberOfDelimiterChars int
	mutator                Mutator
	appendMutator          Mutator
	rulePosition           rulePosition
	stateName              string
	pop                    bool
	popState               bool
	pushState              bool
}

// Pops rule from state-stack and replaces the rule with the previous rule
func popRule(rule ruleReplacingConfig) MutatorFunc {
	return func(state *LexerState) error {
		stackName := genStackName(rule.stateName, rule.rulePosition)

		stack, ok := state.Get(stackName).([]ruleReplacingConfig)

		if ok && len(stack) > 0 {
			// Pop from stack
			stack = stack[:len(stack)-1]
			lastRule := stack[len(stack)-1]
			lastRule.pushState = false
			lastRule.popState = false
			lastRule.pop = true
			state.Set(stackName, stack)

			// Call replaceRule to use the last rule
			err := replaceRule(lastRule)(state)
			if err != nil {
				panic(err)
			}
		}

		return nil
	}
}

// Replaces a state's rule based on the rule config and position
func replaceRule(rule ruleReplacingConfig) MutatorFunc {
	return func(state *LexerState) error {
		stateName := rule.stateName
		stackName := genStackName(rule.stateName, rule.rulePosition)

		stack, ok := state.Get(stackName).([]ruleReplacingConfig)
		if !ok {
			stack = []ruleReplacingConfig{}
		}

		// If state-stack is empty fill it with the placeholder rule
		if len(stack) == 0 {
			stack = []ruleReplacingConfig{
				{
					// Placeholder, will be overwritten by mutators, DO NOT REMOVE!
					pattern:      `\A\z`,
					tokenType:    nil,
					mutator:      nil,
					stateName:    stateName,
					rulePosition: rule.rulePosition,
				},
			}
			state.Set(stackName, stack)
		}

		var mutator Mutator
		mutators := []Mutator{}

		switch {
		case rule.rulePosition == topRule && rule.mutator == nil:
			// Default mutator for top rule
			mutators = []Mutator{Pop(1), popRule(rule)}
		case rule.rulePosition == topRule && rule.mutator != nil:
			// Default mutator for top rule, when rule.mutator is set
			mutators = []Mutator{rule.mutator, popRule(rule)}
		case rule.mutator != nil:
			mutators = []Mutator{rule.mutator}
		}

		if rule.appendMutator != nil {
			mutators = append(mutators, rule.appendMutator)
		}

		if len(mutators) > 0 {
			mutator = Mutators(mutators...)
		} else {
			mutator = nil
		}

		ruleConfig := ruleMakingConfig{
			pattern:                rule.pattern,
			delimiter:              rule.delimiter,
			numberOfDelimiterChars: rule.numberOfDelimiterChars,
			tokenType:              rule.tokenType,
			mutator:                mutator,
		}

		cRule := makeRule(ruleConfig)

		switch rule.rulePosition {
		case topRule:
			state.Rules[stateName][0] = cRule
		case bottomRule:
			state.Rules[stateName][len(state.Rules[stateName])-1] = cRule
		}

		// Pop state name from stack if asked. State should be popped first before Pushing
		if rule.popState {
			err := Pop(1)(state)
			if err != nil {
				panic(err)
			}
		}

		// Push state name to stack if asked
		if rule.pushState {
			err := Push(stateName)(state)
			if err != nil {
				panic(err)
			}
		}

		if !rule.pop {
			state.Set(stackName, append(stack, rule))
		}

		return nil
	}
}

// Generates rule replacing stack using state name and rule position
func genStackName(stateName string, rulePosition rulePosition) (stackName string) {
	switch rulePosition {
	case topRule:
		stackName = stateName + `-top-stack`
	case bottomRule:
		stackName = stateName + `-bottom-stack`
	}
	return
}

// Makes a compiled rule and returns it
func makeRule(config ruleMakingConfig) *CompiledRule {
	var rePattern string

	if len(config.delimiter) > 0 {
		delimiter := string(config.delimiter)

		if config.numberOfDelimiterChars > 1 {
			delimiter = strings.Repeat(delimiter, config.numberOfDelimiterChars)
		}

		rePattern = `(?<!(?<!\\)\\)` + regexp2.Escape(delimiter)
	} else {
		rePattern = config.pattern
	}

	regex := regexp2.MustCompile(rePattern, regexp2.None)

	cRule := &CompiledRule{
		Rule:   Rule{rePattern, config.tokenType, config.mutator},
		Regexp: regex,
	}

	return cRule
}

// Emitter for colon pairs, changes token state based on key and brackets
func colonPair(tokenClass TokenType) Emitter {
	return EmitterFunc(func(groups []string, state *LexerState) Iterator {
		iterators := []Iterator{}
		tokens := []Token{
			{Punctuation, state.NamedGroups[`colon`]},
			{Punctuation, state.NamedGroups[`opening_delimiters`]},
			{Punctuation, state.NamedGroups[`closing_delimiters`]},
		}

		// Append colon
		iterators = append(iterators, Literator(tokens[0]))

		if tokenClass == NameAttribute {
			iterators = append(iterators, Literator(Token{NameAttribute, state.NamedGroups[`key`]}))
		} else {
			var keyTokenState string
			keyre := regexp.MustCompile(`^\d+$`)
			if keyre.MatchString(state.NamedGroups[`key`]) {
				keyTokenState = "common"
			} else {
				keyTokenState = "Q"
			}

			// Use token state to Tokenise key
			if keyTokenState != "" {
				iterator, err := state.Lexer.Tokenise(
					&TokeniseOptions{
						State:  keyTokenState,
						Nested: true,
					}, state.NamedGroups[`key`])

				if err != nil {
					panic(err)
				} else {
					// Append key
					iterators = append(iterators, iterator)
				}
			}
		}

		// Append punctuation
		iterators = append(iterators, Literator(tokens[1]))

		var valueTokenState string

		switch state.NamedGroups[`opening_delimiters`] {
		case "(", "{", "[":
			valueTokenState = "root"
		case "<<", "«":
			valueTokenState = "ww"
		case "<":
			valueTokenState = "Q"
		}

		// Use token state to Tokenise value
		if valueTokenState != "" {
			iterator, err := state.Lexer.Tokenise(
				&TokeniseOptions{
					State:  valueTokenState,
					Nested: true,
				}, state.NamedGroups[`value`])

			if err != nil {
				panic(err)
			} else {
				// Append value
				iterators = append(iterators, iterator)
			}
		}
		// Append last punctuation
		iterators = append(iterators, Literator(tokens[2]))

		return Concaterator(iterators...)
	})
}

// Emitter for quoting constructs, changes token state based on quote name and adverbs
func quote(groups []string, state *LexerState) Iterator {
	keyword := state.NamedGroups[`keyword`]
	adverbsStr := state.NamedGroups[`adverbs`]
	iterators := []Iterator{}
	tokens := []Token{
		{Keyword, keyword},
		{StringAffix, adverbsStr},
		{Text, state.NamedGroups[`ws`]},
		{Punctuation, state.NamedGroups[`opening_delimiters`]},
		{Punctuation, state.NamedGroups[`closing_delimiters`]},
	}

	// Append all tokens before dealing with the main string
	iterators = append(iterators, Literator(tokens[:4]...))

	var tokenStates []string

	// Set tokenStates based on adverbs
	adverbs := strings.Split(adverbsStr, ":")
	for _, adverb := range adverbs {
		switch adverb {
		case "c", "closure":
			tokenStates = append(tokenStates, "Q-closure")
		case "qq":
			tokenStates = append(tokenStates, "qq")
		case "ww":
			tokenStates = append(tokenStates, "ww")
		case "s", "scalar", "a", "array", "h", "hash", "f", "function":
			tokenStates = append(tokenStates, "Q-variable")
		}
	}

	var tokenState string

	switch {
	case keyword == "qq" || contains(tokenStates, "qq"):
		tokenState = "qq"
	case adverbsStr == "ww" || contains(tokenStates, "ww"):
		tokenState = "ww"
	case contains(tokenStates, "Q-closure") && contains(tokenStates, "Q-variable"):
		tokenState = "qq"
	case contains(tokenStates, "Q-closure"):
		tokenState = "Q-closure"
	case contains(tokenStates, "Q-variable"):
		tokenState = "Q-variable"
	default:
		tokenState = "Q"
	}

	iterator, err := state.Lexer.Tokenise(
		&TokeniseOptions{
			State:  tokenState,
			Nested: true,
		}, state.NamedGroups[`value`])

	if err != nil {
		panic(err)
	} else {
		iterators = append(iterators, iterator)
	}

	// Append the last punctuation
	iterators = append(iterators, Literator(tokens[4]))

	return Concaterator(iterators...)
}

// Emitter for pod config, tokenises the properties with "colon-pair-attribute" state
func podConfig(groups []string, state *LexerState) Iterator {
	// Tokenise pod config
	iterator, err := state.Lexer.Tokenise(
		&TokeniseOptions{
			State:  "colon-pair-attribute",
			Nested: true,
		}, groups[0])

	if err != nil {
		panic(err)
	} else {
		return iterator
	}
}

// Emitter for pod code, tokenises the code based on the lang specified
func podCode(groups []string, state *LexerState) Iterator {
	iterators := []Iterator{}
	tokens := []Token{
		{Comment, state.NamedGroups[`ws`]},
		{Keyword, state.NamedGroups[`keyword`]},
		{Keyword, state.NamedGroups[`ws2`]},
		{Keyword, state.NamedGroups[`name`]},
		{StringDoc, state.NamedGroups[`value`]},
		{Comment, state.NamedGroups[`ws3`]},
		{Keyword, state.NamedGroups[`end_keyword`]},
		{Keyword, state.NamedGroups[`ws4`]},
		{Keyword, state.NamedGroups[`name`]},
	}

	// Append all tokens before dealing with the pod config
	iterators = append(iterators, Literator(tokens[:4]...))

	// Tokenise pod config
	iterators = append(iterators, podConfig([]string{state.NamedGroups[`config`]}, state))

	langMatch := regexp.MustCompile(`:lang\W+(\w+)`).FindStringSubmatch(state.NamedGroups[`config`])
	var lang string
	if len(langMatch) > 1 {
		lang = langMatch[1]
	}

	// Tokenise code based on lang property
	sublexer := internal.Get(lang)
	if sublexer != nil {
		iterator, err := sublexer.Tokenise(nil, state.NamedGroups[`value`])

		if err != nil {
			panic(err)
		} else {
			iterators = append(iterators, iterator)
		}
	} else {
		iterators = append(iterators, Literator(tokens[4]))
	}

	// Append the rest of the tokens
	iterators = append(iterators, Literator(tokens[5:]...))

	return Concaterator(iterators...)
}
