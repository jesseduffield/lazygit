package utils

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/gui/style"
)

func TestDecolorise(t *testing.T) {
	tests := []struct {
		input  string
		output string
	}{
		{
			input:  "",
			output: "",
		},
		{
			input:  "hello",
			output: "hello",
		},
		{
			input:  "hello\x1b[31m",
			output: "hello",
		},
		{
			input:  "hello\x1b[31mworld",
			output: "helloworld",
		},
		{
			input:  "hello\x1b[31m\x1b[32mworld",
			output: "helloworld",
		},
		{
			input:  "hello\x1b[31m\x1b[32m\x1b[33mworld",
			output: "helloworld",
		},
		{
			input:  "hello\x1b[31m\x1b[32m\x1b[33m\x1b[34mworld",
			output: "helloworld",
		},
		{
			input:  "hello\x1b[31m\x1b[32m\x1b[33m\x1b[34m\x1b[35mworld",
			output: "helloworld",
		},
		{
			input:  "hello\x1b[31m\x1b[32m\x1b[33m\x1b[34m\x1b[35m\x1b[36mworld",
			output: "helloworld",
		},
		{
			input:  "hello\x1b[31m\x1b[32m\x1b[33m\x1b[34m\x1b[35m\x1b[36m\x1b[37mworld",
			output: "helloworld",
		},
		{
			input:  "hello\x1b[31m\x1b[32m\x1b[33m\x1b[34m\x1b[35m\x1b[36m\x1b[37mworld",
			output: "helloworld",
		},
		{
			input:  "\x1b[38;2;47;228;2mJD\x1b[0m",
			output: "JD",
		},
		{
			input:  "\x1b[38;2;160;47;213mRy\x1b[0m",
			output: "Ry",
		},
		{
			input:  "\x1b[38;2;179;217;72mSB\x1b[0m",
			output: "SB",
		},
		{
			input:  "\x1b[38;2;48;34;214mMK\x1b[0m",
			output: "MK",
		},
		{
			input:  "\x1b[38;2;28;152;222mAŁ\x1b[0m",
			output: "AŁ",
		},
		{
			input:  "\x1b[38;2;237;230;56mHH\x1b[0m",
			output: "HH",
		},
		{
			input:  "\x1b[38;2;63;232;69mmj\x1b[0m",
			output: "mj",
		},
		{
			input:  "\x1b[38;2;111;207;16mbl\x1b[0m",
			output: "bl",
		},
		{
			input:  "\x1b[38;2;250;31;163msa\x1b[0m",
			output: "sa",
		},
		{
			input:  "\x1b[38;2;195;10;54mbt\x1b[0m",
			output: "bt",
		},
		{
			input:  "\x1b[38;2;232;147;68mco\x1b[0m",
			output: "co",
		},
		{
			input:  "\x1b[38;2;116;180;35mDY\x1b[0m",
			output: "DY",
		},
		{
			input:  "\x1b[38;2;232;1;195mDB\x1b[0m",
			output: "DB",
		},
		{
			input:  "\x1b[38;2;245;101;55mLi\x1b[0m",
			output: "Li",
		},
		{
			input:  "\x1b[38;2;47;4;217mRy\x1b[0m",
			output: "Ry",
		},
		{
			input:  "\x1b[38;2;252;197;1mEl\x1b[0m",
			output: "El",
		},
		{
			input:  "\x1b[38;2;41;131;237mMG\x1b[0m",
			output: "MG",
		},
		{
			input:  "\x1b[38;2;65;240;62mDP\x1b[0m",
			output: "DP",
		},
		{
			input:  "\x1b[38;2;29;201;139mFM\x1b[0m",
			output: "FM",
		},
		{
			input:  "\x1b[38;2;141;20;198mEB\x1b[0m",
			output: "EB",
		},
		{
			input:  "\x1b[38;2;60;215;140mDM\x1b[0m",
			output: "DM",
		},
		{
			input:  "\x1b[38;2;247;63;38mDE\x1b[0m",
			output: "DE",
		},
		{
			input:  "\x1b[38;2;67;210;17mCB\x1b[0m",
			output: "CB",
		},
		{
			input:  "\x1b[38;2;220;190;84mST\x1b[0m",
			output: "ST",
		},
		{
			input:  "\x1b[38;2;137;239;6mER\x1b[0m",
			output: "ER",
		},
		{
			input:  "\x1b[38;2;47;249;225mAY\x1b[0m",
			output: "AY",
		},
		{
			input:  "\x1b[38;2;215;16;195mca\x1b[0m",
			output: "ca",
		},
		{
			input:  "\x1b[38;2;73;215;122mRV\x1b[0m",
			output: "RV",
		},
		{
			input:  "\x1b[38;2;118;15;221mJP\x1b[0m",
			output: "JP",
		},
		{
			input:  "\x1b[38;2;186;163;39mHJ\x1b[0m",
			output: "HJ",
		},
		{
			input:  "\x1b[38;2;54;222;111mDD\x1b[0m",
			output: "DD",
		},
		{
			input:  "\x1b[38;2;56;209;108mPZ\x1b[0m",
			output: "PZ",
		},
		{
			input:  "\x1b[38;2;9;179;216mPM\x1b[0m",
			output: "PM",
		},
		{
			input:  "\x1b[38;2;157;205;18mta\x1b[0m",
			output: "ta",
		},
		{
			input:  "a_" + style.PrintSimpleHyperlink("xyz") + "_b",
			output: "a_xyz_b",
		},
	}

	for _, test := range tests {
		output := Decolorise(test.input)
		if output != test.output {
			t.Errorf("Decolorise(%s) = %s, want %s", test.input, output, test.output)
		}
	}
}
