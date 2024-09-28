if not git diff --quiet; or not git diff --cached --quiet
  echo "Has uncommitted changes"
else
  echo "No uncommitted changes"
end
