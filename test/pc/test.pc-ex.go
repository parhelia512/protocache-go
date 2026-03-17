package pc

import (
	"github.com/peterrk/protocache-go"
)

type SmallEX struct {
	__    protocache.MessageEX
	fI32  int32
	fFlag bool
	fStr  string
	fJunk int64
}

func TO_SmallEX(data []byte) *SmallEX {
	out := &SmallEX{}
	out.__.Init(data)
	return out
}

func (m *SmallEX) HasBase() bool { return m.__.HasBase() }

func (m *SmallEX) IsValid() bool { return m.HasBase() }

func (m *SmallEX) Serialize() ([]byte, error) {
	words, err := m.serializeWords()
	if err != nil {
		return nil, err
	}
	return protocache.WordsToBytes(words), nil
}

func (m *SmallEX) serializeWords() ([]uint32, error) {
	if m == nil {
		return []uint32{0}, nil
	}
	parts := make([][]uint32, 5)
	if m.__.IsVisited(_FIELD_Small_i32, _FIELD_TOTAL_Small) {
		parts[0] = protocache.EncodeInt32(m.fI32)
	} else if raw := m.__.RawPart(_FIELD_Small_i32); len(raw) != 0 {
		parts[0] = protocache.BytesToWords(raw)
	}
	if m.__.IsVisited(_FIELD_Small_flag, _FIELD_TOTAL_Small) {
		parts[1] = protocache.EncodeBool(m.fFlag)
	} else if raw := m.__.RawPart(_FIELD_Small_flag); len(raw) != 0 {
		parts[1] = protocache.BytesToWords(raw)
	}
	if raw := m.__.RawPart(2); len(raw) != 0 {
		parts[2] = protocache.BytesToWords(raw)
	}
	if m.__.IsVisited(_FIELD_Small_str, _FIELD_TOTAL_Small) {
		part, err := protocache.EncodeString(m.fStr)
		if err != nil {
			return nil, err
		}
		parts[3] = part
	} else if raw := m.__.RawPart(_FIELD_Small_str); len(raw) != 0 {
		parts[3] = protocache.BytesToWords(raw)
	}
	if m.__.IsVisited(_FIELD_Small_junk, _FIELD_TOTAL_Small) {
		parts[4] = protocache.EncodeInt64(m.fJunk)
	} else if raw := m.__.RawPart(_FIELD_Small_junk); len(raw) != 0 {
		parts[4] = protocache.BytesToWords(raw)
	}
	return protocache.EncodeMessageParts(parts)
}

func (m *SmallEX) GetI32() int32 {
	if !m.__.IsVisited(_FIELD_Small_i32, _FIELD_TOTAL_Small) {
		field := m.__.RawField(_FIELD_Small_i32)
		m.fI32 = field.GetInt32()
		m.__.Visit(_FIELD_Small_i32, _FIELD_TOTAL_Small)
	}
	return m.fI32
}

func (m *SmallEX) SetI32(v int32) {
	m.fI32 = v
	m.__.Visit(_FIELD_Small_i32, _FIELD_TOTAL_Small)
}

func (m *SmallEX) GetFlag() bool {
	if !m.__.IsVisited(_FIELD_Small_flag, _FIELD_TOTAL_Small) {
		field := m.__.RawField(_FIELD_Small_flag)
		m.fFlag = field.GetBool()
		m.__.Visit(_FIELD_Small_flag, _FIELD_TOTAL_Small)
	}
	return m.fFlag
}

func (m *SmallEX) SetFlag(v bool) {
	m.fFlag = v
	m.__.Visit(_FIELD_Small_flag, _FIELD_TOTAL_Small)
}

func (m *SmallEX) GetStr() string {
	if !m.__.IsVisited(_FIELD_Small_str, _FIELD_TOTAL_Small) {
		field := m.__.RawField(_FIELD_Small_str)
		m.fStr = field.GetString()
		m.__.Visit(_FIELD_Small_str, _FIELD_TOTAL_Small)
	}
	return m.fStr
}

