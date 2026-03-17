//go:build linux

package test

import (
	"os"
	"testing"

	pc "github.com/peterrk/protocache-go"
	"github.com/peterrk/protocache-go/reflect"
	"github.com/peterrk/protocache-go/reflect/compiler"
)

func TestReflection(t *testing.T) {
	raw, err := os.ReadFile("test.proto")
	assert(t, err == nil)

	proto, err := compiler.ParseProto(raw)
	assert(t, err == nil)

	var pool reflect.DescriptorPool
	assert(t, pool.Register(proto))

	root := pool.Find("test.Main")
	assert(t, root != nil)

	field := root.Lookup("f64")
	assert(t, field != nil)
	assert(t, !field.IsRepeated())
	assert(t, field.Value() == reflect.TYPE_FLOAT64)

	field = root.Lookup("strv")
	assert(t, field != nil)
	assert(t, field.IsRepeated())
	assert(t, field.Value() == reflect.TYPE_STRING)

	field = root.Lookup("mode")
	assert(t, field != nil)
	assert(t, !field.IsRepeated())
	assert(t, field.Value() == reflect.TYPE_ENUM)

	field = root.Lookup("object")
	assert(t, field != nil)
	assert(t, !field.IsRepeated())
	assert(t, field.Value() == reflect.TYPE_MESSAGE)
	object := pool.Find(field.ValueType())
	assert(t, object != nil)
	assert(t, object.Alias() == nil)
	field = object.Lookup("flag")
	assert(t, field != nil)
	assert(t, !field.IsRepeated())
	assert(t, field.Value() == reflect.TYPE_BOOL)

	field = root.Lookup("index")
	assert(t, field != nil)
	assert(t, field.IsRepeated())
	assert(t, field.IsMap())

	field = root.Lookup("matrix")
	assert(t, field != nil)
	assert(t, !field.IsRepeated())
	assert(t, field.Value() == reflect.TYPE_MESSAGE)
	object = field.ValueDescriptor()
	assert(t, object != nil)
	field = object.Alias()
	assert(t, field != nil)
	assert(t, field.IsRepeated())
	assert(t, !field.IsMap())
	assert(t, field.Value() == reflect.TYPE_MESSAGE)
	object = field.ValueDescriptor()
	assert(t, object != nil)
	field = object.Alias()
	assert(t, field != nil)
	assert(t, field.IsRepeated())
	assert(t, !field.IsMap())
	assert(t, field.Value() == reflect.TYPE_FLOAT32)

	field = root.Lookup("arrays")
	assert(t, field != nil)
	assert(t, !field.IsRepeated())
	assert(t, field.Value() == reflect.TYPE_MESSAGE)
	object = field.ValueDescriptor()
	assert(t, object != nil)
	field = object.Alias()
	assert(t, field != nil)
	assert(t, field.IsRepeated())
	assert(t, field.IsMap())
	assert(t, field.Key() == reflect.TYPE_STRING)
	assert(t, field.Value() == reflect.TYPE_MESSAGE)
	object = field.ValueDescriptor()
	assert(t, object != nil)
	field = object.Alias()
	assert(t, field != nil)
	assert(t, field.IsRepeated())
	assert(t, !field.IsMap())
	assert(t, field.Value() == reflect.TYPE_FLOAT32)
}

/*
	func TestBenchmarkReflect(t *testing.T) {
		raw, err := os.ReadFile("test.proto")
		assert(t, err == nil)

		proto, err := compiler.ParseProto(raw)
		assert(t, err == nil)

		var pool reflect.DescriptorPool
		assert(t, pool.Register(proto))

		descriptor := pool.Find("test.Main")
		assert(t, descriptor != nil)

		var junk Junk
		raw, err = os.ReadFile("test.pc")
		if err != nil {
			fmt.Println(err)
			return
		}
		root := pc.AsMessage(raw)
		junk.traversePcMessage(descriptor, root)
		fmt.Println(junk.fuse())
	}
*/

func BenchmarkProtoCacheReflect(b *testing.B) {
	raw, err := os.ReadFile("test.proto")
	if err != nil {
		b.Fatal(err)
	}
	proto, err := compiler.ParseProto(raw)
	if err != nil {
		b.Fatal(err)
	}
	var pool reflect.DescriptorPool
	if !pool.Register(proto) {
		b.Fatal("fail to register schema")
	}
	descriptor := pool.Find("test.Main")
	if descriptor == nil {
		b.Fatal("fail to get root")
	}
	raw, err = os.ReadFile("test.pc")
	if err != nil {
		b.Fatal(err)
	}
	var junk Junk
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		root := pc.AsMessage(raw)
		junk.traversePcMessage(descriptor, root)
	}
	benchmarkFuse = junk.fuse()
}

