.PHONY: test

# Runs the suite twice: as shipped, and with the VM debug assertions enabled
# (see debugStackCheck in internal/vm).
test:
	go test ./...
	go test -tags bag3ldebug ./...
