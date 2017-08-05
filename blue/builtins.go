package blue

import (
    "github.com/Magnus9/blue/errpkg"
    "github.com/Magnus9/blue/objects"
)
var builtins map[string]objects.BlObject

var blBuiltinMethods = []objects.BlGFunctionObject{
    objects.NewBlGFunction("len", builtinLen, objects.GFUNC_VARARGS),
    objects.NewBlGFunction("err", builtinErr, objects.GFUNC_VARARGS),
}

func builtinLen(obj objects.BlObject,
                args ...objects.BlObject) objects.BlObject {
    var arg objects.BlObject
    if objects.BlParseArguments("o", args, &arg) == -1 {
        return nil
    }
    typeobj := arg.BlType()
    if seq := typeobj.Sequence; seq != nil {
        if seq.SqSize != nil {
            return objects.NewBlInt(int64(seq.SqSize(arg)))
        }
    }
    errpkg.SetErrmsg("'%s' object is not a sequence",
                     typeobj.Name)
    return nil
}

/*
 * Since there is no exception handling at the
 * moment, this subroutine can be used to error
 * the interpreter and give a stack-trace.
 */
func builtinErr(obj objects.BlObject,
                args ...objects.BlObject) objects.BlObject {
    var errmsg string
    if objects.BlParseArguments("s", args, &errmsg) == -1 {
        return nil
    }
    errpkg.SetErrmsg(errmsg)
    return nil
}

func blInitBuiltins() {
    mod := blInitModule("builtins", blBuiltinMethods)
    mod.Locals["string"] = &objects.BlStringType
    mod.Locals["float" ] = &objects.BlFloatType
    mod.Locals["list"  ] = &objects.BlListType
    mod.Locals["file"  ] = &objects.BlFileType
    mod.Locals["bool"  ] = &objects.BlBoolType
    mod.Locals["int"   ] = &objects.BlIntType
    mod.Locals["socket"] = &objects.BlSocketType
    builtins = mod.Locals
}