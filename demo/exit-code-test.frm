inputs:
	"/":
		type: "tar"
		hash: "uJRF46th6rYHt0zt_n3fcDuBfGFVPS6lzRZla5hv6iDoh5DVVzxUTMMzENfPoboL"
		silo: "http+ca://repeatr.s3.amazonaws.com/assets/"
action:
  env:
    EXIT_CODE: "1"
	command:
		- "/bin/bash"
		- "-c"
		- |
			set -euo pipefail
			exit ${EXIT_CODE}
