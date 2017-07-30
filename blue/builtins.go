package blue

import (
	"github.com/Magnus9/blue/errpkg"
	"github.com/Magnus9/blue/objects"
)
var builtins map[string]objects.BlObject

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

func builtinLen(obj objects.BlObject,
				args ...objects.BlObject) objects.BlObject {
	var arg objects.BlObject
	if objects.BlParseArguments("o", args, &arg) == -1 {
		return nil
	}
	typeobj := arg.BlType()
	if seq := typeobj.Sequence; seq != nil {
		if seq.SeqSize != nil {
			return objects.NewBlInt(int64(seq.SeqSize(arg)))
		}
	}
	errpkg.SetErrmsg("'%s' object is not a sequence",
					 typeobj.Name)
	return nil
}

var blBuiltinMethods = []objects.BlGFunctionObject{
	objects.NewBlGFunction("len", builtinLen, objects.GFUNC_VARARGS),
}