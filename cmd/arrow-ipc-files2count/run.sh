#!/bin/sh

sample1=./sample1.file.ipc
sample2=./sample2.file.ipc

gencsv1(){
	echo timestamp,severity,status,body
	echo 2025-10-19T23:39:18.012345Z,INFO,200,apt update done
	echo 2025-10-18T23:39:18.012345Z,WARN,500,apt update failure
}

gencsv2(){
	echo timestamp,severity,status,body
	echo 2025-10-17T23:39:18.012345Z,INFO,200,apt update done
	echo 2025-10-16T23:39:18.012345Z,WARN,500,apt update failure
	echo 2025-10-15T23:39:18.012345Z,WARN,500,apt update failure
}

geninput(){
	echo generating sample ipc file...

	which ccsv2arrow2ipc | fgrep -q csv2arrow2ipc || exec sh -c '
		echo csv2arrow2ipc missing.
		echo see github.com/takanoriyanagitani/go-csv2arrow2ipc to install it.
		exit 1
	'

	gencsv1 | csv2arrow2ipc > "${sample1}"
	gencsv2 | csv2arrow2ipc > "${sample2}"
}

test -f "${sample1}" || geninput
test -f "${sample2}" || geninput

ls "${sample1}" "${sample2}" | ./arrow-ipc-files2count
