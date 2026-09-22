# Terminal settings shared by all demo recordings. The tape that
# demo/record_demo.sh generates for a demo sources this file.

# Use the non-Mono Nerd Font variant; the Mono variant renders the icons too
# small.
Set FontFamily "SauceCodePro NF"
Set FontSize 24
Set LineHeight 1.0
Set Padding 20

# vhs sizes the terminal in pixels rather than in cells, so these two numbers
# are what produce a 120x35 grid at the font settings above. If you change the
# font, the font size or the padding, record a tape that runs `stty size` and
# adjust them until the grid is 120x35 again.
Set Width 1866
Set Height 1140

# There is no frame rate setting here because vhs ignores `Set Framerate` when
# it writes a video; it captures at 25 fps either way.

# The frame around the focused view is bold green, and xterm.js draws bold text
# in the bright variant of a colour, so brightGreen is the one those frames end
# up using.

Set Theme {
    "background": "#1d1d1d",
    "foreground": "#dddad6",
    "cursor": "#c7c7c7",
    "selection": "#44475a",
    "black": "#7a7a7a",
    "red": "#fc4384",
    "green": "#3fb950",
    "yellow": "#ffa727",
    "blue": "#102895",
    "magenta": "#c930c7",
    "cyan": "#00c5c7",
    "white": "#c7c7c7",
    "brightBlack": "#676767",
    "brightRed": "#ff7fac",
    "brightGreen": "#56d364",
    "brightYellow": "#ebdf86",
    "brightBlue": "#6871ff",
    "brightMagenta": "#ff76ff",
    "brightCyan": "#5ffdff",
    "brightWhite": "#fffefe"
}
