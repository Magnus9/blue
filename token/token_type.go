
package token

const (
    // LITERALS
    STRING = iota; INTEGER; FLOAT; TRUE; FALSE
    NIL

    // RESERVED WORDS
    DEF; IF; ELIF; ELSE; DO; END; FOR; WHILE
    SWITCH; CASE; DEFAULT; IN; RETURN; THEN
    PRINT; CONTINUE; BREAK; IMPORT; FROM; CLASS;
    EXTENDS; NEW

    // NON ASSIGNING SYMBOLS
    LT; LTEQ; GT; GTEQ; LEFTSHIFT; RIGHTSHIFT; DOT
    DOTDOT; PLUS; MINUS; STAR; SLASH; PERCENT
    BANG; TILDE; LPAREN; RPAREN; LBRACK; RBRACK
    LBRACE; RBRACE; COMMA; SEMICOLON; EQEQ
    BANGEQ; PIPE; PIPEPIPE; AMP; AMPAMP; CARET; EQGT
    NEWLINE; COLON
    
    // ASSIGNING SYMBOLS
    EQ; PIPEEQ; CARETEQ; AMPEQ; LEFTSHIFTEQ
    RIGHTSHIFTEQ; PLUSEQ; MINUSEQ; STAREQ; SLASHEQ
    PERCENTEQ

    // NAME; EOF
    NAME; EOF

    // IMAGINARY TOKENS
    BLOCK; LIST; HASH; HASH_ELEM; CALL; MAKE_CLASS
    SUBSCRIPT; NEGATE; AUGASSIGN; COMP_OP; MAKE_FUNC
    MAKE_INSTANCE; PATH; SLICE

    CLASSBLOCK; PARAMETERS; ARGUMENTS; LE; GE; MEMBER
    RANGE; ADD; SUB; MUL; DIV; MODULO; COMPL; ASSIGN
    NE; LOGICAL_OR; LOGICAL_AND; BITWISE_OR; BITWISE_AND
    XOR; NOT

    ASS_BITWISE_OR; ASS_BITWISE_AND; ASS_XOR; ASS_LEFTSHIFT
    ASS_RIGHTSHIFT; ASS_ADD; ASS_SUB; ASS_MUL; ASS_DIV
    ASS_MODULO

    // INPUT FROM FILE, INTERACTIVE.
    FILE_INPUT; INTERACTIVE
)

var RES_WORDS = map[string]int{
    "true"    : TRUE,
    "false"   : FALSE,
    "nil"     : NIL,
    "def"     : DEF,
    "if"      : IF,
    "elif"    : ELIF,
    "else"    : ELSE,
    "do"      : DO,
    "then"    : THEN,
    "end"     : END,
    "for"     : FOR,
    "while"   : WHILE,
    "switch"  : SWITCH,
    "case"    : CASE,
    "default" : DEFAULT,
    "in"      : IN,
    "return"  : RETURN,
    "print"   : PRINT,
    "continue": CONTINUE,
    "break"   : BREAK,
    "import"  : IMPORT,
    "from"    : FROM,
    "class"   : CLASS,
    "extends" : EXTENDS,
    "new"     : NEW,
}