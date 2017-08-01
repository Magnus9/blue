
package blue

import (
    "os"
    "fmt"
    "bytes"
    "strings"
    "path/filepath"
    "github.com/Magnus9/blue/parser"
    "github.com/Magnus9/blue/interm"
    "github.com/Magnus9/blue/errpkg"
    "github.com/Magnus9/blue/objects"
)

func blInitModule(
name string,
funcs []objects.BlGFunctionObject) *objects.BlModuleObject {
    modules := GetModuleMap()
    mod, ok := modules[name]
    if ok {
        return mod
    }
    /*
     * The module object wasnt found in the global
     * module dictionary. Create a fresh one and
     * insert the functions passed.
     */
    mod = objects.NewBlModule(name, "builtin")
    if funcs != nil {
        objects.BlInsertFunctions(funcs, mod.Locals)
    }
    modules[name] = mod

    return mod
}

/*
 * The entry point of a module import. The chain
 * of routines is very simple because of how we do
 * imports. There is one rule only: If the import
 * path is a directory, there must be a file within
 * with the same name as the directory.

 * @param NODE, the non-terminal PATH node.
 */
func blImportModule(NODE *interm.Node) (string, objects.BlObject) {
    var buf bytes.Buffer
    siz := len(NODE.Children)
    for i, p := range NODE.Children {
        buf.WriteString(p.Str)
        if i != siz - 1 {
            buf.WriteByte(filepath.Separator)
        }
    }
    path := buf.String()
    base := filepath.Base(path)
    return base, blLocateModule(base, path)
}

func blLocateModule(name, path string) objects.BlObject {
    modules := GetModuleMap()
    // If the module is already in the module map
    // we just return it.
    mod, ok := modules[path]
    if ok {
        return mod
    }
    mod, ok = modules["system"]
    if !ok {
        errpkg.InternError("system module not added at" +
                           " initialization")
        return nil
    }
    lobj := mod.Locals["path"].(*objects.BlListObject)
    for _, obj := range lobj.GetList() {
        /*
         * The initialization stage only adds BlStringObject's
         * into system.path. But if somone were stupid enough
         * to actually append a non-string onto it on runtime
         * we ought to check for it so we dont crash here :).
         */
        sobj, ok := obj.(*objects.BlStringObject)
        if !ok {
            // Just continue on non-strings.
            continue
        }
        fullpath := fmt.Sprintf("%s%c%s", sobj.Value,
                                filepath.Separator, path)
        stat, err := os.Stat(fullpath)
        var f *os.File
        if err == nil && stat.IsDir() {
            f, err = os.Open(fmt.Sprintf("%s%c%s.bl",
                             fullpath, filepath.Separator,
                             name))
            if err != nil {
                break
            }
        } else {
            f, err = os.Open(fullpath + ".bl")
        }
        if f != nil {
            return blLoadModule(f, name, fullpath)
        }
    }
    pstr := strings.Replace(path, string(filepath.Separator),
                            ".", -1)
    errpkg.SetErrmsg("failed to load module '%s'", pstr)
    return nil
}

func blLoadModule(fdesc *os.File,
                  name, fullpath string) objects.BlObject {
    ast := parser.ParseFromFile(fullpath, fdesc)
    mod := blAddModule(name, fullpath)
    return blExecModule(mod, ast)
}

func blExecModule(mod *objects.BlModuleObject,
                  ast *interm.Node) objects.BlObject {
    /*
     * Spawn a new instance of the VM passing in
     * mod.Locals as globals.
     */
    runtime := New(mod.Path, ast)
    runtime.Run(mod.Locals)
    // If the evaluation succeeded mod.Locals is filled.
    return mod
}

/*
 * Caution: This does not check if a module with 'name'
 * already exists in the module map, since its supposed
 * to be called from the import stmt, and module existence
 * is checked in blLocateModule. Builtin module objects
 * should refer to blInitModule.
 */
func blAddModule(name, fullpath string) *objects.BlModuleObject {
    modules := GetModuleMap()
    mod := objects.NewBlModule(name, fullpath)
    modules[name] = mod
    
    return mod
}