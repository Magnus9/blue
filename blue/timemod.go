
package blue

import (
    "time"
    "github.com/Magnus9/blue/objects"
)

var blTimeMethods = []objects.BlGFunctionObject{
    objects.NewBlGFunction("sleep", timeSleep,
                           objects.GFUNC_VARARGS),
    objects.NewBlGFunction("ctime", timeCtime,
                           objects.GFUNC_VARARGS),
    objects.NewBlGFunction("time", timeTime,
                           objects.GFUNC_NOARGS),
}

func timeSleep(obj objects.BlObject,
               args ...objects.BlObject) objects.BlObject {
    var msecs int64
    if objects.BlParseArguments("i", args, &msecs) == -1 {
        return nil
    }
    time.Sleep(time.Duration(msecs) * time.Millisecond)
    return objects.BlNil
}

func timeTime(obj objects.BlObject,
              args ...objects.BlObject) objects.BlObject {
    msecs := time.Duration(time.Now().UnixNano()) /
             time.Millisecond
    return objects.NewBlInt(int64(msecs))
}

func timeCtime(obj objects.BlObject,
               args ...objects.BlObject) objects.BlObject {
    var msecs int64
    if objects.BlParseArguments("i", args, &msecs) == -1 {
        return nil
    }
    tstr := time.Unix(msecs / 1000, 0).Format(time.ANSIC)
    return objects.NewBlString(tstr)
}

func blInitTime() objects.BlObject {
    mod := blInitModule("time", blTimeMethods)
    return mod
}