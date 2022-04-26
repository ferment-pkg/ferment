ARCH=$(uname -m)
OS=$(uname -s)
if [ "$ARCH" = "x86_64" ]; then
  ARCH="amd64"
fi
if [ "$ARCH" = "aarch64" ]; then
  ARCH="arm64"
fi
if [ "$OS" = "Linux" ]; then
  echo "Linux Is Not Supported As The Default Package Manager Is Almost Always The Best Option"
  exit 1
fi
echo "Identified architecture As: $ARCH"
echo "Identified operating system As: $OS"
echo "Checking If Dependencies Are Installed"
PYTHONEXE=$(which python3)
if [ "$PYTHONEXE" = "" ]; then
  echo "Python3 Is Not Installed And Python2 and below are not supported"
  exit 1
fi
echo "Python3 Is Installed"
echo "Adding Project To PATH"
cp bin/$ARCH/ferment .
echo export PATH='$PATH':$PWD >> $HOME/.zshrc
echo "Updated Path in .zshrc"
echo "Install Completed"
exit 0
