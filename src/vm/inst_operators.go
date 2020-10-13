package vm

import (
	. "github.com/tdkr/go-luavm/src/api"
)

/*
二元算术运算指令（iABC模式），对两个寄存器或常量值（索引由操作数B和C指定）进行运算，将结果放入另一个寄存器（索引由操作数A指定）。
如果用RK(N)表示寄存器或者常量值，那么二元算术运算指令可以用如下伪代码表示:
R(A) := RK(B) op RK(C)
*/
func _binaryArith(i Instruction, vm LuaVM, op ArithOp) {
	a, b, c := i.ABC()
	a += 1
	vm.GetRK(b)
	vm.GetRK(c)
	vm.Arith(op)
	vm.Replace(a)
}

/*
一元算术运算指令（iABC模式），对操作数B所指定的寄存器里的值进行运算，然后把结果放入操作数A所指定的寄存器中，操作数C没用。
一元算术指令可以用如下伪代码表示:
R(A) = op R(B)
*/
func _unaryArith(i Instruction, vm LuaVM, op ArithOp) {
	a, b, _ := i.ABC()
	a += 1
	b += 1
	vm.PushValue(b)
	vm.Arith(op)
	vm.Replace(a)
}

func add(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPADD) }  // +
func sub(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPSUB) }  // -
func mul(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPMUL) }  // *
func mod(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPMOD) }  // %
func pow(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPPOW) }  // ^
func div(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPDIV) }  // /
func idiv(i Instruction, vm LuaVM) { _binaryArith(i, vm, LUA_OPIDIV) } // //
func band(i Instruction, vm LuaVM) { _binaryArith(i, vm, LUA_OPBAND) } // &
func bor(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPBOR) }  // |
func bxor(i Instruction, vm LuaVM) { _binaryArith(i, vm, LUA_OPBXOR) } // ~
func shl(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPSHL) }  // <<
func shr(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPSHR) }  // >>
func unm(i Instruction, vm LuaVM)  { _unaryArith(i, vm, LUA_OPUNM) }   // -
func bnot(i Instruction, vm LuaVM) { _unaryArith(i, vm, LUA_OPBNOT) }  // ~

/*
LEN指令（iABC模式）进行的操作和一元算术运算指令类似，可以用伪代码表示为:
R(A) := length of R(B)
*/
func _len(i Instruction, vm LuaVM) {
	a, b, _ := i.ABC()
	a += 1
	b += 1
	vm.Len(b)
	vm.Replace(a)
}

func length(i Instruction, vm LuaVM) {
	_len(i, vm)
}

/*
CONCAT指令（iABC模式），将连续n个寄存器（起止索引分别由操作数B和C指定）里的值拼接，将结果放入另一个寄存器（索引由操作数A指定）。
CONCAT指令可以用如下伪代码表示:
R(A) := R(B).. ... ..R(C)
*/
func concat(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a += 1
	b += 1
	c += 1
	n := c - b + 1
	vm.CheckStack(n)
	for i := b; i <= c; i++ {
		vm.PushValue(i)
	}
	vm.Concat(n)
	vm.Replace(a)
}

/*
比较指令（iABC模式），比较寄存器或常量表里的两个值（索引分别由操作数B和C指定），如果比较结果和操作数A（转换为布尔值）匹配，则跳过下一条指令。
比较指令不改变寄存器状态，可以用如下伪代码表示:
if (RK(B) op RK(C)) != A then pc++
*/
func _compare(i Instruction, vm LuaVM, op CompareOp) {
	a, b, c := i.ABC()
	vm.GetRK(b)
	vm.GetRK(c)
	if vm.Compare(-2, -1, op) != (a != 0) {
		vm.AddPC(1)
	}
	vm.Pop(2)
}

func eq(i Instruction, vm LuaVM) { _compare(i, vm, LUA_OPEQ) } // ==
func lt(i Instruction, vm LuaVM) { _compare(i, vm, LUA_OPLT) } // <
func le(i Instruction, vm LuaVM) { _compare(i, vm, LUA_OPLE) } // <=

/*
NOT指令（iABC模式）进行的操作和一元算术运算指令类似
可以用如下伪代码表示:
R(A) := not R(B)
*/
func not(i Instruction, vm LuaVM) {
	a, b, _ := i.ABC()
	a += 1
	b += 1
	vm.PushBoolean(!vm.ToBoolean(b))
	vm.Replace(b)
}

/*
TESTSET指令（iABC模式），判断寄存器B（索引由操作数B指定）中的值转换为布尔值之后是否和操作数C表示的布尔值一致，如果一致则将寄存器B中的值复制到寄存器A（索引由操作数A指定）中，否则跳过下一条指令。
TESTSET指令可以用如下伪代码表示（<=>表示按布尔值比较）:
if (R(B) <=> R(C)) then R(A) := R(B) else pc++
*/
func testSet(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a += 1
	b += 1
	if vm.ToBoolean(b) == (c != 0) {
		vm.Copy(b, a)
	} else {
		vm.AddPC(1)
	}
}

/*
TEST指令（iABC模式），判断寄存器A（索引由操作数A指定）中的值转换为布尔值之后是否和操作数C表示的布尔值一致，如果一致，则跳过下一条指令。
TEST指令不使用操作数B，也不改变寄存器状态，可以用以下伪代码表示:
if not (R(A) <=> C) then pc++
*/
func test(i Instruction, vm LuaVM) {
	a, _, c := i.ABC()
	a += 1
	if vm.ToBoolean(a) == (c != 0) {
		vm.AddPC(1)
	}
}
