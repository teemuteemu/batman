package vm

import (
	"fmt"

	"github.com/robertkrimen/otto"

	"github.com/teemuteemu/batman/pkg/client"
	"github.com/teemuteemu/batman/pkg/env"
)

type VM struct {
	engine *otto.Otto
	env    env.Env
}

func New(env *env.Env) (*VM, error) {
	engine := otto.New()

	vm := &VM{
		engine: engine,
		env:    *env,
	}

	err := engine.Set("setEnv", vm.SetEnv)
	if err != nil {
		return nil, err
	}

	err = engine.Set("getEnv", vm.GetEnv)
	if err != nil {
		return nil, err
	}

	err = engine.Set("stop", vm.Stop)
	if err != nil {
		return nil, err
	}

	err = engine.Set("assert", vm.Assert)
	if err != nil {
		return nil, err
	}

	err = engine.Set("goto", vm.Goto)
	if err != nil {
		return nil, err
	}

	return vm, nil
}

func setupResponse(vm *otto.Otto, resp *client.Response) error {
	responseObject, err := vm.Object(`response = {}`)
	if err != nil {
		return err
	}

	err = responseObject.Set("statusCode", resp.StatusCode)
	if err != nil {
		return err
	}

	bodyObj, err := vm.Object(fmt.Sprintf("body = %s", resp.Body))
	if err != nil {
		return err
	}

	err = responseObject.Set("body", bodyObj)
	if err != nil {
		return err
	}

	err = responseObject.Set("header", resp.Header)
	if err != nil {
		return err
	}

	err = vm.Set("response", responseObject)
	if err != nil {
		return err
	}

	return nil
}

func (vm *VM) ExecuteScript(script *string, resp *client.Response) (string, error) {
	if resp != nil {
		err := setupResponse(vm.engine, resp)
		if err != nil {
			return "", err
		}
	}

	result, err := vm.engine.Run(*script)
	if err != nil {
		return "", err
	}

	if result.IsObject() {
		resultObject := result.Object()

		gotoVal, err := resultObject.Get("goto")
		if err != nil {
			return "", err
		}

		gotoStep := gotoVal.String()
		if err != nil {
			return "", err
		}

		return gotoStep, nil
	}

	return "", nil
}
