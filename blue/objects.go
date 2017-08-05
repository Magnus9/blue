
package blue

import (
    "fmt"
    "github.com/Magnus9/blue/token"
    "github.com/Magnus9/blue/errpkg"
    "github.com/Magnus9/blue/objects"
)

func blGetMember(obj objects.BlObject,
                 name string) objects.BlObject {
    typeobj := obj.BlType()
    if fn := typeobj.GetMember; fn != nil {
        return fn(obj, name)
    }
    errpkg.SetErrmsg("cant retrieve members from the '%s'" +
                     " object", typeobj.Name)
    return nil
}

func blSetMember(obj, value objects.BlObject,
                 name string) int {
    typeobj := obj.BlType()
    if fn := typeobj.SetMember; fn != nil {
        return fn(obj, name, value)
    }
    errpkg.SetErrmsg("cant set members on the '%s' object",
                     typeobj.Name)
    return -1
}

func blGetSeqItem(obj objects.BlObject,
                  key objects.BlObject) objects.BlObject {
    typeobj := obj.BlType()
    /*
     * Start by validating that key is type
     * integer.
     */
    iobj, ok := key.(*objects.BlIntObject)
    if !ok {
        errpkg.SetErrmsg("'%s' indices must be integers",
                         typeobj.Name)
        return nil
    }
    siz := int(iobj.Value)
    if seq := typeobj.Sequence; seq != nil {
        if seq.SqItem != nil {
            if siz < 0 && seq.SqSize != nil {
                siz += seq.SqSize(obj)
            }
            return seq.SqItem(obj, siz)
        }
    }
    errpkg.SetErrmsg("'%s' object is not subscriptable",
                     typeobj.Name)
    return nil
}

func blSetSeqItem(obj, value, key objects.BlObject) int {
    typeobj := obj.BlType()
    /*
     * Start by validating that key is type
     * integer.
     */
    iobj, ok := key.(*objects.BlIntObject)
    if !ok {
        errpkg.SetErrmsg("'%s' indices must be integers",
                         typeobj.Name)
        return -1
    }
    siz := int(iobj.Value)
    if seq := typeobj.Sequence; seq != nil {
        if seq.SqAssItem != nil {
            if siz < 0 && seq.SqSize != nil {
                siz += seq.SqSize(obj)
            }
            return seq.SqAssItem(obj, value, siz)
        }
    }
    errpkg.SetErrmsg("'%s' object does not support item" +
                     " assignment", typeobj.Name)
    return -1
}

func blGetSlice(obj objects.BlObject, s, e int) objects.BlObject {
    typeobj := obj.BlType()
    if seq := typeobj.Sequence; seq != nil {
        if seq.SqSlice != nil {
            if s < 0 || e < 0 && seq.SqSize != nil {
                siz := seq.SqSize(obj)
                if s < 0 {
                    s += siz
                }
                if e < 0 {
                    e += siz
                }
            }
            return seq.SqSlice(obj, s, e)
        }
    }
    errpkg.SetErrmsg("'%s' object is not slicable",
                     typeobj.Name)
    return nil
}

func blSetSlice(obj, value objects.BlObject, s, e int) int {
    typeobj := obj.BlType()
    if seq := typeobj.Sequence; seq != nil {
        if seq.SqAssSlice != nil {
            if s < 0 || e < 0 && seq.SqSize != nil {
                siz := seq.SqSize(obj)
                if s < 0 {
                    s += siz
                }
                if e < 0 {
                    e += siz
                }
            }
            return seq.SqAssSlice(obj, value, s, e)
        }
    }
    errpkg.SetErrmsg("'%s' object does not support slice" +
                     " assignment", typeobj.Name)
    return -1
}

func blNumNegate(obj objects.BlObject) objects.BlObject {
    typeobj := obj.BlType()
    if typeobj.Numbers != nil {
        if fn := typeobj.Numbers.NumNeg; fn != nil {
            return fn(obj)
        }
    }
    errpkg.SetErrmsg("bad operand type for '-'")
    return nil
}

func blNumCompl(obj objects.BlObject) objects.BlObject {
    typeobj := obj.BlType()
    if typeobj.Numbers != nil {
        if fn := typeobj.Numbers.NumCompl; fn != nil {
            return fn(obj)
        }
    }
    errpkg.SetErrmsg("bad operand type for '~'")
    return nil
}

func blNumOr(a, b objects.BlObject, op string) objects.BlObject {
    if a.BlType().Numbers != nil {
        if objects.BlNumCoerce(&a, &b) == -1 {
            goto err
        }
        typeobj := a.BlType()
        if typeobj.Numbers != nil {
            if fn := typeobj.Numbers.NumOr; fn != nil {
                return fn(a, b)
            }
        }
    }
err:
    errpkg.SetErrmsg("bad operand types for '%s'",
                      op)
    return nil
}

func blNumAnd(a, b objects.BlObject, op string) objects.BlObject {
    if a.BlType().Numbers != nil {
        if objects.BlNumCoerce(&a, &b) == -1 {
            goto err
        }
        typeobj := a.BlType()
        if typeobj.Numbers != nil {
            if fn := typeobj.Numbers.NumAnd; fn != nil {
                return fn(a, b)
            }
        }
    }
err:
    errpkg.SetErrmsg("bad operand types for '%s'",
                      op)
    return nil
}

