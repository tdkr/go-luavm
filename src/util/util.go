package util

import (
	"bytes"
	"fmt"
	"github.com/tdkr/go-luavm/src/api"
	"github.com/tdkr/go-luavm/src/binchunk"
	. "github.com/tdkr/go-luavm/src/vm"
)

func DumpProto(f *binchunk.Prototype) string {
	buf := bytes.NewBuffer([]byte{})
	buf.WriteString(DumpHeader(f))
	buf.WriteString(DumpCode(f))
	buf.WriteString(DumpDetail(f))
	for _, p := range f.Protos {
		buf.WriteString(DumpProto(p))
	}
}

func DumpHeader(f *binchunk.Prototype) string {
	funcType := "main"
	if f.LineDefined > 0 {
		funcType = "function"
	}

	varargFlag := ""
	if f.IsVararg > 0 {
		varargFlag = "+"
	}

	buf := bytes.NewBuffer([]byte{})

	buf.WriteString(fmt.Sprintf("%s <%s:%d,%d> (%d instructions)\n",
		funcType, f.Source, f.LineDefined, f.LastLineDefined, len(f.Code)))

	buf.WriteString(fmt.Sprintf("%d%s params, %d slots, %d upvalues, ",
		f.NumParams, varargFlag, f.MaxStackSize, len(f.Upvalues)))

	buf.WriteString(fmt.Sprintf("%d locals, %d constants, %d functions\n",
		len(f.LocVars), len(f.Constants), len(f.Protos)))

	return buf.String()
}

func DumpCode(f *binchunk.Prototype) string {
	buf := bytes.NewBuffer([]byte{})
	for pc, c := range f.Code {
		line := "-"
		if len(f.LineInfo) > 0 {
			line = fmt.Sprintf("%d", f.LineInfo[pc])
		}

		i := Instruction(c)
		buf.WriteString(fmt.Sprintf("\t%d\t[%s]\t%s \t", pc+1, line, i.OpName()))
		DumpOperands(i)
		buf.WriteString("\n")
	}
	return buf.String()
}

func DumpOperands(i Instruction) string {
	buf := bytes.NewBuffer([]byte{})
	switch i.OpMode() {
	case IABC:
		a, b, c := i.ABC()

		buf.WriteString(fmt.Sprintf("%d", a))
		if i.BMode() != OpArgN {
			if b > 0xFF {
				buf.WriteString(fmt.Sprintf(" %d", -1-b&0xFF))
			} else {
				buf.WriteString(fmt.Sprintf(" %d", b))
			}
		}
		if i.CMode() != OpArgN {
			if c > 0xFF {
				buf.WriteString(fmt.Sprintf(" %d", -1-c&0xFF))
			} else {
				buf.WriteString(fmt.Sprintf(" %d", c))
			}
		}
	case IABx:
		a, bx := i.ABx()

		buf.WriteString(fmt.Sprintf("%d", a))
		if i.BMode() == OpArgK {
			buf.WriteString(fmt.Sprintf(" %d", -1-bx))
		} else if i.BMode() == OpArgU {
			buf.WriteString(fmt.Sprintf(" %d", bx))
		}
	case IAsBx:
		a, sbx := i.AsBx()
		buf.WriteString(fmt.Sprintf("%d %d", a, sbx))
	case IAx:
		ax := i.Ax()
		buf.WriteString(fmt.Sprintf("%d", -1-ax))
	}
}

func DumpDetail(f *binchunk.Prototype) string {
	buf := bytes.NewBuffer([]byte{})
	buf.WriteString(fmt.Sprintf("constants (%d):\n", len(f.Constants)))
	for i, k := range f.Constants {
		buf.WriteString(fmt.Sprintf("\t%d\t%s\n", i+1, constantToString(k)))
	}

	buf.WriteString(fmt.Sprintf("locals (%d):\n", len(f.LocVars)))
	for i, locVar := range f.LocVars {
		buf.WriteString(fmt.Sprintf("\t%d\t%s\t%d\t%d\n",
			i, locVar.VarName, locVar.StartPC+1, locVar.EndPC+1))
	}

	buf.WriteString(fmt.Sprintf("upvalues (%d):\n", len(f.Upvalues)))
	for i, upval := range f.Upvalues {
		buf.WriteString(fmt.Sprintf("\t%d\t%s\t%d\t%d\n",
			i, upvalName(f, i), upval.Instack, upval.Idx))
	}
}

func constantToString(k interface{}) string {
	switch k.(type) {
	case nil:
		return "nil"
	case bool:
		return fmt.Sprintf("%t", k)
	case float64:
		return fmt.Sprintf("%g", k)
	case int64:
		return fmt.Sprintf("%d", k)
	case string:
		return fmt.Sprintf("%q", k)
	default:
		return "?"
	}
}

func upvalName(f *binchunk.Prototype, idx int) string {
	if len(f.UpvalueNames) > 0 {
		return f.UpvalueNames[idx]
	}
	return "-"
}

func DumpStack(ls api.LuaState) string {
	buf := bytes.NewBuffer([]byte{})
	top := ls.GetTop()
	for i := 1; i <= top; i++ {
		t := ls.Type(i)
		switch t {
		case api.LUA_TBOOLEAN:
			buf.WriteString(fmt.Sprintf("[%t]", ls.ToBoolean(i)))
		case api.LUA_TNUMBER:
			buf.WriteString(fmt.Sprintf("[%g]", ls.ToNumber(i)))
		case api.LUA_TSTRING:
			buf.WriteString(fmt.Sprintf("[%q]", ls.ToString(i)))
		default: // other values
			buf.WriteString(fmt.Sprintf("[%s]", ls.TypeName(t)))
		}
	}
	fmt.Println()
}
