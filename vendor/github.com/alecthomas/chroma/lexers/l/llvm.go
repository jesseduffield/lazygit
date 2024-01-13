package l

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Llvm lexer.
var Llvm = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "LLVM",
		Aliases:   []string{"llvm"},
		Filenames: []string{"*.ll"},
		MimeTypes: []string{"text/x-llvm"},
	},
	llvmRules,
))

func llvmRules() Rules {
	return Rules{
		"root": {
			Include("whitespace"),
			{`([-a-zA-Z$._][\w\-$.]*|"[^"]*?")\s*:`, NameLabel, nil},
			Include("keyword"),
			{`%([-a-zA-Z$._][\w\-$.]*|"[^"]*?")`, NameVariable, nil},
			{`@([-a-zA-Z$._][\w\-$.]*|"[^"]*?")`, NameVariableGlobal, nil},
			{`%\d+`, NameVariableAnonymous, nil},
			{`@\d+`, NameVariableGlobal, nil},
			{`#\d+`, NameVariableGlobal, nil},
			{`!([-a-zA-Z$._][\w\-$.]*|"[^"]*?")`, NameVariable, nil},
			{`!\d+`, NameVariableAnonymous, nil},
			{`c?"[^"]*?"`, LiteralString, nil},
			{`0[xX][a-fA-F0-9]+`, LiteralNumber, nil},
			{`-?\d+(?:[.]\d+)?(?:[eE][-+]?\d+(?:[.]\d+)?)?`, LiteralNumber, nil},
			{`[=<>{}\[\]()*.,!]|x\b`, Punctuation, nil},
		},
		"whitespace": {
			{`(\n|\s)+`, Text, nil},
			{`;.*?\n`, Comment, nil},
		},
		"keyword": {
			{Words(``, `\b`, `begin`, `end`, `true`, `false`, `declare`, `define`, `global`, `constant`, `private`, `linker_private`, `internal`, `available_externally`, `linkonce`, `linkonce_odr`, `weak`, `weak_odr`, `appending`, `dllimport`, `dllexport`, `common`, `default`, `hidden`, `protected`, `extern_weak`, `external`, `thread_local`, `zeroinitializer`, `undef`, `null`, `to`, `tail`, `target`, `triple`, `datalayout`, `volatile`, `nuw`, `nsw`, `nnan`, `ninf`, `nsz`, `arcp`, `fast`, `exact`, `inbounds`, `align`, `addrspace`, `section`, `alias`, `module`, `asm`, `sideeffect`, `gc`, `dbg`, `linker_private_weak`, `attributes`, `blockaddress`, `initialexec`, `localdynamic`, `localexec`, `prefix`, `unnamed_addr`, `ccc`, `fastcc`, `coldcc`, `x86_stdcallcc`, `x86_fastcallcc`, `arm_apcscc`, `arm_aapcscc`, `arm_aapcs_vfpcc`, `ptx_device`, `ptx_kernel`, `intel_ocl_bicc`, `msp430_intrcc`, `spir_func`, `spir_kernel`, `x86_64_sysvcc`, `x86_64_win64cc`, `x86_thiscallcc`, `cc`, `c`, `signext`, `zeroext`, `inreg`, `sret`, `nounwind`, `noreturn`, `noalias`, `nocapture`, `byval`, `nest`, `readnone`, `readonly`, `inlinehint`, `noinline`, `alwaysinline`, `optsize`, `ssp`, `sspreq`, `noredzone`, `noimplicitfloat`, `naked`, `builtin`, `cold`, `nobuiltin`, `noduplicate`, `nonlazybind`, `optnone`, `returns_twice`, `sanitize_address`, `sanitize_memory`, `sanitize_thread`, `sspstrong`, `uwtable`, `returned`, `type`, `opaque`, `eq`, `ne`, `slt`, `sgt`, `sle`, `sge`, `ult`, `ugt`, `ule`, `uge`, `oeq`, `one`, `olt`, `ogt`, `ole`, `oge`, `ord`, `uno`, `ueq`, `une`, `x`, `acq_rel`, `acquire`, `alignstack`, `atomic`, `catch`, `cleanup`, `filter`, `inteldialect`, `max`, `min`, `monotonic`, `nand`, `personality`, `release`, `seq_cst`, `singlethread`, `umax`, `umin`, `unordered`, `xchg`, `add`, `fadd`, `sub`, `fsub`, `mul`, `fmul`, `udiv`, `sdiv`, `fdiv`, `urem`, `srem`, `frem`, `shl`, `lshr`, `ashr`, `and`, `or`, `xor`, `icmp`, `fcmp`, `phi`, `call`, `trunc`, `zext`, `sext`, `fptrunc`, `fpext`, `uitofp`, `sitofp`, `fptoui`, `fptosi`, `inttoptr`, `ptrtoint`, `bitcast`, `addrspacecast`, `select`, `va_arg`, `ret`, `br`, `switch`, `invoke`, `unwind`, `unreachable`, `indirectbr`, `landingpad`, `resume`, `malloc`, `alloca`, `free`, `load`, `store`, `getelementptr`, `extractelement`, `insertelement`, `shufflevector`, `getresult`, `extractvalue`, `insertvalue`, `atomicrmw`, `cmpxchg`, `fence`, `allocsize`, `amdgpu_cs`, `amdgpu_gs`, `amdgpu_kernel`, `amdgpu_ps`, `amdgpu_vs`, `any`, `anyregcc`, `argmemonly`, `avr_intrcc`, `avr_signalcc`, `caller`, `catchpad`, `catchret`, `catchswitch`, `cleanuppad`, `cleanupret`, `comdat`, `convergent`, `cxx_fast_tlscc`, `deplibs`, `dereferenceable`, `dereferenceable_or_null`, `distinct`, `exactmatch`, `externally_initialized`, `from`, `ghccc`, `hhvm_ccc`, `hhvmcc`, `ifunc`, `inaccessiblemem_or_argmemonly`, `inaccessiblememonly`, `inalloca`, `jumptable`, `largest`, `local_unnamed_addr`, `minsize`, `musttail`, `noduplicates`, `none`, `nonnull`, `norecurse`, `notail`, `preserve_allcc`, `preserve_mostcc`, `prologue`, `safestack`, `samesize`, `source_filename`, `swiftcc`, `swifterror`, `swiftself`, `webkit_jscc`, `within`, `writeonly`, `x86_intrcc`, `x86_vectorcallcc`), Keyword, nil},
			{Words(``, ``, `void`, `half`, `float`, `double`, `x86_fp80`, `fp128`, `ppc_fp128`, `label`, `metadata`, `token`), KeywordType, nil},
			{`i[1-9]\d*`, Keyword, nil},
		},
	}
}
