# Allowing cloning inside lazygit

## TODO
- Add it to the recent repos list
- Destination is not working properly
    "~" is not exapnded to home

## QOL

### Repo URL (QoL)

valid URLs
- https://github.com/willparsons/astrovim-config.git
- git@github.com:willparsons/astrovim-config.git
- willparsons/astrovim-config

user can enter entire clone link (https or ssh) or just the author/repo

if just the author/repo is given we need to then ask for https or ssh
then we can fill in the blanks

### Options
Some repos need more options when cloning like
    - `--depth 1`
    - `--bare`

### Suggestions
- Search suggestions?
