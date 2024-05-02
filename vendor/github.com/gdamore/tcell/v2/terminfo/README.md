This package represents the parent for all terminals.

In older versions of tcell we had (a couple of) different
external file formats for the terminal database.  Those are
now removed.  All terminal definitions are supplied by
one of two methods:

1. Compiled Go code

2. For systems with terminfo and infocmp, dynamically
   generated at runtime.

The Go code can be generated using the mkinfo utility in
this directory.  The database entry should be generated
into a package in a directory named as the first character
of the package name.  (This permits us to group them all
without having a huge directory of little packages.)

It may be desirable to add new packages to the extended
package, or -- rarely -- the base package.

Applications which want to have the large set of terminal
descriptions built into the binary can simply import the
extended package.  Otherwise a smaller reasonable default
set (the base package) will be included instead.
