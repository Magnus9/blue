
package objects

import (
    "fmt"
    "github.com/Magnus9/blue/errpkg"
)

type BlModuleObject struct {
    header blHeader
    // The module name
    name   string
    // The full path of the module (file).
    Path   string
    // Map filled with symbols.
    Locals map[string]BlObject
}
func (bmo *BlModuleObject) BlType() *BlTypeObject {
    return bmo.header.typeobj
}
var BlModuleType BlTypeObject

func NewBlModule(name, path string) *BlModuleObject {
    return &BlModuleObject{
        header: blHeader{&BlModuleType},
        name  : name,
        Path  : path,
        Locals: make(map[string]BlObject, 0),
    }
}

func blModuleRepr(obj BlObject) *BlStringObject {
    mobj := obj.(*BlModuleObject)
    return NewBlString(fmt.Sprintf("<module '%s', path='%s'>",
                       mobj.name, mobj.Path))
}

func blModuleGetMember(obj BlObject, name string) BlObject {
    mobj := obj.(*BlModuleObject)
    retv, ok := mobj.Locals[name]
    if ok {
        return retv
    }
    errpkg.SetErrmsg("'module' object has no member '%s'",
                     name)
    return nil
}

func blModuleSetMember(obj BlObject, name string,
                       value BlObject) int {
    obj.(*BlModuleObject).Locals[name] = value
    return 0
}

func blInitModule() {
    BlModuleType = BlTypeObject{
        header: blHeader{&BlTypeType},
        Name     : "module",
        Repr     : blModuleRepr,
        GetMember: blModuleGetMember,
        SetMember: blModuleSetMember,
    }
}