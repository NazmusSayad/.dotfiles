local ignoredApps = {
  ["Finder"] = true,
  ["Hammerspoon"] = true
}

local windowFilter = hs.window.filter.new(function(window)
  local subrole = window:subrole()
  return subrole == "AXStandardWindow" or subrole == "AXDialog"
end)

local function quitIfNoWindows(appName, app)
  if ignoredApps[appName] then return end
  if not app or not app:isRunning() then return end

  local pid = app:pid()
  local windows = windowFilter:getWindows()

  for _, window in ipairs(windows) do
    local windowApp = window:application()
    if windowApp and windowApp:pid() == pid then return end
  end

  app:kill()
end

windowFilter:subscribe(hs.window.filter.windowDestroyed, function(window, appName)
  local app = window and window:application() or hs.application.get(appName)
  quitIfNoWindows(appName, app)
end)
