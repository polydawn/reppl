inputs:
	"/":            {tag: "base"}
	"/lib/asset-A": {tag: "asset-A"}
	"/app/go":      {tag: "go"}
action:
	command:
		- "/bin/bash"
		- "-c"
		- |
			set -euo pipefail ; set -x
			export PATH="$PATH:/app/go/go/bin"
			export GOROOT="/app/go/go"
			cat >main.go <<EOF
				package main
				
				import "fmt"
				
				func main() {
					fmt.Println("Hello, 世界")
				}
			EOF
			mkdir bin
			go build -o bin/hello .
outputs:
	"/task/bin":
		tag: "hellogopher"
		type: "tar"
		silo: "file+ca://./wares/"
