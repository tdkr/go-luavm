package state

import (
	. "github.com/tdkr/go-luavm/src/api"
	"github.com/tdkr/go-luavm/src/binchunk"
)

type luaState struct {
	stack *luaStack
	proto *binchunk.Prototype
	pc    int
}

var _ LuaState = New(10, nil)

func New(stackSize int, proto *binchunk.Prototype) *luaState {
	return &luaState{
		stack: newLuaStack(stackSize),
		proto: proto,
		pc:    0,
	}
}
