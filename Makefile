.PHONY: package clean

# Default target - clean first, then package
all: clean package

# Package target - builds the Go binary and runs dxt pack
package:
	go build
	chmod +x golang-mcp-testing
	dxt pack

# Optional clean target to remove built artifacts
clean:
	rm -f golang-mcp-testing
	rm -f golang-mcp-testing.dxt
