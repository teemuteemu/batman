package vm

import (
	"fmt"
	"os"

	"github.com/robertkrimen/otto"
)

func (vm *VM) SetEnv(call otto.FunctionCall) otto.Value {
	key := call.Argument(0).String()
	val := call.Argument(1).String()

	vm.env[key] = val

	return otto.Value{}
}

func (vm *VM) GetEnv(call otto.FunctionCall) otto.Value {
	key := call.Argument(0).String()

	val, err := vm.engine.ToValue(vm.env[key])
	if err != nil {
		fmt.Printf("GetEnv: Bad value %v", val)
	}

	return val
}

func (vm *VM) Stop(call otto.FunctionCall) otto.Value {
	os.Exit(1)

	return otto.Value{}
}

func (vm *VM) Assert(call otto.FunctionCall) otto.Value {
	value, err := call.Argument(0).ToBoolean()

	if err != nil {
		fmt.Println("Assert: Asserted value is not boolean")
	}

	if value {
		fmt.Println("✅")
	} else {
		fmt.Println("❌")
	}

	return otto.Value{}
}

func (vm *VM) Goto(call otto.FunctionCall) otto.Value {
	step := call.Argument(0).String()

	var stepValue interface{}

	if step == "undefined" {
		stepValue = true
	} else {
		stepValue = step
	}

	resultObject, err := vm.engine.Object(`result = {}`)
	if err != nil {
		fmt.Println("Goto: Bad result object")
	}

	err = resultObject.Set("goto", stepValue)
	if err != nil {
		fmt.Println("Goto: Could not set result object")
	}

	return resultObject.Value()
}
