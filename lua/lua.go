package lua

import (
	"context"
	"errors"
	"time"

	glua "github.com/yuin/gopher-lua"
)

type ValueType int

const (
	Undefined ValueType = iota
	String
	Float64
	Table
)

var (
	ErrUnsupportedTableValue   = errors.New("unsupported config table value")
	ErrUnsupportedTableKey     = errors.New("unsupported config table key")
	ErrNilRequestedGlobalBAlue = errors.New("nil requested global value")
)

func extractTable(luaVM *glua.LState, table *glua.LTable) (config map[string]interface{}, err error) {
	defer func() {
		e, isError := recover().(error)
		if isError {
			err = e
		}
	}()
	config = map[string]interface{}{}
	table.ForEach(func(k, v glua.LValue) {
		var key string
		var val interface{}
		switch {
		case glua.LVCanConvToString(k):
			key = glua.LVAsString(k)
		default:
			if glua.LVCanConvToString(k) {
				key = glua.LVAsString(luaVM.ToStringMeta(k))
				break
			}
			panic(ErrUnsupportedTableKey)
		}
		switch v.Type() {
		case glua.LTBool:
			val = glua.LVAsBool(v)
		case glua.LTNumber:
			val = float64(glua.LVAsNumber(v))
		case glua.LTString:
			val = glua.LVAsString(v)
		case glua.LTTable:
			if glua.LVCanConvToString(v) {
				val = glua.LVAsString(luaVM.ToStringMeta(v))
				break
			}
			m, e := extractTable(luaVM, v.(*glua.LTable))
			if e != nil {
				panic(e)
			}
			val = m
		default:
			panic(ErrUnsupportedTableValue)
		}
		config[key] = val
	})
	return config, nil
}

func EvalConfig(configStr string, timeout time.Duration) (map[string]interface{}, error) {
	luaVM := glua.NewState(glua.Options{
		CallStackSize: 120,
		RegistrySize:  1024 * 100,
	})
	defer luaVM.Close()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	luaVM.SetContext(ctx)
	err := luaVM.DoString(configStr)
	if err != nil {
		return nil, err
	}
	configTable, ok := luaVM.GetGlobal("config").(*glua.LTable)
	if !ok {
		return nil, ErrNilRequestedGlobalBAlue
	}
	config, err := extractTable(luaVM, configTable)
	return config, err
}
