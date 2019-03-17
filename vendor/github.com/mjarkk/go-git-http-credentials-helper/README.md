# go-git-http-credentials-helper
A Go library that adds a secure password/username ask function to a git push/pull/fetch command over http

## The problem:
When executing git fetch/push/pull it might ask for credentials but git craches when it detects it's ran in a tty.  
This libary fixes that problem and makes it possible to fillin the username/password fields if there they are asked for.

## How to use:
```go
// 1. Import the package:
import (
  gitcredentialhelper "github.com/mjarkk/go-git-http-credentials-helper"
)

// 2: Add the SetupClient to the top of the application
// IMPORTANT: Make sure the application doesn't write to stdout before this function!
func main() {
  gitcredentialhelper.SetupClient()
}

// 3: run a git command that needs credentials
func pushToTheRepo() {
  // Create a exec.Cmd
  cmd := exec.Command("git", "push")
  
  // The output is the same as cmd.CombinedOutput()
  out, err := gitcredentialhelper.Run(cmd, func(question string) string {
    // question is "username" or "password"
    // and it expects to get back the username or password depending on the question
    fmt.Println("git asked for:", question)
		return "my username or password"
  })
  
  // ....
}
```

If git commands take forever you might need to set the Options argument:  
```go
// ....

// "your-app-name" needs to be the appname as snown in the terminal title or in task manager this does NOT contain any / or .
options := gitcredentialhelper.Options{AppName: "your-app-name"}
gitcredentialhelper.Run(cmd, askFunction, options)

// ....
```

If you have a logger and want to log more errors:  
```go
func main() {
  gitcredentialhelper.SetupClient(func(err error) {
    // DO NOT PRINT THE ERROR!
    // If you do that the error will be seen as password/username for git
    log.Error(err)
  })
}
```


## How it works:  
1. Your program calles `git push` in a http repo
2. This libary creates a webserver for username/password questions and adds the needed shell variables to the process
3. Git sees the `GIT_ASKPASS` and runs your program with as arugment: `username for repo "https://example.com/you/yourprogram"`
4. Now the proccess tree looks a bit like this: `yourProgram -> git -> yourProgram` (in reality this look more like `... -> git -> git -> git -> ...`)
5. This libary detects if `yourProgram` is started by git and if so runs it's own startup script (this is the `SetupClient()` function) 
6. This libary will create a keypair and do a network request to the main program that spawned git in the first place
7. The main program will pickup the network request and check if it is not a fake request. After that it will ask the ask function for a password/username and will encrypt the data with a public key from `yourProgram` and send the encrypted message back to the program that was created by git.
8. `yourProgram` spawned by git will decrypt the message and printout the contents of the message, after that it will exit the program.  
(There are still some things i've not mentioned but those things are mostly useless to know)  

## Q and A:
> Why all this if you can just use a pty?  

Mostly because of Windows.  
Windows does have a dll to create a PTY but there are no inplementations yet and i would need to included the dll because a lot of users don't have the dll.  
Also PTY support on the Windows 10 subsystem *(linux on windows)* as non-root user is completely broken.  
Beside all of that this is the offical way to do these kinds of things with git so..

> Is this secure?  

TL;DR Yes.  
The long answer is no, it's probebly possible to break this libary. Though i could not succeed in breaking this and if someone did break it would be still a 50% change to get any input due to security measures. It's easier for someone to add a keylogger in your terminal than to break this.

> Why so menny functions to inplement this?  

Most things are the result of security measures to check if it's really your program that is asking for someones password.

## Known issues:
- The credential manager on windows can break this if the popup is closed by the user. As user the best thing to do in that case is to remove the credential manager from git: `git config --system --unset credential.helper`
- This does not work with the `go run ...` 
