
package parser

import "os"
import "fmt"
import "bytes"
import "github.com/Magnus9/blue/token"
import "github.com/Magnus9/blue/interm"

type Parser struct {
    scanner  Scanner
    current  token.Token
    next     token.Token
    pathname string
}

func readFp(fp *os.File) string {
    stat, err := fp.Stat()
    if err != nil {
        panic(err)
    }
    buf := make([]byte, stat.Size())
    fp.Read(buf)

    return string(buf)
}

func ParseFromFile(pathname string, fp *os.File) *interm.Node {
    scanner := newScanner(readFp(fp), pathname)
    p := Parser{
        scanner : scanner,
        current : scanner.nextToken(),
        next    : scanner.nextToken(),
        pathname: pathname,
    }
    root := interm.New("FILE_INPUT", "", token.FILE_INPUT,
                       0)
    return p.Program(root)
}

func ParseFromRepl(pathname, program string) *interm.Node {
    scanner := newScanner(program, pathname)
    p := Parser{
        scanner : scanner,
        current : scanner.nextToken(),
        next    : scanner.nextToken(),
        pathname: pathname,
    }
    root := interm.New("INTERACTIVE", "", token.INTERACTIVE,
                       0)
    return p.Program(root)
}

func (p *Parser) createNode(str string,
                            nodeType int) *interm.Node {
    return interm.New(str, p.current.Line, nodeType,
                      p.current.LineNum)
}

func (p *Parser) postError(message string) {
    var buf bytes.Buffer
    buf.WriteString(fmt.Sprintf("%s:%d => ", p.pathname,
                    p.current.LineNum))

    tokenType := p.current.TokenType
    if tokenType == token.NEWLINE {
        buf.WriteString("unexpected newline, ")
    } else if tokenType == token.EOF {
        buf.WriteString("unexpected end-of-file, ")
    } else if tokenType == token.NAME {
        buf.WriteString("unexpected name near '" +
                        p.current.Str + "', ")
    } else if tokenType >= token.STRING &&
              tokenType <= token.NIL {
        buf.WriteString("unexpected literal near '" +
                        p.current.Str + "', ")
    } else if tokenType >= token.DEF &&
              tokenType <= token.NEW {
        buf.WriteString("unexpected keyword near '" +
                        p.current.Str + ", ")
    } else {
        buf.WriteString("unexpected symbol near '" +
                        p.current.Str + "', ")
    }
    buf.WriteString(message)
    buf.WriteString("\n   ")
    buf.WriteString(p.current.Line)

    panic(buf.String())
}

func (p *Parser) matchToken(TokenType int,
                            message string) {
    if p.peekCurrent() != TokenType {
        p.postError(message)
    }
    p.nextToken()
}

func (p *Parser) nextToken() {
    p.current = p.next
    p.next = p.scanner.nextToken()
}

func (p *Parser) peekCurrent() int {
    return p.current.TokenType
}

func (p *Parser) peekNext() int {
    return p.next.TokenType
}

func (p *Parser) skipNL() {
    for p.peekCurrent() == token.NEWLINE {
        p.nextToken()
    }
}

func (p *Parser) nextAndSkipNL() {
    p.nextToken()
    p.skipNL()
}

func (p *Parser) matchNewline(message string) {
    p.matchToken(token.NEWLINE, message)
    p.skipNL()
}

func (p *Parser) isFactor() bool {
    tokenType := p.peekCurrent()

    return (tokenType == token.MINUS ||
            tokenType == token.BANG  ||
            tokenType == token.TILDE)
}

func (p *Parser) stmtTrailer() {
    tokenType := p.peekCurrent()
    if tokenType == token.SEMICOLON ||
       tokenType == token.NEWLINE {
        p.nextAndSkipNL()
    } else {
        p.matchToken(token.EOF, "expected end-of-file")
    }
}

func (p *Parser) Program(node *interm.Node) *interm.Node {
    p.skipNL()
    tokenType := p.peekCurrent()

    for tokenType != token.EOF {
        switch tokenType {
            case token.CLASS:
                node.Add(p.classStmt())
            case token.DEF:
                node.Add(p.defStmt())
            default:
                node.Add(p.stmt())
        }
        p.stmtTrailer()
        tokenType = p.peekCurrent()
    }
    return node
}

