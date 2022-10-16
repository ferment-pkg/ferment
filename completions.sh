echo "Looking for .zshrc"
content="$(cat ~/.zshrc)"
if [[ $content == *"autoload -U compinit; compinit"* ]]; then
  echo "Completion is already enabled"
else
	echo "autoload -U compinit; compinit" >> ~/.zshrc
	echo "Completion Enabled"
fi
echo "Creating completions directory"
mkdir ~/.completions
cd ~/.completions
ferment completion zsh > _ferment
echo "Run exec zsh for changes to take effect"
