
# go-osc52

<p>
    <a href="https://github.com/aymanbagabas/go-osc52/releases"><img src="https://img.shields.io/github/release/aymanbagabas/go-osc52.svg" alt="Latest Release"></a>
    <a href="https://pkg.go.dev/github.com/aymanbagabas/go-osc52?tab=doc"><img src="https://godoc.org/github.com/golang/gddo?status.svg" alt="GoDoc"></a>
</p>

A terminal Go library to copy text to clipboard from anywhere. It does so using [ANSI OSC52](https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h3-Operating-System-Commands). The `Copy()` function defaults to copying text from terminals running locally.

To use this over SSH, using [gliderlabs/ssh](https://github.com/gliderlabs/ssh), use `NewOutput(sshSession, sshSession.Environ())` and make sure you pass the `TERM` environment variable in your SSH connection.

```sh
ssh -o SendEnv=TERM <host>
```

Tmux users need to pass an additional environment variable `TMUX`.

```sh
ssh -o SendEnv=TERM -o SendEnv=TMUX <host>
```

# Credits

* [vim-oscyank](https://github.com/ojroques/vim-oscyank) this is heavily inspired by vim-oscyank.