func (p *Parser) stmt() *interm.Node {
    switch p.peekCurrent() {
        case token.WHILE:
            return p.controlStmt()
        case token.FOR:
            return p.forStmt()
        case token.CONTINUE:
            fallthrough
        case token.BREAK:
            node := p.createNode(p.current.Str, p.current.TokenType)
            p.nextToken()

            return node
        case token.IF:
            return p.ifStmt()
        case token.RETURN:
            return p.returnStmt()
        case token.IMPORT:
            node := p.createNode(p.current.Str, p.current.TokenType)
            for {
                node.Add(p.importPath())
                if p.peekCurrent() != token.COMMA {
                    break
                }
            }
            return node
        case token.PRINT:
            node := p.createNode(p.current.Str, p.current.TokenType)
            p.nextToken()
            node.Add(p.expr())

            return node
        default:
            return p.exprStmt()
    }
}

func (p *Parser) classStmt() *interm.Node {
    root := p.createNode("MAKE_CLASS", token.MAKE_CLASS)
    
    p.nextToken()
    if p.peekCurrent() != token.NAME {
        p.postError("expected name after 'class'")
    }
    nameNode := p.createNode(p.current.Str, p.current.TokenType)
    root.Add(nameNode)
    extendsNode := p.createNode("EXTENDS", token.EXTENDS)
    root.Add(extendsNode)

    p.nextToken()
    if p.peekCurrent() == token.COLON {
        p.baseClass(extendsNode)
    }
    root.Add(p.classBlock())

    return root
}

func (p *Parser) baseClass(node *interm.Node) {
    p.nextToken()
    if p.peekCurrent() != token.NAME {
        p.postError("expected name after ':'")
    }
    nameNode := p.createNode(p.current.Str, p.current.TokenType)
    node.Add(nameNode)
    p.nextToken()
}

func (p *Parser) classBlock() *interm.Node {
    p.matchNewline("expected newline to open class")
    root := p.createNode("CLASSBLOCK", token.CLASSBLOCK)

    tokenType := p.peekCurrent()
    for tokenType != token.END && tokenType != token.EOF {
        if tokenType == token.DEF {
            root.Add(p.defStmt())
        } else {
            root.Add(p.stmt())
        }
        tokenType = p.peekCurrent()
        if tokenType == token.SEMICOLON {
            p.nextAndSkipNL()
        } else {
            p.matchNewline("expected newline")
        }
        tokenType = p.peekCurrent()
    }
    p.matchToken(token.END, "expected 'end' to close class")

    return root
}

func (p *Parser) defStmt() *interm.Node {
    root := p.createNode("MAKE_FUNC", token.MAKE_FUNC)

    p.nextToken()
    if p.peekCurrent() != token.NAME {
        p.postError("expected name after 'def'")
    }
    nameNode := p.createNode(p.current.Str, p.current.TokenType)
    root.Add(nameNode)
    p.nextToken()

    if p.peekCurrent() != token.LPAREN {
        p.postError("expected '(' to open parameter list")
    }
    paramsNode := p.createNode("PARAMETERS", token.PARAMETERS)
    root.Add(paramsNode)
    if p.peekNext() == token.RPAREN {
        p.nextToken()
        p.nextToken()
    } else {
        p.defParams(paramsNode)
        p.matchToken(token.RPAREN, "expected ')' to close parameter" +
                     " list")
    }
    p.matchNewline("expected newline")
    root.Add(p.stmtBlock())
    p.matchToken(token.END, "expected 'end' to close function")

    return root
}

func (p *Parser) defParams(root *interm.Node) {
    for (true) {
        p.nextAndSkipNL()
        if (root.Flags & interm.FLAG_STARPARAM) != 0 {
            p.postError("star parameter must be the last param")
        }
        if p.peekCurrent() == token.STAR {
            root.Flags |= interm.FLAG_STARPARAM
            p.nextToken()
        }
        if p.peekCurrent() != token.NAME {
            p.postError("expected name as argument")
        }
        nameNode := p.createNode(p.current.Str, p.current.TokenType)
        root.Add(nameNode)
        
        p.nextAndSkipNL()
        if p.peekCurrent() != token.COMMA {
            break
        }
    }
}