func (m *SmallEX) SetStr(v string) {
	m.fStr = v
	m.__.Visit(_FIELD_Small_str, _FIELD_TOTAL_Small)
}

type Vec2D_Vec1DEX []float32

func TO_Vec2D_Vec1DEX(data []byte) Vec2D_Vec1DEX {
	arr := protocache.AsFloat32Array(data)
	return append(Vec2D_Vec1DEX(nil), arr.Raw()...)
}

func (x Vec2D_Vec1DEX) Serialize() ([]byte, error) {
	words, err := serializeVec2D_Vec1DEX(x)
	if err != nil {
		return nil, err
	}
	return protocache.WordsToBytes(words), nil
}

func serializeVec2D_Vec1DEX(x Vec2D_Vec1DEX) ([]uint32, error) {
	if len(x) == 0 {
		return []uint32{1}, nil
	}
	return protocache.EncodeFloat32Array([]float32(x))
}

type Vec2DEX []Vec2D_Vec1DEX

func TO_Vec2DEX(data []byte) Vec2DEX {
	arr := protocache.AsArray(data)
	out := make(Vec2DEX, int(arr.Size()))
	for i := uint32(0); i < arr.Size(); i++ {
		elem := arr.Get(i)
		out[i] = TO_Vec2D_Vec1DEX(elem.GetObject())
	}
	return out
}

func (x Vec2DEX) Serialize() ([]byte, error) {
	words, err := serializeVec2DEX(x)
	if err != nil {
		return nil, err
	}
	return protocache.WordsToBytes(words), nil
}

func serializeVec2DEX(x Vec2DEX) ([]uint32, error) {
	if len(x) == 0 {
		return []uint32{1}, nil
	}
	return protocache.EncodeObjectArray(len(x), func(i int) ([]uint32, error) {
		return serializeVec2D_Vec1DEX(x[i])
	})
}

type ArrMap_ArrayEX []float32

func TO_ArrMap_ArrayEX(data []byte) ArrMap_ArrayEX {
	arr := protocache.AsFloat32Array(data)
	return append(ArrMap_ArrayEX(nil), arr.Raw()...)
}

func (x ArrMap_ArrayEX) Serialize() ([]byte, error) {
	words, err := serializeArrMap_ArrayEX(x)
	if err != nil {
		return nil, err
	}
	return protocache.WordsToBytes(words), nil
}

func serializeArrMap_ArrayEX(x ArrMap_ArrayEX) ([]uint32, error) {
	if len(x) == 0 {
		return []uint32{1}, nil
	}
	return protocache.EncodeFloat32Array([]float32(x))
}

type ArrMapEX map[string]ArrMap_ArrayEX

func TO_ArrMapEX(data []byte) ArrMapEX {
	pack := protocache.AsMap(data)
	out := make(ArrMapEX, int(pack.Size()))
	for i := uint32(0); i < pack.Size(); i++ {
		keyField := pack.Key(i)
		valField := pack.Value(i)
		out[keyField.GetString()] = TO_ArrMap_ArrayEX(valField.GetObject())
	}
	return out
}

func (x ArrMapEX) Serialize() ([]byte, error) {
	words, err := serializeArrMapEX(x)
	if err != nil {
		return nil, err
	}
	return protocache.WordsToBytes(words), nil
}

func serializeArrMapEX(x ArrMapEX) ([]uint32, error) {
	if len(x) == 0 {
		return []uint32{5 << 28}, nil
	}
	keys := make([][]uint32, 0, len(x))
	vals := make([][]uint32, 0, len(x))
	for k, v := range x {
		keyPart, err := protocache.EncodeString(k)
		if err != nil {
			return nil, err
		}
		valPart, err := serializeArrMap_ArrayEX(v)
		if err != nil {
			return nil, err
		}
		keys = append(keys, keyPart)
		vals = append(vals, valPart)
	}
	return protocache.EncodeMapParts(keys, vals, true)
}

