<!-- use this template to generate the contributor docs with the following command: `$ lingo run docs --template CONTRIBUTING_TEMPLATE.md  --output CONTRIBUTING.md` -->

# Contributing


♥ We love pull requests from everyone !


When contributing to this repository, please first discuss the change you wish
to make via issue, email, or any other method with the owners of this repository
before making a change. 

## So all code changes happen through Pull Requests
Pull requests are the best way to propose changes to the codebase. We actively
welcome your pull requests:

1. Fork the repo and create your branch from `master`.
2. If you've added code that should be tested, add tests.
3. If you've added code that need documentation, update the documentation.
4. Make sure your code follows the [effective go](https://golang.org/doc/effective_go.html) guidelines as much as possible.
5. Be sure to test your modifications.
6. Write a [good commit message](http://tbaggery.com/2008/04/19/a-note-about-git-commit-messages.html).
7. Issue that pull request!

## Code of conduct
Please note by participating in this project, you agree to abide by the [code of conduct].

[code of conduct]: https://github.com/jesseduffield/lazygit/blob/master/CODE-OF-CONDUCT.md

## Any contributions you make will be under the MIT Software License
In short, when you submit code changes, your submissions are understood to be
under the same [MIT License](http://choosealicense.com/licenses/mit/) that
covers the project. Feel free to contact the maintainers if that's a concern.

## Report bugs using Github's [issues](https://github.com/jesseduffield/lazygit/issues)
We use GitHub issues to track public bugs. Report a bug by [opening a new
issue](https://github.com/jesseduffield/lazygit/issues/new); it's that easy!

## Code Review Comments and Effective Go Guidelines
[CodeLingo](https://codelingo.io) automatically checks every pull request against the following guidelines from [Effective Go](https://golang.org/doc/effective_go.html) and [Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments).


## Go Error Format
Error strings should not be capitalized (unless beginning with proper nouns 
or acronyms) or end with punctuation, since they are usually printed following
other context. That is, use fmt.Errorf("something bad") not fmt.Errorf("Something bad"),
so that log.Printf("Reading %s: %v", filename, err) formats without a spurious 
capital letter mid-message. This does not apply to logging, which is implicitly
line-oriented and not combined inside other messages.


## Use Crypto Rand
Do not use package math/rand to generate keys, even 
throwaway ones. Unseeded, the generator is completely predictable. 
Seeded with time.Nanoseconds(), there are just a few bits of entropy. 
Instead, use crypto/rand's Reader, and if you need text, print to 
hexadecimal or base64


## Avoid Renaming Imports
Avoid renaming imports except to avoid a name collision; good package names
should not require renaming. In the event of collision, prefer to rename the
most local or project-specific import.


## Context as First Argument
Values of the context.Context type carry security credentials, tracing information, 
deadlines, and cancellation signals across API and process boundaries. Go programs 
pass Contexts explicitly along the entire function call chain from incoming RPCs 
and HTTP requests to outgoing requests.

Most functions that use a Context should accept it as their first parameter.


## Do Not Discard Errors
Do not discard errors using _ variables. If a function returns an error, 
check it to make sure the function succeeded. Handle the error, return it, or, 
in truly exceptional situations, panic.


## Avoid Annotations in Comments
Comments do not need extra formatting such as banners of stars. The generated output
may not even be presented in a fixed-width font, so don't depend on spacing for alignment—godoc, 
like gofmt, takes care of that. The comments are uninterpreted plain text, so HTML and other 
annotations such as _this_ will reproduce verbatim and should not be used. One adjustment godoc 
does do is to display indented text in a fixed-width font, suitable for program snippets. 
The package comment for the fmt package uses this to good effect.


## Comment First Word as Subject
Doc comments work best as complete sentences, which allow a wide variety of automated presentations.
The first sentence should be a one-sentence summary that starts with the name being declared.


## Package Comment
Every package should have a package comment, a block comment preceding the package clause. 
For multi-file packages, the package comment only needs to be present in one file, and any one will do. 
The package comment should introduce the package and provide information relevant to the package as a 
whole. It will appear first on the godoc page and should set up the detailed documentation that follows.


## Single Method Interface Name
By convention, one-method interfaces are named by the method name plus an -er suffix 
or similar modification to construct an agent noun: Reader, Writer, Formatter, CloseNotifier etc.

There are a number of such names and it's productive to honor them and the function names they capture. 
Read, Write, Close, Flush, String and so on have canonical signatures and meanings. To avoid confusion, 
don't give your method one of those names unless it has the same signature and meaning. Conversely, 
if your type implements a method with the same meaning as a method on a well-known type, give it the 
same name and signature; call your string-converter method String not ToString.

