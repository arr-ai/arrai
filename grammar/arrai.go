package grammar

//nolint:unused,varcheck,deadcode
var arraiGrammar = `
expr       -> "&"* arrow; // TODO: What's with NewFunction("-", expr) in the Go source?

arrow      -> (nest | unnest)* binop_tail arrow_tail;
arrow_tail -> nest* (unnest arrow)?;
nest       -> "nest" name_list ident;
unnest     -> "unnest" ident;

expr       -> expr:("with"|"without")
            ^ expr:("||")
            ^ expr:("&&")
            ^ expr:/(!?(?:<>?=?|>=?))/
            ^ "if" expr2 "then" expr "else" expr;
expr2      -> expr2:(/([&])/ | JOIN)
            ^ expr2:/([-+|])/
            ^ expr2:/([\*/%]|-%|//)/
            ^ expr2<:"**"
            ^ /([-+!]|\*\*?)/ touch;
touch      -> suffix touch_tail;
touch_tail -> ("->*" )
suffix     -> 
dot        -> atom dot_tail ("." (("&" IDENT) | STRING | "*"))*;
atom       -> ""
			| "(" (IDENT? ":" expr):","! ")"
			| "(" expr ")"
			| "{" expr:","! "}"
			| array
			| xml
			| "." "*"?
			| "\\\\" expr
			| "\\" ident expr
			| NUMBER | IDENT | STRING;
attr_expr  -> "(" attr_term:"," ","? ")"
name_list  -> "|" ident:"," "|"

NUMBER     -> /((?:[0-9]+(?:\.[0-9]*)?|\.[0-9]+)(?:[Ee][-+]?[0-9]+)?)/;
IDENT      -> /([$@A-Za-z_][0-9$@A-Za-z_]*)/;
STRING     -> /("(?:\\.|[^\\"])*")/
JOIN       -> /[-<][-&][->]/;
`
