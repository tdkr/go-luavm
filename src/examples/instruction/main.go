package main

import (
	"fmt"
	"github.com/tdkr/go-luavm/src/api"
	"github.com/tdkr/go-luavm/src/binchunk"
	"github.com/tdkr/go-luavm/src/state"
	"github.com/tdkr/go-luavm/src/util"
	. "github.com/tdkr/go-luavm/src/vm"
	"io/ioutil"
	"os"
	"path"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadFile(path.Join(wd, "src/examples/instruction/luac.out"))
	if err != nil {
		panic(err)
	}

	proto := binchunk.Undump(data)
	fmt.Println("********* list proto **********")
	list(proto)
	fmt.Println("********* exec proto **********")
	luaMain(proto)
}

func luaMain(proto *binchunk.Prototype) {
	nRegs := int(proto.MaxStackSize)
	vm := state.New().(api.LuaVM)
	vm.SetTop(nRegs)
	for {
		pc := vm.PC()
		inst := Instruction(vm.Fetch())
		if inst.Opcode() != OP_RETURN {
			inst.Execute(vm)
			fmt.Printf("[%02d] %s ", pc+1, inst.OpName())
			util.PrintStack(vm)
		} else {
			break
		}
	}
}

func list(f *binchunk.Prototype) {
	util.PrintHeader(f)
	util.PrintCode(f)
	util.PrintDetail(f)
	for _, p := range f.Protos {
		list(p)
	}
}
