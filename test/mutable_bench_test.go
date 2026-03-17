package test

import (
	"os"
	"testing"

	"github.com/peterrk/protocache-go"
	"github.com/peterrk/protocache-go/test/pb"
	"github.com/peterrk/protocache-go/test/pc"
	"google.golang.org/protobuf/encoding/protojson"
)

var (
	benchmarkMutableBytes []byte
	benchmarkMutableFuse  uint64
)

func loadMutableRaw(b *testing.B) []byte {
	b.Helper()
	raw, err := os.ReadFile("test.json")
	if err != nil {
		b.Fatal(err)
	}
	message := &pb.Main{}
	if err := protojson.Unmarshal(raw, message); err != nil {
		b.Fatal(err)
	}
	raw, err = protocache.Serialize(message)
	if err != nil {
		b.Fatal(err)
	}
	return raw
}

func loadPBMain(b *testing.B) *pb.Main {
	b.Helper()
	raw, err := os.ReadFile("test.json")
	if err != nil {
		b.Fatal(err)
	}
	message := &pb.Main{}
	if err := protojson.Unmarshal(raw, message); err != nil {
		b.Fatal(err)
	}
	return message
}

func (p *Junk) traverseMutableSmall(root *pc.SmallEX) {
	if root == nil {
		return
	}
	p.u32 += uint32(root.GetI32())
	p.consumeBool(root.GetFlag())
	p.consumeString(root.GetStr())
}

func (p *Junk) traverseMutableMain(root *pc.MainEX) {
	p.u32 += uint32(root.GetI32())
	for _, v := range root.GetI32V() {
		p.u32 += uint32(v)
	}
	p.consumeString(root.GetStr())
	p.consumeBytes(root.GetData())
	for _, v := range root.GetStrv() {
		p.consumeString(v)
	}

	p.traverseMutableSmall(root.GetObject())
	for _, v := range root.GetObjectv() {
		p.traverseMutableSmall(v)
	}

	for key, val := range root.GetIndex() {
		p.consumeString(key)
		p.u32 += uint32(val)
	}

	for key, val := range root.GetObjects() {
		p.u32 += uint32(key)
		p.traverseMutableSmall(val)
	}

	for _, one := range root.GetMatrix() {
		for _, v := range one {
			p.f32 += v
		}
	}

	for _, one := range root.GetVector() {
		for key, val := range one {
			p.consumeString(key)
			for _, v := range val {
				p.f32 += v
			}
		}
	}

	for key, val := range root.GetArrays() {
		p.consumeString(key)
		for _, v := range val {
			p.f32 += v
		}
	}

	for _, v := range root.GetModev() {
		p.u32 += uint32(v)
	}
}

func (p *Junk) traverseMutableRootScalars(root *pc.MainEX) {
	p.u32 += uint32(root.GetI32())
	p.consumeString(root.GetStr())
	p.consumeBytes(root.GetData())
}

func BenchmarkProtoCacheEX(b *testing.B) {
	raw := loadMutableRaw(b)
	var junk Junk

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		root := pc.TO_MainEX(raw)
		junk.traverseMutableMain(root)
	}
	benchmarkMutableFuse = junk.fuse()
}

func BenchmarkFullSerializeEX(b *testing.B) {
	raw := loadMutableRaw(b)
	var junk Junk
	root := pc.TO_MainEX(raw)
	junk.traverseMutableMain(root)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		out, err := root.Serialize()
		if err != nil {
			b.Fatal(err)
		}
		benchmarkMutableBytes = out
	}
	benchmarkMutableFuse = junk.fuse()
}

func BenchmarkPartlySerializeEX(b *testing.B) {
	raw := loadMutableRaw(b)
	var junk Junk
	root := pc.TO_MainEX(raw)
	junk.traverseMutableRootScalars(root)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		out, err := root.Serialize()
		if err != nil {
			b.Fatal(err)
		}
		benchmarkMutableBytes = out
	}
	benchmarkMutableFuse = junk.fuse()
}

func BenchmarkSerialize(b *testing.B) {
	message := loadPBMain(b)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		out, err := protocache.Serialize(message)
		if err != nil {
			b.Fatal(err)
		}
		benchmarkMutableBytes = out
	}
}
