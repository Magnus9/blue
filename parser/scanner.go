
package parser

import (
    "fmt"
    "bytes"
    "github.com/Magnus9/blue/token"
)

const EOF = 0

type Scanner struct {
    sourceProgram string
    pathname      string
    lineBuf       bytes.Buffer
    charPointer   byte
    sourcePos     int
    lineNum       int
}

func newScanner(sourceProgram, pathname string) Scanner {
    scanner := Scanner{
        sourceProgram: sourceProgram,
        pathname     : pathname,
        sourcePos    : -1,
        lineNum      :  1,
    }
    // Read the first line / until EOF.
    scanner.readLine()
    // Call nextChar() to initialize charPointer.
    scanner.nextChar()

    return scanner
}

func (s *Scanner) postError(message string) {
    buf := fmt.Sprintf("%s:%d => %s\n   %s", s.pathname,
                       s.lineNum, message,
                       s.lineBuf.String())
    panic(buf)
}

func (s *Scanner) makeSymToken(str string,
                               ttype int) token.Token {
    switch len(str) {
        case 1:
            s.nextCharx(1)
        case 2:
            s.nextCharx(2)
        case 3:
            s.nextCharx(3)
    }
    return token.New(str, s.lineBuf.String(),
                     ttype, s.lineNum)
}

func (s *Scanner) makeToken(str string,
                            ttype int) token.Token {
    return token.New(str, s.lineBuf.String(),
                     ttype, s.lineNum)
}

func (s *Scanner) readLine() {
    s.lineBuf.Reset()
    var i int
    // First for loop removes leading whitespace.
    for i = s.sourcePos + 1; i < len(s.sourceProgram);
        i++ {
        ch := s.sourceProgram[i]
        if ch != ' ' && ch != '\t' {
            break
        }
    }
    for ; i < len(s.sourceProgram); i++ {
        ch := s.sourceProgram[i]
        if ch == '\n' {
            break
        }
        s.lineBuf.WriteByte(ch)
    }
}

func (s *Scanner) nextChar() byte {
    s.sourcePos++
    if s.sourcePos >= len(s.sourceProgram) {
        s.charPointer = EOF
    } else {
        s.charPointer = s.sourceProgram[s.sourcePos]
    }
    return s.charPointer
}

func (s *Scanner) peekChar(num int) byte {
    newPos := s.sourcePos + num
    if newPos >= len(s.sourceProgram) {
        return EOF
    }
    return s.sourceProgram[newPos]
}

func (s *Scanner) nextCharx(num int) {
    for i := 0; i < num; i++ { s.nextChar() }
}

func (s *Scanner) getSlice(spos int) string {
    return s.sourceProgram[spos:s.sourcePos]
}

func (s *Scanner) nextToken() token.Token {
    var ch byte
    for s.charPointer != EOF {
        if s.charPointer == '=' && s.peekChar(1) == '=' &&
           s.peekChar(2) == '=' {
            s.longComment()
            continue
        }
        switch s.charPointer {
            // List of one-character symbols.
            case '~': return s.makeSymToken("~", token.TILDE)
            case '(': return s.makeSymToken("(", token.LPAREN)
            case ')': return s.makeSymToken(")", token.RPAREN)
            case '[': return s.makeSymToken("[", token.LBRACK)
            case ']': return s.makeSymToken("]", token.RBRACK)
            case '{': return s.makeSymToken("{", token.LBRACE)
            case '}': return s.makeSymToken("}", token.RBRACE)
            case ',': return s.makeSymToken(",", token.COMMA)
            case ';': return s.makeSymToken(";", token.SEMICOLON)
            case ':': return s.makeSymToken("(", token.COLON)
            case '\n':
                token := s.makeToken("NL", token.NEWLINE)
                s.readLine()
                s.nextChar()

                s.lineNum++
                return token
            case ' ':
                for s.charPointer == ' ' || s.charPointer == '\r' ||
                    s.charPointer == '\t' {
                    s.nextChar()
                }
                continue
            case '#':
                for s.charPointer != '\n' && s.charPointer != EOF {
                    s.nextChar()
                }
            case '0':
                nextChar := s.peekChar(1)
                if nextChar == 'X' || nextChar == 'x' {
                    return s.parseHex()
                }
                return s.parseNumber()
            case '"':
                fallthrough
            case '\'':
                return s.parseString()
            case '\\':
                if s.peekChar(1) != '\n' {
                    s.postError("missing newline after line-continuation" +
                                " character")
                }
                s.nextCharx(2)

            // List of two-character symbols.
            case '.':
                if s.peekChar(1) == '.' {
                    return s.makeSymToken("..", token.DOTDOT)
                }
                return s.makeSymToken(".", token.DOT)
            case '+':
                if s.peekChar(1) == '=' {
                    return s.makeSymToken("+=", token.PLUSEQ)
                }
                return s.makeSymToken("+", token.PLUS)
            case '-':
                if s.peekChar(1) == '=' {
                    return s.makeSymToken("-=", token.MINUSEQ)
                }
                return s.makeSymToken("-", token.MINUS)
            case '*':
                if s.peekChar(1) == '=' {
                    return s.makeSymToken("*=", token.STAREQ)
                }
                return s.makeSymToken("*", token.STAR)
            case '/':
                if s.peekChar(1) == '=' {
                    return s.makeSymToken("/=", token.SLASHEQ)
                }
                return s.makeSymToken("/", token.SLASH)
            case '%':
                if s.peekChar(1) == '=' {
                    return s.makeSymToken("%=", token.PERCENTEQ)
                }
                return s.makeSymToken("%", token.PERCENT)
            case '!':
                if s.peekChar(1) == '=' {
                    return s.makeSymToken("!=", token.BANGEQ)
                }
                return s.makeSymToken("!", token.BANG)
            case '=':
                ch = s.peekChar(1)
                if ch == '=' {
                    return s.makeSymToken("==", token.EQEQ)
                } else if ch == '>' {
                    return s.makeSymToken("=>", token.EQGT)
                }
                return s.makeSymToken("=", token.EQ)
            case '|':
                ch = s.peekChar(1)
                if ch == '|' {
                    return s.makeSymToken("||", token.PIPEPIPE)
                } else if ch == '=' {
                    return s.makeSymToken("|=", token.PIPEEQ)
                }
                return s.makeSymToken("|", token.PIPE)
            case '&':
                ch = s.peekChar(1)
                if ch == '&' {
                    return s.makeSymToken("&&", token.AMPAMP)
                } else if ch == '=' {
                    return s.makeSymToken("&=", token.AMPEQ)
                }
                return s.makeSymToken("&", token.AMP)
            case '^':
                if s.peekChar(1) == '=' {
                    return s.makeSymToken("^=", token.CARETEQ)
                }
                return s.makeSymToken("^", token.CARET)

            // List of three-character symbols.
            case '<':
                ch = s.peekChar(1)
                if ch == '<' {
                    if s.peekChar(2) == '=' {
                        return s.makeSymToken("<<=", token.LEFTSHIFTEQ)
                    }
                    return s.makeSymToken("<<", token.LEFTSHIFT)
                } else if ch == '=' {
                    return s.makeSymToken("<=", token.LTEQ)
                }
                return s.makeSymToken("<", token.LT)
            case '>':
                ch = s.peekChar(1)
                if ch == '>' {
                    if s.peekChar(2) == '=' {
                        return s.makeSymToken(">>=", token.RIGHTSHIFTEQ)
                    }
                    return s.makeSymToken(">>", token.RIGHTSHIFT)
                } else if ch == '=' {
                    return s.makeSymToken(">=", token.GTEQ)
                }
                return s.makeSymToken(">", token.GT)
            
            default:
                if s.isLetter() || s.charPointer == '_' {
                    return s.parseWord()
                } else if s.isDigit() {
                    return s.parseNumber()
                } else {
                    s.postError("unrecognized character '" +
                                 string(s.charPointer) + "'")
                }

        }
    }
    return s.makeSymToken("EOF", token.EOF)
}

