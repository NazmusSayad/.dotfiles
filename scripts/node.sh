cat ~/temp-custom-scripts/main.sh > ~/.bashrc
rm ~/temp-custom-scripts -rf

rm ~/.nodejs /opt/nodejs -rf
curl https://raw.githubusercontent.com/creationix/nvm/master/install.sh | bash
nvm install node --use
nvm alias default node
npm i -g npm@latest live-server gitignore

clear
