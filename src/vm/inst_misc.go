package vm

import . "github.com/tdkr/go-luavm/src/api"

/*
移动和跳转指令
*/

/*
MOVE指令（iABC模式）把源寄存器（索引由操作数B指定）里的值移动到目标寄存器（索引由操作数A指定）里。
我们先解码指令，得到目标寄存器和源寄存器索引，然后把它们转换成栈索引，最后调用Lua API提供的Copy()方法拷贝栈值。
*/
func move(i Instruction, vm LuaVM) {
	a, b, _ := i.ABC()
	a += 1
	b += 1
	vm.Copy(b, a)
}

/*
JMP指令（iAsBx模式）执行无条件跳转。
*/
func jmp(i Instruction, vm LuaVM) {
	a, sBx := i.AsBx()
	vm.AddPC(sBx)
	if a != 0 {
		panic("todo!")
	}
}