/*
 * This subroutine checks if one of the tokens
 * that can terminate a block is the lookahead
 * token.
 */
func (p *Parser) blockFollows() bool {
    tokenType := p.peekCurrent()

    return (tokenType == token.ELIF ||
            tokenType == token.ELSE ||
            tokenType == token.END  ||
            tokenType == token.EOF)
}

func (p *Parser) stmtBlock() *interm.Node {
    root := p.createNode("BLOCK", token.BLOCK)
    p.skipNL()

    for !p.blockFollows() {
        root.Add(p.stmt())
        if p.peekCurrent() == token.SEMICOLON {
            p.nextAndSkipNL()
        } else {
            p.matchNewline("expected newline")
        }
    }
    return root
}

func (p *Parser) controlStmt() *interm.Node {
    root := p.createNode(p.current.Str, p.current.TokenType)
    p.nextToken()

    root.Add(p.expr())
    p.matchToken(token.DO, "expected 'do' to open block")

    root.Add(p.stmtBlock())
    p.matchToken(token.END, "expected 'end' to close block")

    return root
}

func (p *Parser) forStmt() *interm.Node {
    root := p.createNode(p.current.Str, p.current.TokenType)
    p.nextToken()

    argsNode := p.createNode("ARGUMENTS", token.ARGUMENTS)
    for i := 0; i < 2; i++ {
        if p.peekCurrent() != token.NAME {
            p.postError("expected name")
        }
        nameNode := p.createNode(p.current.Str, p.current.TokenType)
        argsNode.Add(nameNode)
        p.nextToken()
        if p.peekCurrent() != token.COMMA {
            break
        }
        p.nextToken()
    }
    root.Add(argsNode)
    p.matchToken(token.IN, "expected 'in'")
    root.Add(p.expr())
    p.matchToken(token.DO, "expected 'do' to open block")
    root.Add(p.stmtBlock())
    p.matchToken(token.END, "expected 'end' to close block")

    return root
}

func (p *Parser) ifStmt() *interm.Node {
    root := p.createNode(p.current.Str, p.current.TokenType)
    p.nextToken()

    root.Add(p.expr())
    p.matchToken(token.DO, "expected 'do' to open block")

    root.Add(p.stmtBlock())
    for p.peekCurrent() == token.ELIF {
        elifNode := p.createNode(p.current.Str, p.current.TokenType)

        root.Add(elifNode)
        p.nextToken()
        root.Add(p.expr())
        p.matchToken(token.DO, "expected 'do' to open block")
        root.Add(p.stmtBlock())
    }
    if p.peekCurrent() == token.ELSE {
        p.nextToken()
        root.Add(p.stmtBlock())
    }
    p.matchToken(token.END, "expected 'end' to close block")

    return root
}

func (p *Parser) returnStmt() *interm.Node {
    root := p.createNode(p.current.Str, p.current.TokenType)
    p.nextToken()

    tokenType := p.peekCurrent()
    if tokenType != token.NEWLINE &&
       tokenType != token.SEMICOLON &&
       tokenType != token.EOF {
        root.Add(p.expr())
    }
    return root
}

func (p *Parser) importPath() *interm.Node {
    root := p.createNode("PATH", token.PATH)
    p.nextToken()

    for {
        if p.peekCurrent() != token.NAME {
            p.postError("expected name")
        }
        root.Add(p.createNode(p.current.Str,
                 p.current.TokenType))
        p.nextToken()
        /*
         * Check for a dot '.' token. If we find it
         * we skip the token and continue the for loop,
         * otherwise we break the loop.
         */
        if p.peekCurrent() == token.DOT {
            p.nextToken()
            continue
        }
        break
    }
    return root
}

/*
 * The beginning of the expression-chain routines.
 * The deeper the chain, the higher precedence. Since
 * there is no code-generator at the moment, denial of LHS
 * types for assignments is done here.
 */
func (p *Parser) exprStmt() *interm.Node {
    LHS := p.expr()
    tokenType := p.peekCurrent()
    if tokenType >= token.PIPEEQ && tokenType <= token.PERCENTEQ {
        p.checkLHS(LHS.NodeType)
        LHS = p.augAssign(LHS)
    } else if tokenType == token.EQ {
        p.checkLHS(LHS.NodeType)
        opNode := p.createNode(p.current.Str, token.ASSIGN)
        LHS = LHS.GiveRootTo(opNode)

        p.nextToken()
        LHS.Add(p.expr())
    }
    return LHS
}

