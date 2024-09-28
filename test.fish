if git log --branches --not --remotes ":/*" >/dev/null 2>&1
  echo "Has changes"
end