func (p *Junk) traversePcMessage(descriptor *reflect.Descriptor, root pc.Message) {
	descriptor.Traverse(func(name string, field *reflect.Field) bool {
		unit := root.GetField(field.Id())
		if unit.IsValid() {
			p.traversePcField(field, unit)
		}
		return true
	})
}

func (p *Junk) traversePcField(descriptor *reflect.Field, field pc.Field) {
	if descriptor.IsMap() {
		pack := field.GetMap()
		for i := uint32(0); i < pack.Size(); i++ {
			key := pack.Key(i)
			switch descriptor.Key() {
			case reflect.TYPE_STRING:
				p.consumeString(key.GetString())
			case reflect.TYPE_UINT64:
				p.u64 += key.GetUint64()
			case reflect.TYPE_UINT32:
				p.u32 += key.GetUint32()
			case reflect.TYPE_INT64:
				p.u64 += uint64(key.GetInt64())
			case reflect.TYPE_INT32:
				p.u32 += uint32(key.GetInt32())
			}
			p.accessPcField(descriptor, pack.Value(i))
		}
	} else if descriptor.IsRepeated() {
		switch descriptor.Value() {
		case reflect.TYPE_MESSAGE:
			array := field.GetArray()
			subtype := descriptor.ValueDescriptor()
			if alias := subtype.Alias(); alias != nil {
				for i := uint32(0); i < array.Size(); i++ {
					p.traversePcField(alias, array.Get(i))
				}
			} else {
				for i := uint32(0); i < array.Size(); i++ {
					unit := array.Get(i)
					p.traversePcMessage(subtype, unit.GetMessage())
				}
			}
		case reflect.TYPE_BYTES:
			array := field.GetArray()
			for i := uint32(0); i < array.Size(); i++ {
				unit := array.Get(i)
				p.consumeBytes(unit.GetBytes())
			}
		case reflect.TYPE_STRING:
			array := field.GetArray()
			for i := uint32(0); i < array.Size(); i++ {
				unit := array.Get(i)
				p.consumeString(unit.GetString())
			}
		case reflect.TYPE_FLOAT64:
			for _, v := range field.GetFloat64Array() {
				p.f64 += v
			}
		case reflect.TYPE_FLOAT32:
			for _, v := range field.GetFloat32Array() {
				p.f32 += v
			}
		case reflect.TYPE_UINT64:
			for _, v := range field.GetUint64Array() {
				p.u64 += v
			}
		case reflect.TYPE_UINT32:
			for _, v := range field.GetUint32Array() {
				p.u32 += v
			}
		case reflect.TYPE_INT64:
			for _, v := range field.GetInt64Array() {
				p.u64 += uint64(v)
			}
		case reflect.TYPE_INT32:
			for _, v := range field.GetInt32Array() {
				p.u32 += uint32(v)
			}
		case reflect.TYPE_BOOL:
			for _, v := range field.GetBoolArray() {
				p.consumeBool(v)
			}
		case reflect.TYPE_ENUM:
			for _, v := range field.GetEnumValueArray() {
				p.u32 += uint32(v)
			}
		}

	} else {
		p.accessPcField(descriptor, field)
	}
}

func (p *Junk) accessPcField(descriptor *reflect.Field, field pc.Field) {
	switch descriptor.Value() {
	case reflect.TYPE_MESSAGE:
		subtype := descriptor.ValueDescriptor()
		if alias := subtype.Alias(); alias != nil {
			p.traversePcField(alias, field)
		} else {
			p.traversePcMessage(subtype, field.GetMessage())
		}
	case reflect.TYPE_BYTES:
		p.consumeBytes(field.GetBytes())
	case reflect.TYPE_STRING:
		p.consumeString(field.GetString())
	case reflect.TYPE_FLOAT64:
		p.f64 += field.GetFloat64()
	case reflect.TYPE_FLOAT32:
		p.f32 += field.GetFloat32()
	case reflect.TYPE_UINT64:
		p.u64 += field.GetUint64()
	case reflect.TYPE_UINT32:
		p.u32 += field.GetUint32()
	case reflect.TYPE_INT64:
		p.u64 += uint64(field.GetInt64())
	case reflect.TYPE_INT32:
		p.u32 += uint32(field.GetInt32())
	case reflect.TYPE_BOOL:
		p.consumeBool(field.GetBool())
	case reflect.TYPE_ENUM:
		p.u32 += uint32(field.GetEnumValue())
	}
}