func (p *Parser) augAssign(node *interm.Node) *interm.Node {
    root := p.createNode("AUGASSIGN", token.AUGASSIGN)
    var nodeType int
    switch p.current.TokenType {
        case token.PIPEEQ:
            nodeType = token.ASS_BITWISE_OR
        case token.AMPEQ:
            nodeType = token.ASS_BITWISE_AND
        case token.CARETEQ:
            nodeType = token.ASS_XOR
        case token.LEFTSHIFTEQ:
            nodeType = token.ASS_LEFTSHIFT
        case token.RIGHTSHIFTEQ:
            nodeType = token.ASS_RIGHTSHIFT
        case token.PLUSEQ:
            nodeType = token.ASS_ADD
        case token.MINUSEQ:
            nodeType = token.ASS_SUB
        case token.STAREQ:
            nodeType = token.ASS_MUL
        case token.SLASHEQ:
            nodeType = token.ASS_DIV
        case token.PERCENTEQ:
            nodeType = token.ASS_MODULO
    }
    opNode := p.createNode(p.current.Str, nodeType)
    node = node.GiveRootTo(opNode)
    p.nextToken()
    node.Add(p.expr())

    root.Add(node)

    return root
}

func (p *Parser) checkLHS(tokenType int) {
    switch tokenType {
        case token.LIST:
            p.postError("cant assign to list")
        case token.HASH:
            p.postError("cant assign to hash")
        case token.CALL:
            p.postError("cant assign to function call")
        case token.NEW:
            p.postError("cant assign to constructor")
        case token.NAME:
            fallthrough
        case token.SUBSCRIPT:
            fallthrough
        case token.MEMBER:
        default:
            if tokenType >= token.STRING &&
               tokenType <= token.NIL {
                p.postError("cant assign to literal")
            }
            p.postError("cant assign to operator")
    }
}

func (p *Parser) expr() *interm.Node {
    return p.rangeExpr()
}

func (p *Parser) rangeExpr() *interm.Node {
    root := p.orExpr()
    for p.peekCurrent() == token.DOTDOT {
        opNode := p.createNode(p.current.Str, token.RANGE)
        root = root.GiveRootTo(opNode)

        p.nextToken()
        root.Add(p.orExpr())
    }
    return root
}

func (p *Parser) orExpr() *interm.Node {
    root := p.andExpr()
    for p.peekCurrent() == token.PIPEPIPE {
        opNode := p.createNode(p.current.Str, token.LOGICAL_OR)
        root = root.GiveRootTo(opNode)

        p.nextToken()
        root.Add(p.andExpr())
    }
    return root
}

func (p *Parser) andExpr() *interm.Node {
    root := p.equalExpr()
    for p.peekCurrent() == token.AMPAMP {
        opNode := p.createNode(p.current.Str, token.LOGICAL_AND)
        root = root.GiveRootTo(opNode)

        p.nextToken()
        root.Add(p.equalExpr())
    }
    return root
}

func (p *Parser) equalExpr() *interm.Node {
    root := p.compExpr()
    tokenType := p.peekCurrent()

    for tokenType == token.BANGEQ || tokenType == token.EQEQ {
        var opNode *interm.Node
        if tokenType == token.BANGEQ {
            opNode = p.createNode(p.current.Str, token.NE)
        } else {
            opNode = p.createNode(p.current.Str, token.EQ)
        }
        root = root.GiveRootTo(opNode)

        p.nextToken()
        root.Add(p.compExpr())

        compNode := p.createNode("COMP_OP", token.COMP_OP)
        root = root.GiveRootTo(compNode)

        tokenType = p.peekCurrent()
    }
    return root
}

