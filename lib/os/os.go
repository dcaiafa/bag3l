package os

import (
	"os"

	nitro "github.com/dcaiafa/bag3l"
)

func home_dir0(vm *nitro.VM) (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return dir, nil
}

func get_workdir0(vm *nitro.VM) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return dir, nil
}

func set_workdir0(vm *nitro.VM, dir string) error {
	err := os.Chdir(dir)
	if err != nil {
		return err
	}
	return nil
}
