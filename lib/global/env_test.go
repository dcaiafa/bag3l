package global_test

import (
	"os"
	"testing"

	"github.com/dcaiafa/bag3l/internal/btesting"
)

func TestEnv(t *testing.T) {
	// Set a known environment variable for testing.
	testEnvName := "BAG3L_TEST_ENV_VAR"
	testEnvValue := "test_value_12345"
	os.Setenv(testEnvName, testEnvValue)
	defer os.Unsetenv(testEnvName)

	btesting.RunSubO(t, "env_get_specific", `
		print(env("BAG3L_TEST_ENV_VAR"))
	`, testEnvValue)

	btesting.RunSubO(t, "env_get_nonexistent", `
		print(env("BAG3L_NONEXISTENT_VAR_12345"))
	`, "")

	btesting.RunSubO(t, "env_get_all", `
		var e = env()
		print(type(e))
	`, "list")

	btesting.RunSubO(t, "env_contains_test_var", `
		var all = env()
		var found = first(filter(all, &e -> e == "BAG3L_TEST_ENV_VAR=test_value_12345"))
		print(found)
	`, "BAG3L_TEST_ENV_VAR=test_value_12345")
}
