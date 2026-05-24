package runtime

import "github.com/dcaiafa/bag3l/internal/vm"

//go:generate go run ../../internal/stub/stubgen runtime.stubgen

type scriptDirUserDataKey struct{}

func SetScriptDir(m *vm.VM, dir string) {
	m.SetUserData(scriptDirUserDataKey{}, dir)
}

func script_dir0(m *vm.VM) (string, error) {
	scriptDir, ok := m.GetUserData(scriptDirUserDataKey{}).(string)
	if !ok {
		panic("scriptDir is not set")
	}
	return scriptDir, nil
}
