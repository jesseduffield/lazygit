# _Tcell_ on Plan 9

> [!NOTE]
> Plan 9 is supported on a best-effort basis, as the main _Tcell_ development team does not have a Plan 9 environment.

The Plan 9 backend opens `/dev/cons` for I/O, enables raw mode by writing `rawon`/`rawoff` to `/dev/consctl`.
It watches `/dev/wctl` for resize notifications.

The default mode for `vt((1)` is VT100, which will only provide basic monochrome text, and few additional features.
In this case, it is expected that `TERM=vt100` is set.

It may be possible to emulate more modern terminals using `-2` (VT220), `-a` (ANSI), or `-x` (XTerm) flags to `vt`.
While this has not been tested, the use of `-x` to get xterm like features, combinerd with a `TERM=xterm` may yield superior results,
including possibly color and mouse support.

Note that if _Tcell_ does not find a suitable value for `TERM` in the environment, it will assume XTerm like functionality.
