__bash_prompt() {
  local gitbranch='`\
      if [ "$(git config --get codespaces-theme.hide-status 2>/dev/null)" != 1 ]; then \
          export BRANCH=$(git symbolic-ref --short HEAD 2>/dev/null || git rev-parse --short HEAD 2>/dev/null); \
          if [ "${BRANCH}" != "" ]; then \
              echo -n "\[\033[0;36m\](\[\033[1;31m\]${BRANCH}" \
              && if git ls-files --error-unmatch -m --directory --no-empty-directory -o --exclude-standard ":/*" > /dev/null 2>&1; then \
                      echo -n " \[\033[1;33m\]✗"; \
              fi \
              && echo -n "\[\033[0;36m\]) "; \
          fi; \
      fi`'
  local magenta='\e[35m'
  local lightblue='\[\033[1;34m\]'
  local resetcolor='\e[0m'
  PS1="${magenta}➜ ${lightblue}\w ${gitbranch}\n${magenta}$ ${resetcolor}"
  unset -f __bash_prompt
}
__bash_prompt
[[ -z $REMOTE_GITHUB_TOKEN ]] || export GITHUB_TOKEN=$REMOTE_GITHUB_TOKEN

alias y="yarn"
alias ya="yarn add"
alias yad="yarn add -D"
alias ni="touch"
alias md="mkdir"
alias cls="clear"

alias n="node"
alias nw="node --watch"
alias nd="node --inspect"
alias ndw="node --inspect --watch"
alias tn="ts-node"
alias tnw="ts-node-dev"
alias nm="nodemon"
alias live="live-server"

ghp() {
  local msg=${*:-`git status --short --no-renames`}
  
  git status --short
  echo

  git add -A >> /dev/null &&
  git commit -m "$msg" &&
  echo &&
  git push --quiet
}

ghr() {
  local hash="$1"
  if [ -z "$hash" ]; then
    echo "You must need to give a <commit_id>"
    return
  fi
  
  git checkout $hash . &&
  ghp "$hash restored; ${*:2}"
}

ghhr() {
  local branch=${1:-"master"}
 
  git checkout --orphan latest_branch;
  git add -A;
  git commit -am "initial commit";
  git branch -D $branch;
  git branch -m $branch;
  git push -f origin $branch;
}

ghc() {
  if [[ $1 == http* ]]; then
    git clone $*
  else
    gh repo clone $*
  fi
}
