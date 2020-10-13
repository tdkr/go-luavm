package vm

import . "github.com/tdkr/go-luavm/src/api"

/*
加载指令
*/

/*
LOADNIL指令（i ABC模式）用于给连续n个寄存器放置nil值。寄存器的起始索引由操作数A指定，寄存器数量则由操作数B指定，操作数C没有用。
LOADNIL指令可以用如下伪代码表示.
R(A), R(A+1), ..., R(A+B) := nil
Lua编译器在编译函数生成指令表时，会把指令执行阶段所需要的寄存器数量预先算好，保存在函数原型里。
这里假定虚拟机在执行第一条指令前，已经根据这一信息调用SetTop()方法保留了必要数量的栈空间。
有了这个假设，我们就可以先调用PushNil()方法往栈顶推入一个nil值，然后连续调用Copy()方法将nil值复制到指定寄存器中，最后调用Pop()方法把一开始推入栈顶的那个nil值弹出，让栈顶指针恢复原状。
*/
func loadNil(i Instruction, vm LuaVM) {
	a, b, _ := i.ABC()
	a += 1
	vm.PushNil()
	for i := a; i <= a+b; i++ {
		vm.Copy(-1, i)
	}
	vm.Pop(1)
}

/*
LOADBOOL指令（iABC模式）给单个寄存器设置布尔值。寄存器索引由操作数A指定，布尔值由寄存器B指定（0代表false，非0代表true），如果寄存器C非0则跳过下一条指令。LOADBOOL指令可以用如下伪代码表示。
R(A) := (bool)B; if (C) pc++
*/
func loadBool(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a += 1
	vm.PushBoolean(b != 0)
	vm.Replace(a)
	if c != 0 {
		vm.AddPC(1)
	}
}

/*
LOADK指令（iABx模式）将常量表里的某个常量加载到指定寄存器，寄存器索引由操作数A指定，常量表索引由操作数Bx指定。
如果用Kst(N)表示常量表中的第N个常量，那么LOADK指令可以用以下伪代码表示。R(A) := Kst(Bx)
*/
func loadK(i Instruction, vm LuaVM) {
	a, bx := i.ABx()
	a += 1
	vm.GetConst(bx)
	vm.Replace(a)
}

/*
LOADKX指令（也是iABx模式）需要和EXTRAARG指令（iAx模式）搭配使用，用后者的Ax操作数来指定常量索引。Ax操作数占26个比特，可以表达的最大无符号整数是67108864，可以满足大部分情况了。
*/
func loadKx(i Instruction, vm LuaVM) {
	a, _ := i.ABx()
	a += 1
	ax := Instruction(vm.Fetch()).Ax()
	vm.GetConst(ax)
	vm.Replace(a)
}
