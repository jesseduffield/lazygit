# _Tcell_ on Windows

Windows is supported starting from either Windows 10 version 1703 (the Creators Update, released in April 2017)
or Windows Server 2016.

> [NOTE!]
> Windows 8 and earlier are _not_ supported!

On Windows, _Tcell_ uses the modern VT API for console applications, including the _win32-input-mode_ to give
advanced key reporting. However, in order to receive resize notifications, the `ReadConsoleInput` API is
used even while in VT input mode.

On Windows we use 24-bit color by default.

## Terminal Emulators

Our preferred terminal emulator on Windows 11 is actually the modern _Windows Terminal_, which can also be
used as a portable application.  This is the main test target we use, and the one we recommend if problems are
encountered using third party applications.

We have tested _Alacritty_ and found it to work reasonably well, including support for modern
key reporting, mouse events, and bracketed paste, although its support support for advanced Unicode features
is limited.

We have also tested _WezTerm_, _Termius_, _Putty_, and _MobaXterm_.
These appear to work reasonably for remote sessions, but do not support the newer keyboard protocols, and
may not work well for local sessions at all.

_Termius_ in particular gave a poor experience when used for local sessions, and we would not recommend it.

While historically popular, we cannot recommend _ConEmu_ or _mintty_.  These applications have not kept pace
with modern APIs nor modern terminal standards, and we found their experience suboptimal.

## SSH From Windows

Unfortunately, we have found that some features (particularly the rich keyboard support) are degraded when
using SSH from a Windows Terminal (and probably other terminals as well!) to a remote host. It appears that
the full _win32-input-mode_ is captured and refactored into legacy VT style encodings, which can result in
old-school issues such as the inability to to distinguish CTRL-I from TAB.  This appears to be something
done by the Windows SSH or terminal application and is nothing we can fix.

It's likely that WSL will suffer these limitations as well.

Further, you'll probably need to request 24-bit color explicitly by setting `COLORTERM=truecolor` in your
environment, as this is not typically done for Windows terminals as it is for Posix terminals.

Note that _Alacritty_ has support for the full key bindings both locally and remotely. Others may as well.
_WezTerm_, while not great locally, performed quite well as an SSH client, with complete support for all
various modern keyboard featuers, bracketed paste, and good Unicode support.

Another option might be to run an X11 server (such as _MobaXterm_) and then remotely display a Linux terminal
application such as _Ghostty_ or _Kitty_.
