package test

import (
	//"fmt"
	"math"
	"os"
	"testing"
	"unsafe"

	"github.com/peterrk/protocache-go"
	//"github.com/peterrk/protocache-go/test/fb"
	"github.com/peterrk/protocache-go/test/pb"
	"github.com/peterrk/protocache-go/test/pc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var (
	benchmarkBytes []byte
	benchmarkFuse  uint64
)

func BenchmarkCompress(b *testing.B) {
	raw, err := os.ReadFile("test.pc")
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkBytes = protocache.Compress(raw)
	}
}

func BenchmarkDecompress(b *testing.B) {
	raw, err := os.ReadFile("test.pc")
	if err != nil {
		b.Fatal(err)
	}
	raw = protocache.Compress(raw)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		out, err := protocache.Decompress(raw)
		if err != nil {
			b.Fatal(err)
		}
		benchmarkBytes = out
	}
}

/*
	func TestProtoCacheBenchmark(t *testing.T) {
		var junk Junk
		raw, err := os.ReadFile("test.pc")
		if err != nil {
			t.Fatal(err)
		}
		root := pc.AS_Main(raw)
		junk.traversePcMain(root)
		fmt.Println(junk.fuse())
	}

	func TestProtobufBenchmark(t *testing.T) {
		var junk Junk
		raw, err := os.ReadFile("test.pb")
		if err != nil {
			t.Fatal(err)
		}
		root := &pb.Main{}
		err = proto.Unmarshal(raw, root)
		if err != nil {
			t.Fatal(err)
		}
		junk.traversePbMain(root)
		fmt.Println(junk.fuse())

		junk = Junk{}
		junk.traversePbMessage(root.ProtoReflect())
		fmt.Println(junk.fuse())
	}

	func TestFlatbuffersBenchmark(t *testing.T) {
		var junk Junk
		raw, err := os.ReadFile("test.fb")
		if err != nil {
			t.Fatal(err)
		}
		root := fb.GetRootAsMain(raw, 0)
		junk.traverseFbMain(root)
		fmt.Println(junk.fuse())
	}
*/

func BenchmarkProtobufVT(b *testing.B) {
	raw, err := os.ReadFile("test.pb")
	if err != nil {
		b.Fatal(err)
	}
	var junk Junk
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		root := &pb.Main{}
		err = root.UnmarshalVT(raw)
		if err != nil {
			b.Fatal(err)
		}
		junk.traversePbMain(root)
	}
	benchmarkFuse = junk.fuse()
}

func BenchmarkProtobufVTReflect(b *testing.B) {
	raw, err := os.ReadFile("test.pb")
	if err != nil {
		b.Fatal(err)
	}
	var junk Junk
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		root := &pb.Main{}
		err = root.UnmarshalVT(raw)
		if err != nil {
			b.Fatal(err)
		}
		junk.traversePbMessage(root.ProtoReflect())
	}
	benchmarkFuse = junk.fuse()
}

func BenchmarkProtobuf(b *testing.B) {
	raw, err := os.ReadFile("test.pb")
	if err != nil {
		b.Fatal(err)
	}
	var junk Junk
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		root := &pb.Main{}
		err = proto.Unmarshal(raw, root)
		if err != nil {
			b.Fatal(err)
		}
		junk.traversePbMain(root)
	}
	benchmarkFuse = junk.fuse()
}

func BenchmarkProtobufReflect(b *testing.B) {
	raw, err := os.ReadFile("test.pb")
	if err != nil {
		b.Fatal(err)
	}
	var junk Junk
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		root := &pb.Main{}
		err = proto.Unmarshal(raw, root)
		if err != nil {
			b.Fatal(err)
		}
		junk.traversePbMessage(root.ProtoReflect())
	}
	benchmarkFuse = junk.fuse()
}

func BenchmarkProtoCache(b *testing.B) {
	raw, err := os.ReadFile("test.pc")
	if err != nil {
		b.Fatal(err)
	}
	var junk Junk
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		root := pc.AS_Main(raw)
		junk.traversePcMain(root)
	}
	benchmarkFuse = junk.fuse()
}

