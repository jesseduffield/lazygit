#!/bin/bash

set -e

for element in ./styles/examples/*.md; do
    echo "Generating screenshot for element ${element}"
    basename="`basename -s .md ${element}`"
    stylename="${basename}.style"
    filename="${basename}.png"

    # take screenshot
    ./termshot -o ./styles/examples/ -f "$filename" glow -s ./styles/examples/${stylename} ${element}

    # add border
    convert -bordercolor black -border 16x16 "./styles/examples/$filename" "./styles/examples/$filename"

    # optimize filesize
    pngcrush -ow "./styles/examples/$filename"
done
