#!/bin/bash
# For testing subprocesses that require input
# Ask the user for login details
read -p 'Username: ' user
read -sp 'Password: ' pass
echo
echo Hello $user