func (p *Parser) compExpr() *interm.Node {
    root := p.bitwiseOrExpr()
    tokenType := p.peekCurrent()

    for tokenType >= token.LT && tokenType <= token.GTEQ {
        var opNode *interm.Node
        switch tokenType {
            case token.LT:
                opNode = p.createNode(p.current.Str, token.LT)
            case token.LTEQ:
                opNode = p.createNode(p.current.Str, token.LE)
            case token.GT:
                opNode = p.createNode(p.current.Str, token.GT)
            case token.GTEQ:
                opNode = p.createNode(p.current.Str, token.GE)
        }
        root = root.GiveRootTo(opNode)

        p.nextToken()
        root.Add(p.bitwiseOrExpr())

        compNode := p.createNode("COMP_OP", token.COMP_OP)
        root = root.GiveRootTo(compNode)
        
        tokenType = p.peekCurrent()
    }
    return root
}

func (p *Parser) bitwiseOrExpr() *interm.Node {
    root := p.bitwiseXorExpr()
    for p.peekCurrent() == token.PIPE {
        opNode := p.createNode(p.current.Str, token.BITWISE_OR)
        root = root.GiveRootTo(opNode)

        p.nextToken()
        root.Add(p.bitwiseXorExpr())
    }
    return root
}

func (p *Parser) bitwiseXorExpr() *interm.Node {
    root := p.bitwiseAndExpr()
    for p.peekCurrent() == token.CARET {
        opNode := p.createNode(p.current.Str, token.XOR)
        root = root.GiveRootTo(opNode)

        p.nextToken()
        root.Add(p.bitwiseAndExpr())
    }
    return root
}

func (p *Parser) bitwiseAndExpr() *interm.Node {
    root := p.bitwiseShiftExpr()
    for p.peekCurrent() == token.AMP {
        opNode := p.createNode(p.current.Str, token.BITWISE_AND)
        root = root.GiveRootTo(opNode)

        p.nextToken()
        root.Add(p.bitwiseShiftExpr())
    }
    return root
}

func (p *Parser) bitwiseShiftExpr() *interm.Node {
    root := p.arithExpr()
    tokenType := p.peekCurrent()

    for tokenType == token.LEFTSHIFT ||
        tokenType == token.RIGHTSHIFT {
        opNode := p.createNode(p.current.Str, token.LEFTSHIFT)
        root = root.GiveRootTo(opNode)

        p.nextToken()
        root.Add(p.arithExpr())
        tokenType = p.peekCurrent()
    }
    return root
}

func (p *Parser) arithExpr() *interm.Node {
    root := p.termExpr()
    tokenType := p.peekCurrent()

    for tokenType == token.PLUS || tokenType == token.MINUS {
        var opNode *interm.Node
        if tokenType == token.PLUS {
            opNode = p.createNode(p.current.Str, token.ADD)
        } else {
            opNode = p.createNode(p.current.Str, token.SUB)
        }
        root = root.GiveRootTo(opNode)

        p.nextToken()
        root.Add(p.termExpr())
        tokenType = p.peekCurrent()
    }
    return root
}

func (p *Parser) termExpr() *interm.Node {
    root := p.factorExpr()
    tokenType := p.peekCurrent()

    for tokenType >= token.STAR && tokenType <= token.PERCENT {
        var opNode *interm.Node
        switch tokenType {
            case token.STAR:
                opNode = p.createNode(p.current.Str, token.MUL)
            case token.SLASH:
                opNode = p.createNode(p.current.Str, token.DIV)
            case token.PERCENT:
                opNode = p.createNode(p.current.Str, token.MODULO)
        }
        root = root.GiveRootTo(opNode)

        p.nextToken()
        root.Add(p.factorExpr())
        tokenType = p.peekCurrent()
    }
    return root
}

func (p *Parser) factorExpr() *interm.Node {
    if p.isFactor() {
        var nodeType int
        switch p.current.TokenType {
            case token.MINUS:
                nodeType = token.NEGATE
            case token.BANG:
                nodeType = token.NOT
            case token.TILDE:
                nodeType = token.COMPL
        }
        root := p.createNode(p.current.Str, nodeType)
        p.nextToken()
        if p.isFactor() {
            root.Add(p.factorExpr())
        } else {
            root.Add(p.trailerExpr())
        }
        return root
    }
    return p.trailerExpr()
}

func (p *Parser) trailerExpr() *interm.Node {
    root := p.atom()
out:
    for true {
        switch tokenType := p.peekCurrent(); tokenType {
            case token.LBRACK:
                root = p.subscript(root)
            case token.LPAREN:
                root = p.callFunction(root)
            case token.DOT:
                root = p.instanceAttr(root)
            default:
                break out
        }
    }
    return root
}

