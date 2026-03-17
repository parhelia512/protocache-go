package pc

import (
	"github.com/peterrk/protocache-go"
)

type Mode protocache.EnumValue

const (
	Mode_MODE_A Mode = 0
	Mode_MODE_B Mode = 1
	Mode_MODE_C Mode = 2
)

const (
	_FIELD_Small_i32  uint16 = 0
	_FIELD_Small_flag uint16 = 1
	_FIELD_Small_str  uint16 = 3
	_FIELD_Small_junk uint16 = 4
	_FIELD_TOTAL_Small uint16 = 5
)

type Small struct{ core protocache.Message }

func AS_Small(data []byte) Small { return Small{core: protocache.AsMessage(data)} }

func (m *Small) IsValid() bool { return m.core.IsValid() }

func (m *Small) GetI32() int32 {
	field := m.core.GetField(_FIELD_Small_i32)
	return field.GetInt32()
}

func (m *Small) GetFlag() bool {
	field := m.core.GetField(_FIELD_Small_flag)
	return field.GetBool()
}

func (m *Small) GetStr() string {
	field := m.core.GetField(_FIELD_Small_str)
	return field.GetString()
}

func (m *Small) GetJunk() int64 {
	field := m.core.GetField(_FIELD_Small_junk)
	return field.GetInt64()
}

type Vec2D_Vec1D = protocache.Float32Array

func AS_Vec2D_Vec1D(data []byte) Vec2D_Vec1D {
	return protocache.AsFloat32Array(data)
}

type Vec2D = ARRAY_Vec2D_Vec1D

func AS_Vec2D(data []byte) Vec2D {
	return AS_ARRAY_Vec2D_Vec1D(data)
}

type ArrMap_Array = protocache.Float32Array

func AS_ArrMap_Array(data []byte) ArrMap_Array {
	return protocache.AsFloat32Array(data)
}

type ArrMap = MAP_string_ArrMap_Array

func AS_ArrMap(data []byte) ArrMap {
	return AS_MAP_string_ArrMap_Array(data)
}

const (
	_FIELD_Main_i32     uint16 = 0
	_FIELD_Main_u32     uint16 = 1
	_FIELD_Main_i64     uint16 = 2
	_FIELD_Main_u64     uint16 = 3
	_FIELD_Main_flag    uint16 = 4
	_FIELD_Main_mode    uint16 = 5
	_FIELD_Main_str     uint16 = 6
	_FIELD_Main_data    uint16 = 7
	_FIELD_Main_f32     uint16 = 8
	_FIELD_Main_f64     uint16 = 9
	_FIELD_Main_object  uint16 = 10
	_FIELD_Main_i32v    uint16 = 11
	_FIELD_Main_u64v    uint16 = 12
	_FIELD_Main_strv    uint16 = 13
	_FIELD_Main_datav   uint16 = 14
	_FIELD_Main_f32v    uint16 = 15
	_FIELD_Main_f64v    uint16 = 16
	_FIELD_Main_flags   uint16 = 17
	_FIELD_Main_objectv uint16 = 18
	_FIELD_Main_t_u32   uint16 = 19
	_FIELD_Main_t_i32   uint16 = 20
	_FIELD_Main_t_s32   uint16 = 21
	_FIELD_Main_t_u64   uint16 = 22
	_FIELD_Main_t_i64   uint16 = 23
	_FIELD_Main_t_s64   uint16 = 24
	_FIELD_Main_index   uint16 = 25
	_FIELD_Main_objects uint16 = 26
	_FIELD_Main_matrix  uint16 = 27
	_FIELD_Main_vector  uint16 = 28
	_FIELD_Main_arrays  uint16 = 29
	_FIELD_Main_modev   uint16 = 31
	_FIELD_TOTAL_Main   uint16 = 32
)

type Main struct{ core protocache.Message }

func AS_Main(data []byte) Main { return Main{core: protocache.AsMessage(data)} }

func (m *Main) IsValid() bool { return m.core.IsValid() }

func (m *Main) GetI32() int32 {
	field := m.core.GetField(_FIELD_Main_i32)
	return field.GetInt32()
}

func (m *Main) GetU32() uint32 {
	field := m.core.GetField(_FIELD_Main_u32)
	return field.GetUint32()
}

func (m *Main) GetI64() int64 {
	field := m.core.GetField(_FIELD_Main_i64)
	return field.GetInt64()
}

