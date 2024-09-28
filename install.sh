cp ./scripts/.bashrc ~/.bashrc -f
echo "Bashrc added..."

cp ./scripts/config.fish ~/.config/fish/config.fish -f
cp ./scripts/fish_theme.fish ~/.config/fish/fish_theme.fish -f
cp ./scripts/fish_aliases.fish ~/.config/fish/fish_aliases.fish -f
cp ./scripts/fish_helpers.fish ~/.config/fish/fish_helpers.fish -f
echo "Fish config added..."

cp ./scripts/default.gitignore ~/default.gitignore -f
echo "Default gitignore added..."

git config --global user.name "Nazmus Sayad"
git config --global user.email "87106526+NazmusSayad@users.noreply.github.com"
git config --global core.autocrlf false
git config --global init.defaultBranch main
git config --global --add safe.directory '*'
git config --global core.excludesfile ~/default.gitignore
git config --global --add --bool push.autoSetupRemote true
echo "Git config added..."

npm config set ignore-scripts true
echo "Npm config added..."
