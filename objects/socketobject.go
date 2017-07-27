
/*
 * File represents the socket type. It is based
 * on doing internal syscalls instead of using
 * the high-level network primitives of golang,
 * since its easier to deal with a file-desc for
 * this task (it yields less code). New socket classes
 * should wrap this one in blue itself.
 */
package objects

import (
    "fmt"
    "net"
    "os"
    "syscall"
    "github.com/Magnus9/blue/errpkg"
)
const (
    AF_INET     = int64(2)
    AF_INET6    = int64(10)
    AF_UNIX     = int64(1)
    SOCK_STREAM = int64(1)
    SOCK_DGRAM  = int64(2)
    IPPROTO_TCP = int64(6)
    IPPROTO_UDP = int64(17)
)
type BlSocketObject struct {
    header blHeader
    // Our file desc. Points to the socket.
    fd     int
    // The domain the fd is attached to.
    domain int64
    // The type that describes semantics.
    stype  int64
    // File object that takes care of IO.
    f      *os.File
    // The abstract socket address.
    saddr  syscall.Sockaddr
}
func (bso *BlSocketObject) BlType() *BlTypeObject {
    return bso.header.typeobj
}

var blSocketMethods = []BlGFunctionObject{
    NewBlGFunction("connect",  socketConnect,    GFUNC_VARARGS),
    NewBlGFunction("bind",     socketBind,       GFUNC_VARARGS),
    NewBlGFunction("listen",   socketListen,     GFUNC_VARARGS),
    NewBlGFunction("accept",   socketAccept,     GFUNC_NOARGS ),
    NewBlGFunction("read",     socketRead,       GFUNC_VARARGS),
    NewBlGFunction("write",    socketWrite,      GFUNC_VARARGS),
    NewBlGFunction("writeall", socketWriteAll,   GFUNC_VARARGS),
    NewBlGFunction("getaddr",  socketGetAddress, GFUNC_NOARGS ),
    NewBlGFunction("close",    socketClose,      GFUNC_VARARGS),
}
var blSocketFields = []BlFields{
    {"AF_INET", T_INT, AF_INET},
    {"AF_INET6", T_INT, AF_INET6},
    {"AF_UNIX", T_INT, AF_UNIX},
    {"SOCK_STREAM", T_INT, SOCK_STREAM},
    {"SOCK_DGRAM", T_INT, SOCK_DGRAM},
    {"IPPROTO_TCP", T_INT, IPPROTO_TCP},
    {"IPPROTO_UDP", T_INT, IPPROTO_UDP},
}
var BlSocketType BlTypeObject

func NewBlSocket(fd int, domain, stype int64,
                 saddr syscall.Sockaddr) *BlSocketObject {
    return &BlSocketObject{
        header: blHeader{&BlSocketType},
        fd    : fd,
        domain: domain,
        stype : stype,
        f     : os.NewFile(uintptr(fd), "socket"),
        saddr : saddr,
    }
}

func blSocketRepr(obj BlObject) *BlStringObject {
    self := obj.(*BlSocketObject)
    msg := fmt.Sprintf("<socket object, domain=%d," +
                       " streamtype=%d>", self.domain,
                       self.stype)
    return NewBlString(msg)
}

func blSocketGetMember(obj BlObject, name string) BlObject {
    return genericGetMember(obj.BlType(), name, obj)
}

func blSocketInit(obj *BlTypeObject, args ...BlObject) BlObject {
    var domain, stype, proto int64
    if blParseArguments("iii", args, &domain, &stype,
                        &proto) == -1 {
        return nil
    }
    /*
     * We start off a bit small and only accept the
     * AF_INET, AF_INET6 and AF_UNIX domains in
     * the beginning.
     */
    var saddr syscall.Sockaddr
    switch domain {
        case AF_UNIX:
            saddr = &syscall.SockaddrUnix{}
        case AF_INET:
            saddr = &syscall.SockaddrInet4{}
        case AF_INET6:
            saddr = &syscall.SockaddrInet6{}
        default:
            errpkg.SetErrmsg("domain (%d) not supported",
                             domain)
            return nil
    }
    fd, err := syscall.Socket(int(domain), int(stype),
                              int(proto))
    if err != nil {
        errpkg.SetErrmsg(err.Error())
        return nil
    }
    return NewBlSocket(fd, domain, stype, saddr)
}

func socketAFUNIXConnect(self *BlSocketObject,
                         arg BlObject) int {
    sobj, ok := arg.(*BlStringObject)
    if !ok {
        errpkg.SetErrmsg("expected string")
        return -1
    }
    self.saddr.(*syscall.SockaddrUnix).Name = sobj.Value
    err := syscall.Connect(self.fd, self.saddr)
    if err != nil {
        errpkg.SetErrmsg(err.Error())
        return -1
    }
    return 0
}

