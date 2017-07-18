
package blue

import (
    "fmt"
    "bytes"
    "github.com/Magnus9/blue/errpkg"
    "github.com/Magnus9/blue/objects"
)

func genericTraceFunc(frame *objects.BlFrame) {
    var buf bytes.Buffer
    for f := frame; f != nil; f = f.Prev {
        str := fmt.Sprintf("in %s:%d\n   %s", f.Pathname,
                           f.Node.LineNum, f.Node.Line)
        buf.WriteString(str)
    }
    buf.WriteByte('\n')
    buf.WriteString(errpkg.Errmsg)

    panic(buf.String())
}