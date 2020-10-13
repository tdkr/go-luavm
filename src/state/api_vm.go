package state

func (self *luaState) PC() int {
	return self.pc
}

func (self *luaState) AddPC(n int) {
	self.pc += n
}

func (self *luaState) Fetch() uint32 {
	i := self.proto.Code[self.pc]
	self.pc++
	return i
}

func (self *luaState) GetConst(idx int) {
	val := self.proto.Constants[idx]
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
