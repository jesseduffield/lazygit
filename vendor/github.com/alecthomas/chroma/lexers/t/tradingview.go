package t

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// TradingView lexer
var TradingView = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "TradingView",
		Aliases:   []string{"tradingview", "tv"},
		Filenames: []string{"*.tv"},
		MimeTypes: []string{"text/x-tradingview"},
		DotAll:    true,
		EnsureNL:  true,
	},
	tradingViewRules,
))

func tradingViewRules() Rules {
	return Rules{
		"root": {
			{`[^\S\n]+|\n|[()]`, Text, nil},
			{`(//.*?)(\n)`, ByGroups(CommentSingle, Text), nil},
			{`>=|<=|==|!=|>|<|\?|-|\+|\*|\/|%|\[|\]`, Operator, nil},
			{`[:,.]`, Punctuation, nil},
			{`=`, KeywordPseudo, nil},
			{`"(\\\\|\\"|[^"\n])*["\n]`, LiteralString, nil},
			{`'\\.'|'[^\\]'`, LiteralString, nil},
			{`[0-9](\.[0-9]*)?([eE][+-][0-9]+)?`, LiteralNumber, nil},
			{`#[a-fA-F0-9]{8}|#[a-fA-F0-9]{6}|#[a-fA-F0-9]{3}`, LiteralStringOther, nil},
			{`(abs|acos|alertcondition|alma|asin|atan|atr|avg|barcolor|barssince|bgcolor|cci|ceil|change|cog|color\.new|correlation|cos|crossover|crossunder|cum|dev|ema|exp|falling|fill|fixnan|floor|heikinashi|highest|highestbars|hline|iff|kagi|label\.(delete|get_text|get_x|get_y|new|set_color|set_size|set_style|set_text|set_textcolor|set_x|set_xloc|set_xy|set_y|set_yloc)|line\.(new|delete|get_x1|get_x2|get_y1|get_y2|set_color|set_width|set_style|set_extend|set_xy1|set_xy2|set_x1|set_x2|set_y1|set_y2|set_xloc)|linebreak|linreg|log|log10|lowest|lowestbars|macd|max|max_bars_back|min|mom|nz|percentile_(linear_interpolation|nearest_rank)|percentrank|pivothigh|pivotlow|plot|plotarrow|plotbar|plotcandle|plotchar|plotshape|pointfigure|pow|renko|rising|rma|roc|round|rsi|sar|security|sign|sin|sma|sqrt|stdev|stoch|study|sum|swma|tan|timestamp|tostring|tsi|valuewhen|variance|vwma|wma|strategy\.(cancel|cancel_all|close|close_all|entry|exit|order|risk\.(allow_entry_in|max_cons_loss_days|max_drawdown|max_intraday_filled_orders|max_intraday_loss|max_position_size)))\b`, NameFunction, nil},
			{`\b(bool|color|cross|dayofmonth|dayofweek|float|hour|input|int|label|line|minute|month|na|offset|second|strategy|string|tickerid|time|tr|vwap|weekofyear|year)(\()`, ByGroups(NameFunction, Text), nil}, // functions that can also be variable
			{`(accdist|adjustment\.(dividends|none|splits)|aqua|area|areabr|bar_index|black|blue|bool|circles|close|columns|currency\.(AUD|CAD|CHF|EUR|GBP|HKD|JPY|NOK|NONE|NZD|RUB|SEK|SGD|TRY|USD|ZAR)|color\.(aqua|black|blue|fuchsia|gray|green|lime|maroon|navy|olive|orange|purple|red|silver|teal|white|yellow)|dashed|dotted|dayofweek\.(monday|tuesday|wednesday|thursday|friday|saturday|sunday)|extend\.(both|left|right|none)|float|format\.(inherit|price|volume)|friday|fuchsia|gray|green|high|histogram|hl2|hlc3|hline\.style_(dotted|solid|dashed)|input\.(bool|float|integer|resolution|session|source|string|symbol)|integer|interval|isdaily|isdwm|isintraday|ismonthly|isweekly|label\.style_(arrowdown|arrowup|circle|cross|diamond|flag|labeldown|labelup|none|square|triangledown|triangleup|xcross)|lime|line\.style_(dashed|dotted|solid|arrow_both|arrow_left|arrow_right)|linebr|location\.(abovebar|absolute|belowbar|bottom|top)|low|maroon|monday|n|navy|ohlc4|olive|open|orange|period|plot\.style_(area|areabr|circles|columns|cross|histogram|line|linebr|stepline)|purple|red|resolution|saturday|scale\.(left|none|right)|session|session\.(extended|regular)|silver|size\.(auto|huge|large|normal|small|tiny)|solid|source|stepline|string|sunday|symbol|syminfo\.(mintick|pointvalue|prefix|root|session|ticker|tickerid|timezone)|teal|thursday|ticker|timeframe\.(isdaily|isdwm|isintraday|ismonthly|isweekly|multiplier|period)|timenow|tuesday|volume|wednesday|white|yellow|strategy\.(cash|closedtrades|commission\.(cash_per_contract|cash_per_order|percent)|direction\.(all|long|short)|equity|eventrades|fixed|grossloss|grossprofit|initial_capital|long|losstrades|max_contracts_held_(all|long|short)|max_drawdown|netprofit|oca\.(cancel|none|reduce)|openprofit|opentrades|percent_of_equity|position_avg_price|position_entry_name|position_size|short|wintrades)|shape\.(arrowdown|arrowup|circle|cross|diamond|flag|labeldown|labelup|square|triangledown|triangleup|xcross)|barstate\.is(first|history|last|new|realtime)|barmerge\.(gaps_on|gaps_off|lookahead_on|lookahead_off)|xloc\.bar_(index|time)|yloc\.(abovebar|belowbar|price))\b`, NameVariable, nil},
			{`(cross|dayofmonth|dayofweek|hour|minute|month|na|second|tickerid|time|tr|vwap|weekofyear|year)(\b[^\(])`, ByGroups(NameVariable, Text), nil}, // variables that can also be function
			{`(int|float|bool|color|string|label|line)(\b[^\(=.])`, ByGroups(KeywordType, Text), nil},                                                      // types that can also be a function
			{`(var)\b`, KeywordType, nil},
			{`(true|false)\b`, KeywordConstant, nil},
			{`(and|or|not|if|else|for|to)\b`, OperatorWord, nil},
			{`@?[_a-zA-Z]\w*`, Text, nil},
		},
	}
}
