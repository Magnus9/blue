
package objects

import (
    "fmt"
)
type BlInstanceObject struct {
    header  blHeader
    class   *BlClassObject
    fields  map[string]struct{}
    methods map[string]struct{}
    members map[string]BlObject
}
func (bio *BlInstanceObject) BlType() *BlTypeObject {
    return bio.header.typeobj
}
var BlInstanceType BlTypeObject

func NewBlInstance(class *BlClassObject) *BlInstanceObject {
    return &BlInstanceObject{
        header : blHeader{&BlInstanceType},
        class  : class,
        fields : make(map[string]struct{}),
        methods: make(map[string]struct{}),
        members: make(map[string]BlObject, 0),
    }
}

func blInstanceRepr(obj BlObject) *BlStringObject {
    iobj := obj.(*BlInstanceObject)
    str := fmt.Sprintf("<class '%s' instance>", iobj.class.name)
    return NewBlString(str)
}

func blInstanceGetMember(obj BlObject,
                         name string) BlObject {
    iobj := obj.(*BlInstanceObject)
    ret, ok := iobj.members[name]
    if !ok {
        ret = iobj.class.BlType().GetMember(iobj.class, name)
    }
    if ret == nil {
        return nil
    }
    switch t := ret.(type) {
        case *BlFunctionObject:
            return NewBlMethod(iobj.class, iobj, t)
        case *BlMethodObject:
            t.Self = iobj
            return t
    }
    return nil
}

func blInstanceSetMember(obj BlObject, name string,
                         value BlObject) int {
    iobj := obj.(*BlInstanceObject)
    iobj.members[name] = value
    return 0
}

func blInstanceEvalCond(obj BlObject) bool {
    iobj := obj.(*BlInstanceObject)
    if len(iobj.members) > 0 {
        return true
    }
    return false
}

func blInitInstance() {
    BlInstanceType = BlTypeObject{
        header   : blHeader{&BlTypeType},
        Name     : "instance",
        Repr     : blInstanceRepr,
        GetMember: blInstanceGetMember,
        SetMember: blInstanceSetMember,
        EvalCond : blInstanceEvalCond,
    }
}