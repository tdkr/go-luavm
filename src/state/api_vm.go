package state

func (self *luaState) PC() int {
	return self.stack.pc
}

func (self *luaState) AddPC(n int) {
	self.stack.pc += n
}

func (self *luaState) Fetch() uint32 {
	i := self.stack.closure.proto.Code[self.stack.pc]
	self.stack.pc++
	return i
}

func (self *luaState) GetConst(idx int) {
	val := self.stack.closure.proto.Constants[idx]
	self.stack.push(val)
}

/*
传递给GetRK()方法的参数实际上是iABC模式指令里的OpArgK类型参数。
这种类型的参数一共占9个比特。如果最高位是1，那么参数里存放的是常量表索引，把最高位去掉就可以得到索引值；
否则最高位是0，参数里存放的就是寄存器索引值。但是Lua虚拟机指令操作数里携带的寄存器索引是从0开始的，而Lua API里的栈索引是从1开始的，所以当需要把寄存器索引当成栈索引使用时，要对寄存器索引加1。
*/
func (self *luaState) GetRK(rk int) {
	if rk > 0xff { // constant
		self.GetConst(rk & 0xff)
	} else { // register
		self.PushValue(rk + 1)
	}
}

func (self *luaState) RegisterCount() int {
	return int(self.stack.closure.proto.MaxStackSize)
}

func (self *luaState) LoadVararg(n int) {
	if n < 0 {
		n = len(self.stack.varargs)
	}
	self.stack.check(n)
	self.stack.pushN(self.stack.varargs, n)
}

func (self *luaState) LoadProto(idx int) {
	stack := self.stack
	subProto := stack.closure.proto.Protos[idx]
	closure := newLuaClosure(subProto)
	stack.push(closure)

	/*
		我们需要根据函数原型里的Upvalue表来初始化闭包的Upvalue值。对于每个Upvalue，又有两种情况需要考虑：
		如果某一个Upvalue捕获的是当前函数的局部变量（Instack == 1），那么我们只要访问当前函数的局部变量即可；
		如果某一个Upvalue捕获的是更外围的函数中的局部变量（Instack ==0），该Upvalue已经被当前函数捕获，我们只要把该Upvalue传递给闭包即可。
	*/
	for i, uvInfo := range subProto.Upvalues {
		uvIdx := int(uvInfo.Idx)
		if uvInfo.Instack == 1 {
			if stack.openuvs == nil {
				stack.openuvs = map[int]*upvalue{}
			}
			if openuv, ok := stack.openuvs[uvIdx]; ok {
				closure.upvals[i] = openuv
			} else {
				closure.upvals[i] = &upvalue{&stack.slots[uvIdx]}
				stack.openuvs[uvIdx] = closure.upvals[i]
			}
		} else {
			closure.upvals[i] = stack.closure.upvals[uvIdx]
		}
	}
}

/*
todo:???
处于开启状态的Upvalue引用了还在寄存器里的Lua值，我们把这些Lua值从寄存器里复制出来，然后更新Upvalue，这样就将其改为了闭合状态
*/
func (self *luaState) CloseUpvalues(a int) {
	for i, openuv := range self.stack.openuvs {
		if i >= a-1 {
			val := *openuv.val
			openuv.val = &val
			delete(self.stack.openuvs, i)
		}
	}
}
