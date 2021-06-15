# Code Generation

In order to make our code easy to test, we need to mock certain structs so that we're not e.g. running actual commands on the machine's OS as part of our unit tests. This means where we might have otherwise had a pointer to a Git struct as a field in one of our structs, we now have an IGit instead. IGit is simply the interface of the Git struct; using this interface allows us to pass in a mock struct which implements the interface but stubs the methods and performs assertions on which methods get called.

Unfortunately, golang has no easy way of saying that some interface is derived from a struct, so it's up to us to keep our interfaces in-sync with the structs they represent. We make use of two commands to help us keep things in sync: ifacemaker and counterfeiter

ifacemaker --file="pkg/commands/\*.go" --struct=Git --iface=IGit --pkg=commands -o pkg/commands/igit.go --doc false --comment="\$(cat pkg/commands/auto-generation-message.txt)"

counterfeiter pkg/commands IGitConfig
