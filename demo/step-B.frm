inputs:
	"/":
		tag:  "base"
		silo:
			- "file+ca://./wares/"
			- "http+ca://repeatr.s3.amazonaws.com/assets/"
	"/lib/asset-A":
		tag:  "asset-A"
		silo: "file+ca://./wares/"
action:
	command:
		- "/bin/bash"
		- "-c"
		- |
			set -euo pipefail ; set -x
			mkdir -p output/one output/two
			cp -a /lib/asset-A output/one/lib-A
			cp -a /lib/asset-A output/two/lib-A
			echo "+processB1" > output/one/b.sh
			echo "+processB2" > output/two/b.sh
outputs:
	"/task/output/one":
		tag: "asset-B1"
		type: "tar"
		silo: "file+ca://./wares/"
	"/task/output/two":
		tag: "asset-B2"
		type: "tar"
		silo: "file+ca://./wares/"