func (m *Main) GetU64() uint64 {
	field := m.core.GetField(_FIELD_Main_u64)
	return field.GetUint64()
}

func (m *Main) GetFlag() bool {
	field := m.core.GetField(_FIELD_Main_flag)
	return field.GetBool()
}

func (m *Main) GetMode() Mode {
	field := m.core.GetField(_FIELD_Main_mode)
	return Mode(field.GetEnumValue())
}

func (m *Main) GetStr() string {
	field := m.core.GetField(_FIELD_Main_str)
	return field.GetString()
}

func (m *Main) GetData() []byte {
	field := m.core.GetField(_FIELD_Main_data)
	return field.GetBytes()
}

func (m *Main) GetF32() float32 {
	field := m.core.GetField(_FIELD_Main_f32)
	return field.GetFloat32()
}

func (m *Main) GetF64() float64 {
	field := m.core.GetField(_FIELD_Main_f64)
	return field.GetFloat64()
}

func (m *Main) GetObject() Small {
	field := m.core.GetField(_FIELD_Main_object)
	return AS_Small(field.GetObject())
}
func (m *Main) HasObject() bool {
	return m.core.HasField(_FIELD_Main_object)
}

func (m *Main) GetI32V() []int32 {
	field := m.core.GetField(_FIELD_Main_i32v)
	return field.GetInt32Array()
}

func (m *Main) GetU64V() []uint64 {
	field := m.core.GetField(_FIELD_Main_u64v)
	return field.GetUint64Array()
}

func (m *Main) GetStrv() protocache.StringArray {
	field := m.core.GetField(_FIELD_Main_strv)
	return protocache.AsStringArray(field.GetObject())
}

func (m *Main) GetDatav() protocache.BytesArray {
	field := m.core.GetField(_FIELD_Main_datav)
	return protocache.AsBytesArray(field.GetObject())
}

func (m *Main) GetF32V() []float32 {
	field := m.core.GetField(_FIELD_Main_f32v)
	return field.GetFloat32Array()
}

func (m *Main) GetF64V() []float64 {
	field := m.core.GetField(_FIELD_Main_f64v)
	return field.GetFloat64Array()
}

func (m *Main) GetFlags() []bool {
	field := m.core.GetField(_FIELD_Main_flags)
	return field.GetBoolArray()
}

func (m *Main) GetObjectv() ARRAY_Small {
	field := m.core.GetField(_FIELD_Main_objectv)
	return AS_ARRAY_Small(field.GetObject())
}

func (m *Main) GetTU32() uint32 {
	field := m.core.GetField(_FIELD_Main_t_u32)
	return field.GetUint32()
}

func (m *Main) GetTI32() int32 {
	field := m.core.GetField(_FIELD_Main_t_i32)
	return field.GetInt32()
}

func (m *Main) GetTS32() int32 {
	field := m.core.GetField(_FIELD_Main_t_s32)
	return field.GetInt32()
}

func (m *Main) GetTU64() uint64 {
	field := m.core.GetField(_FIELD_Main_t_u64)
	return field.GetUint64()
}

func (m *Main) GetTI64() int64 {
	field := m.core.GetField(_FIELD_Main_t_i64)
	return field.GetInt64()
}

func (m *Main) GetTS64() int64 {
	field := m.core.GetField(_FIELD_Main_t_s64)
	return field.GetInt64()
}

func (m *Main) GetIndex() MAP_string_int32 {
	field := m.core.GetField(_FIELD_Main_index)
	return AS_MAP_string_int32(field.GetObject())
}

func (m *Main) GetObjects() MAP_int32_Small {
	field := m.core.GetField(_FIELD_Main_objects)
	return AS_MAP_int32_Small(field.GetObject())
}

func (m *Main) GetMatrix() Vec2D {
	field := m.core.GetField(_FIELD_Main_matrix)
	return AS_Vec2D(field.GetObject())
}

func (m *Main) GetVector() ARRAY_ArrMap {
	field := m.core.GetField(_FIELD_Main_vector)
	return AS_ARRAY_ArrMap(field.GetObject())
}

func (m *Main) GetArrays() ArrMap {
	field := m.core.GetField(_FIELD_Main_arrays)
	return AS_ArrMap(field.GetObject())
}

func (m *Main) GetModev() []Mode {
	field := m.core.GetField(_FIELD_Main_modev)
	return protocache.CastEnumArray[Mode](field.GetEnumValueArray())
}

