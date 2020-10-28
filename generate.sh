set -e

log () {
    echo `date +"%m/%d/%Y %H:%M:%S"` "$@"
}
cleanup() {
    rm -f filtrite >> /dev/null 2>&1 &
}

cleanup

log "Building"

go build -v -o filtrite


# Last setup steps
chmod +x filtrite
chmod +x deps/ruleset_converter
mkdir -p dist

log "Start generating bromite-default"
./filtrite lists/bromite-default.txt dist/bromite-default.dat

log "Start generating bromite-extended"
./filtrite lists/bromite-extended.txt dist/bromite-extended.dat

cleanup
