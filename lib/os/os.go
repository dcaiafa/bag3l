package os

import (
	"os"

	"github.com/dcaiafa/bag3l/internal/vm"
)

func home_dir0(m *vm.VM) (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return dir, nil
}

func get_workdir0(m *vm.VM) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return dir, nil
}

func set_workdir0(m *vm.VM, dir string) error {
	err := os.Chdir(dir)
	if err != nil {
		return err
	}
	return nil
}
