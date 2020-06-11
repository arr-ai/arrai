// AUTOGENERATED. DO NOT EDIT.
package syntax

import (
	"strings"

	"github.com/arr-ai/wbnf/wbnf"
)

func unfakeBackquote(s string) string {
	return strings.ReplaceAll(s, "‵", "`")
}

var arraiParsers = wbnf.MustCompile(unfakeBackquote(`
default -> expr;
expr   -> C* amp="&"* @ C* arrow=(
              nest |
              unnest |
              ARROW @ |
              binding="->" C* "\\" C* pattern C* %%bind C* @ |
              binding="->" C* %%bind @
          )* C*
        > C* unop=/{:>|=>|>>}* @ C*
        > C* @:binop=("without" | "with") C*
        > C* @:binop="||" C*
        > C* @:binop="&&" C*
        > C* @:compare=/{!?(?:<:|=|<=?|>=?|\((?:<=?|>=?|<>=?)\))} C*
        > C* @ if=("if" t=expr ("else" f=expr)?)* C*
        > C* @:binop=/{\+\+|[+|]|-%?} C*
        > C* @:binop=/{&~|&|~~?|[-<][-&][->]} C*
        > C* @:binop=/{//|[*/%]|\\} C*
        > C* @:rbinop="^" C*
        > C* unop=/{[-+!*^]}* @ C*
        > C* @:binop=">>>" C*
        > C* @ postfix=/{count|single}? C* touch? C*
        > C* (get | @) tail_op=(
            safe_tail=(first_safe=(tail "?") ops=(safe=(tail "?") | tail)* ":" fall=@)
            | tail
          )* C*
        > %!patternterms(expr)
        | C* cond=("cond" "{" (key=@ ":" value=@):SEQ_COMMENT,? "}") C*
        | C* cond=("cond" controlVar=expr "{" (condition=pattern ":" value=@):SEQ_COMMENT,? "}") C*
        | C* "{:" C* embed=(macro=@ rule? ":" subgrammar=%%ast) ":}" C*
        | C* op="\\\\" @ C*
        | C* fn="\\" IDENT @ C*
        | C* "//" pkg=( "{" dot="."? PKGPATH "}" | std=IDENT?)
        | C* "(" @ ")" C*
        | C* let=("let" C* rec="rec"? pattern C* "=" C* @ %%bind C* ";" C* @) C*
        | C* xstr C*
        | C* IDENT C*
        | C* STR C*
        | C* NUM C*
        | C* CHAR C*;
rule   -> C* "[" C* name C* "]" C*;
nest   -> C* "nest" names IDENT C*;
unnest -> C* "unnest" IDENT C*;
touch  -> C* ("->*" ("&"? IDENT | STR))+ "(" expr:"," ","? ")" C*;
get    -> C* dot="." ("&"? IDENT | STR | "*") C*;
names  -> C* "|" C* IDENT:"," C* "|" C*;
name   -> C* IDENT C* | C* STR C*;
xstr   -> C* quote=/{\$"\s*} part=( sexpr | fragment=/{(?: \\. | \$[^{"] | [^\\"$] )+} )* '"' C*
        | C* quote=/{\$'\s*} part=( sexpr | fragment=/{(?: \\. | \$[^{'] | [^\\'$] )+} )* "'" C*
        | C* quote=/{\$‵\s*} part=( sexpr | fragment=/{(?: ‵‵  | \$[^{‵] | [^‵  $] )+} )* "‵" C*;
sexpr  -> "${"
          C* expr C*
          control=/{ (?: : [-+#*\.\_0-9a-z]* (?: : (?: \\. | [^\\:}] )* ){0,2} )? }
          close=/{\}\s*};
tail   -> get
          | call=("("
                arg=(
                    expr (":" end=expr? (":" step=expr)?)?
                    |     ":" end=expr  (":" step=expr)?
                ):SEQ_COMMENT,
            ")");
pattern -> extra | %!patternterms(pattern|expr) | IDENT | NUM | C* "(" exprpattern=expr:SEQ_COMMENT,? ")" C* | C* exprpattern=STR C*;
extra -> ("..." ident=IDENT?);

ARROW  -> /{:>|=>|>>|orderby|order|rank|where|sum|max|mean|median|min};
IDENT  -> /{ \. | [$@A-Za-z_][0-9$@A-Za-z_]* };
PKGPATH -> /{ (?: \\ | [^\\}] )* };
STR    -> /{ " (?: \\. | [^\\"] )* "
           | ' (?: \\. | [^\\'] )* '
           | ‵ (?: ‵‵  | [^‵  ] )* ‵
           };
NUM    -> /{ (?: \d+(?:\.\d*)? | \.\d+ ) (?: [Ee][-+]?\d+ )? };
CHAR   -> /{%(\\.|.)};
C      -> /{ # .* $ };
SEQ_COMMENT -> "," C*;

.wrapRE -> /{\s*()\s*};

.macro patternterms(top) {
    C* "{" C* rel=(names tuple=("(" v=top:SEQ_COMMENT, ")"):SEQ_COMMENT,?) "}" C*
  | C* "{" C* set=(elt=top:SEQ_COMMENT,?) "}" C*
  | C* "{" C* dict=((ext=extra|key=expr ":" value=top):SEQ_COMMENT,?) "}" C*
  | C* "[" C* array=(%!sparse_sequence(top)?) C* "]" C*
  | C* "<<" C* bytes=(item=top:SEQ_COMMENT,?) C* ">>" C*
  | C* "(" tuple=(pairs=(extra|name? ":" v=top):SEQ_COMMENT,?) ")" C*
  | C* "(" identpattern=IDENT ")" C*
};

.macro sparse_sequence(top) {
  first_item=(top) C* (SEQ_COMMENT:(item=(top|empty=\s*)),)?
}

`), nil)
