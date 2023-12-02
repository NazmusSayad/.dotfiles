_____SHELL_PATH() {
  local dir="$PWD"
  if [[ "$OSTYPE" == "cygwin" || "$OSTYPE" == "msys" ]]; then
    local win_dir=$dir
    local first_part="${win_dir:1:1}"
    local first_part_uppercase="${first_part^^}"
    local second_part="${win_dir:3}"
    dir="${first_part_uppercase}: ${second_part}"
  fi
  echo $dir
}

_____SHELL_PROMPT_COMMAND() {
  local resetcolor='\e[0m'
  local magenta='\e[35m'
  local lightblue='\[\033[1;34m\]'
  local gray='\e[0m'
  
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

  PS1="${magenta}⌠ ${lightblue}$(_____SHELL_PATH) ${gitbranch}\n${magenta}$ ${resetcolor}"
}

PROMPT_COMMAND="_____SHELL_PROMPT_COMMAND"
[[ -z $REMOTE_GITHUB_TOKEN ]] || export GITHUB_TOKEN=$REMOTE_GITHUB_TOKEN

alias ps="powershell"
alias ni="touch"
alias md="mkdir"
alias cls="clear"

alias n="node"
alias w="nodemon"

alias x="npx"
alias r="npm run"
alias i="npm install"
alias u="npm uninstall"
alias di="npm install --save-dev"
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
