package vm

import . "github.com/tdkr/go-luavm/src/api"

/*
CLOSURE指令（iBx模式）把当前Lua函数的子函数原型实例化为闭包，放入由操作数A指定的寄存器中。
子函数原型来自于当前函数原型的子函数原型表，索引由操作数Bx指定。
CLOSURE指令可以用如下伪代码表示:
R(A) := closure(KPROTO[Bx])
*/
func closure(i Instruction, vm LuaVM) {
	a, bx := i.ABx()
	a += 1

	vm.LoadProto(bx)
	vm.Replace(a)
}

/*
CALLCALL指令（iABC模式）调用Lua函数。其中被调函数位于寄存器中，索引由操作数A指定。
需要传递给被调函数的参数值也在寄存器中，紧挨着被调函数，数量由操作数B指定。
函数调用结束后，原先存放函数和参数值的寄存器会被返回值占据，具体有多少个返回值则由操作数C指定。
CALL指令可以用如下伪代码表示:
R(A), ..., R(A+C-2) := R(A)(R(A+1), ..., R(A+B-1))
*/
func call(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a += 1

	nArgs := _pushFuncAndArgs(a, b, vm)
	vm.Call(nArgs, c-1)
	_popResults(a, c, vm)
}

/*
如果操作数B大于0，需要传递的参数是B -1个，循环调用PushValue()方法把函数和参数值推入栈顶即可。
如果操作数B等于0，表示需要接受子函数的所有返回值作为参数
*/
func _pushFuncAndArgs(a, b int, vm LuaVM) (nArgs int) {
	if b >= 1 { // b-1 args
		vm.CheckStack(b)
		for i := a; i < a+b; i++ {
			vm.PushValue(i)
		}
		return b - 1
	} else {
		_fixStack(a, vm)
		return vm.GetTop() - vm.RegisterCount() - 1
	}
}

/*
如果操作数C大于1，则返回值数量是C-1，循环调用Replace()方法把栈顶返回值移动到相应寄存器即可；
如果操作数C等于1，则返回值数量是0，不需要任何处理；
如果C等于0，那么需要把被调函数的返回值全部返回。
对于最后这种情况，干脆就把这些返回值先留在栈顶，反正后面也是要把它们再推入栈顶的。我们往栈顶推入一个整数值，标记这些返回值原本是要移动到哪些寄存器中。
*/
func _popResults(a, c int, vm LuaVM) {
	if c == 1 { // no results

	} else if c > 1 { // c-1 results
		for i := a + c - 2; i >= a; i-- {
			vm.Replace(i)
		}
	} else {
		// leave results on stack
		vm.CheckStack(1)
		vm.PushInteger(int64(a))
	}
}

func _fixStack(a int, vm LuaVM) {
	x := int(vm.ToInteger(-1))
	vm.Pop(1)

	vm.CheckStack(x - a)
	for i := a; i < x; i++ {
		vm.PushValue(i)
	}
	vm.Rotate(vm.RegisterCount()+1, x-a)
}

/*
RETURN指令（iABC模式）把存放在连续多个寄存器里的值返回给主调函数。其中第一个寄存器的索引由操作数A指定，寄存器数量由操作数B指定，操作数C没用。
RETURN指令可以用如下伪代码表示:
return R(A, ..., R(A+B-2)
我们需要将返回值推入栈顶。如果操作数B等于1，则不需要返回任何值；
如果操作数B大于1，则需要返回B -1个值，这些值已经在寄存器里了，循环调用PushValue()方法复制到栈顶即可。
如果操作数B等于0，则一部分返回值已经在栈顶了，调用_fixStack()函数把另一部分也推入栈顶。
*/
func _return(i Instruction, vm LuaVM) {
	a, b, _ := i.ABC()
	a += 1

	if b == 1 { // no return values

	} else if b > 1 {
		vm.CheckStack(b - 1)
		for i := a; i <= a+b-2; i++ {
			vm.PushValue(i)
		}
	} else {
		_fixStack(a, vm)
	}
}

/*
VARARG指令（iABC模式）把传递给当前函数的变长参数加载到连续多个寄存器中。其中第一个寄存器的索引由操作数A指定，寄存器数量由操作数B指定，操作数C没有用。
VARARG指令可以用如下伪代码表示:
R(A), ..., R(A+B-2) = vararg
操作数B若大于1，表示把B-1个vararg参数复制到寄存器；
否则只能等于0，表示把全部vararg参数复制到寄存器。
对于这两种情况，我们统一调用LoadVararg()方法把vararg参数推入栈顶，剩下的工作交给_popResults()函数就可以了。
*/
func vararg(i Instruction, vm LuaVM) {
	a, b, _ := i.ABC()
	a += 1

	if b != 1 {
		vm.LoadVararg(b - 1)
		_popResults(a, b, vm)
	}
}

/*
TAILCALL指令（iABC模式）可以用如下伪代码表示:
return R(A)(R(A+1), ..., R(A+B-1))
todo:尾递归调用优化
*/
func tailCall(i Instruction, vm LuaVM) {
	a, b, _ := i.ABC()
	a += 1
	c := 0
	nArgs := _pushFuncAndArgs(a, b, vm)
	vm.Call(nArgs, c-1)
	_popResults(a, c, vm)
}

/*
SELF指令（iABC模式）把对象和方法拷贝到相邻的两个目标寄存器中。对象在寄存器中，索引由操作数B指定。方法名在常量表里，索引由操作数C指定。目标寄存器索引由操作数A指定。
SELF指令可以用如下伪代码表示:
R(A+1) := R(B); R(A) := R(B)[RK(C)]
*/
func self(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a += 1
	b += 1

	vm.Copy(b, a+1)
	vm.GetRK(c)
	vm.GetTable(b)
	vm.Replace(a)
}
