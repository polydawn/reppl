inputs:
	"/":
		tag:  "base"
		silo:
			- "file+ca://./wares/"
			- "http+ca://repeatr.s3.amazonaws.com/assets/"
action:
	command:
		- "/bin/bash"
		- "-c"
		- |
			set -euo pipefail ; set -x
			mkdir output
			echo "wheeeee asset A" > output/wow
outputs:
	"/task/output":
		tag: "asset-A"
		type: "tar"
		silo: "file+ca://./wares/"
