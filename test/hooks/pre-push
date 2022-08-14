#!/bin/bash

# test pre-push hook for testing the lazygit credentials view
#
# to enable, use:
# chmod +x .git/hooks/pre-push
#
# this will hang if you're using git from the command line, so only enable this
# when you are testing the credentials view in lazygit

exec < /dev/tty

echo -n "Username for 'github': "
read username

echo -n "Password for 'github': "
# this will print the password to the log view but real git won't do that.
# We could use read -s but that's not POSIX compliant.
read password

if [ "$username" = "username" -a "$password" = "password" ]; then
  echo "success"
  exit 0
fi

>&2 echo "incorrect username/password"
exit 1