func socketAFINETConnect(self *BlSocketObject,
                         ips []net.IP, port int) int {
    saddr := self.saddr.(*syscall.SockaddrInet4)
    saddr.Port = port
    var err error
    for _, e := range ips {
        if ip := e.To4(); ip != nil {
            copy(saddr.Addr[:], ip)
            err = syscall.Connect(self.fd, self.saddr)
            if err == nil {
                return 0
            }
        }
    }
    if err != nil {
        errpkg.SetErrmsg(err.Error())
    } else {
        errpkg.SetErrmsg("failed to find an ipv4" +
                         " address")
    }
    return -1
}

func socketAFINET6Connect(self *BlSocketObject,
                          ips []net.IP, port int) int {
    saddr := self.saddr.(*syscall.SockaddrInet6)
    saddr.Port = port
    var err error
    for _, e := range ips {
        if ip := e.To16(); ip != nil {
            copy(saddr.Addr[:], ip)
            err = syscall.Connect(self.fd, self.saddr)
            if err == nil {
                return 0
            }
        }
    }
    if err != nil {
        errpkg.SetErrmsg(err.Error())
    } else {
        errpkg.SetErrmsg("failed to find an ipv6" +
                         " address")
    }
    return -1
}

func socketConnect(obj BlObject, args ...BlObject) BlObject {
    var arg BlObject
    if blParseArguments("o", args, &arg) == -1 {
        return nil
    }
    self := obj.(*BlSocketObject)
    switch self.domain {
        case AF_UNIX:
            if socketAFUNIXConnect(self, arg) == -1 {
                return nil
            }
        case AF_INET:
            fallthrough
        case AF_INET6:
            lobj, ok := arg.(*BlListObject)
            if !ok || lobj.lsize != 2 {
                errpkg.SetErrmsg("expected list with addr-port" +
                                 " pair")
                return nil
            }
            args = []BlObject{lobj.list[0],
                              lobj.list[1]}
            var host string
            var port int64
            if blParseArguments("si", args, &host, &port) == -1 {
                return nil
            }
            ips, err := net.LookupIP(host)
            if err != nil {
                errpkg.SetErrmsg(err.Error())
                return nil
            }
            if self.domain == AF_INET {
                if socketAFINETConnect(self, ips, int(port)) == -1 {
                    return nil
                }
            } else {
                if socketAFINET6Connect(self, ips, int(port)) == -1 {
                    return nil
                }
            }
    }
    return BlNil
}

func socketAFUNIXBind(self *BlSocketObject,
                      arg BlObject) int {
    sobj, ok := arg.(*BlStringObject)
    if !ok {
        errpkg.SetErrmsg("expected string")
        return -1
    }
    self.saddr.(*syscall.SockaddrUnix).Name = sobj.Value
    err := syscall.Bind(self.fd, self.saddr)
    if err != nil {
        errpkg.SetErrmsg(err.Error())
        return -1
    }
    return 0
}

func socketAFINETBind(self *BlSocketObject,
                      ips []net.IP, port int) int {
    saddr := self.saddr.(*syscall.SockaddrInet4)
    saddr.Port = port
    for _, e := range ips {
        if ip := e.To4(); ip != nil {
            err := syscall.Bind(self.fd, self.saddr)
            if err != nil {
                errpkg.SetErrmsg(err.Error())
                return -1
            }
            return 0
        }
    }
    errpkg.SetErrmsg("failed to find an ipv4 address")
    return -1
}

func socketAFINET6Bind(self *BlSocketObject,
                       ips []net.IP, port int) int {
    saddr := self.saddr.(*syscall.SockaddrInet6)
    saddr.Port = port
    for _, e := range ips {
        if ip := e.To16(); ip != nil {
            err := syscall.Bind(self.fd, self.saddr)
            if err != nil {
                errpkg.SetErrmsg(err.Error())
                return -1
            }
            return 0
        }
    }
    errpkg.SetErrmsg("failed to find an ipv6 address")
    return -1
}