func (p *Parser) atom() *interm.Node {
    var node *interm.Node
    tokenType := p.peekCurrent()

    if tokenType >= token.STRING && tokenType <= token.NIL ||
       tokenType == token.NAME {
        node = p.createNode(p.current.Str, p.current.TokenType)
        p.nextToken()
    } else if tokenType == token.LBRACK {
        node = p.arrayLiteral()
    } else if tokenType == token.LBRACE {
        node = p.hashLiteral()
    } else if tokenType == token.LPAREN {
        p.nextToken()
        node = p.expr()
        p.matchToken(token.RPAREN, "expected ')' to close group")
    } else if tokenType == token.NEW {
        node = p.newStmt()
    } else {
        p.postError("expected expression")
    }
    return node
}

func (p *Parser) arrayLiteral() *interm.Node {
    root := p.createNode("LIST", token.LIST)
    p.nextAndSkipNL()
    p.expressionList(root, token.RBRACK)
    p.skipNL()
    p.matchToken(token.RBRACK, "expected ']' to close array" +
                 " literal")
    return root
}

func (p *Parser) hashLiteral() *interm.Node {
    root := p.createNode("HASH", token.HASH)
    if p.peekNext() == token.RBRACE {
        p.nextToken()
        p.nextToken()
        return root
    }
    for true {
        hashElem := p.createNode("HASH_ELEM", token.HASH_ELEM)
        p.nextAndSkipNL()

        hashElem.Add(p.expr())
        p.skipNL()
        p.matchToken(token.EQGT, "expected '=>' between key" +
                     " and value")
        p.skipNL()
        hashElem.Add(p.expr())
        root.Add(hashElem)

        p.skipNL()
        if p.peekCurrent() != token.COMMA {
            break
        }
    }
    p.matchToken(token.RBRACE, "expected '}' to close hash literal")
    return root
}

func (p *Parser) newStmt() *interm.Node {
    root := p.createNode(p.current.Str, token.MAKE_INSTANCE)

    p.nextToken()
    if p.peekCurrent() != token.NAME {
        p.postError("expected name after 'new'")
    }
    nameNode := p.createNode(p.current.Str, p.current.TokenType)
    root.Add(nameNode)
    p.nextToken()

    p.matchToken(token.LPAREN, "expected '('")
    argsNode := p.createNode("ARGUMENTS", token.ARGUMENTS)

    p.expressionList(argsNode, token.RPAREN)
    root.Add(argsNode)

    p.matchToken(token.RPAREN, "expected ')'")
    return root
}

func (p *Parser) subscript(node *interm.Node) *interm.Node {
    root := p.createNode("SUBSCRIPT", token.SUBSCRIPT)
    root = node.GiveRootTo(root)

    p.nextAndSkipNL()
    root.Add(p.expr())
    p.skipNL()
    p.matchToken(token.RBRACK, "expected ']' to close subscript")

    return root
}

func (p *Parser) callFunction(node *interm.Node) *interm.Node {
    root := p.createNode("CALL", token.CALL)
    root = node.GiveRootTo(root)

    p.nextAndSkipNL()
    argsNode := p.createNode("ARGUMENTS", token.ARGUMENTS)
    root.Add(argsNode)

    p.expressionList(argsNode, token.RPAREN)
    p.skipNL()
    p.matchToken(token.RPAREN, "expected ')' to close func call")

    return root
}

func (p *Parser) instanceAttr(node *interm.Node) *interm.Node {
    root := p.createNode(p.current.Str, token.MEMBER)
    root = node.GiveRootTo(root)

    p.nextToken()
    if p.peekCurrent() != token.NAME {
        p.postError("expected name after '.'")
    }
    nameNode := p.createNode(p.current.Str, p.current.TokenType)
    
    root.Add(nameNode)
    p.nextToken()

    return root
}

func (p *Parser) expressionList(node *interm.Node,
                                end int) {
    if p.peekCurrent() == end {
        return
    }
    for true {
        node.Add(p.expr())
        p.skipNL()

        if p.peekCurrent() != token.COMMA {
            break
        }
        p.nextAndSkipNL()
    }
}