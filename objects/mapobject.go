
package objects

import (
    "bytes"
    "github.com/Magnus9/blue/errpkg"
)

type MapPair struct {
    key BlObject
    val BlObject
}

type BlMapObject struct {
    header blHeader
    m      map[int64]MapPair
    mlen   int
}
func (bmo *BlMapObject) BlType() *BlTypeObject {
    return bmo.header.typeobj
}

var blMapMapping = BlMappingMethods{
    MpSize   : blMapSize,
    MpItem   : blMapItem,
    MpAssItem: blMapAssItem,
}

var BlMapType BlTypeObject

func NewBlMap() *BlMapObject {
    return &BlMapObject{
        header: blHeader{&BlMapType},
        m     : make(map[int64]MapPair),
        mlen  : 0,
    }
}

func blMapSize(obj BlObject) int {
    return obj.(*BlMapObject).mlen
}

func blMapItem(obj, key BlObject) BlObject {
    hash := blObjectHash(key)
    if hash == -1 {
        return nil
    }
    mobj := obj.(*BlMapObject)
    pair, ok := mobj.m[hash]
    if !ok {
        errpkg.SetErrmsg("key not found")
        return nil
    }
    return pair.val
}

func blMapAssItem(obj, value, key BlObject) int {
    hash := blObjectHash(key)
    if hash == -1 {
        return -1
    }
    mobj := obj.(*BlMapObject)
    mobj.m[hash] = MapPair{
        key: key,
        val: value,
    }
    return 0
}

func blMapRepr(obj BlObject) *BlStringObject {
    mobj := obj.(*BlMapObject)

    var buf bytes.Buffer
    buf.WriteByte('{')
    i := 0
    for _, pair := range mobj.m {
        if i > 0 {
            buf.WriteString(", ")
        }
        key := pair.key.BlType()
        val := pair.val.BlType()
        buf.WriteString(key.Repr(pair.key).Value)
        buf.WriteString("=>")
        buf.WriteString(val.Repr(pair.val).Value)
        i++
    }
    buf.WriteByte('}')
    return NewBlString(buf.String())
}

func blInitMap() {
    BlMapType = BlTypeObject{
        header : blHeader{&BlTypeType},
        Name   : "map",
        Repr   : blMapRepr,
        Mapping: &blMapMapping,
    }
    blTypeFinish(&BlMapType)
}