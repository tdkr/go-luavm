package binchunk

/*
Lua的一个执行单元被称作Chunk，一个Chunk就是一个语句（赋值、控制结构、函数调用、变量声明）序列，它们会按次序执行。每个语句可以以一个分号结束：
当一个 chunk 被执行，首先它会被预编译成虚拟机中的指令序列， 然后被虚拟机解释运行这些指令。
二进制chunk内部使用的数据类型大致可以分为数字、字符串和列表三种。
数字类型主要包括字节、C语言整型（后文简称cint）、C语言size_t类型（简称size_t）、Lua整数、Lua浮点数五种。
*/

const (
	LUA_SIGNATURE    = "\x1bLua" //a fixed magic number(0x1B4C7561), lua signature, ascii code to "EscLua"
	LUAC_VERSION     = 0x53
	LUAC_FORMAT      = 0
	LUAC_DATA        = "\x19\x93\r\n\x1a\n"
	CINT_SIZE        = 4
	CSIZET_SIZE      = 8
	INSTRUCTION_SIZE = 4
	LUA_INTEGER_SIZE = 8
	LUA_NUMBER_SIZE  = 8
	LUAC_INT         = 0x5678
	LUAC_NUM         = 370.5
)

const (
	TAG_NIL       = 0x00
	TAG_BOOLEAN   = 0x01
	TAG_NUMBER    = 0x03
	TAG_INTEGER   = 0x13
	TAG_SHORT_STR = 0x04
	TAG_LONG_STR  = 0x14
)

type header struct {
	signature       [4]byte
	version         byte
	format          byte
	luacData        [6]byte
	cintSize        byte
	sizetSize       byte
	instructionSize byte
	luaIntegerSize  byte
	luaNumberSize   byte
	luacInt         int64
	luacNum         int64
}

type Upvalue struct {
	Instack byte
	Idx     byte
}

type LocVar struct {
	VarName string
	StartPC uint32
	EndPC   uint32
}

type Prototype struct {
	Source          string
	LineDefined     uint32
	LastLineDefined uint32
	NumParams       byte
	IsVararg        byte
	MaxStackSize    byte
	Code            []uint32
	Constants       []interface{}
	Upvalues        []Upvalue
	Protos          []*Prototype
	LineInfo        []uint32
	LocVars         []LocVar
	UpvalueNames    []string
}

type binaryChunk struct {
	header                  // 头部
	sizeUpvalues byte       // 头函数 upvalue 数量
	mainFunc     *Prototype //主函数原型
}

func Undump(data []byte) *Prototype {
	reader := &reader{data}
	reader.checkHeader()
	reader.readByte()
	return reader.readProto("")
}
