package vm

import (
	. "github.com/tdkr/go-luavm/src/api"
)

/*
NEWTABLE指令（iABC模式）创建空表，并将其放入指定寄存器。寄存器索引由操作数A指定，表的初始数组容量和哈希表容量分别由操作数B和C指定。
NEWTABLE指令可以用如下伪代码表示:
R(A) := {} (size = B, C)
*/
func newTable(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a += 1
	vm.CreateTable(Fb2int(b), Fb2int(c))
	vm.Replace(a)
}

/*
GETTABLE指令（iABC模式）根据键从表里取值，并放入目标寄存器中。其中表位于寄存器中，索引由操作数B指定；键可能位于寄存器中，也可能在常量表里，索引由操作数C指定；目标寄存器索引则由操作数A指定。
GETTABLE指令可以用如下伪代码表示:
R(A) := R(B)[RK(C)]
*/
func getTable(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a += 1
	b += 1
	vm.GetRK(c)
	vm.GetTable(b)
	vm.Replace(a)
}