/*
	func BenchmarkFlatbuffers(b *testing.B) {
		raw, err := os.ReadFile("test.fb")
		if err != nil {
			b.Fatal(err)
		}
		var junk Junk
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			root := fb.GetRootAsMain(raw, 0)
			junk.traverseFbMain(root)
		}
		benchmarkFuse = junk.fuse()
	}
*/
type Junk struct {
	u32 uint32
	f32 float32
	u64 uint64
	f64 float64
}

func (p *Junk) fuse() uint64 {
	return p.u64 + uint64(p.u32) + math.Float64bits(p.f64+float64(p.f32))
}

func (p *Junk) consumeBool(flag bool) {
	if flag {
		p.u32++
	}
}

func (p *Junk) consumeString(str string) {
	view := unsafe.Slice((*uint32)(unsafe.Pointer(unsafe.StringData(str))), len(str)/4)
	for _, v := range view {
		p.u32 += v
	}
}

func (p *Junk) consumeBytes(data []byte) {
	view := unsafe.Slice((*uint32)(unsafe.Pointer(unsafe.SliceData(data))), len(data)/4)
	for _, v := range view {
		p.u32 += v
	}
}

func (p *Junk) traversePbSmall(root *pb.Small) {
	p.u32 += uint32(root.I32)
	p.consumeBool(root.Flag)
	p.consumeString(root.Str)
}

func (p *Junk) traversePbVec2D(root *pb.Vec2D) {
	for _, one := range root.X {
		for _, v := range one.X {
			p.f32 += v
		}
	}
}

func (p *Junk) traversePbArrMap(root *pb.ArrMap) {
	for key, val := range root.X {
		p.consumeString(key)
		for _, v := range val.X {
			p.f32 += v
		}
	}
}

func (p *Junk) traversePbMain(root *pb.Main) {
	p.u32 += uint32(root.I32) + root.U32 + uint32(root.Mode)
	p.consumeBool(root.Flag)
	p.u32 += uint32(root.TI32) + uint32(root.TS32) + root.TU32
	for _, v := range root.I32V {
		p.u32 += uint32(v)
	}
	p.u64 += uint64(root.I64) + root.U64 +
		uint64(root.TI64) + uint64(root.TS64) + root.TU64
	for _, v := range root.U64V {
		p.u64 += v
	}
	for _, v := range root.Flags {
		p.consumeBool(v)
	}
	p.consumeString(root.Str)
	p.consumeBytes(root.Data)
	for _, v := range root.Strv {
		p.consumeString(v)
	}
	for _, v := range root.Datav {
		p.consumeBytes(v)
	}

	p.f32 += root.F32
	for _, v := range root.F32V {
		p.f32 += v
	}
	p.f64 += root.F64
	for _, v := range root.F64V {
		p.f64 += v
	}

	p.traversePbSmall(root.Object)
	for _, v := range root.Objectv {
		p.traversePbSmall(v)
	}

	for key, val := range root.Index {
		p.consumeString(key)
		p.u32 += uint32(val)
	}

	for key, val := range root.Objects {
		p.u32 += uint32(key)
		p.traversePbSmall(val)
	}

	p.traversePbVec2D(root.Matrix)

	for _, one := range root.Vector {
		p.traversePbArrMap(one)
	}
	p.traversePbArrMap(root.Arrays)
}

func (p *Junk) traversePcSmall(root pc.Small) {
	p.u32 += uint32(root.GetI32())
	p.consumeBool(root.GetFlag())
	p.consumeString(root.GetStr())
}