const (
	_FIELD_CyclicA_value  uint16 = 0
	_FIELD_CyclicA_cyclic uint16 = 1
	_FIELD_TOTAL_CyclicA uint16 = 2
)

type CyclicA struct{ core protocache.Message }

func AS_CyclicA(data []byte) CyclicA { return CyclicA{core: protocache.AsMessage(data)} }

func (m *CyclicA) IsValid() bool { return m.core.IsValid() }

func (m *CyclicA) GetValue() int32 {
	field := m.core.GetField(_FIELD_CyclicA_value)
	return field.GetInt32()
}

func (m *CyclicA) GetCyclic() CyclicB {
	field := m.core.GetField(_FIELD_CyclicA_cyclic)
	return AS_CyclicB(field.GetObject())
}
func (m *CyclicA) HasCyclic() bool {
	return m.core.HasField(_FIELD_CyclicA_cyclic)
}

const (
	_FIELD_CyclicB_value  uint16 = 0
	_FIELD_CyclicB_cyclic uint16 = 1
	_FIELD_TOTAL_CyclicB uint16 = 2
)

type CyclicB struct{ core protocache.Message }

func AS_CyclicB(data []byte) CyclicB { return CyclicB{core: protocache.AsMessage(data)} }

func (m *CyclicB) IsValid() bool { return m.core.IsValid() }

func (m *CyclicB) GetValue() int32 {
	field := m.core.GetField(_FIELD_CyclicB_value)
	return field.GetInt32()
}

func (m *CyclicB) GetCyclic() CyclicA {
	field := m.core.GetField(_FIELD_CyclicB_cyclic)
	return AS_CyclicA(field.GetObject())
}
func (m *CyclicB) HasCyclic() bool {
	return m.core.HasField(_FIELD_CyclicB_cyclic)
}

const (
	_FIELD_Deprecated_Valid_val uint16 = 0
	_FIELD_TOTAL_Deprecated_Valid uint16 = 1
)

type Deprecated_Valid struct{ core protocache.Message }

func AS_Deprecated_Valid(data []byte) Deprecated_Valid {
	return Deprecated_Valid{core: protocache.AsMessage(data)}
}

func (m *Deprecated_Valid) IsValid() bool { return m.core.IsValid() }

func (m *Deprecated_Valid) GetVal() int32 {
	field := m.core.GetField(_FIELD_Deprecated_Valid_val)
	return field.GetInt32()
}

const (
	_FIELD_Deprecated_junk uint16 = 0
	_FIELD_TOTAL_Deprecated uint16 = 1
)

type Deprecated struct{ core protocache.Message }

func AS_Deprecated(data []byte) Deprecated { return Deprecated{core: protocache.AsMessage(data)} }

func (m *Deprecated) IsValid() bool { return m.core.IsValid() }

func (m *Deprecated) GetJunk() int32 {
	field := m.core.GetField(_FIELD_Deprecated_junk)
	return field.GetInt32()
}

type ModeDict_Value = protocache.EnumArray[Mode]

func AS_ModeDict_Value(data []byte) ModeDict_Value {
	return protocache.AsEnumArray[Mode](data)
}

type ModeDict = MAP_int32_ModeDict_Value

func AS_ModeDict(data []byte) ModeDict {
	return AS_MAP_int32_ModeDict_Value(data)
}

type ARRAY_ArrMap struct{ core protocache.Array }

func AS_ARRAY_ArrMap(data []byte) ARRAY_ArrMap {
	return ARRAY_ArrMap{core: protocache.AsArray(data)}
}

func (x *ARRAY_ArrMap) Get(i uint32) ArrMap {
	field := x.core.Get(i)
	return AS_ArrMap(field.GetObject())
}

func (x *ARRAY_ArrMap) IsValid() bool { return x.core.IsValid() }

func (x *ARRAY_ArrMap) Size() uint32 { return x.core.Size() }

type ARRAY_Small struct{ core protocache.Array }

func AS_ARRAY_Small(data []byte) ARRAY_Small {
	return ARRAY_Small{core: protocache.AsArray(data)}
}

func (x *ARRAY_Small) Get(i uint32) Small {
	field := x.core.Get(i)
	return AS_Small(field.GetObject())
}

func (x *ARRAY_Small) IsValid() bool { return x.core.IsValid() }

