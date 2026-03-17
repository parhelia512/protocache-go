package protocache

type MessageEX struct {
	base    Message
	visited []uint64
}

func (m *MessageEX) Init(data []byte) {
	m.base = AsMessage(data)
	m.visited = nil
}

func (m *MessageEX) HasBase() bool {
	return m.base.IsValid()
}

// IsValid is kept as a compatibility alias for HasBase.
func (m *MessageEX) IsValid() bool {
	return m.HasBase()
}

func (m *MessageEX) HasField(id uint16) bool {
	return m.base.HasField(id)
}

func (m *MessageEX) IsVisited(id uint16, fields uint16) bool {
	if id >= fields {
		return false
	}
	idx := int(id >> 6)
	if idx >= len(m.visited) {
		return false
	}
	return (m.visited[idx] & (uint64(1) << (id & 63))) != 0
}

func (m *MessageEX) Visit(id uint16, fields uint16) {
	if id >= fields {
		return
	}
	if len(m.visited) == 0 && fields != 0 {
		m.visited = make([]uint64, (uint32(fields)+63)/64)
	}
	idx := int(id >> 6)
	m.visited[idx] |= uint64(1) << (id & 63)
}

func (m *MessageEX) RawField(id uint16) Field {
	return m.base.GetField(id)
}

func (m *MessageEX) RawPart(id uint16) []byte {
	return m.base.RawPart(id)
}
