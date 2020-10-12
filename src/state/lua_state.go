package state

import (
	"fmt"
	. "github.com/tdkr/go-luavm/src/api"
)

type luaState struct {
	stack *luaStack
}

func (l luaState) GetTop() int {
	return l.stack.top
}

func (l luaState) AbsIndex(idx int) int {
	return l.stack.absIndex(idx)
}

func (l luaState) CheckStack(n int) bool {
	l.stack.check(n)
	return true
}

func (l luaState) Pop(n int) {
	for i := 0; i < n; i++ {
		l.stack.pop()
	}
}

func (l luaState) Copy(fromIdx, toIdx int) {
	val := l.stack.get(fromIdx)
	l.stack.set(toIdx, val)
}

/* PushValue()方法把指定索引处的值推入栈顶。*/
func (l luaState) PushValue(idx int) {
	val := l.stack.get(idx)
	l.stack.push(val)
}

/* Replace()是PushValue()的反操作：将栈顶值弹出，然后写入指定位置。 */
func (l luaState) Replace(idx int) {
	val := l.stack.pop()
	l.stack.set(idx, val)
}

/* Insert()方法将栈顶值弹出，然后插入指定位置。 */
func (l luaState) Insert(idx int) {
	l.Rotate(idx, 1)
}

func (l luaState) Remove(idx int) {
	l.Rotate(idx, -1)
	l.Pop(1)
}

/* Rotate()方法将[idx, top]索引区间内的值朝栈顶方向循环移动n个位置。如果n是负数，那么实际效果就是朝栈底方向移动。 */
func (l luaState) Rotate(idx, n int) {
	// 三次反转法
	t := l.stack.top - 1
	p := l.stack.absIndex(idx) - 1
	var m int
	if n >= 0 {
		m = t - n
	} else {
		m = p - n - 1
	}
	l.stack.reverse(p, m)
	l.stack.reverse(m+1, t)
	l.stack.reverse(p, t)
}

/*
SetTop()方法将栈顶索引设置为指定值。如果指定值小于当前栈顶索引，效果则相当于弹出操作（指定值为0相当于清空栈）
如果指定值大于当前栈顶索引，则效果相当于推入多个nil值。
*/
func (l luaState) SetTop(idx int) {
	top := l.AbsIndex(idx)
	if top < 0 {
		panic("stack underflow")
	}
	if n := l.stack.top - top; n > 0 {
		for i := 0; i < n; i++ {
			l.stack.pop()
		}
	} else {
		for i := n; i < 0; i++ {
			l.stack.push(nil)
		}
	}
}

func (l luaState) TypeName(tp LuaType) string {
	switch tp {
	case LUA_TNONE:
		return "no value"
	case LUA_TNIL:
		return "nil"
	case LUA_TBOOLEAN:
		return "boolean"
	case LUA_TNUMBER:
		return "number"
	case LUA_TSTRING:
		return "string"
	case LUA_TTABLE:
		return "table"
	case LUA_TFUNCTION:
		return "function"
	case LUA_TTHREAD:
		return "thread"
	default:
		return "userdata"
	}
}

func (l luaState) Type(idx int) LuaType {
	if l.stack.isValid(idx) {
		val := l.stack.get(idx)
		return typeOf(val)
	}
	return LUA_TNONE
}

func (l luaState) IsNone(idx int) bool {
	return l.Type(idx) == LUA_TNONE
}

func (l luaState) IsNil(idx int) bool {
	return l.Type(idx) == LUA_TNIL
}

func (l luaState) IsNoneOrNil(idx int) bool {
	return l.Type(idx) <= LUA_TNIL
}

func (l luaState) IsBoolean(idx int) bool {
	return l.Type(idx) == LUA_TBOOLEAN
}

func (l luaState) IsInteger(idx int) bool {
	val := l.stack.get(idx)
	_, ok := val.(int64)
	return ok
}

func (l luaState) IsNumber(idx int) bool {
	return l.Type(idx) == LUA_TNUMBER
}

func (l luaState) IsString(idx int) bool {
	return l.Type(idx) == LUA_TNUMBER || l.Type(idx) == LUA_TSTRING
}

func (l luaState) ToBoolean(idx int) bool {
	val := l.stack.get(idx)
	return convertToBoolean(val)
}

func (l luaState) ToInteger(idx int) int64 {
	n, _ := l.ToIntegerX(idx)
	return n
}

func (l luaState) ToIntegerX(idx int) (int64, bool) {
	val := l.stack.get(idx)
	return convertToInteger(val)
}

func (l luaState) ToNumber(idx int) float64 {
	n, _ := l.ToNumberX(idx)
	return n
}

func (l luaState) ToNumberX(idx int) (float64, bool) {
	val := l.stack.get(idx)
	return convertToFloat(val)
}

func (l luaState) ToString(idx int) string {
	s, _ := l.ToStringX(idx)
	return s
}

func (l luaState) ToStringX(idx int) (string, bool) {
	val := l.stack.get(idx)
	switch x := val.(type) {
	case string:
		return x, true
	case int64, float64:
		s := fmt.Sprintf("%v", x)
		l.stack.set(idx, s)
		return s, true
	default:
		return "", false
	}
}

func (l luaState) PushNil() {
	l.stack.push(nil)
}

func (l luaState) PushBoolean(b bool) {
	l.stack.push(b)
}

func (l luaState) PushInteger(n int64) {
	l.stack.push(n)
}

func (l luaState) PushNumber(n float64) {
	l.stack.push(n)
}

func (l luaState) PushString(s string) {
	l.stack.push(s)
}

var _ LuaState = New()

func New() *luaState {
	return &luaState{
		stack: newLuaStack(20),
	}
}
