package runtime

import "github.com/dcaiafa/bag3l"

//go:generate go run ../../internal/stub/stubgen runtime.stubgen

type scriptDirUserDataKey struct{}

func SetScriptDir(vm *bag3l.VM, dir string) {
	vm.SetUserData(scriptDirUserDataKey{}, dir)
}

func script_dir0(vm *bag3l.VM) (string, error) {
	scriptDir, ok := vm.GetUserData(scriptDirUserDataKey{}).(string)
	if !ok {
		panic("scriptDir is not set")
	}
	return scriptDir, nil
}
