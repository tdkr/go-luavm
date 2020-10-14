package vm

import . "github.com/tdkr/go-luavm/src/api"

/*
如果当前闭包的某个Upvalue是表，则GETTABUP指令（iABC模式）可以根据键从该表里取值，然后把值放入目标寄存器中。
其中目标寄存器索引由操作数A指定，Upvalue索引由操作数B指定，键（可能在寄存器中也可能在常量表中）索引由操作数C指定。
GETTABUP指令相当于GETUPVAL和GETTABLE这两条指令的组合，不过前者的效率明显要高一些。
GETTABUP指令可以用如下伪代码表示:
R(A) := UpValue[B](RK[C])
*/
func getTabUp(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a += 1
	b += 1

	vm.GetRK(c)
	vm.GetTable(LuaUpvalueIndex(b))
	vm.Replace(a)
}

/*
如果当前闭包的某个Upvalue是表，则SETTABUP指令（iABC模式）可以根据键往该表里写入值。
其中Upvalue索引由操作数A指定，键和值可能在寄存器中也可能在常量表中，索引分别由操作数B和C指定。
和GETTABUP指令类似，SETTABUP指令相当于GETUPVAL和SETTABLE这两条指令的组合，不过一条指令的效率要高一些。
SETTABUP指令可以用如下伪代码表示:
Upvalue[A][RK(B)] := RK(C)
*/
func setTabUp(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a += 1

	vm.GetRK(b)
	vm.GetRK(c)
	vm.SetTable(LuaUpvalueIndex(a))
}

/*
GETUPVAL指令（iABC模式），把当前闭包的某个Upvalue值拷贝到目标寄存器中。
其中目标寄存器的索引由操作数A指定，Upvalue索引由操作数B指定，操作数C没用。
GETUPVAL指令可以用如下伪代码表示:
R(A) := UpValue[B]
*/
func getUpval(i Instruction, vm LuaVM) {
	a, b, _ := i.ABC()
	a += 1
	b += 1
	vm.Copy(LuaUpvalueIndex(b), a)
}

/*
SETUPVAL指令（iABC模式），使用寄存器中的值给当前闭包的Upvalue赋值。
其中寄存器索引由操作数A指定，Upvalue索引由操作数B指定，操作数C没用。
SETUPVAL指令可以用如下伪代码表示:
Upvalue[B] := R(A)
*/
func setUpval(i Instruction, vm LuaVM) {
	a, b, _ := i.ABC()
	a += 1
	b += 1
	vm.Copy(a, LuaUpvalueIndex(b))
}
