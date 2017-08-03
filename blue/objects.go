
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
    tobj := obj.BlType()
    if tobj.Sequence == nil {
        errpkg.SetErrmsg("'%s' object is not subscriptable",
                         tobj.Name)
        return nil
    }
    seq := tobj.Sequence
    if seq.SeqItem == nil {
        errpkg.SetErrmsg("'%s' object has no subscript getter",
                         tobj.Name)
        return nil
    }
    iobj, ok := key.(*objects.BlIntObject)
    if !ok {
        errpkg.SetErrmsg("'%s' indices must be integers",
                         tobj.Name)
        return nil
    }
    return seq.SeqItem(obj, int(iobj.Value))
}

func blSetSeqItem(obj, value, key objects.BlObject) int {
    tobj := obj.BlType()
    if tobj.Sequence == nil {
        errpkg.SetErrmsg("'%s' object is not subscriptable",
                         tobj.Name)
        return -1
    }
    seq := tobj.Sequence
    if seq.SeqAssItem == nil {
        errpkg.SetErrmsg("'%s' object has no subscript setter",
                         tobj.Name)
        return -1
    }
    iobj, ok := key.(*objects.BlIntObject)
    if !ok {
        errpkg.SetErrmsg("'%s' indices must be integers",
                         tobj.Name)
        return -1
    }
    return seq.SeqAssItem(obj, value, int(iobj.Value))
}

func blGetSlice(obj objects.BlObject, s, e int) objects.BlObject {
    typeobj := obj.BlType()
    if typeobj.Sequence != nil {
        if fn := typeobj.Sequence.SeqSlice; fn != nil {
            return fn(obj, s, e)
        }
    }
    errpkg.SetErrmsg("'%s' object is not subscriptable",
                     typeobj.Name)
    return nil
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
        if fn := typeobj.Sequence.SeqConcat; fn != nil {
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
        if fn := typeobj.Sequence.SeqRepeat; fn != nil {
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