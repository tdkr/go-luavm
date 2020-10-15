package state

import . "github.com/tdkr/go-luavm/src/api"

type luaStack struct {
	state   *luaState
	slots   []luaValue
	top     int
	prev    *luaStack
	closure *closure
	varargs []luaValue
	pc      int
	openuvs map[int]*upvalue
}

func newLuaStack(size int, state *luaState) *luaStack {
	return &luaStack{
		slots: make([]luaValue, size),
		top:   0,
		state: state,
	}
}

func (self *luaStack) check(n int) {
	free := len(self.slots) - self.top
	for i := free; i < n; i++ {
		self.slots = append(self.slots, nil)
	}
}

func (self *luaStack) push(val luaValue) {
	if self.top == len(self.slots) {
		panic("lua stack overflow")
	}
	self.slots[self.top] = val
	self.top++
}

func (self *luaStack) pushN(vals []luaValue, n int) {
	nVals := len(vals)
	if n < 0 {
		n = nVals
	}
	for i := 0; i < n; i++ {
		if i < nVals {
			self.push(vals[i])
		} else {
			self.push(nil)
		}
	}
}

func (self *luaStack) pop() luaValue {
	if self.top < 1 {
		panic("lua stack underflow")
	}
	self.top--
	val := self.slots[self.top]
	return val
}

func (self *luaStack) popN(n int) []luaValue {
	arr := make([]luaValue, n)
	for i := n - 1; i >= 0; i-- {
		arr[i] = self.pop()
	}
	return arr
}

func (self *luaStack) absIndex(idx int) int {
	if idx >= 0 || idx <= LUA_REGISTRYINDEX {
		return idx
	}
	return self.top + idx + 1
}

func (self *luaStack) isValid(idx int) bool {
	if idx < LUA_REGISTRYINDEX { /* upvalues */
		uvIdx := LUA_REGISTRYINDEX - idx - 1
		c := self.closure
		return c != nil && uvIdx < len(c.upvals)
	}
	if idx == LUA_REGISTRYINDEX {
		return true
	}
	idx = self.absIndex(idx)
	return idx > 0 && idx <= self.top
}

func (self *luaStack) get(idx int) luaValue {
	if idx < LUA_REGISTRYINDEX { /* upvalues */
		uvIdx := LUA_REGISTRYINDEX - idx - 1
		c := self.closure
		if c == nil || uvIdx >= len(c.upvals) {
			return nil
		}
		return *(c.upvals[uvIdx].val)
	}
	if idx == LUA_REGISTRYINDEX {
		return self.state.registry
	}
	absIdx := self.absIndex(idx)
	if absIdx > 0 && absIdx <= self.top {
		return self.slots[absIdx-1]
	}
	return nil
}

func (self *luaStack) set(idx int, val luaValue) {
	if idx < LUA_REGISTRYINDEX { /* upvalues */
		uvIdx := LUA_REGISTRYINDEX - idx - 1
		c := self.closure
		if c != nil && uvIdx < len(c.upvals) {
			*(c.upvals[uvIdx].val) = val
		}
		return
	}
	if idx == LUA_REGISTRYINDEX {
		self.state.registry = val.(*luaTable)
		return
	}
	absIdx := self.absIndex(idx)
	if absIdx > 0 && absIdx <= self.top {
		self.slots[absIdx-1] = val
		return
	}
	panic("invalid index!")
}

func (self *luaStack) reverse(from, to int) {
	slots := self.slots
	for from < to {
		slots[from], slots[to] = slots[to], slots[from]
		from++
		to--
	}
}
