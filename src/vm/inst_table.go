package vm

import (
	. "github.com/tdkr/go-luavm/src/api"
)

const LFIELDS_PER_FLUSH = 50

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

/*
SETTABLE指令（iABC模式）根据键往表里赋值。其中表位于寄存器中，索引由操作数A指定；键和值可能位于寄存器中，也可能在常量表里，索引分别由操作数B和C指定。
SETTABLE指令可以用如下伪代码表示:
R(A)[RK(B)] := RK(C)
*/
func setTable(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a += 1
	vm.GetRK(b)
	vm.GetRK(c)
	vm.SetTable(a)
}

/*
SETTABLE是通用指令，每次只处理一个键值对，具体操作交给表去处理，并不关心实际写入的是表的哈希部分还是数组部分。
SETLIST指令（iABC模式）则是专门给数组准备的，用于按索引批量设置数组元素。
其中数组位于寄存器中，索引由操作数A指定；
需要写入数组的一系列值也在寄存器中，紧挨着数组，数量由操作数B指定；
数组起始索引则由操作数C指定。
当表构造器的最后一个元素是函数调用或者vararg表达式时，Lua会把它们产生的所有值都收集起来，以供SETLIST指令使用。
SETLIST指令稍微有一点复杂，可以用如下伪代码表示：
R(A)[(C-1)*FPF+i] := R(A+i), 1 <= i <= B
*/
func setList(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a += 1
	if c > 0 {
		c = c - 1
	} else {
		c = Instruction(vm.Fetch()).Ax()
	}

	bIsZero := b == 0
	if bIsZero {
		b = int(vm.ToInteger(-1)) - a - 1
		vm.Pop(1)
	}

	vm.CheckStack(1)
	idx := int64(c * LFIELDS_PER_FLUSH)
	for j := 1; j <= b; j++ {
		idx++
		vm.PushValue(a + j)
		vm.SetI(a, idx)
	}

	if bIsZero {
		for j := vm.RegisterCount() + 1; j <= vm.GetTop(); j++ {
			idx++
			vm.PushValue(j)
			vm.SetI(a, idx)
		}

		// clear stack
		vm.SetTop(vm.RegisterCount())
	}
}
