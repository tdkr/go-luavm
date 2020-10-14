package main

import (
	"fmt"
	"github.com/tdkr/go-luavm/src/api"
	"github.com/tdkr/go-luavm/src/binchunk"
	"github.com/tdkr/go-luavm/src/state"
	"github.com/tdkr/go-luavm/src/util"
	"io/ioutil"
	"os"
	"path"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadFile(path.Join(wd, "src/examples/closure/luac.out"))
	if err != nil {
		panic(err)
	}

	util.PrintProto(binchunk.Undump(data))
	ls := state.New()
	ls.Register("print", print)
	ls.Register("fail", fail)
	ls.Load(data, "test.lua", "b")
	ls.Call(0, 0)
}

func print(ls api.LuaState) int {
	nArgs := ls.GetTop()
	for i := 1; i <= nArgs; i++ {
		if ls.IsBoolean(i) {
			fmt.Printf("%t", ls.ToBoolean(i))
		} else if ls.IsString(i) {
			fmt.Printf("%s", ls.ToString(i))
		} else {
			fmt.Print(ls.TypeName(ls.Type(i)))
		}
		if i < nArgs {
			fmt.Print("\t")
		}
	}
	fmt.Println()
	return 0
}

func fail(ls api.LuaState) int {
	fmt.Printf("fail lua state :%+v\n", ls)
	return 0
}