type MainEX struct {
	__       protocache.MessageEX
	fI32     int32
	fU32     uint32
	fI64     int64
	fU64     uint64
	fFlag    bool
	fMode    Mode
	fStr     string
	fData    []byte
	fF32     float32
	fF64     float64
	fObject  *SmallEX
	fI32V    []int32
	fU64V    []uint64
	fStrv    []string
	fDatav   [][]byte
	fF32V    []float32
	fF64V    []float64
	fFlags   []bool
	fObjectv []*SmallEX
	fTU32    uint32
	fTI32    int32
	fTS32    int32
	fTU64    uint64
	fTI64    int64
	fTS64    int64
	fIndex   map[string]int32
	fObjects map[int32]*SmallEX
	fMatrix  Vec2DEX
	fVector  []ArrMapEX
	fArrays  ArrMapEX
	fModev   []Mode
}

func TO_MainEX(data []byte) *MainEX {
	out := &MainEX{}
	out.__.Init(data)
	return out
}

func (m *MainEX) HasBase() bool { return m.__.HasBase() }

func (m *MainEX) IsValid() bool { return m.HasBase() }

func (m *MainEX) Serialize() ([]byte, error) {
	words, err := m.serializeWords()
	if err != nil {
		return nil, err
	}
	return protocache.WordsToBytes(words), nil
}

