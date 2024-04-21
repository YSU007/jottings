#!/bin/sh

CURR_DIR=$(pwd)
IN_PATH=$CURR_DIR
OUT_PATH=$CURR_DIR/gen
PROTO_PATH=$CURR_DIR/../../proto

gen() {
    if [ -d "$OUT_PATH"  ];then
      rm -fr "$OUT_PATH"
    fi
    mkdir "$OUT_PATH"

    python3.9 proto_go.py "$IN_PATH" "$OUT_PATH"
    protoc --proto_path="$IN_PATH" --go_out="$OUT_PATH" "$IN_PATH"/*.proto
    for f in $OUT_PATH/*;do
      go fmt "$f"
    done
}

clean() {
  rm -vfr "$OUT_PATH"
}

copy() {
  if [ -d "$OUT_PATH"  ];then
        files=$(ls "$OUT_PATH"|grep -v handle.go)
        cnt=$(echo "$files"|wc -l|sed s/[[:space:]]//g)
        if [  "$cnt" -eq 0  ]; then
          echo "no files"
          return
        fi
        for f in $files ; do
          mv -v "$OUT_PATH"/"$f" "$PROTO_PATH"
        done
  fi
}

if [ "$1" = "gen" ]||[ "$1" = "" ]; then
    gen
fi

if [ "$1" = "clean" ]; then
    clean
fi

if [ "$1" = "copy" ]; then
    copy
fi

if [ "$1" = "-h" ]||[ "$1" = "--help" ]; then
    echo "Parameter:"
    echo "    'gen' or empty gen proto."
    echo "    'clean' clean gen files."
    echo "    'copy' mv proto files to project."
fi