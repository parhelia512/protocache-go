package test

import (
	"bytes"
	"os"
	"testing"

	"github.com/peterrk/protocache-go"
	"github.com/peterrk/protocache-go/test/pb"
	"github.com/peterrk/protocache-go/test/pc"
	"google.golang.org/protobuf/encoding/protojson"
)

func loadMain(t *testing.T) []byte {
	raw, err := os.ReadFile("test.json")
	assert(t, err == nil)

	message := &pb.Main{}
	err = protojson.Unmarshal(raw, message)
	assert(t, err == nil)

	raw, err = protocache.Serialize(message)
	assert(t, err == nil)
	return raw
}

func TestEXReplayWithoutAccess(t *testing.T) {
	raw := loadMain(t)

	root := pc.TO_MainEX(raw)
	out, err := root.Serialize()
	assert(t, err == nil)
	assert(t, bytes.Equal(raw, out))
}

func TestEXModifyOneFieldPreservesUntouchedRawParts(t *testing.T) {
	raw := loadMain(t)

	root := pc.TO_MainEX(raw)
	root.SetStr("patched")
	out, err := root.Serialize()
	assert(t, err == nil)

	before := protocache.AsMessage(raw)
	after := protocache.AsMessage(out)
	assert(t, bytes.Equal(before.RawPart(13), after.RawPart(13)))

	view := pc.AS_Main(out)
	obj := view.GetObject()
	assert(t, view.GetStr() == "patched")
	assert(t, obj.GetStr() == "tmp")
}

func TestEXModifyNestedMessage(t *testing.T) {
	raw := loadMain(t)

	root := pc.TO_MainEX(raw)
	obj := root.GetObject()
	obj.SetStr("nested")
	obj.SetI32(777)

	out, err := root.Serialize()
	assert(t, err == nil)

	view := pc.AS_Main(out)
	objView := view.GetObject()
	assert(t, objView.GetStr() == "nested")
	assert(t, objView.GetI32() == 777)
	assert(t, len(view.GetI32V()) == 2)
}

func TestEXModifyRepeatedScalarAndEnum(t *testing.T) {
	raw := loadMain(t)

	root := pc.TO_MainEX(raw)
	root.SetI32V([]int32{9, 8, 7})
	root.SetStrv([]string{"left", "right"})
	root.SetModev([]pc.Mode{pc.Mode_MODE_A, pc.Mode_MODE_B})

	out, err := root.Serialize()
	assert(t, err == nil)

	view := pc.AS_Main(out)
	i32v := view.GetI32V()
	assert(t, len(i32v) == 3)
	assert(t, i32v[0] == 9 && i32v[2] == 7)

	strv := view.GetStrv()
	assert(t, strv.Size() == 2)
	assert(t, strv.Get(0) == "left")
	assert(t, strv.Get(1) == "right")

	modev := view.GetModev()
	assert(t, len(modev) == 2)
	assert(t, modev[0] == pc.Mode_MODE_A)
	assert(t, modev[1] == pc.Mode_MODE_B)
}

func TestMessageEXVisitUsesExplicitFieldTotal(t *testing.T) {
	var msg protocache.MessageEX
	msg.Init(nil)

	assert(t, !msg.HasBase())
	assert(t, !msg.IsValid())

	assert(t, !msg.IsVisited(2, 2))
	msg.Visit(2, 2)
	assert(t, !msg.IsVisited(2, 2))

	assert(t, !msg.IsVisited(1, 2))
	msg.Visit(1, 2)
	assert(t, msg.IsVisited(1, 2))
}

func TestMessageEXHasBaseTracksInitData(t *testing.T) {
	var empty protocache.MessageEX
	assert(t, !empty.HasBase())
	assert(t, !empty.IsValid())

	var withBase protocache.MessageEX
	withBase.Init([]byte{0, 0, 0, 0})
	assert(t, withBase.HasBase())
	assert(t, withBase.IsValid())
}