func (m *MainEX) serializeWords() ([]uint32, error) {
	if m == nil {
		return []uint32{0}, nil
	}
	parts := make([][]uint32, 32)
	if m.__.IsVisited(_FIELD_Main_i32, _FIELD_TOTAL_Main) {
		parts[0] = protocache.EncodeInt32(m.fI32)
	} else if raw := m.__.RawPart(_FIELD_Main_i32); len(raw) != 0 {
		parts[0] = protocache.BytesToWords(raw)
	}
	if m.__.IsVisited(_FIELD_Main_u32, _FIELD_TOTAL_Main) {
		parts[1] = protocache.EncodeUint32(m.fU32)
	} else if raw := m.__.RawPart(_FIELD_Main_u32); len(raw) != 0 {
		parts[1] = protocache.BytesToWords(raw)
	}
	if m.__.IsVisited(_FIELD_Main_i64, _FIELD_TOTAL_Main) {
		parts[2] = protocache.EncodeInt64(m.fI64)
	} else if raw := m.__.RawPart(_FIELD_Main_i64); len(raw) != 0 {
		parts[2] = protocache.BytesToWords(raw)
	}
	if m.__.IsVisited(_FIELD_Main_u64, _FIELD_TOTAL_Main) {
		parts[3] = protocache.EncodeUint64(m.fU64)
	} else if raw := m.__.RawPart(_FIELD_Main_u64); len(raw) != 0 {
		parts[3] = protocache.BytesToWords(raw)
	}
	if m.__.IsVisited(_FIELD_Main_flag, _FIELD_TOTAL_Main) {
		parts[4] = protocache.EncodeBool(m.fFlag)
	} else if raw := m.__.RawPart(_FIELD_Main_flag); len(raw) != 0 {
		parts[4] = protocache.BytesToWords(raw)
	}
	if m.__.IsVisited(_FIELD_Main_mode, _FIELD_TOTAL_Main) {
		parts[5] = protocache.EncodeInt32(int32(m.fMode))
	} else if raw := m.__.RawPart(_FIELD_Main_mode); len(raw) != 0 {
		parts[5] = protocache.BytesToWords(raw)
	}
	if m.__.IsVisited(_FIELD_Main_str, _FIELD_TOTAL_Main) {
		part, err := protocache.EncodeString(m.fStr)
		if err != nil {
			return nil, err
		}
		parts[6] = part
	} else if raw := m.__.RawPart(_FIELD_Main_str); len(raw) != 0 {
		parts[6] = protocache.BytesToWords(raw)
	}
	if m.__.IsVisited(_FIELD_Main_data, _FIELD_TOTAL_Main) {
		part, err := protocache.EncodeBytes(m.fData)
		if err != nil {
			return nil, err
		}
		parts[7] = part
	} else if raw := m.__.RawPart(_FIELD_Main_data); len(raw) != 0 {
		parts[7] = protocache.BytesToWords(raw)
	}
	if m.__.IsVisited(_FIELD_Main_f32, _FIELD_TOTAL_Main) {
		parts[8] = protocache.EncodeFloat32(m.fF32)
	} else if raw := m.__.RawPart(_FIELD_Main_f32); len(raw) != 0 {
		parts[8] = protocache.BytesToWords(raw)
	}
	if m.__.IsVisited(_FIELD_Main_f64, _FIELD_TOTAL_Main) {
		parts[9] = protocache.EncodeFloat64(m.fF64)
	} else if raw := m.__.RawPart(_FIELD_Main_f64); len(raw) != 0 {
		parts[9] = protocache.BytesToWords(raw)
	}
	if m.__.IsVisited(_FIELD_Main_object, _FIELD_TOTAL_Main) {
		if m.fObject != nil {
			part, err := m.fObject.serializeWords()
			if err != nil {
				return nil, err
			}
			if len(part) > 1 {
				parts[10] = part
			}
		}
	} else if raw := m.__.RawPart(_FIELD_Main_object); len(raw) != 0 {
		parts[10] = protocache.BytesToWords(raw)
	}
	if m.__.IsVisited(_FIELD_Main_i32v, _FIELD_TOTAL_Main) {
		if len(m.fI32V) != 0 {
			part, err := protocache.EncodeInt32Array(m.fI32V)
			if err != nil {
				return nil, err
			}
			parts[11] = part
		}
	} else if raw := m.__.RawPart(_FIELD_Main_i32v); len(raw) != 0 {
		parts[11] = protocache.BytesToWords(raw)
	}
	if m.__.IsVisited(_FIELD_Main_u64v, _FIELD_TOTAL_Main) {
		if len(m.fU64V) != 0 {
			part, err := protocache.EncodeUint64Array(m.fU64V)
			if err != nil {
				return nil, err
			}
			parts[12] = part
		}
	} else if raw := m.__.RawPart(_FIELD_Main_u64v); len(raw) != 0 {
		parts[12] = protocache.BytesToWords(raw)
	}
	if m.__.IsVisited(_FIELD_Main_strv, _FIELD_TOTAL_Main) {
		if len(m.fStrv) != 0 {
			part, err := protocache.EncodeStringArray(m.fStrv)
			if err != nil {
				return nil, err
			}
			parts[13] = part
		}
	} else if raw := m.__.RawPart(_FIELD_Main_strv); len(raw) != 0 {
		parts[13] = protocache.BytesToWords(raw)
	}
	if m.__.IsVisited(_FIELD_Main_datav, _FIELD_TOTAL_Main) {
		if len(m.fDatav) != 0 {
			part, err := protocache.EncodeBytesArray(m.fDatav)
			if err != nil {
				return nil, err
			}
			parts[14] = part
		}
	} else if raw := m.__.RawPart(_FIELD_Main_datav); len(raw) != 0 {
		parts[14] = protocache.BytesToWords(raw)
	}
	if m.__.IsVisited(_FIELD_Main_f32v, _FIELD_TOTAL_Main) {
		if len(m.fF32V) != 0 {
			part, err := protocache.EncodeFloat32Array(m.fF32V)
			if err != nil {
				return nil, err
			}
			parts[15] = part
		}
	} else if raw := m.__.RawPart(_FIELD_Main_f32v); len(raw) != 0 {
		parts[15] = protocache.BytesToWords(raw)
	}
	if m.__.IsVisited(_FIELD_Main_f64v, _FIELD_TOTAL_Main) {
		if len(m.fF64V) != 0 {
			part, err := protocache.EncodeFloat64Array(m.fF64V)
			if err != nil {
				return nil, err
			}
			parts[16] = part
		}
	} else if raw := m.__.RawPart(_FIELD_Main_f64v); len(raw) != 0 {
		parts[16] = protocache.BytesToWords(raw)
	}
	if m.__.IsVisited(_FIELD_Main_flags, _FIELD_TOTAL_Main) {
		if len(m.fFlags) != 0 {
			part, err := protocache.EncodeBoolArray(m.fFlags)
			if err != nil {
				return nil, err
			}
			parts[17] = part
		}
	} else if raw := m.__.RawPart(_FIELD_Main_flags); len(raw) != 0 {
		parts[17] = protocache.BytesToWords(raw)
	}
	if m.__.IsVisited(_FIELD_Main_objectv, _FIELD_TOTAL_Main) {
		if len(m.fObjectv) != 0 {
			part, err := protocache.EncodeObjectArray(len(m.fObjectv), func(i int) ([]uint32, error) {
				if m.fObjectv[i] == nil {
					return []uint32{0}, nil
				}
				return m.fObjectv[i].serializeWords()
			})
			if err != nil {
				return nil, err
			}
			parts[18] = part
		}
	} else if raw := m.__.RawPart(_FIELD_Main_objectv); len(raw) != 0 {
		parts[18] = protocache.BytesToWords(raw)
	}
	if m.__.IsVisited(_FIELD_Main_t_u32, _FIELD_TOTAL_Main) {
		parts[19] = protocache.EncodeUint32(m.fTU32)
	} else if raw := m.__.RawPart(_FIELD_Main_t_u32); len(raw) != 0 {
		parts[19] = protocache.BytesToWords(raw)
	}
	if m.__.IsVisited(_FIELD_Main_t_i32, _FIELD_TOTAL_Main) {
		parts[20] = protocache.EncodeInt32(m.fTI32)
	} else if raw := m.__.RawPart(_FIELD_Main_t_i32); len(raw) != 0 {
		parts[20] = protocache.BytesToWords(raw)
	}
	if m.__.IsVisited(_FIELD_Main_t_s32, _FIELD_TOTAL_Main) {
		parts[21] = protocache.EncodeInt32(m.fTS32)
	} else if raw := m.__.RawPart(_FIELD_Main_t_s32); len(raw) != 0 {
		parts[21] = protocache.BytesToWords(raw)
	}
	if m.__.IsVisited(_FIELD_Main_t_u64, _FIELD_TOTAL_Main) {
		parts[22] = protocache.EncodeUint64(m.fTU64)
	} else if raw := m.__.RawPart(_FIELD_Main_t_u64); len(raw) != 0 {
		parts[22] = protocache.BytesToWords(raw)
	}
	if m.__.IsVisited(_FIELD_Main_t_i64, _FIELD_TOTAL_Main) {
		parts[23] = protocache.EncodeInt64(m.fTI64)
	} else if raw := m.__.RawPart(_FIELD_Main_t_i64); len(raw) != 0 {
		parts[23] = protocache.BytesToWords(raw)
	}
	if m.__.IsVisited(_FIELD_Main_t_s64, _FIELD_TOTAL_Main) {
		parts[24] = protocache.EncodeInt64(m.fTS64)
	} else if raw := m.__.RawPart(_FIELD_Main_t_s64); len(raw) != 0 {
		parts[24] = protocache.BytesToWords(raw)
	}
	if m.__.IsVisited(_FIELD_Main_index, _FIELD_TOTAL_Main) {
		if len(m.fIndex) != 0 {
			keys := make([][]uint32, 0, len(m.fIndex))
			vals := make([][]uint32, 0, len(m.fIndex))
			for k, v := range m.fIndex {
				keyPart, err := protocache.EncodeString(k)
				if err != nil {
					return nil, err
				}
				keys = append(keys, keyPart)
				vals = append(vals, protocache.EncodeInt32(v))
			}
			part, err := protocache.EncodeMapParts(keys, vals, true)
			if err != nil {
				return nil, err
			}
			parts[25] = part
		}
	} else if raw := m.__.RawPart(_FIELD_Main_index); len(raw) != 0 {
		parts[25] = protocache.BytesToWords(raw)
	}
	if m.__.IsVisited(_FIELD_Main_objects, _FIELD_TOTAL_Main) {
		if len(m.fObjects) != 0 {
			keys := make([][]uint32, 0, len(m.fObjects))
			vals := make([][]uint32, 0, len(m.fObjects))
			for k, v := range m.fObjects {
				keys = append(keys, protocache.EncodeInt32(k))
				if v == nil {
					vals = append(vals, nil)
					continue
				}
				valPart, err := v.serializeWords()
				if err != nil {
					return nil, err
				}
				if len(valPart) <= 1 {
					valPart = nil
				}
				vals = append(vals, valPart)
			}
			part, err := protocache.EncodeMapParts(keys, vals, false)
			if err != nil {
				return nil, err
			}
			parts[26] = part
		}
	} else if raw := m.__.RawPart(_FIELD_Main_objects); len(raw) != 0 {
		parts[26] = protocache.BytesToWords(raw)
	}
	if m.__.IsVisited(_FIELD_Main_matrix, _FIELD_TOTAL_Main) {
		part, err := serializeVec2DEX(m.fMatrix)
		if err != nil {
			return nil, err
		}
		if len(part) > 1 {
			parts[27] = part
		}
	} else if raw := m.__.RawPart(_FIELD_Main_matrix); len(raw) != 0 {
		parts[27] = protocache.BytesToWords(raw)
	}
	if m.__.IsVisited(_FIELD_Main_vector, _FIELD_TOTAL_Main) {
		if len(m.fVector) != 0 {
			part, err := protocache.EncodeObjectArray(len(m.fVector), func(i int) ([]uint32, error) {
				return serializeArrMapEX(m.fVector[i])
			})
			if err != nil {
				return nil, err
			}
			parts[28] = part
		}
	} else if raw := m.__.RawPart(_FIELD_Main_vector); len(raw) != 0 {
		parts[28] = protocache.BytesToWords(raw)
	}
	if m.__.IsVisited(_FIELD_Main_arrays, _FIELD_TOTAL_Main) {
		part, err := serializeArrMapEX(m.fArrays)
		if err != nil {
			return nil, err
		}
		if len(part) > 1 {
			parts[29] = part
		}
	} else if raw := m.__.RawPart(_FIELD_Main_arrays); len(raw) != 0 {
		parts[29] = protocache.BytesToWords(raw)
	}
	if m.__.IsVisited(_FIELD_Main_modev, _FIELD_TOTAL_Main) {
		if len(m.fModev) != 0 {
			part, err := protocache.EncodeEnumArray(m.fModev)
			if err != nil {
				return nil, err
			}
			parts[31] = part
		}
	} else if raw := m.__.RawPart(_FIELD_Main_modev); len(raw) != 0 {
		parts[31] = protocache.BytesToWords(raw)
	}
	for _, id := range []uint16{30} {
		if raw := m.__.RawPart(id); len(raw) != 0 {
			parts[id] = protocache.BytesToWords(raw)
		}
	}
	return protocache.EncodeMessageParts(parts)
}

