package protocache

import (
	"unsafe"
)

func init() {
	buf := [2]byte{0xff, 0}
	if *(*uint16)(unsafe.Pointer(&buf[0])) != 0xff {
		panic("little-endian only")
	}
}

func sizeof[T any]() int {
	var t T
	return int(unsafe.Sizeof(t))
}

func castStrToBytes(str string) []byte {
	return unsafe.Slice(unsafe.StringData(str), len(str))
}

func castBytesToStr(raw []byte) string {
	return *(*string)(unsafe.Pointer(&raw))
}

type pod interface {
	~bool | ~int8 | ~uint8 | ~int16 | ~uint16 | ~int32 | ~uint32 |
		~int64 | ~uint64 | ~float32 | ~float64
}

func upCast[A, B pod](src []A) []B {
	a, b := sizeof[A](), sizeof[B]()
	if b%a != 0 {
		panic("illegal up cast")
	}
	m := b / a
	return unsafe.Slice((*B)(unsafe.Pointer(unsafe.SliceData(src))), len(src)/m)
}

func downCast[A, B pod](src []A) []B {
	a, b := sizeof[A](), sizeof[B]()
	if a%b != 0 {
		panic("illegal down cast")
	}
	m := a / b
	return unsafe.Slice((*B)(unsafe.Pointer(unsafe.SliceData(src))), len(src)*m)
}

func castToBytes[T pod](src []T) []byte {
	return downCast[T, byte](src)
}

func castBytesToWords(src []byte) []uint32 {
	return upCast[byte, uint32](src)
}

func getInt32(data []byte) int32 {
	return *(*int32)(unsafe.Pointer(&data[0]))
}

func getUint32(data []byte) uint32 {
	return *(*uint32)(unsafe.Pointer(&data[0]))
}

func getUint16(data []byte) uint16 {
	return *(*uint16)(unsafe.Pointer(&data[0]))
}

func getInt64(data []byte) int64 {
	return *(*int64)(unsafe.Pointer(&data[0]))
}

func getUint64(data []byte) uint64 {
	return *(*uint64)(unsafe.Pointer(&data[0]))
}

func getFloat32(data []byte) float32 {
	return *(*float32)(unsafe.Pointer(&data[0]))
}

func getFloat64(data []byte) float64 {
	return *(*float64)(unsafe.Pointer(&data[0]))
}

func putInt32(data []byte, val int32) {
	*(*int32)(unsafe.Pointer(&data[0])) = val
}

func putUint32(data []byte, val uint32) {
	*(*uint32)(unsafe.Pointer(&data[0])) = val
}

func putUint16(data []byte, val uint16) {
	*(*uint16)(unsafe.Pointer(&data[0])) = val
}

func putInt64(data []byte, val int64) {
	*(*int64)(unsafe.Pointer(&data[0])) = val
}

func putUint64(data []byte, val uint64) {
	*(*uint64)(unsafe.Pointer(&data[0])) = val
}

func putFloat32(data []byte, val float32) {
	*(*float32)(unsafe.Pointer(&data[0])) = val
}

func putFloat64(data []byte, val float64) {
	*(*float64)(unsafe.Pointer(&data[0])) = val
}
