package binchunk

import (
	"encoding/binary"
	"math"
)

type reader struct {
	data []byte
}

func (self *reader) readByte() byte {
	v := self.data[0]
	self.data = self.data[1:]
	return v
}
func (self *reader) readBytes(length uint) []byte {
	bytes := self.data[:length]
	self.data = self.data[length:]
	return bytes
}

func (self *reader) readUint32() uint32 {
	v := binary.LittleEndian.Uint32(self.data)
	self.data = self.data[4:]
	return v
}

func (self *reader) readUint64() uint64 {
	v := binary.LittleEndian.Uint64(self.data)
	self.data = self.data[8:]
	return v
}

func (self *reader) readLuaInteger() int64 {
	return int64(self.readUint64())
}

func (self *reader) readLuaNumber() float64 {
	return math.Float64frombits(self.readUint64())
}

func (self *reader) readString() string {
	size := uint(self.readByte())
	if size == 0 {
		return ""
	}
	if size == 0xff { //长字符串
		size = uint(self.readUint64())
	}
	bytes := self.readBytes(size - 1)
	return string(bytes)
}

func (self *reader) checkHeader() {
	if string(self.readBytes(4)) != LUA_SIGNATURE {
		panic("lua signature not match")
	} else if self.readByte() != LUAC_VERSION {
		panic("luac version not match")
	} else if self.readByte() != LUAC_FORMAT {
		panic("luac format not match")
	} else if string(self.readBytes(6)) != LUAC_DATA {
		panic("luac data not match")
	} else if self.readByte() != CINT_SIZE {
		panic("cint size not match")
	} else if self.readByte() != CSIZET_SIZE {
		panic("csizet size not match")
	} else if self.readByte() != INSTRUCTION_SIZE {
		panic("instruction size not match")
	} else if self.readByte() != LUA_INTEGER_SIZE {
		panic("lua integer size not match")
	} else if self.readByte() != LUA_NUMBER_SIZE {
		panic("lua number size not match")
	} else if self.readLuaInteger() != LUAC_INT {
		panic("luac int not match")
	} else if self.readLuaNumber() != LUAC_NUM {
		panic("luac num not match")
	}
}

func (self *reader) readProto(parentSource string) *Prototype {
	src := self.readString()
	if src == "" {
		src = parentSource
	}
	return &Prototype{
		Source:          src,
		LineDefined:     self.readUint32(),
		LastLineDefined: self.readUint32(),
		NumParams:       self.readByte(),
		IsVararg:        self.readByte(),
		MaxStackSize:    self.readByte(),
		Code:            self.readCode(),
		Constants:       self.readConstants(),
		Upvalues:        self.readUpvalues(),
		Protos:          self.readProtos(src),
		LineInfo:        self.readLineInfo(),
		LocVars:         self.readLocVars(),
		UpvalueNames:    self.readUpvalueNames(),
	}
}

// 读取指令表
func (self *reader) readCode() []uint32 {
	code := make([]uint32, self.readUint32())
	for i := range code {
		code[i] = self.readUint32()
	}
	return code
}

func (self *reader) readConstants() []interface{} {
	constants := make([]interface{}, self.readUint32())
	for i := range constants {
		constants[i] = self.readConstant()
	}
	return constants
}

func (self *reader) readConstant() interface{} {
	switch self.readByte() { // tag
	case TAG_NIL:
		return nil
	case TAG_BOOLEAN:
		return self.readByte() != 0
	case TAG_INTEGER:
		return self.readLuaInteger()
	case TAG_NUMBER:
		return self.readLuaNumber()
	case TAG_SHORT_STR, TAG_LONG_STR:
		return self.readString()
	default:
		panic("invalid tag")
	}
}

func (self *reader) readUpvalues() []Upvalue {
	upvalues := make([]Upvalue, self.readUint32())
	for i := range upvalues {
		upvalues[i] = Upvalue{
			Instack: self.readByte(),
			Idx:     self.readByte(),
		}
	}
	return upvalues
}

func (self *reader) readProtos(parentSource string) []*Prototype {
	protos := make([]*Prototype, self.readUint32())
	for i := range protos {
		protos[i] = self.readProto(parentSource)
	}
	return protos
}

func (self *reader) readLineInfo() []uint32 {
	lineInfo := make([]uint32, self.readUint32())
	for i := range lineInfo {
		lineInfo[i] = self.readUint32()
	}
	return lineInfo
}

func (self *reader) readLocVars() []LocVar {
	locVars := make([]LocVar, self.readUint32())
	for i := range locVars {
		locVars[i] = LocVar{
			VarName: self.readString(),
			StartPC: self.readUint32(),
			EndPC:   self.readUint32(),
		}
	}
	return locVars
}

func (self *reader) readUpvalueNames() []string {
	names := make([]string, self.readUint32())
	for i := range names {
		names[i] = self.readString()
	}
	return names
}
