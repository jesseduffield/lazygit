"""Make a copy of the Flog Symbols font whose lines fill a terminal cell.

The demo recordings draw the commit graph with the branch drawing symbols, and
the terminal that vhs records takes them from the Flog Symbols Demo font in
demo/fonts. This script made its regular and bold faces from FlogSymbols.ttf of
https://github.com/rbong/flog-symbols (see LICENSE-FlogSymbols):

    pip install fonttools
    python3 demo/fonts/fit_flog_symbols.py FlogSymbols.ttf "Flog Symbols Demo" \\
        demo/fonts

Flog Symbols draws its lines for a cell that is 620 units wide and reaches from
-206 to 1006 units. A terminal that takes the symbols from a fallback font
draws them in the cells of its main font, and if those are larger, the lines
stop short of the cell edges and leave gaps between neighbouring cells.

So move the ends of the strokes that run to an edge of the cell out to the
terminal's cell edges, and centre everything else in the cell. The circles and
bends in the middle of the cell keep their shape.
"""

import os
import sys

from fontTools.ttLib import TTFont

# The cell of the symbols, in font units
LEFT, RIGHT, BOTTOM, TOP = -6, 626, -206, 1006

# The width of the terminal cells in the recordings, in font units: at the font
# size in demo/settings.tape they are 16 pixels wide (SauceCodePro's 14.4
# pixels, plus the pixel of letter spacing that vhs adds, rounded up), which is
# 16/24 of an em. Regenerate the font if the font size changes.
CELL_WIDTH = 667

# How far the horizontal strokes reach into the neighbouring cells, in font
# units. xterm.js doesn't clip a character to its cell horizontally, so this
# has to stay short of where the bends in a neighbouring cell begin. In the
# recordings, less than 30 leaves a dim line where two cells meet, and 35 or
# more opens a dark gap there.
HORIZONTAL_OVERLAP = 30

# How far past the top and bottom of the symbols' cell the vertical strokes
# reach, in font units. xterm.js clips a character to its row, so this only
# has to be more than the terminal's row sticks out beyond the symbols' cell.
VERTICAL_REACH = 400

# Points this close to an edge of the symbols' cell belong to the end of a
# stroke that runs to that edge
EDGE_TOLERANCE = 30


def main():
    flog_path, family, output_dir = sys.argv[1:4]
    font = TTFont(flog_path)
    fit_to_cell(font)

    # The bold face has the same outlines. Without one, the browser makes the
    # symbols of bold text bold itself by thickening them, and that leaves gaps
    # where they meet.
    for style in ("Regular", "Bold"):
        set_style(font, family, style)
        font.save(os.path.join(output_dir, f"{family.replace(' ', '')}-{style}.ttf"))


def fit_to_cell(font):
    glyf = font["glyf"]

    # Centre the symbols in the terminal's cell
    dx = round((CELL_WIDTH - (RIGHT + LEFT)) / 2)

    for name in font.getGlyphOrder():
        glyph = glyf[name]
        if glyph.numberOfContours <= 0:
            continue
        coordinates = glyph.coordinates
        for i, (x, y) in enumerate(coordinates):
            if x <= LEFT + EDGE_TOLERANCE:
                x = -HORIZONTAL_OVERLAP
            elif x >= RIGHT - EDGE_TOLERANCE:
                x = CELL_WIDTH + HORIZONTAL_OVERLAP
            else:
                x += dx
            if y <= BOTTOM + EDGE_TOLERANCE:
                y -= VERTICAL_REACH
            elif y >= TOP - EDGE_TOLERANCE:
                y += VERTICAL_REACH
            coordinates[i] = (x, y)
        glyph.recalcBounds(glyf)
        font["hmtx"][name] = (CELL_WIDTH, glyph.xMin)


def set_style(font, family, style):
    bold = style == "Bold"
    postscript_name = f"{family.replace(' ', '')}-{style}"
    full_name = family if not bold else f"{family} {style}"
    for record in font["name"].names:
        if record.nameID == 1:
            record.string = family
        elif record.nameID == 2:
            record.string = style
        elif record.nameID == 4:
            record.string = full_name
        elif record.nameID in (3, 6):
            record.string = postscript_name

    os2 = font["OS/2"]
    os2.usWeightClass = 700 if bold else 400
    fs_bold, fs_regular = 1 << 5, 1 << 6
    os2.fsSelection &= ~(fs_bold | fs_regular)
    os2.fsSelection |= fs_bold if bold else fs_regular
    font["head"].macStyle = 1 if bold else 0


if __name__ == "__main__":
    main()
