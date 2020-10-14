package state

import (
	. "github.com/tdkr/go-luavm/src/api"
)

type luaState struct {
	stack    *luaStack
	registry *luaTable
}

var _ LuaState = New()

func New() *luaState {
	registry := newLuaTable(0, 0)
	registry.put(LUA_RIDX_GLOBALS, newLuaTable(0, 0)) //全局huanjing
	ls := &luaState{registry: registry}
	ls.pushLuaStack(newLuaStack(LUA_MINSTACK, ls))
	return ls
}

func (self *luaState) pushLuaStack(stack *luaStack) {
	stack.prev = self.stack
	self.stack = stack
}

func (self *luaState) popLuaStack() {
	stack := self.stack
	self.stack = stack.prev
	stack.prev = nil
}
