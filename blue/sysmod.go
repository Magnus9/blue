
/*
 * The system module works as the core piece
 * of the language. The most important data piece
 * it contains at this moment is the path list which
 * contains paths of where blue should look for
 * modules to be imported.
 */
package blue

import (
	"os"
	"path/filepath"
	"github.com/Magnus9/blue/objects"
)
var (
	BLPATH_LIST = []string{
		"." + string(filepath.Separator) + "modules",
		".",
	}
)

var blSystemMethods = []objects.BlGFunctionObject{
	objects.NewBlGFunction("exit", systemExit,
						   objects.GFUNC_VARARGS),
}

func systemExit(obj objects.BlObject,
				args ...objects.BlObject) objects.BlObject {
	var exitCode int64
	if objects.BlParseArguments("i", args, &exitCode) == -1 {
		return nil
	}
	os.Exit(int(exitCode))
	return nil
}

func makePathObject(paths []string) objects.BlObject {
	lobj := objects.NewBlList(0)
	for _, s := range paths {
		lobj.Append(objects.NewBlString(s))
	}
	return lobj
}

func makeArgvObject(argv []string) objects.BlObject {
	lobj := objects.NewBlList(0)
	for _, s := range argv {
		lobj.Append(objects.NewBlString(s))
	}
	return lobj
}

func blInitSystem(argv []string) objects.BlObject {
	mod := blInitModule("system", blSystemMethods)
	mod.Locals["stdin" ] = objects.NewBlFile(os.Stdin, "r")
	mod.Locals["stdout"] = objects.NewBlFile(os.Stdout, "w")
	mod.Locals["stderr"] = objects.NewBlFile(os.Stderr, "w")
	mod.Locals["path"  ] = makePathObject(BLPATH_LIST)
	mod.Locals["argv"  ] = makeArgvObject(argv)
	return mod
}