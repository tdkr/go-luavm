package state

import "github.com/tdkr/go-luavm/src/api"

func (self *luaState) PushNil() {
	self.stack.push(nil)
}

func (self *luaState) PushBoolean(b bool) {
	self.stack.push(b)
}

func (self *luaState) PushInteger(n int64) {
	self.stack.push(n)
}

func (self *luaState) PushNumber(n float64) {
	self.stack.push(n)
}

func (self *luaState) PushString(s string) {
	self.stack.push(s)
}

func (self *luaState) PushGoFunction(f api.GoFunction) {
	self.stack.push(newGoClosure(f, 0))
}

func (self *luaState) PushGlobalTable() {
	g := self.registry.get(api.LUA_RIDX_GLOBALS)
	self.stack.push(g)
}

func (self *luaState) GetGlobal(name string) api.LuaType {
	g := self.registry.get(api.LUA_RIDX_GLOBALS)
	return self.getTable(g, name, false)
}

func (self *luaState) SetGlobal(name string) {
	g := self.registry.get(api.LUA_RIDX_GLOBALS)
	v := self.stack.pop()
	self.setTable(g, name, v, false)
}

func (self *luaState) Register(name string, f api.GoFunction) {
	self.PushGoFunction(f)
	self.SetGlobal(name)
}
