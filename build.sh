echo "This Is A Tool To Build The Project"
echo "Do you want to rebuild the project?(y/n)"
read answer
if [ $answer != "y" ];then
     echo "Abort!"
     exit 1
fi
echo "Checking for GO"
GOEXEC=$(which go)
if [ "$GOEXEC" = "" ]; then
    echo "GO Is Not Installed"
    exit 1
fi
echo "GO Found"
echo "Checking for Python3"
PYTHONEXE=$(which python3)
if [ "$PYTHONEXE" = "" ]; then
    echo "Python3 Is Not Installed"
    exit 1
fi
echo "Python3 Found"
echo "Building For amd64..."
GOARCH="amd64" go build -o bin/amd64/ferment-amd64
echo "Build Completed"
echo "Building For arm64..."
GOARCH="arm64" go build -o bin/arm64/ferment-arm64
echo "Done"
exit 0