func TestEXModifyRepeatedMessage(t *testing.T) {
	raw := loadMain(t)

	root := pc.TO_MainEX(raw)
	items := root.GetObjectv()
	assert(t, len(items) == 3)
	items[0].SetStr("first")
	items[2].SetI32(303)
	root.SetObjectv(items[:2])

	out, err := root.Serialize()
	assert(t, err == nil)

	view := pc.AS_Main(out)
	objs := view.GetObjectv()
	assert(t, objs.Size() == 2)
	first := objs.Get(0)
	second := objs.Get(1)
	assert(t, first.GetStr() == "first")
	assert(t, second.GetFlag())
}

func TestEXModifyMapFields(t *testing.T) {
	raw := loadMain(t)

	root := pc.TO_MainEX(raw)
	index := root.GetIndex()
	index["abc-1"] = 101
	index["new-key"] = 202
	root.SetIndex(index)

	objects := root.GetObjects()
	obj := objects[2]
	obj.SetStr("updated")
	objects[9] = pc.TO_SmallEX(nil)
	objects[9].SetI32(909)
	root.SetObjects(objects)

	out, err := root.Serialize()
	assert(t, err == nil)

	view := pc.AS_Main(out)
	gotIndex := view.GetIndex()
	val, found := gotIndex.Find("new-key")
	assert(t, found && val == 202)
	val, found = gotIndex.Find("abc-1")
	assert(t, found && val == 101)

	gotObjects := view.GetObjects()
	obj2, found := gotObjects.Find(2)
	assert(t, found && obj2.GetStr() == "updated")
	obj9, found := gotObjects.Find(9)
	assert(t, found && obj9.GetI32() == 909)
}

func TestEXModifyAliasFields(t *testing.T) {
	raw := loadMain(t)

	root := pc.TO_MainEX(raw)

	matrix := root.GetMatrix()
	matrix[0][0] = 101
	matrix = append(matrix, pc.Vec2D_Vec1DEX{10, 11})
	root.SetMatrix(matrix)

	vector := root.GetVector()
	vector[0]["lv1"] = pc.ArrMap_ArrayEX{111, 112}
	vector = append(vector, pc.ArrMapEX{
		"lv9": pc.ArrMap_ArrayEX{901, 902},
	})
	root.SetVector(vector)

	arrays := root.GetArrays()
	arrays["lv4"] = pc.ArrMap_ArrayEX{401, 402, 403}
	arrays["lv6"] = pc.ArrMap_ArrayEX{601, 602}
	root.SetArrays(arrays)

	out, err := root.Serialize()
	assert(t, err == nil)

	view := pc.AS_Main(out)

	gotMatrix := view.GetMatrix()
	assert(t, gotMatrix.Size() == 4)
	row0 := gotMatrix.Get(0)
	assert(t, row0.Size() == 3)
	assert(t, row0.Get(0) == 101)
	row3 := gotMatrix.Get(3)
	assert(t, row3.Size() == 2)
	assert(t, row3.Get(1) == 11)

	gotVector := view.GetVector()
	assert(t, gotVector.Size() == 3)
	vec0 := gotVector.Get(0)
	lv1, found := vec0.Find("lv1")
	assert(t, found)
	assert(t, lv1.Size() == 2)
	assert(t, lv1.Get(0) == 111)
	vec2 := gotVector.Get(2)
	lv9, found := vec2.Find("lv9")
	assert(t, found)
	assert(t, lv9.Size() == 2)
	assert(t, lv9.Get(1) == 902)

	gotArrays := view.GetArrays()
	lv4, found := gotArrays.Find("lv4")
	assert(t, found)
	assert(t, lv4.Size() == 3)
	assert(t, lv4.Get(2) == 403)
	lv6, found := gotArrays.Find("lv6")
	assert(t, found)
	assert(t, lv6.Size() == 2)
	assert(t, lv6.Get(0) == 601)
}