func (s *Scanner) isLetter() bool {
    return (s.charPointer >= 'A' && s.charPointer <= 'Z' ||
            s.charPointer >= 'a' && s.charPointer <= 'z')
}

func (s *Scanner) isDigit() bool {
    return (s.charPointer >= '0' && s.charPointer <= '9')
}

func (s *Scanner) isxDigit(ch byte) bool {
    if ch >= '0' && ch <= '9' {
        return true
    }
    if ch >= 'A' && ch <= 'F' {
        return true
    }
    if ch >= 'a' && ch <= 'f' {
        return true
    }
    return false
}

func (s *Scanner) parseWord() token.Token {
    ch := s.peekChar(1)
    if (s.charPointer == 'R' || s.charPointer == 'r') &&
        (ch == '"' || ch == '\'') {
        return s.parseString()
    }
    pos := s.sourcePos
    for s.isLetter() || s.isDigit() || s.charPointer == '_' {
        s.nextChar()
    }
    buf := s.getSlice(pos)
    v, ok := token.RES_WORDS[buf]

    if (ok) {
        return s.makeToken(buf, v)
    }
    return s.makeToken(buf, token.NAME)
}

func (s *Scanner) parseNumber() token.Token {
    pos := s.sourcePos

    for (s.isDigit()) {
        s.nextChar()
    }
    if s.charPointer == '.' && s.peekChar(1) != '.' {
        s.nextChar()
        for (s.isDigit()) {
            s.nextChar()
        }
        return s.makeToken(s.getSlice(pos), token.FLOAT)
    }
    return s.makeToken(s.getSlice(pos), token.INTEGER)
}

func (s *Scanner) parseHex() token.Token {
    pos := s.sourcePos

    s.nextCharx(2)
    for (s.isxDigit(s.charPointer)) {
        s.nextChar()
    }
    return s.makeToken(s.getSlice(pos), token.INTEGER)
}

func (s *Scanner) parseString() token.Token {
    pos := s.sourcePos

    if s.charPointer == 'R' || s.charPointer == 'r' {
        s.nextChar()
    }
    quote := s.charPointer
    s.nextChar()

    for (s.charPointer != quote && s.charPointer != EOF) {
        if s.charPointer == '\\' {
            s.nextChar()
        }
        s.nextChar()
    }
    if s.charPointer == EOF {
        s.postError("unterminated string literal")
    }
    s.nextChar()

    return s.makeToken(s.getSlice(pos), token.STRING)
}

func (s *Scanner) longComment() {
    s.nextCharx(3)
    for s.charPointer != EOF {
        if s.charPointer == '=' {
            if s.peekChar(1) == '=' && s.peekChar(2) == '=' {
                break
            }
        }
        if s.charPointer == '\n' {
            s.lineNum++
        }
        s.nextChar()
    }
    if s.charPointer == EOF {
        s.postError("unterminated long comment")
    }
    s.nextCharx(3)
}