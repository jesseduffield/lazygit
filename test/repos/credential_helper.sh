#!/bin/bash
set -ex; rm -rf repo; mkdir repo; cd repo

helper_type=$1

git init

if [[ -n "$helper_type" ]]; then
	if [[ "$helper_type" == "credential" ]]; then
		git config credential.helper cache
	elif [[ "$helper_type" == "gitAskPass" ]]; then
		export GIT_ASKPASS="/bin/false"
	elif [[ "$helper_type" == "coreAskPass" ]]; then
		git config core.askPass /bin/false
	elif [[ "$helper_type" == "sshAskPass" ]]; then
		export SSH_ASKPASS="/bin/false"
	fi
else
	echo "No argument specified"
fi
