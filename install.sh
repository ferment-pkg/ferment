ARCH=$(uname -m)
OS=$(uname -s)
echo "THIS SCRIPT EXPECTS YOU TO FOLLOW THE DOCUMENTATION IN GITHUB TO WORK PROPERLY"
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
mkdir -p /usr/local/bin
ln -sf bin/$ARCH/ferment-$ARCH ferment
zshrcOut=$(cat ~/.zshrc|grep /usr/local/ferment)
if [ "$zshrcOut" = "" ]; then
  echo "Adding ferment to your zshrc"
  echo export PATH='$PATH':$PWD >> $HOME/.zshrc
fi
echo "Updated Path in .zshrc"
echo "Run source ~/.zshrc to update PATH"
echo "Install Completed"
exit 0
