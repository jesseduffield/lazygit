## Change Foreground Color

This command should enable a blue foreground color:

```bash
echo -ne "\033]10;#0000ff\007"
```

## Change Background Color

This command should enable a green background color:

```bash
echo -ne "\033]11;#00ff00\007"
```

## Change Cursor Color

This command should enable a red cursor color:

```bash
echo -ne "\033]12;#ff0000\007"
```

## Query Color Scheme

These two commands should print out the currently active color scheme:

```bash
echo -ne "\033]10;?\033\\"
echo -ne "\033]11;?\033\\"
```

## Query Cursor Position

This command should print out the current cursor position:

```bash
echo -ne "\033[6n"
```

## Set Window Title

This command should set the window title to "Test":

```bash
echo -ne "\033]2;Test\007" && sleep 10
```

## Bracketed paste

Enter this command, then paste a word from the clipboard. The text
displayed on the terminal should contain the codes `200~` and `201~`:

```bash
echo -ne "\033[?2004h" && sleep 10
```
