if test (git log --branches --not --remotes | wc -l) -gt 0
  echo "Has changes"
end
