
package objects

import (
    "fmt"
    "github.com/Magnus9/blue/errpkg"
)
type BlClassObject struct {
    header  blHeader
    name    string
    // Either builtin or Blue class.
    base    BlObject
    fields  map[string]struct{}
    methods map[string]struct{}
    members map[string]BlObject
}
func (bco *BlClassObject) BlType() *BlTypeObject {
    return bco.header.typeobj
}
var BlClassType BlTypeObject

func NewBlClass(name string, base BlObject) BlObject {
    return &BlClassObject{
        header : blHeader{&BlClassType},
        name   : name,
        base   : base,
        fields : make(map[string]struct{}),
        methods: make(map[string]struct{}),
        members: make(map[string]BlObject),
    }
}

func blClassRepr(obj BlObject) *BlStringObject {
    cobj := obj.(*BlClassObject)
    str := fmt.Sprintf("<class '%s'>", cobj.name)
    
    return NewBlString(str)
}

func blClassGetMember(obj BlObject,
                      name string) BlObject {
    cobj := obj.(*BlClassObject)
    ret, ok := cobj.members[name]
    if !ok {
        if cobj.base != nil {
            switch t := cobj.base.(type) {
                case *BlClassObject:
                    ret = t.BlType().GetMember(cobj.base, name)
                case *BlTypeObject:
                    ret = t.BlType().GetMember(cobj.base, name)
            }
        }
    }
    if ret == nil {
        errpkg.SetErrmsg("'%s' object has no member '%s'",
                         cobj.name, name)
        return nil
    }
    f, ok := ret.(*BlFunctionObject)
    if ok {
        return NewBlMethod(cobj, nil, f)
    }
    return ret
}

func blClassSetMember(obj BlObject, name string,
                      value BlObject) int {
    cobj := obj.(*BlClassObject)
    size := len(name)
    if name[0] == '_' && name[size - 1] == '_' {
        switch name[1:size - 1] {
            case "_init_":
                switch value.(type) {
                    case *BlFunctionObject:
                    case *BlMethodObject:
                default:
                    errpkg.SetErrmsg(
                        "'__init__' must be a blue" +
                        " function or blue method")
                    return -1
                }
        }
    }
    cobj.members[name] = value
    return 0
}

func blClassEvalCond(obj BlObject) bool {
    cobj := obj.(*BlClassObject)
    if len(cobj.members) > 0 {
        return true
    }
    return false
}

func blInitClass() {
    BlClassType = BlTypeObject{
        header   : blHeader{&BlTypeType},
        Name     : "class",
        Repr     : blClassRepr,
        GetMember: blClassGetMember,
        SetMember: blClassSetMember,
        EvalCond : blClassEvalCond,
    }
}