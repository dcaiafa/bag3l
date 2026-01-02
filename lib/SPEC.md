# stubgen conversion

Convert all functions in `GlobalPackage` @/Users/dcaiafa/src/my/bag3l/lib/lib.go#23 to using stubgen.

- Spawn a clean agent starting with empty context for each function to convert.
- Use `GlobalPackage` in lib.go to track which functions to update; Process them
  in order.
- Skip the functions starting with "$" for now.
- If it is not clear how a function should be converted, add the comment `//
  REVISIT` in from of the function line in `lib.go` and move to the next
  function.
- For each function, use `git mv` for the implementation to move the
  implementation and tests (if existing) to the `lib/global` directory.
- Update the package name from `lib` to `global`.
- If the source file includes multiple bag3l functions, it is OK to move them
  all in a single operation.
- Refer to the stubgen documentation @/Users/dcaiafa/src/my/bag3l/internal/stub/README.md for instructions on how to convert the function.
- After the conversion, remove the function from `GlobalPackage`.
- Run all tests in `lib` using `go test ./...`.
- If the function includes no tests at all, write some tests.
- After tests pass, commit the changes; don't forget to `git add` new test
  files. There should be one commit per function (or source file) that was
  converted.

