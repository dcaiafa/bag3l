package lib

import (
	"fmt"
	"os"
	"testing"

	"github.com/dcaiafa/bag3l/internal/stub"
	"github.com/dcaiafa/bag3l/lib/exec"
)

func TestExec(t *testing.T) {
	RunSubO(t, `complex-arg`,
		"e`go run ./testexec/testexec.go -print-args single-{\"arg\"}-yeah last` |"+`
			io.out
	`, `
[single-arg-yeah]
[last]
  `)

	RunSubErr(t, `complex-arg-expand-err`,
		"e`go run ./testexec/testexec.go -print-args single-{[1,2]...}-yeah last` |"+`
			io.out
	`, nil)

	RunSubO(t, `literals`,
		"e`go run ./testexec/testexec.go -print-args ./some/path   \\other\\path 123 !@#$%^&*()[]| ` |\n"+`
			io.out
	`, `
[./some/path]
[\other\path]
[123]
[!@#$%^&*()[]|]
  `)

	RunSubO(t, `exprs`, `
  var pets = [
    { name: "Chewie", alive: false },
    { name: "Deedee", alive: false },
    { name: "Ollie", alive: true },
  ]
`+"e`go run ./testexec/testexec.go -print-args {1} {[\"hello\", \"world\"] | str.join(\" \")} {pets | filter(&p->not p.alive) | map(&p->p.name)...}` |\n"+`
			io.out
	`, `
[1]
[hello world]
[Chewie]
[Deedee]
  `)

	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}
	RunSubO(t, `expr_home`,
		"e`go run ./testexec/testexec.go -print-args ~ a=~ b=~/foo` | io.out",
		fmt.Sprintf("[%v]\n[a=%v]\n[b=%v/foo]", homeDir, homeDir, homeDir))

	RunSubO(t, `expr_nil`, `
`+"e`go run ./testexec/testexec.go -print-args a {nil} b` |\n"+`
			io.out
	`, `
[a]
[b]
  `)

	RunSubO(t, `expr_iterator`, `
`+"e`go run ./testexec/testexec.go -print-args a {range(5) | filter(&n->mod(n,2)==0)...} b` |\n"+`
			io.out
	`, `
[a]
[0]
[2]
[4]
[b]
  `)

	RunSubO(t, `expr_iterable`, `
  var l = ["hi", 123, nil, 3.1415]
`+"e`go run ./testexec/testexec.go -print-args a {l...} b` |\n"+`
			io.out
	`, `
[a]
[hi]
[123]
[3.1415]
[b]
  `)

	RunSubO(t, `no_input`, `
		exec.exec(["go", "run", "./testexec/testexec.go", "-range", "1024", "-range-stdout"]) |
			str.lines() |
			map(&e -> int.parse(e)) |
			reduce(sum) |
			print()
	`, `523776`)

	RunSubO(t, `input`, `
		range(100000) | 
			map(str.into) |
			stream |
			exec.exec(["go", "run", "./testexec/testexec.go", "-echo-to-stdout"]) |
				read() |
				str.lines() |
				count() |
				print()
  `, `100000`)

	RunSubO(t, `pipe`, `
		exec.exec(["go", "run", "./testexec/testexec.go", "-range", "100000", "-range-stdout"]) |
		exec.exec(["go", "run", "./testexec/testexec.go", "-echo-to-stdout"]) |
		str.lines() |
		count() |
		print() 
	`, `100000`)

	RunSubO(t, `input2`, `
		range(100001) |
			map(str.into) |
			stream |
			exec.exec(["go", "run", "./testexec/testexec.go", "-echo-to-stdout"]) |
			str.lines() |
			map(&e -> int.parse(e)) |
			reduce(sum) |
			print()
	`, `5000050000`)

	RunSubO(t, `capture_stderr`, `
			var err_buf = buf.new()
			var out = range(2049) |
				map(str.into) |
				stream |
				exec.exec(["go", "run", "./testexec/testexec.go", "-echo-to-stderr", "-range", "11", "-range-stdout"]) |
        exec.with_stderr(err_buf)
	    var out_sum = out | str.lines() | map(&l -> int.parse(l)) | reduce(sum)
			var err_sum = err_buf | str.lines() | map(&l -> int.parse(l)) | reduce(sum)
			print(out_sum, err_sum)
		`, `55 2098176`)

	RunSubO(t, `process_returns_non_zero`, `
		try {
			"this will go in the error\n" |
				exec.exec(["go", "run", "./testexec/testexec.go", "-echo-to-stderr", "-exit-code=128"]) |
				read()
		} catch e {
			print(e.error)
		}
	`, `
exit status 1
this will go in the error
exit status 128
`)

	RunSubO(t, `abort`, `
		exec.exec(["go", "run", "./testexec/testexec.go", "-range", "1024", "-range-stdout"]) |
			str.lines() |
			take(10) |
			map(&e -> int.parse(e)) |
			reduce(sum) |
			print()
	`, `45`)

	RunSubO(t, `nil_stdin`, `
		nil |
			exec.exec(["go", "run", "./testexec/testexec.go", "-echo-to-stdout"]) |
			read() |
			len() |
			print()
	`, `0`)

	RunSubO(t, `stdin_is_file`, `
		var tmp = file.create_temp()
		defer file.remove(tmp)
		range(100000) | 
			map(str.into) |
			stream |
			tmp
		file.seek(tmp, 0)
		tmp |
			exec.exec(["go", "run", "./testexec/testexec.go", "-echo-to-stdout"]) |
		  read |
			str.lines |
			count |
			print
`, `100000`)

	RunSubO(t, `stderr_is_file`, `
      var tmp = file.create_temp()
      defer file.remove(tmp)
      range(100000) |
        map(str.into) |
				stream |
        exec.exec(["go", "run", "./testexec/testexec.go", "-echo-to-stderr"]) |
        exec.with_stderr(tmp) |
        discard
      file.seek(tmp, 0)
      read(tmp) |
        str.lines() |
        count() |
        print()
   `, `100000`)

	RunSubO(t, `combined_output`, `
		range(11) |
			map(str.into) |
			stream |
			exec.exec(["go", "run", "./testexec/testexec.go", "-echo-to-stderr", "-range", "11", "-range-stdout"]) |
			exec.with_combined_output() |
			str.lines() |
			map(&l -> int.parse(l)) |
			reduce(sum) |
			print()
	`, `110`)

	RunSubO(t, `combined_output_nonzero_exit`, `
		var out = buf.new()
		try {
			"this goes to the reader\n" |
				exec.exec(["go", "run", "./testexec/testexec.go", "-echo-to-stderr", "-exit-code=3"]) |
				exec.with_combined_output() |
				out
		} catch e {
			print(e.error)
		}
		print(read(out))
	`, `
exit status 1
this goes to the reader
exit status 3
`)

	RunSubO(t, `combined_output_conflicts_with_stderr`, `
		try {
			exec.exec(["go", "run", "./testexec/testexec.go"]) |
				exec.with_stderr(buf.new()) |
				exec.with_combined_output() |
				read()
		} catch e {
			print(e.error)
		}
	`, `with_stderr and with_combined_output are mutually exclusive`)

	RunSubO(t, `read_all`, `
		var out, code = range(5) |
			map(str.into) |
			stream |
			exec.exec(["go", "run", "./testexec/testexec.go", "-echo-to-stdout"]) |
			exec.read_all()
		print(out | str.lines() | count(), code)
	`, `5 0`)

	RunSubO(t, `read_all_nonzero_exit`, `
		// go run exits with 1 when the program exits with non-zero.
		var out, code = "this goes to the reader\n" |
			exec.exec(["go", "run", "./testexec/testexec.go", "-echo-to-stderr", "-exit-code=3"]) |
			exec.with_combined_output() |
			exec.read_all()
		print(code)
		print(out)
	`, `
1
this goes to the reader
exit status 3
`)

	RunSubO(t, `read_all_executable_not_found`, `
	try {
		exec.exec(["./foo_bar_doesnt_exist"]) | exec.read_all()
	} catch {
		print("failed as expected")
	}
`, `failed as expected`)

	RunSubO(t, `with_partial_env`, `
		exec.exec(["go", "run", "./testexec/testexec.go", "-print-env", "BAG3L_TEST_VAR"]) |
			exec.with_partial_env("BAG3L_TEST_VAR=hello") |
			io.out
	`, `hello`)

	RunSubO(t, `with_env_no_vars`, `
		exec.exec(["go", "run", "./testexec/testexec.go", "-print-env", "BAG3L_TEST_VAR"]) |
			exec.with_env() |
			exec.with_partial_env("BAG3L_TEST_VAR=hello") |
			io.out
	`, `hello`)

	RunSubErr(t, `with_env_non_string`, `
		exec.exec(["go", "run", "./testexec/testexec.go"]) |
			exec.with_env(1) |
			read()
	`, nil)

	RunSubErr(t, `no_args`, `exec.exec()`, stub.ErrInsufficientArgs)

	RunSubErr(t, `stdin_without_command`, `"input" | exec.exec()`, stub.ErrInsufficientArgs)

	RunSubErr(t, `empty_command`, `exec.exec([]) | read()`, exec.ErrEmptyCommand)

	RunSubErr(t, `nil_only_command`, "e`{nil}` | read()", exec.ErrEmptyCommand)

	RunSubO(t, `executable_not_found`, `
	try {
		exec.exec(["./foo_bar_doesnt_exist"]) | read()
	} catch {
  	print("failed as expected")
	}
`, `failed as expected`)
}
