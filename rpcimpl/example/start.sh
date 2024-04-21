#!/usr/bin/env sh

alias terminal='/Applications/iTerm.app/Contents/MacOS/iTerm2'

curr_dir=$(pwd)
server_dir="$curr_dir/server"
client_dir="$curr_dir/client"

ls "$server_dir"/* | grep -v '\.go$' | xargs rm -vfr
ls "$client_dir"/* | grep -v '\.go$' | xargs rm -vfr
