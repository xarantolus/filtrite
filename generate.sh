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

log "Start generating bromite-extended"

./filtrite lists/bromite-extended.txt dist/bromite-extended.dat

cleanup
