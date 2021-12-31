echo "::group::Init"

set -e

log () {
    echo `date +"%m/%d/%Y %H:%M:%S"` "$@"
}

cleanup() {
    rm -f filtrite >> /dev/null 2>&1
}
cleanup

# Make sure all dependencies are installed
sudo apt-get install unzip wget || true

echo "::endgroup::"

echo "::group::Build executable"
log "Building"
go build -v -o filtrite
echo "::endgroup::"

echo "::group::Downloading latest subresource_filter_tools build"
wget -O "subresource_filter_tools_linux.zip" "https://github.com/xarantolus/subresource_filter_tools/releases/latest/download/subresource_filter_tools_linux-x64.zip"

mkdir -p deps
unzip -u "subresource_filter_tools_linux.zip" -d deps

rm "subresource_filter_tools_linux.zip"
echo "::endgroup::"


echo "::group::Other setup steps"
chmod +x filtrite
chmod +x deps/ruleset_converter
mkdir -p dist
mkdir -p logs
echo "::endgroup::"

# If the default list file exists, we overwrite it with the actual official list
if [[ -f "lists/bromite-default.txt" ]]; then
    echo "::group::Downloading official list"
    wget -O "lists/bromite-default.txt" "https://raw.githubusercontent.com/bromite/filters/master/lists.txt"
    echo "::endgroup::"
fi

# Now that everything is set up, we can start actually generating filter lists
./filtrite 

echo "::group::Cleanup"
cleanup

# Reset the downloaded list to the previous text, in case this is run locally
if [[ -f "lists/bromite-default.txt" ]]; then
    git restore lists/bromite-default.txt
fi

echo "::endgroup::"