func (p *Junk) traversePcMain(root pc.Main) {
	p.u32 += uint32(root.GetI32()) + root.GetU32() + uint32(root.GetMode())
	p.consumeBool(root.GetFlag())
	p.u32 += uint32(root.GetTI32()) + uint32(root.GetTS32()) + root.GetTU32()
	for _, v := range root.GetI32V() {
		p.u32 += uint32(v)
	}
	p.u64 += uint64(root.GetI64()) + root.GetU64() +
		uint64(root.GetTI64()) + uint64(root.GetTS64()) + root.GetTU64()
	for _, v := range root.GetU64V() {
		p.u64 += v
	}
	for _, v := range root.GetFlags() {
		p.consumeBool(v)
	}
	p.consumeString(root.GetStr())
	p.consumeBytes(root.GetData())
	strv := root.GetStrv()
	for i := uint32(0); i < strv.Size(); i++ {
		p.consumeString(strv.Get(i))
	}
	datv := root.GetDatav()
	for i := uint32(0); i < datv.Size(); i++ {
		p.consumeBytes(datv.Get(i))
	}

	p.f32 += root.GetF32()
	for _, v := range root.GetF32V() {
		p.f32 += v
	}
	p.f64 += root.GetF64()
	for _, v := range root.GetF64V() {
		p.f64 += v
	}

	p.traversePcSmall(root.GetObject())
	objs := root.GetObjectv()
	for i := uint32(0); i < objs.Size(); i++ {
		p.traversePcSmall(objs.Get(i))
	}

	map1 := root.GetIndex()
	for i := uint32(0); i < map1.Size(); i++ {
		p.consumeString(map1.Key(i))
		p.u32 += uint32(map1.Value(i))
	}

	map2 := root.GetObjects()
	for i := uint32(0); i < map2.Size(); i++ {
		p.u32 += uint32(map2.Key(i))
		p.traversePcSmall(map2.Value(i))
	}

	matrix := root.GetMatrix()
	for i := uint32(0); i < matrix.Size(); i++ {
		w := matrix.Get(i)
		for _, v := range w.Raw() {
			p.f32 += v
		}
	}

	vector := root.GetVector()
	for i := uint32(0); i < vector.Size(); i++ {
		map3 := vector.Get(i)
		for i := uint32(0); i < map3.Size(); i++ {
			p.consumeString(map3.Key(i))
			w := map3.Value(i)
			for _, v := range w.Raw() {
				p.f32 += v
			}
		}
	}

	map4 := root.GetArrays()
	for i := uint32(0); i < map4.Size(); i++ {
		p.consumeString(map4.Key(i))
		w := map4.Value(i)
		for _, v := range w.Raw() {
			p.f32 += v
		}
	}
}