func (x *ARRAY_Small) Size() uint32 { return x.core.Size() }

type ARRAY_Vec2D_Vec1D struct{ core protocache.Array }

func AS_ARRAY_Vec2D_Vec1D(data []byte) ARRAY_Vec2D_Vec1D {
	return ARRAY_Vec2D_Vec1D{core: protocache.AsArray(data)}
}

func (x *ARRAY_Vec2D_Vec1D) Get(i uint32) Vec2D_Vec1D {
	field := x.core.Get(i)
	return AS_Vec2D_Vec1D(field.GetObject())
}

func (x *ARRAY_Vec2D_Vec1D) IsValid() bool { return x.core.IsValid() }

func (x *ARRAY_Vec2D_Vec1D) Size() uint32 { return x.core.Size() }

type MAP_int32_ModeDict_Value struct{ core protocache.Map }

func AS_MAP_int32_ModeDict_Value(data []byte) MAP_int32_ModeDict_Value {
	return MAP_int32_ModeDict_Value{core: protocache.AsMap(data)}
}

func (x *MAP_int32_ModeDict_Value) Key(i uint32) int32 {
	field := x.core.Key(i)
	return field.GetInt32()
}

func (x *MAP_int32_ModeDict_Value) Value(i uint32) ModeDict_Value {
	field := x.core.Value(i)
	return AS_ModeDict_Value(field.GetObject())
}

func (x *MAP_int32_ModeDict_Value) Find(key int32) (ModeDict_Value, bool) {
	field := x.core.FindByInt32(key)
	return AS_ModeDict_Value(field.GetObject()), field.IsValid()
}

func (x *MAP_int32_ModeDict_Value) IsValid() bool { return x.core.IsValid() }

func (x *MAP_int32_ModeDict_Value) Size() uint32 { return x.core.Size() }

type MAP_int32_Small struct{ core protocache.Map }

func AS_MAP_int32_Small(data []byte) MAP_int32_Small {
	return MAP_int32_Small{core: protocache.AsMap(data)}
}

func (x *MAP_int32_Small) Key(i uint32) int32 {
	field := x.core.Key(i)
	return field.GetInt32()
}

func (x *MAP_int32_Small) Value(i uint32) Small {
	field := x.core.Value(i)
	return AS_Small(field.GetObject())
}

func (x *MAP_int32_Small) Find(key int32) (Small, bool) {
	field := x.core.FindByInt32(key)
	return AS_Small(field.GetObject()), field.IsValid()
}

func (x *MAP_int32_Small) IsValid() bool { return x.core.IsValid() }

func (x *MAP_int32_Small) Size() uint32 { return x.core.Size() }

type MAP_string_ArrMap_Array struct{ core protocache.Map }

func AS_MAP_string_ArrMap_Array(data []byte) MAP_string_ArrMap_Array {
	return MAP_string_ArrMap_Array{core: protocache.AsMap(data)}
}

func (x *MAP_string_ArrMap_Array) Key(i uint32) string {
	field := x.core.Key(i)
	return field.GetString()
}

func (x *MAP_string_ArrMap_Array) Value(i uint32) ArrMap_Array {
	field := x.core.Value(i)
	return AS_ArrMap_Array(field.GetObject())
}

func (x *MAP_string_ArrMap_Array) Find(key string) (ArrMap_Array, bool) {
	field := x.core.FindByString(key)
	return AS_ArrMap_Array(field.GetObject()), field.IsValid()
}

func (x *MAP_string_ArrMap_Array) IsValid() bool { return x.core.IsValid() }

func (x *MAP_string_ArrMap_Array) Size() uint32 { return x.core.Size() }

type MAP_string_int32 struct{ core protocache.Map }

func AS_MAP_string_int32(data []byte) MAP_string_int32 {
	return MAP_string_int32{core: protocache.AsMap(data)}
}

func (x *MAP_string_int32) Key(i uint32) string {
	field := x.core.Key(i)
	return field.GetString()
}

func (x *MAP_string_int32) Value(i uint32) int32 {
	field := x.core.Value(i)
	return field.GetInt32()
}

func (x *MAP_string_int32) Find(key string) (int32, bool) {
	field := x.core.FindByString(key)
	return field.GetInt32(), field.IsValid()
}

func (x *MAP_string_int32) IsValid() bool { return x.core.IsValid() }

func (x *MAP_string_int32) Size() uint32 { return x.core.Size() }
