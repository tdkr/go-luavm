package state

type luaStack struct {
	slots []luaValue
	top   int
}

func newLuaStack(size int) *luaStack {
	return &luaStack{
		slots: make([]luaValue, size),
		top:   0,
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

func (self *luaStack) pop() luaValue {
	if self.top == 0 {
		panic("lua stack underflow")
	}
	self.top--
	val := self.slots[self.top]
	return val
}

func (self *luaStack) absIndex(idx int) int {
	if idx >= 0 {
		return idx
	}
	return self.top + idx + 1
}

func (self *luaStack) isValid(idx int) bool {
	idx = self.absIndex(idx)
	return idx > 0 && idx <= self.top
}

func (self *luaStack) get(idx int) luaValue {
	idx = self.absIndex(idx)
	if self.isValid(idx) {
		return self.slots[idx-1]
	}
	return nil
}

func (self *luaStack) set(idx int, val luaValue) {
	idx = self.absIndex(idx)
	if self.isValid(idx) {
		self.slots[idx] = val
	} else {
		panic("invalid idx")
	}
}

func (self *luaStack) reverse(from, to int) {
	slots := self.slots
	for from < to {
		slots[from], slots[to] = slots[to], slots[from]
		from++
		to--
	}
}