package state

import (
	"github.com/tdkr/go-luavm/src/api"
	"github.com/tdkr/go-luavm/src/binchunk"
	"github.com/tdkr/go-luavm/src/vm"
)

/*
Load()方法接收三个参数。
其中第一个参数是字节数组，给出要加载的chunk数据；
第二个参数是字符串，指定chunk的名字，供加载错误或调试时使用；
第三个参数也是字符串，指定加载模式，可选值是"b"、"t"或者"bt"。
如果加载模式是"b"，那么第一个参数必须是二进制chunk数据，否则加载失败，
如果加载模式是"t"，那么第一个参数必须是文本chunk数据，否则加载失败，
如果加载模式是"bt"，那么第一个参数可以是二进制或者文本chunk数据，Load()方法会根据实际的数据格式进行处理。
如果Load()方法无法成功加载chunk，需要在栈顶留下一条错误消息。Load()方法会返回一个状态码，0表示加载成功，非0表示加载失败。
*/
func (self *luaState) Load(chunk []byte, chunkName, mode string) int {
	proto := binchunk.Undump(chunk) // todo
	c := newLuaClosure(proto)
	self.stack.push(c)
	if len(proto.Upvalues) > 0 {
		env := self.registry.get(api.LUA_RIDX_GLOBALS)
		c.upvals[0] = &upvalue{&env}
	}
	return api.LUA_OK
}

/*
Call()方法对Lua函数进行调用。在执行Call()方法之前，必须先把被调函数推入栈顶，然后把参数值依次推入栈顶。
Call()方法结束之后，参数值和函数会被弹出栈顶，取而代之的是指定数量的返回值。
Call()方法接收两个参数：
第一个参数指定准备传递给被调函数的参数数量，同时也隐含给出了被调函数在栈里的位置；
第二个参数指定需要的返回值数量（多退少补），如果是-1，则被调函数的返回值会全部留在栈顶。
当我们试图调用一个非函数类型的值时，Lua会看这个值是否有__call元方法，如果有，Lua会以该值为第一个参数，后跟原方法调用的其他参数，来调用元方法，以元方法返回值为返回值。
*/
func (self *luaState) Call(nArgs, nResults int) {
	val := self.stack.get(-(nArgs + 1))

	c, ok := val.(*closure)
	if !ok {
		if mf := getMetaField(val, "__call", self); mf != nil {
			if c, ok = mf.(*closure); ok {
				self.stack.push(val)
				self.Insert(-(nArgs + 2))
				nArgs += 1
			}
		}
	}

	if ok {
		if c.proto != nil {
			self.callLuaClosure(nArgs, nResults, c)
		} else {
			self.callGoClosure(nArgs, nResults, c)
		}
	} else {
		panic("not a function")
	}
}

func (self *luaState) callLuaClosure(nArgs, nResults int, c *closure) {
	nRegs := int(c.proto.MaxStackSize)
	nParams := int(c.proto.NumParams)
	isVararg := c.proto.IsVararg == 1

	// create new lua stack
	newStack := newLuaStack(nRegs+api.LUA_MINSTACK, self)
	newStack.closure = c

	// pass args, pop func
	funcAndArgs := self.stack.popN(nArgs + 1)
	newStack.pushN(funcAndArgs[1:], nParams)
	newStack.top = nRegs
	if nArgs > nParams && isVararg {
		newStack.varargs = funcAndArgs[nParams+1:]
	}

	self.pushLuaStack(newStack)
	self.runLuaClosure()
	self.popLuaStack()

	if nResults != 0 { // 被调函数运行完毕之后，返回值会留在被调帧的栈顶（寄存器之上）。
		results := newStack.popN(newStack.top - nRegs)
		self.stack.check(len(results))
		self.stack.pushN(results, nResults)
	}
}

func (self *luaState) callGoClosure(nArgs, nResults int, c *closure) {
	newStack := newLuaStack(nArgs+api.LUA_MINSTACK, self)
	newStack.closure = c

	// pass args, pop func
	if nArgs > 0 {
		args := self.stack.popN(nArgs)
		newStack.pushN(args, nArgs)
	}
	self.stack.pop()

	// run closure
	self.pushLuaStack(newStack)
	r := c.goFunc(self)
	self.popLuaStack()

	// return results
	if nResults != 0 {
		results := newStack.popN(r)
		self.stack.check(len(results))
		self.stack.pushN(results, nResults)
	}
}

func (self *luaState) runLuaClosure() {
	for {
		inst := vm.Instruction(self.Fetch())
		inst.Execute(self)
		//fmt.Printf("run instruction name:%s, source:%s, line:%v",
		//	inst.OpName(), self.stack.closure.proto.Source, self.stack.closure.proto.LineDefined)
		//fmt.Print(", operands:")
		//util.PrintOperands(inst)
		//fmt.Print(", stack:")
		//util.PrintStack(self)
		//fmt.Println()
		if inst.Opcode() == vm.OP_RETURN {
			break
		}
	}
}

/*
PCall()会捕获函数调用过程中产生的错误，如果没有错误产生，那么PCall()的行为和Call()完全一致，最后返回LUA_OK。
如果有错误产生，那么PCall()会捕获错误，把错误对象留在栈顶，并且会返回相应的错误码。PCall()的第三个参数用于指定错误处理器，
*/
func (self *luaState) PCall(nArgs, nResult, ksgh int) (status int) {
	caller := self.stack
	status = api.LUA_ERRRUN

	defer func() {
		if r := recover(); r != nil {
			for self.stack != caller {
				self.popLuaStack()
			}
			self.stack.push(r)
		}
	}()

	self.Call(nArgs, nResult)
	status = api.LUA_OK
	return
}
