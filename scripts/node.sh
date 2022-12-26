curl https://raw.githubusercontent.com/creationix/nvm/master/install.sh | bash
rm ~/.nodejs /opt/nodejs -rf
nvm install node --use
nvm alias default node

cat ~/temp-custom-scripts/main.sh > ~/.bashrc
rm ~/temp-custom-scripts -rf

npm i -g npm@latest live-server gitignore

clear