func socketBind(obj BlObject, args ...BlObject) BlObject {
    var arg BlObject
    if blParseArguments("o", args, &arg) == -1 {
        return nil
    }
    self := obj.(*BlSocketObject)
    switch self.domain {
        case AF_UNIX:
            if socketAFUNIXBind(self, arg) == -1 {
                return nil
            }
        case AF_INET:
            fallthrough
        case AF_INET6:
            lobj, ok := arg.(*BlListObject)
            if !ok || lobj.lsize != 2 {
                errpkg.SetErrmsg("expected list with addr-port" +
                                 " pair")
                return nil
            }
            args = []BlObject{lobj.list[0],
                              lobj.list[1]}
            var host string
            var port int64
            if blParseArguments("si", args, &host, &port) == -1 {
                return nil
            }
            ips, err := net.LookupIP(host)
            if err != nil {
                errpkg.SetErrmsg(err.Error())
                return nil
            }
            if self.domain == AF_INET {
                if socketAFINETBind(self, ips, int(port)) == -1 {
                    return nil
                }
            } else {
                if socketAFINET6Bind(self, ips, int(port)) == -1 {
                    return nil 
                }
            }
    }
    return BlNil
}

func socketListen(obj BlObject, args ...BlObject) BlObject {
    var bklog int64
    if blParseArguments("i", args, &bklog) == -1 {
        return nil
    }
    self := obj.(*BlSocketObject)
    err  := syscall.Listen(self.fd, int(bklog))
    if err != nil {
        errpkg.SetErrmsg(err.Error())
        return nil
    }
    return BlNil
}

func socketAccept(obj BlObject, args ...BlObject) BlObject {
    self := obj.(*BlSocketObject)
    fd, saddr, err := syscall.Accept(self.fd)
    if err != nil {
        errpkg.SetErrmsg(err.Error())
        return nil
    }
    return NewBlSocket(fd, self.domain, self.stype,
                       saddr)
}

/*
 * Get the socketaddr attached to a file-desc. For
 * a AF_UNIX socket it returns a string. For AF_INET
 * and AF_INET6 sockets it returns an [addr, port] pair.
 */
func socketGetAddress(obj BlObject, args ...BlObject) BlObject {
    self := obj.(*BlSocketObject)
    switch self.domain {
        case AF_UNIX:
            saddr := self.saddr.(*syscall.SockaddrUnix)
            return NewBlString(saddr.Name)
        case AF_INET:
            lobj := NewBlList(2)
            saddr := self.saddr.(*syscall.SockaddrInet4)
            lobj.list[0] = NewBlString(
                net.IP(saddr.Addr[:]).String())
            lobj.list[1] = NewBlInt(int64(saddr.Port))
            return lobj
        case AF_INET6:
            lobj := NewBlList(2)
            saddr := self.saddr.(*syscall.SockaddrInet6)
            lobj.list[0] = NewBlString(
                net.IP(saddr.Addr[:]).String())
            lobj.list[1] = NewBlInt(int64(saddr.Port))
            return lobj
    }
    /*
     * Never reaches this state (hence the no default case)
     * above. Domains are checked by the initialization
     * routine and declines domains not supported.
     */
    return BlNil
}

func socketRead(obj BlObject, args ...BlObject) BlObject {
    var size int64
    if blParseArguments("i", args, &size) == -1 {
        return nil
    }
    self := obj.(*BlSocketObject)
    var buf []byte = make([]byte, size)
    n, err := self.f.Read(buf)
    if err != nil {
        /*
         * If n == 0 we just return an empty
         * string.
         */
        if n == 0 {
            return NewBlString("")
        }
        errpkg.SetErrmsg(err.Error())
        return nil
    }
    return NewBlString(string(buf))
}

func socketWrite(obj BlObject, args ...BlObject) BlObject {
    var data string
    if blParseArguments("s", args, &data) == -1 {
        return nil
    }
    self := obj.(*BlSocketObject)
    size, err := self.f.WriteString(data)
    if err != nil {
        errpkg.SetErrmsg(err.Error())
        return nil
    }
    return NewBlInt(int64(size))
}

func socketWriteAll(obj BlObject, args ...BlObject) BlObject {
    var data string
    if blParseArguments("s", args, &data) == -1 {
        return nil
    }
    self := obj.(*BlSocketObject)
    dataLen := len(data)
    var pos int
    for {
        siz, err := self.f.WriteString(data[pos:])
        if err != nil {
            errpkg.SetErrmsg(err.Error())
            return nil
        }
        pos += siz
        if pos == dataLen {
            break
        }
    }
    return NewBlInt(int64(dataLen))
}

func socketClose(obj BlObject, args ...BlObject) BlObject {
    err := obj.(*BlSocketObject).f.Close()
    if err != nil {
        errpkg.SetErrmsg(err.Error())
        return nil
    }
    return BlNil
}

func blInitSocket() {
    BlSocketType = BlTypeObject{
        header   : blHeader{&BlTypeType},
        Name     : "socket",
        Repr     : blSocketRepr,
        GetMember: blSocketGetMember,
        Init     : blSocketInit,
        methods  : blSocketMethods,
        fields   : blSocketFields,
    }
    blTypeFinish(&BlSocketType)
}