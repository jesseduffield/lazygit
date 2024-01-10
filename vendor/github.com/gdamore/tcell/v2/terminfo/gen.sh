#!/bin/bash
while read line
do
        case "$line" in
        *'|'*)
                alias=${line#*|}
                line=${line%|*}
                ;;
        *)
                alias=${line%%,*}
                ;;
        esac

        alias=${alias//-/_}
        direc=${alias:0:1}

        mkdir -p ${direc}/${alias}
        go run mkinfo.go -P ${alias} -go ${direc}/${alias}/term.go ${line//,/ }
done < models.txt
