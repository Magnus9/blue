
/*
 * File takes care of parsing literals before
 * they are used as values to Blue objects.
 */
package blue

import (
    "bytes"
    "strconv"
    "github.com/Magnus9/blue/errpkg"
)

func parseInt(value string) int64 {
    size := len(value)
    if size >= 2 {
        if value[1] == 'X' || value[1] == 'x' {
            return parseHex(value[2:], size - 2)
        }
    }
    var prevValue, sum int64 = -1, 0
    for i := 0; i < size && (sum >= prevValue); i++ {
        prevValue = sum
        sum = sum * 10 + int64(value[i]) - 48
    }
    if sum >= prevValue {
        return sum
    } else {
        errpkg.SetErrmsg("number overflow")
        return -1
    }
}

func parseHex(value string, size int) int64 {
    var prevValue, sum int64 = -1, 0
    for i := 0; i < size && (sum >= prevValue); i++ {
        ch := value[i]
        switch {
            case ch >= '0' && ch <= '9':
                ch = ch - '0'
            case ch >= 'A' && ch <= 'F':
                ch = ch - 'A' + 10
            case ch >= 'a' && ch <= 'f':
                ch = ch - 'a' + 10
        }
        prevValue = sum
        sum = sum * 16 + int64(ch)
    }
    if sum >= prevValue {
        return sum
    } else {
        errpkg.SetErrmsg("number overflow")
        return -1
    }
}

func parseFloat(value string) float64 {
    fval, errv := strconv.ParseFloat(value, 64)
    if errv != nil {
        errpkg.SetErrmsg("%s", errv)
        return -1.0
    }
    return fval
}

func parseString(value string) *string {
    size := len(value)
    if size == 2 {
        return new(string)
    }
    if value[0] == 'R' || value[0] == 'r' {
        slice := value[1:size - 1]
        return &slice
    }
    var buf bytes.Buffer
    for spos := 1; spos < size - 1; spos++ {
        if value[spos] == '\\' {
            spos++
            switch value[spos] {
                case '\'': buf.WriteByte('\'')
                case '"':  buf.WriteByte('"')
                case 'n':  buf.WriteByte('\n')
                case 't':  buf.WriteByte('\t')
                case 'r':  buf.WriteByte('\r')
                case 'X':
                    fallthrough
                case 'x':
                    spos++
                    var sum byte
                    ch := value[spos]
                    switch {
                        case ch >= '0' && ch <= '9':
                            sum += ch - '0'
                        case ch >= 'A' && ch <= 'F':
                            sum += ch - 'A' + 10
                        case ch >= 'a' && ch <= 'f':
                            sum += ch - 'a' + 10
                        default:
                            errpkg.SetErrmsg("invalid \\x escape")
                            return nil
                    }
                    sum <<= 4
                    spos++
                    ch = value[spos]
                    switch {
                        case ch >= '0' && ch <= '9':
                            sum += ch - '0'
                        case ch >= 'A' && ch <= 'F':
                            sum += ch - 'A' + 10
                        case ch >= 'a' && ch <= 'f':
                            sum += ch - 'a' + 10
                        default:
                            errpkg.SetErrmsg("invalid \\x escape")
                            return nil
                    }
                    buf.WriteByte(sum)
            }
        } else {
            buf.WriteByte(value[spos])
        }
    }
    bufData := buf.String()
    return &bufData
}