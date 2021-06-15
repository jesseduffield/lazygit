module github.com/jesseduffield/lazygit

go 1.14

require (
	github.com/OpenPeeDeeP/xdg v1.0.0
	github.com/atotto/clipboard v0.1.2
	github.com/aybabtme/humanlog v0.4.1
	github.com/cli/safeexec v1.0.0
	github.com/cloudfoundry/jibber_jabber v0.0.0-20151120183258-bcc4c8345a21
	github.com/creack/pty v1.1.11
	github.com/fatih/color v1.9.0
	github.com/fsnotify/fsnotify v1.4.9
	github.com/gdamore/tcell/v2 v2.3.8 // indirect
	github.com/go-errors/errors v1.1.1
	github.com/go-logfmt/logfmt v0.5.0 // indirect
	github.com/golang-collections/collections v0.0.0-20130729185459-604e922904d3
	github.com/imdario/mergo v0.3.11
	github.com/integrii/flaggy v1.4.0
	github.com/jesseduffield/go-git/v5 v5.1.2-0.20201006095850-341962be15a4
	github.com/jesseduffield/gocui v0.3.1-0.20210417110745-37f79434200d
	github.com/jesseduffield/yaml v2.1.0+incompatible
	github.com/kardianos/osext v0.0.0-20190222173326-2bc1f35cddc0
	github.com/konsorten/go-windows-terminal-sequences v1.0.2 // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-colorable v0.1.7 // indirect
	github.com/mattn/go-runewidth v0.0.13
	github.com/maxbrunsfeld/counterfeiter/v6 v6.4.1
	github.com/mgutz/str v1.2.0
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.13.0
	github.com/sahilm/fuzzy v0.1.0
	github.com/sirupsen/logrus v1.4.2
	github.com/spkg/bom v0.0.0-20160624110644-59b7046e48ad
	github.com/stretchr/testify v1.5.1
	golang.org/x/crypto v0.0.0-20201002170205-7f63de1d35b0 // indirect
	golang.org/x/net v0.0.0-20210525063256-abc453219eb5 // indirect
	golang.org/x/sys v0.0.0-20210603125802-9665404d3644 // indirect
	golang.org/x/term v0.0.0-20210503060354-a79de5458b56 // indirect
)

replace github.com/go-git/go-git/v5 => github.com/jesseduffield/go-git/v5 v5.1.1
