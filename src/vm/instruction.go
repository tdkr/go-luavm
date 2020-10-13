package vm

import (
	"github.com/tdkr/go-luavm/src/api"
)

const MAXARG_Bx = 1<<18 - 1       // 2^18-1 = 262143
const MAXARG_sBx = MAXARG_Bx >> 1 // 262143/2 = 131071

/*
，每条Lua虚拟机指令占用4个字节，共32个比特（可以用Go语言uint32类型表示），其中低6个比特用于操作码，高26个比特用于操作数。
*/
type Instruction uint32

func (self Instruction) Opcode() int {
	return int(self & 0x3f)
}

func (self Instruction) ABC() (a, b, c int) {
	a = int(self >> 6 & 0xff)
	c = int(self >> 14 & 0x1ff)
	b = int(self >> 23 & 0x1ff)
	return
}

func (self Instruction) ABx() (a, bx int) {
	a = int(self >> 6 & 0xff)
	bx = int(self >> 14)
	return
}

func (self Instruction) AsBx() (a, bx int) {
	a, bx = self.ABx()
	bx = bx - MAXARG_sBx
	return
}

func (self Instruction) Ax() int {
	return int(self >> 6)
}

func (self Instruction) OpName() string {
	return opcodes[self.Opcode()].name
}

func (self Instruction) OpMode() byte {
	return opcodes[self.Opcode()].opMode
}

func (self Instruction) BMode() byte {
	return opcodes[self.Opcode()].argBMode
}

func (self Instruction) CMode() byte {
	return opcodes[self.Opcode()].argCMode
}

func (self Instruction) Execute(vm api.LuaVM) {
	action := opcodes[self.Opcode()].action
	if action != nil {
		action(self, vm)
	} else {
		panic(self.OpName())
	}
}
