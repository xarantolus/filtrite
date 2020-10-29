echo "::group::Init"

set -e

log () {
    echo `date +"%m/%d/%Y %H:%M:%S"` "$@"
}
cleanup() {
    rm -f filtrite >> /dev/null 2>&1 &
}

cleanup

echo "::endgroup::"
echo "::group::Build executable"

log "Building"

go build -v -o filtrite

echo "::endgroup::"
echo "::group::Other setup steps"
# Last setup steps
chmod +x filtrite
chmod +x deps/ruleset_converter
mkdir -p dist

echo "::endgroup::"

echo "::group::List: bromite-default"

# Download default bromite filter list
wget -O lists/bromite-default.txt https://raw.githubusercontent.com/bromite/filters/master/lists.txt

log "Start generating bromite-default"
./filtrite lists/bromite-default.txt dist/bromite-default.dat

echo "::endgroup::"

echo "::group::List: bromite-extended"
log "Start generating bromite-extended"
./filtrite lists/bromite-extended.txt dist/bromite-extended.dat

echo "::endgroup::"
echo "::group::Cleanup"

cleanup
echo "::endgroup::"