func (m *MainEX) GetI32() int32 {
	if !m.__.IsVisited(_FIELD_Main_i32, _FIELD_TOTAL_Main) {
		field := m.__.RawField(_FIELD_Main_i32)
		m.fI32 = field.GetInt32()
		m.__.Visit(_FIELD_Main_i32, _FIELD_TOTAL_Main)
	}
	return m.fI32
}

func (m *MainEX) SetI32(v int32) {
	m.fI32 = v
	m.__.Visit(_FIELD_Main_i32, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetStr() string {
	if !m.__.IsVisited(_FIELD_Main_str, _FIELD_TOTAL_Main) {
		field := m.__.RawField(_FIELD_Main_str)
		m.fStr = field.GetString()
		m.__.Visit(_FIELD_Main_str, _FIELD_TOTAL_Main)
	}
	return m.fStr
}

func (m *MainEX) SetStr(v string) {
	m.fStr = v
	m.__.Visit(_FIELD_Main_str, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetData() []byte {
	if !m.__.IsVisited(_FIELD_Main_data, _FIELD_TOTAL_Main) {
		field := m.__.RawField(_FIELD_Main_data)
		if data := field.GetBytes(); data != nil {
			m.fData = append([]byte(nil), data...)
		}
		m.__.Visit(_FIELD_Main_data, _FIELD_TOTAL_Main)
	}
	return m.fData
}

func (m *MainEX) SetData(v []byte) {
	if v == nil {
		m.fData = nil
	} else {
		m.fData = append([]byte(nil), v...)
	}
	m.__.Visit(_FIELD_Main_data, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetObject() *SmallEX {
	if !m.__.IsVisited(_FIELD_Main_object, _FIELD_TOTAL_Main) {
		field := m.__.RawField(_FIELD_Main_object)
		m.fObject = TO_SmallEX(field.GetObject())
		m.__.Visit(_FIELD_Main_object, _FIELD_TOTAL_Main)
	}
	return m.fObject
}

func (m *MainEX) SetObject(v *SmallEX) {
	m.fObject = v
	m.__.Visit(_FIELD_Main_object, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetI32V() []int32 {
	if !m.__.IsVisited(_FIELD_Main_i32v, _FIELD_TOTAL_Main) {
		field := m.__.RawField(_FIELD_Main_i32v)
		m.fI32V = append([]int32(nil), field.GetInt32Array()...)
		m.__.Visit(_FIELD_Main_i32v, _FIELD_TOTAL_Main)
	}
	return m.fI32V
}

func (m *MainEX) SetI32V(v []int32) {
	if v == nil {
		m.fI32V = nil
	} else {
		m.fI32V = append(m.fI32V[:0], v...)
	}
	m.__.Visit(_FIELD_Main_i32v, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetStrv() []string {
	if !m.__.IsVisited(_FIELD_Main_strv, _FIELD_TOTAL_Main) {
		field := m.__.RawField(_FIELD_Main_strv)
		arr := protocache.AsStringArray(field.GetObject())
		m.fStrv = make([]string, int(arr.Size()))
		for i := uint32(0); i < arr.Size(); i++ {
			m.fStrv[i] = arr.Get(i)
		}
		m.__.Visit(_FIELD_Main_strv, _FIELD_TOTAL_Main)
	}
	return m.fStrv
}

func (m *MainEX) SetStrv(v []string) {
	if v == nil {
		m.fStrv = nil
	} else {
		m.fStrv = append(m.fStrv[:0], v...)
	}
	m.__.Visit(_FIELD_Main_strv, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetObjectv() []*SmallEX {
	if !m.__.IsVisited(_FIELD_Main_objectv, _FIELD_TOTAL_Main) {
		field := m.__.RawField(_FIELD_Main_objectv)
		arr := field.GetArray()
		m.fObjectv = make([]*SmallEX, int(arr.Size()))
		for i := uint32(0); i < arr.Size(); i++ {
			elem := arr.Get(i)
			m.fObjectv[i] = TO_SmallEX(elem.GetObject())
		}
		m.__.Visit(_FIELD_Main_objectv, _FIELD_TOTAL_Main)
	}
	return m.fObjectv
}

func (m *MainEX) SetObjectv(v []*SmallEX) {
	if v == nil {
		m.fObjectv = nil
	} else {
		m.fObjectv = append(m.fObjectv[:0], v...)
	}
	m.__.Visit(_FIELD_Main_objectv, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetMatrix() Vec2DEX {
	if !m.__.IsVisited(_FIELD_Main_matrix, _FIELD_TOTAL_Main) {
		field := m.__.RawField(_FIELD_Main_matrix)
		m.fMatrix = TO_Vec2DEX(field.GetObject())
		m.__.Visit(_FIELD_Main_matrix, _FIELD_TOTAL_Main)
	}
	return m.fMatrix
}

func (m *MainEX) SetMatrix(v Vec2DEX) {
	if v == nil {
		m.fMatrix = nil
	} else {
		m.fMatrix = append(m.fMatrix[:0], v...)
	}
	m.__.Visit(_FIELD_Main_matrix, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetVector() []ArrMapEX {
	if !m.__.IsVisited(_FIELD_Main_vector, _FIELD_TOTAL_Main) {
		field := m.__.RawField(_FIELD_Main_vector)
		arr := field.GetArray()
		m.fVector = make([]ArrMapEX, int(arr.Size()))
		for i := uint32(0); i < arr.Size(); i++ {
			elem := arr.Get(i)
			m.fVector[i] = TO_ArrMapEX(elem.GetObject())
		}
		m.__.Visit(_FIELD_Main_vector, _FIELD_TOTAL_Main)
	}
	return m.fVector
}

func (m *MainEX) SetVector(v []ArrMapEX) {
	if v == nil {
		m.fVector = nil
	} else {
		m.fVector = append(m.fVector[:0], v...)
	}
	m.__.Visit(_FIELD_Main_vector, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetArrays() ArrMapEX {
	if !m.__.IsVisited(_FIELD_Main_arrays, _FIELD_TOTAL_Main) {
		field := m.__.RawField(_FIELD_Main_arrays)
		m.fArrays = TO_ArrMapEX(field.GetObject())
		m.__.Visit(_FIELD_Main_arrays, _FIELD_TOTAL_Main)
	}
	return m.fArrays
}

func (m *MainEX) SetArrays(v ArrMapEX) {
	if v == nil {
		m.fArrays = nil
	} else {
		m.fArrays = make(ArrMapEX, len(v))
		for k, one := range v {
			m.fArrays[k] = one
		}
	}
	m.__.Visit(_FIELD_Main_arrays, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetModev() []Mode {
	if !m.__.IsVisited(_FIELD_Main_modev, _FIELD_TOTAL_Main) {
		field := m.__.RawField(_FIELD_Main_modev)
		m.fModev = append([]Mode(nil), protocache.CastEnumArray[Mode](field.GetEnumValueArray())...)
		m.__.Visit(_FIELD_Main_modev, _FIELD_TOTAL_Main)
	}
	return m.fModev
}

func (m *MainEX) SetModev(v []Mode) {
	if v == nil {
		m.fModev = nil
	} else {
		m.fModev = append(m.fModev[:0], v...)
	}
	m.__.Visit(_FIELD_Main_modev, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetIndex() map[string]int32 {
	if !m.__.IsVisited(_FIELD_Main_index, _FIELD_TOTAL_Main) {
		field := m.__.RawField(_FIELD_Main_index)
		pack := field.GetMap()
		m.fIndex = make(map[string]int32, int(pack.Size()))
		for i := uint32(0); i < pack.Size(); i++ {
			keyField := pack.Key(i)
			key := keyField.GetString()
			valField := pack.Value(i)
			m.fIndex[key] = valField.GetInt32()
		}
		m.__.Visit(_FIELD_Main_index, _FIELD_TOTAL_Main)
	}
	return m.fIndex
}

func (m *MainEX) SetIndex(v map[string]int32) {
	if v == nil {
		m.fIndex = nil
	} else {
		m.fIndex = make(map[string]int32, len(v))
		for k, one := range v {
			m.fIndex[k] = one
		}
	}
	m.__.Visit(_FIELD_Main_index, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetObjects() map[int32]*SmallEX {
	if !m.__.IsVisited(_FIELD_Main_objects, _FIELD_TOTAL_Main) {
		field := m.__.RawField(_FIELD_Main_objects)
		pack := field.GetMap()
		m.fObjects = make(map[int32]*SmallEX, int(pack.Size()))
		for i := uint32(0); i < pack.Size(); i++ {
			keyField := pack.Key(i)
			key := keyField.GetInt32()
			valField := pack.Value(i)
			m.fObjects[key] = TO_SmallEX(valField.GetObject())
		}
		m.__.Visit(_FIELD_Main_objects, _FIELD_TOTAL_Main)
	}
	return m.fObjects
}

func (m *MainEX) SetObjects(v map[int32]*SmallEX) {
	if v == nil {
		m.fObjects = nil
	} else {
		m.fObjects = make(map[int32]*SmallEX, len(v))
		for k, one := range v {
			m.fObjects[k] = one
		}
	}
	m.__.Visit(_FIELD_Main_objects, _FIELD_TOTAL_Main)
}
