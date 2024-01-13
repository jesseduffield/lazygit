#!/bin/bash

for style in ./styles/*.json; do
    echo "Generating screenshot for ${style}"
    filename="`basename -s .json ${style}`.png"

    light=""
    if [[ $style == *"light"* ]]; then
        light="-l"
    fi

    # take screenshot
    ./termshot ${light} -o ./styles/gallery/ -f "$filename" glow -s ${style}

    # add border
    convert -bordercolor black -border 16x16 "./styles/gallery/$filename" "./styles/gallery/$filename"

    # optimize filesize
    pngcrush -ow "./styles/gallery/$filename"
done
