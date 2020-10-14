package main

import (
	"fmt"
	"path"
)
import "io/ioutil"
import "os"
import "github.com/tdkr/go-luavm/src/binchunk"
import . "github.com/tdkr/go-luavm/src/api"
import "github.com/tdkr/go-luavm/src/state"
import . "github.com/tdkr/go-luavm/src/vm"

func main() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadFile(path.Join(wd, "src/examples/table/luac.out"))
	if err != nil {
		panic(err)
	}

	proto := binchunk.Undump(data)
	luaMain(proto)
}

func luaMain(proto *binchunk.Prototype) {
	nRegs := int(proto.MaxStackSize)
	ls := state.New(nRegs+8)
	ls.SetTop(nRegs)
	for {
		pc := ls.PC()
		inst := Instruction(ls.Fetch())
		if inst.Opcode() != OP_RETURN {
			inst.Execute(ls)

			fmt.Printf("[%02d] %s ", pc+1, inst.OpName())
			printStack(ls)
		} else {
			break
		}
	}
}

func printStack(ls LuaState) {
	top := ls.GetTop()
	for i := 1; i <= top; i++ {
		t := ls.Type(i)
		switch t {
		case LUA_TBOOLEAN:
			fmt.Printf("[%t]", ls.ToBoolean(i))
		case LUA_TNUMBER:
			fmt.Printf("[%g]", ls.ToNumber(i))
		case LUA_TSTRING:
			fmt.Printf("[%q]", ls.ToString(i))
		default: // other values
			fmt.Printf("[%s]", ls.TypeName(t))
		}
	}
	fmt.Println()
}
