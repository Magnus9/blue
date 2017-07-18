
package token

type Token struct {
    Str       string
    Line      string
    TokenType int
    LineNum   int
}

func New(str string, line string,
         tokenType, lineNum int) Token {
    return Token{
        Str       : str,
        Line      : line,
        TokenType : tokenType,
        LineNum   : lineNum,
    }
}