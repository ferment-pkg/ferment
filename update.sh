v=$(cat VERSION.meta)
echo "This is the ferment-pkg UPDATER"
echo "Checking Current Version..."
echo "Version On System: $v"
echo "Getting Latest Git Pull"
result=$(git pull)

if [ "$result" = "Already up to date." ]; then
    echo "Already Up To Date"
    exit 0
fi
echo "Updating..."
v=$(cat VERSION.meta)
echo "NEW VERSION: $v"
ARCH=$(uname -m)
if [ "$ARCH" = "x86_64" ]; then
    ARCH="amd64"
else 
    ARCH="arm64"
fi
ln -sf $PWD/bin/$ARCH/ferment-$ARCH ferment
echo "DONE"
exit 0

