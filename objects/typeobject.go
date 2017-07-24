
/*
 * Object is used for now to scale out the code
 * of figuring out whether we are dealing with
 * a type-object/an object attached to a type-object
 * or an instance. This handles the first case.
 */
package objects

import (
    "fmt"
    "github.com/Magnus9/blue/errpkg"
)
var BlTypeType BlTypeObject

func blTypeRepr(obj BlObject) *BlStringObject {
    typeobj := obj.(*BlTypeObject)
    str := fmt.Sprintf("<class '%s'>", typeobj.Name)

    return NewBlString(str)
}

func blTypeGetMember(obj BlObject,
                     name string) BlObject {
    typeobj := obj.(*BlTypeObject)
    ret := locate(typeobj, name)
    if ret == nil {
        errpkg.SetErrmsg("'%s' object has no member '%s'",
                         typeobj.Name, name)
        return nil
    }
    f, ok := ret.(*BlGFunctionObject)
    if ok {
        return newBlGMethod(typeobj, nil, f)
    }
    return ret
}

func blTypeInit(obj *BlTypeObject,
                args ...BlObject) BlObject {
    if obj.Init == nil {
        errpkg.SetErrmsg("'%s' object is missing init" +
                         " function", obj.Name)
        return nil
    }
    return obj.Init(obj, args...)
}

func blInitType() {
    BlTypeType = BlTypeObject{
        Name     : "type",
        Repr     : blTypeRepr,
        GetMember: blTypeGetMember,
        Init     : blTypeInit,
    }
}