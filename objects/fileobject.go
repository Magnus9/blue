
package objects

import (
    "fmt"
    "os"
    "github.com/Magnus9/blue/errpkg"
)
type BlFileObject struct {
    header blHeader
    f      *os.File
    mode   string
    open   bool
}
func (bfo *BlFileObject) BlType() *BlTypeObject {
    return bfo.header.typeobj
}
var BlFileType BlTypeObject

func NewBlFile(f *os.File, mode string) *BlFileObject {
    return &BlFileObject{
        header: blHeader{&BlFileType},
        f     : f,
        mode  : mode,
        open  : true,
    }
}

func blFileRepr(obj BlObject) *BlStringObject {
    fobj := obj.(*BlFileObject)
    var str string
    if fobj.open {
        str = "open"
    } else {
        str = "closed"
    }
    msg := fmt.Sprintf("<%s file '%s', mode='%s'>", str,
                       fobj.f.Name(), fobj.mode)
    return NewBlString(msg)
}

func blFileGetMember(obj BlObject,
                     name string) BlObject {
    return genericGetMember(obj.BlType(), name, obj)
}

func setModeBits(mode string) int {
    var flag int
    var rfound, wfound bool
    for i := 0; i < len(mode); i++ {
        switch ch := mode[i]; ch {
            case 'r': rfound = true
            case 'w':
                flag |= os.O_WRONLY
                wfound = true
            case 'a': flag |= os.O_APPEND
            case 't': flag |= os.O_TRUNC
            default:
                errpkg.SetErrmsg("unrecognized file mode char" +
                                 " '%c'", ch)
                return -1
        }
    }
    if rfound && wfound {
        flag &= ^os.O_WRONLY
        flag |= os.O_RDWR
    }
    if wfound {
        flag |= os.O_CREATE
    }
    return flag
}

func blFileInit(obj *BlTypeObject,
                args ...BlObject) BlObject {
    var fpath, mode string = "", "r"
    var perm int64 = 0666
    if blParseArguments("s|si", args, &fpath, &mode,
                        &perm) == -1 {
        return nil
    }
    flag := setModeBits(mode)
    if flag == -1 {
        return nil
    }
    f, err := os.OpenFile(fpath, flag, os.FileMode(perm))
    if err != nil {
        errpkg.SetErrmsg(err.Error())
        return nil
    }
    return NewBlFile(f, mode)
}

func blFileRead(self BlObject, args ...BlObject) BlObject {
    var size int64
    if blParseArguments("i", args, &size) == -1 {
        return nil
    }
    fobj := self.(*BlFileObject)
    data := make([]byte, size)
    num, err := fobj.f.Read(data)
    if err != nil {
        if num != 0 {
            errpkg.SetErrmsg(err.Error())
            return nil
        } else {
            return NewBlString("")
        }
    }
    return NewBlString(string(data))
}

func blFileReadAll(self BlObject, args ...BlObject) BlObject {
    fobj := self.(*BlFileObject)
    finfo, err := fobj.f.Stat()
    if err != nil {
        errpkg.SetErrmsg(err.Error())
        return nil
    }
    siz := finfo.Size()
    if siz == 0 {
        return NewBlString("")
    }
    data := make([]byte, siz)
    num, err := fobj.f.Read(data)
    if err != nil {
        if num != 0 {
            errpkg.SetErrmsg(err.Error())
            return nil
        } else {
            return NewBlString("")
        }
    }
    return NewBlString(string(data))
}

func blFileWrite(self BlObject, args ...BlObject) BlObject {
    var buf string
    if blParseArguments("s", args, &buf) == -1 {
        return nil
    }
    fobj := self.(*BlFileObject)
    num, err := fobj.f.WriteString(buf)
    if err != nil {
        errpkg.SetErrmsg(err.Error())
        return nil
    }
    return NewBlInt(int64(num))
}

func blFileClose(self BlObject, args ...BlObject) BlObject {
    fobj := self.(*BlFileObject)

    err := fobj.f.Close()
    if err != nil {
        errpkg.SetErrmsg(err.Error())
        return nil
    }
    fobj.open = false
    return BlNil
}

var blFileMethods = []BlGFunctionObject{
    NewBlGFunction("read",    blFileRead,    GFUNC_VARARGS),
    NewBlGFunction("readall", blFileReadAll, GFUNC_NOARGS ),
    NewBlGFunction("write",   blFileWrite,   GFUNC_VARARGS),
    NewBlGFunction("close",   blFileClose,   GFUNC_NOARGS ),
}

func blInitFile() {
    BlFileType = BlTypeObject{
        header   : blHeader{&BlTypeType},
        Name     : "file",
        Repr     : blFileRepr,
        GetMember: blFileGetMember,
        Init     : blFileInit,
        methods  : blFileMethods,
    }
    blTypeFinish(&BlFileType)
}