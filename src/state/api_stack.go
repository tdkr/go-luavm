package state

func (self *luaState) GetTop() int {
	return self.stack.top
}

func (self *luaState) AbsIndex(idx int) int {
	return self.stack.absIndex(idx)
}

func (self *luaState) CheckStack(n int) bool {
	self.stack.check(n)
	return true
}

func (self *luaState) Pop(n int) {
	for i := 0; i < n; i++ {
		self.stack.pop()
	}
}

func (self *luaState) Copy(fromIdx, toIdx int) {
	val := self.stack.get(fromIdx)
	self.stack.set(toIdx, val)
}

/* PushValue()方法把指定索引处的值推入栈顶。*/
func (self *luaState) PushValue(idx int) {
	val := self.stack.get(idx)
	self.stack.push(val)
}

/* Replace()是PushValue()的反操作：将栈顶值弹出，然后写入指定位置。 */
func (self *luaState) Replace(idx int) {
	val := self.stack.pop()
	self.stack.set(idx, val)
}

/* Insert()方法将栈顶值弹出，然后插入指定位置。 */
func (self *luaState) Insert(idx int) {
	self.Rotate(idx, 1)
}

func (self *luaState) Remove(idx int) {
	self.Rotate(idx, -1)
	self.Pop(1)
}

/* Rotate()方法将[idx, top]索引区间内的值朝栈顶方向循环移动n个位置。如果n是负数，那么实际效果就是朝栈底方向移动。 */
func (self *luaState) Rotate(idx, n int) {
	// 三次反转法
	t := self.stack.top - 1
	p := self.stack.absIndex(idx) - 1
	var m int
	if n >= 0 {
		m = t - n
	} else {
		m = p - n - 1
	}
	self.stack.reverse(p, m)
	self.stack.reverse(m+1, t)
	self.stack.reverse(p, t)
}

/*
SetTop()方法将栈顶索引设置为指定值。如果指定值小于当前栈顶索引，效果则相当于弹出操作（指定值为0相当于清空栈）
如果指定值大于当前栈顶索引，则效果相当于推入多个nil值。
*/
func (self *luaState) SetTop(idx int) {
	top := self.AbsIndex(idx)
	if top < 0 {
		panic("stack underflow")
	}
	if n := self.stack.top - top; n > 0 {
		for i := 0; i < n; i++ {
			self.stack.pop()
		}
	} else if n < 0 {
		for i := 0; i > n; i-- {
			self.stack.push(nil)
		}
	}
}
