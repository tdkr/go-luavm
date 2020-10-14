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
除了可以进行无条件跳转之外，还兼顾着闭合处于开启状态的Upvalue的责任。
如果某个块内部定义的局部变量已经被嵌套函数捕获，那么当这些局部变量退出作用域（也就是块结束）时，编译器会生成一条JMP指令，指示虚拟机闭合相应的Upvalue。
伪代码: JMP A sBx   pc+=sBx; if (A) close all upvalues >= R(A - 1)
*/
func jmp(i Instruction, vm LuaVM) {
	a, sBx := i.AsBx()

	vm.AddPC(sBx)
	if a != 0 {
		vm.CloseUpvalues(a)
	}
}