/*
	func (p *Junk) traverseFbSmall(root *fb.Small) {
		p.u32 += uint32(root.I32())
		p.consumeBool(root.Flag())
		p.consumeBytes(root.Str())
	}

	func (p *Junk) traverseFbVec2D(root *fb.Vec2D) {
		var unit fb.Vec1D
		for i := 0; i < root.AliasLength(); i++ {
			root.Alias(&unit, i)
			for j := 0; j < unit.AliasLength(); j++ {
				p.f32 += unit.Alias(j)
			}
		}
	}

	func (p *Junk) traverseFbArrMap(root *fb.ArrMap) {
		var unit fb.Array
		var pair fb.ArrMapEntry
		for i := 0; i < root.AliasLength(); i++ {
			root.Alias(&pair, i)
			p.consumeBytes(pair.Key())
			pair.Value(&unit)
			for j := 0; j < unit.AliasLength(); j++ {
				p.f32 += unit.Alias(j)
			}
		}
	}

	func (p *Junk) traverseFbMain(root *fb.Main) {
		p.u32 += uint32(root.I32()) + root.U32() + uint32(root.Mode())
		p.consumeBool(root.Flag())
		p.u32 += uint32(root.TI32()) + uint32(root.TS32()) + root.TU32()
		for i := 0; i < root.I32vLength(); i++ {
			p.u32 += uint32(root.I32v(i))
		}
		p.u64 += uint64(root.I64()) + root.U64() +
			uint64(root.TI64()) + uint64(root.TS64()) + root.TU64()
		for i := 0; i < root.U64vLength(); i++ {
			p.u64 += root.U64v(i)
		}
		for i := 0; i < root.FlagsLength(); i++ {
			p.consumeBool(root.Flags(i))
		}
		p.consumeBytes(root.Str())

		data := make([]byte, root.DataLength())
		for i := 0; i < len(data); i++ {
			data[i] = byte(root.Data(i))
		}
		p.consumeBytes(data)
		for i := 0; i < root.StrvLength(); i++ {
			p.consumeBytes(root.Strv(i))
		}

		var bytes fb.Bytes
		for i := 0; i < root.DatavLength(); i++ {
			root.Datav(&bytes, i)
			data := make([]byte, bytes.AliasLength())
			for i := 0; i < len(data); i++ {
				data[i] = byte(bytes.Alias(i))
			}
			p.consumeBytes(data)
		}

		p.f32 += root.F32()
		for i := 0; i < root.F32vLength(); i++ {
			p.f32 += root.F32v(i)
		}
		p.f64 += root.F64()
		for i := 0; i < root.F64vLength(); i++ {
			p.f64 += root.F64v(i)
		}

		var small fb.Small
		root.Object(&small)
		p.traverseFbSmall(&small)
		for i := 0; i < root.ObjectvLength(); i++ {
			root.Objectv(&small, i)
			p.traverseFbSmall(&small)
		}

		var pair1 fb.Map1Entry
		for i := 0; i < root.IndexLength(); i++ {
			root.Index(&pair1, i)
			p.consumeBytes(pair1.Key())
			p.u32 += uint32(pair1.Value())
		}

		var pair2 fb.Map2Entry
		for i := 0; i < root.ObjectsLength(); i++ {
			root.Objects(&pair2, i)
			p.u32 += uint32(pair2.Key())
			pair2.Value(&small)
			p.traverseFbSmall(&small)
		}

		var vec2d fb.Vec2D
		root.Matrix(&vec2d)
		p.traverseFbVec2D(&vec2d)

		var arrMap fb.ArrMap
		for i := 0; i < root.VectorLength(); i++ {
			root.Vector(&arrMap, i)
			p.traverseFbArrMap(&arrMap)
		}
		root.Arrays(&arrMap)
		p.traverseFbArrMap(&arrMap)
	}
*/
func (p *Junk) traversePbMessage(root protoreflect.Message) {
	consume := func(field protoreflect.FieldDescriptor, value protoreflect.Value) {
		switch field.Kind() {
		case protoreflect.MessageKind:
			p.traversePbMessage(value.Message())
		case protoreflect.BytesKind:
			p.consumeBytes(value.Bytes())
		case protoreflect.StringKind:
			p.consumeString(value.String())
		case protoreflect.DoubleKind:
			p.f64 += value.Float()
		case protoreflect.FloatKind:
			p.f32 += float32(value.Float())
		case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
			p.u64 += value.Uint()
		case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
			p.u64 += uint64(value.Int())
		case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
			p.u32 += uint32(value.Uint())
		case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
			p.u32 += uint32(value.Int())
		case protoreflect.BoolKind:
			p.consumeBool(value.Bool())
		case protoreflect.EnumKind:
			p.u32 += uint32(value.Enum())
		default:
			break
		}
	}

	fields := root.Descriptor().Fields()
	for i := 0; i < fields.Len(); i++ {
		field := fields.Get(i)
		if !root.Has(field) {
			continue
		}
		value := root.Get(field)
		if field.IsMap() {
			kField := field.MapKey()
			vField := field.MapValue()
			value.Map().Range(func(key protoreflect.MapKey, val protoreflect.Value) bool {
				consume(kField, key.Value())
				consume(vField, val)
				return true
			})
		} else if field.IsList() {
			list := value.List()
			switch field.Kind() {
			case protoreflect.MessageKind:
				for j := 0; j < list.Len(); j++ {
					p.traversePbMessage(list.Get(j).Message())
				}
			case protoreflect.BytesKind:
				for j := 0; j < list.Len(); j++ {
					p.consumeBytes(list.Get(j).Bytes())
				}
			case protoreflect.StringKind:
				for j := 0; j < list.Len(); j++ {
					p.consumeString(list.Get(j).String())
				}
			case protoreflect.DoubleKind:
				for j := 0; j < list.Len(); j++ {
					p.f64 += list.Get(j).Float()
				}
			case protoreflect.FloatKind:
				for j := 0; j < list.Len(); j++ {
					p.f32 += float32(list.Get(j).Float())
				}
			case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
				for j := 0; j < list.Len(); j++ {
					p.u64 += list.Get(j).Uint()
				}
			case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
				for j := 0; j < list.Len(); j++ {
					p.u64 += uint64(list.Get(j).Int())
				}
			case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
				for j := 0; j < list.Len(); j++ {
					p.u32 += uint32(list.Get(j).Uint())
				}
			case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
				for j := 0; j < list.Len(); j++ {
					p.u32 += uint32(list.Get(j).Int())
				}
			case protoreflect.BoolKind:
				for j := 0; j < list.Len(); j++ {
					p.consumeBool(list.Get(j).Bool())
				}
			case protoreflect.EnumKind:
				for j := 0; j < list.Len(); j++ {
					p.u32 += uint32(list.Get(j).Enum())
				}
			default:
				break
			}
		} else {
			consume(field, value)
		}
	}
}
