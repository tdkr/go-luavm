package vm

/*
Lua语言的for循环语句有两种形式：数值（Numerical）形式和通用（Generic）形式。数值for循环用于按一定步长遍历某个范围内的数值，通用for循环主要用于遍历表。
数值for循环需要借助两条指令来实现：FORPREP和FORLOOP
*/

import . "github.com/tdkr/go-luavm/src/api"

/*
FORPREP指令可以用以下伪代码表示:
R(A) -= R(A+2); pc += sBx
*/
func forPrep(i Instruction, vm LuaVM) {
	a, sBx := i.AsBx()
	a += 1
	//R(A) -= R(A+2)
	vm.PushValue(a)
	vm.PushValue(a + 2)
	vm.Arith(LUA_OPSUB)
	vm.Replace(a)
	// pc += sBx
	vm.AddPC(sBx)
}

/*
FORLOOP指令可以用以下伪代码表示:
R(A) += R(A+2);
if R(A) <? = R(A+1) then {
	pc += sBx; R(A+3) = R(A)
}
FORLOOP指令伪代码中的“<? =”符号。当步长是正数时，这个符号的含义是“<=”，也就是说继续循环的条件是数值不大于限制；当步长是负数时，这个符号的含义是“>=”，循环继续的条件就变成了数值不小于限制。
*/
func forLoop(i Instruction, vm LuaVM) {
	a, sBx := i.AsBx()
	a += 1

	// R(A) += R(A+2)
	vm.PushValue(a + 2)
	vm.PushValue(a)
	vm.Arith(LUA_OPADD)
	vm.Replace(a)

	// R(A) <? = R(A+1)
	isPositiveStep := vm.ToNumber(a+2) >= 0
	if isPositiveStep && vm.Compare(a, a+1, LUA_OPLE) ||
		!isPositiveStep && vm.Compare(a+1, a, LUA_OPLE) {
		vm.AddPC(sBx)   // pc += sBx
		vm.Copy(a, a+3) // R(A+3) = R(A)
	}
}

/*
R(A+3), ..., R(A+2+C) := R(A)(R(A+1), R(A+2))
*/
func tForCall(i Instruction, vm LuaVM) {
	a, _, c := i.ABC()
	a += 1

	_pushFuncAndArgs(a, 3, vm)
	vm.Call(2, c)
	_popResults(a+3, c+1, vm)
}

/*
if R(A+1) ~= nil then {
	R(A) = R(A+1); pc += sBx
}
*/
func tForLoop(i Instruction, vm LuaVM) {
	a, sBx := i.AsBx()
	a += 1

	if !vm.IsNil(a + 1) {
		vm.Copy(a+1, a)
		vm.AddPC(sBx)
	}
}