func blNumXor(a, b objects.BlObject, op string) objects.BlObject {
    if a.BlType().Numbers != nil {
        if objects.BlNumCoerce(&a, &b) == -1 {
            goto err
        }
        typeobj := a.BlType()
        if typeobj.Numbers != nil {
            if fn := typeobj.Numbers.NumXor; fn != nil {
                return fn(a, b)
            }
        }
    }
err:
    errpkg.SetErrmsg("bad operand types for '%s'",
                      op)
    return nil
}

func blNumLshift(a, b objects.BlObject, op string) objects.BlObject {
    if a.BlType().Numbers != nil {
        if objects.BlNumCoerce(&a, &b) == -1 {
            goto err
        }
        typeobj := a.BlType()
        if typeobj.Numbers != nil {
            if fn := typeobj.Numbers.NumLshift; fn != nil {
                return fn(a, b)
            }
        }
    }
err:
    errpkg.SetErrmsg("bad operand types for '%s'",
                      op)
    return nil
}

func blNumRshift(a, b objects.BlObject, op string) objects.BlObject {
    if a.BlType().Numbers != nil {
        if objects.BlNumCoerce(&a, &b) == -1 {
            goto err
        }
        typeobj := a.BlType()
        if typeobj.Numbers != nil {
            if fn := typeobj.Numbers.NumRshift; fn != nil {
                return fn(a, b)
            }
        }
    }
err:
    errpkg.SetErrmsg("bad operand types for '%s'",
                      op)
    return nil
}

func blNumAddition(a, b objects.BlObject, op string) objects.BlObject {
    typeobj := a.BlType()
    if typeobj.Sequence != nil {
        if fn := typeobj.Sequence.SqConcat; fn != nil {
            return fn(a, b)
        }
    }
    if typeobj.Numbers != nil {
        if objects.BlNumCoerce(&a, &b) == -1 {
            goto err
        }
        typeobj = a.BlType()
        if typeobj.Numbers != nil {
            if fn := typeobj.Numbers.NumAdd; fn != nil {
                return fn(a, b)
            }
        }
    }
err:
    errpkg.SetErrmsg("bad operand types for '%s'",
                      op)
    return nil
}

func blNumSubtract(a, b objects.BlObject, op string) objects.BlObject {
    if a.BlType().Numbers != nil {
        if objects.BlNumCoerce(&a, &b) == -1 {
            goto err
        }
        typeobj := a.BlType()
        if typeobj.Numbers != nil {
            if fn := typeobj.Numbers.NumSub; fn != nil {
                return fn(a, b)
            }
        }
    }
err:
    errpkg.SetErrmsg("bad operand types for '%s'",
                      op)
    return nil
}

func blNumMultiply(a, b objects.BlObject, op string) objects.BlObject {
    typeobj := a.BlType()
    if typeobj.Sequence != nil {
        if fn := typeobj.Sequence.SqRepeat; fn != nil {
            return fn(a, b)
        }
    }
    if typeobj.Numbers != nil {
        if objects.BlNumCoerce(&a, &b) == -1 {
            goto err
        }
        typeobj = a.BlType()
        if typeobj.Numbers != nil {
            if fn := typeobj.Numbers.NumMul; fn != nil {
                return fn(a, b)
            }
        }
    }
err:
    errpkg.SetErrmsg("bad operand types for '%s'",
                      op)
    return nil
}

func blNumDivide(a, b objects.BlObject, op string) objects.BlObject {
    if a.BlType().Numbers != nil {
        if objects.BlNumCoerce(&a, &b) == -1 {
            goto err
        }
        typeobj := a.BlType()
        if typeobj.Numbers != nil {
            if fn := typeobj.Numbers.NumDiv; fn != nil {
                return fn(a, b)
            }
        }
    }
err:
    errpkg.SetErrmsg("bad operand types for '%s'",
                      op)
    return nil
}

func blNumModulo(a, b objects.BlObject, op string) objects.BlObject {
    if a.BlType().Numbers != nil {
        if objects.BlNumCoerce(&a, &b) == -1 {
            goto err
        }
        typeobj := a.BlType()
        if typeobj.Numbers != nil {
            if fn := typeobj.Numbers.NumMod; fn != nil {
                return fn(a, b)
            }
        }
    }
err:
    errpkg.SetErrmsg("bad operand types for '%s'",
                      op)
    return nil
}

func blCmp(a, b objects.BlObject, op int) objects.BlObject {
    value := objects.BlCompare(a, b)
    if op != token.EQ && value == -2 {
        return nil
    }
    var res bool
    switch op {
        case token.EQ: res = value == 0
        case token.NE: res = value != 0
        case token.LT: res = value <  0
        case token.LE: res = value <= 0
        case token.GT: res = value >  0
        case token.GE: res = value >= 0
    }
    if res {
        return objects.BlTrue
    }
    return objects.BlFalse
}

/*
 * Objects that dont have an EvalCond function
 * returns true as default.
 */
func blEvalCondition(obj objects.BlObject) bool {
    fn := obj.BlType().EvalCond
    if fn == nil {
        return true
    }
    return fn(obj)
}

func blPrint(obj objects.BlObject) int {
    typeobj := obj.BlType()
    if typeobj.Repr == nil {
        errpkg.SetErrmsg("'%s' object has no representation",
                         typeobj.Name)
        return -1
    }
    fmt.Println(typeobj.Repr(obj).Value)
    return 0
}