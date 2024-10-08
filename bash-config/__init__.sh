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

alias n="node --no-warnings"
alias nw="node --watch --no-warnings"

alias x="npm exec"
alias r="npm run"
alias rp="run-p --silent"
alias rs="run-s --silent"

alias ni="npm install"
alias nid="npm install --save-dev"
alias nu="npm uninstall"

alias gc="gh repo clone"

gp() {
  local msg=${*:-$(git status --short --no-renames)}

  git status --short
  echo

  git add -A >>/dev/null &&
    git commit -m "$msg" &&
    echo &&
    git push --quiet
}

gr() {
  local hash="$1"
  if [ -z "$hash" ]; then
    echo "You must need to give a <commit_id>"
    return
  fi

  git checkout $hash . &&
    gp "$hash restored; ${*:2}"
}

ghr() {
  local branch=${1:-"master"}

  git checkout --orphan latest_branch
  git add -A
  git commit -am "initial commit"
  git branch -D $branch
  git branch -m $branch
  git push -f origin $branch
}

gpull() {
  local src="$1"

  if [ -z "$src" ]; then
    echo "Error: Please provide the source branch."
    return 1
  fi

  if [ ! -d ".git" ]; then
    echo "Error: Not a Git repository."
    return 1
  fi

  local target=$(git rev-parse --abbrev-ref HEAD)
  git pull origin "$src:$target"
  return $?
}
