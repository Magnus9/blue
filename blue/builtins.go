
package blue

import (
    "github.com/Magnus9/blue/objects"
)
var builtins map[string]objects.BlObject

/*
 * Defined as global because much of the other
 * code might call it from outside the package.
 */
func AddToBuiltins(name string, obj objects.BlObject) {
    builtins[name] = obj
}

func blInitBuiltins() {
    builtins = make(map[string]objects.BlObject)
    
    // Put the string object into builtins.
    AddToBuiltins("string", &objects.BlStringType)
    // Put the float object into builtins (64bit).
    AddToBuiltins("float", &objects.BlFloatType)
    // Put the list object into builtins.
    AddToBuiltins("list", &objects.BlListType)
    // Put the bool object into builtins.
    AddToBuiltins("bool", &objects.BlBoolType)
    // Put the int object into builtins (64bit).
    AddToBuiltins("int", &objects.BlIntType)
    // Put the nil object into builtins (no-val).
    AddToBuiltins("nil", &objects.BlNilType)
}