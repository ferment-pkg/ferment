v=$(cat VERSION.meta)
echo "This is the ferment-pkg UPDATER"
echo "Checking Current Version..."
echo "Version On System: $v"
echo "Getting Latest Git Pull"
result=$(git pull)
function UpdatePackages(){
    echo "Now Updating Packages..."
    r=$(ls|grep Barrells)
    if [ "$r" = "" ]
    then
        echo "Barrells Not Found"
        echo "Cloning Barrells"
        git clone https://github.com/ferment-pkg/Barrells Barrells
    else
        r=$(cd Barrells)
        result=$(git pull -c Barrells)
        if [ "$result" = "fatal: not a git repository (or any of the parent directories): .git" ]
        then
            echo "Packages Are Not Installed"
            echo "Installing Packages..."
            rm -rf $(ls -a)
            git clone https://github.com/ferment-pkg/Barrells .
        fi
    fi

}

if [ "$result" = "Already up to date." ]; then
    echo "Already Up To Date"
    UpdatePackages
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
ln -sf $PWD/bin/$ARCH/ferment-$ARCH /usr/local/ferment/ferment
UpdatePackages
echo "DONE"
exit 0

