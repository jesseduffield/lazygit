git fetch
git checkout %GIT_COMMIT%

SET GOPATH=%CD%\Godeps\_workspace;c:\Users\Administrator\go
c:\Go\bin\go test -v .
