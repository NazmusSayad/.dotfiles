local ignoredApps = {
  ["Finder"] = true,
  ["Hammerspoon"] = true
}

local logPath = hs.configdir .. "/ghost-apps.log"

local function log(message)
  local file = io.open(logPath, "a")
  if not file then return end

  file:write(os.date("%Y-%m-%d %H:%M:%S"), " ", message, "\n")
  file:close()
end

local file = io.open(logPath, "w")
if file then file:close() end
log("configuration loaded")

local windowFilter = hs.window.filter.new(function(window)
  local subrole = window:subrole()
  return subrole == "AXStandardWindow" or subrole == "AXDialog"
end)

local function quitIfNoWindows(appName, app)
  if ignoredApps[appName] then
    log(string.format("skip ignored app=%q", appName))
    return
  end

  if not app or not app:isRunning() then
    log(string.format("skip not-running app=%q", appName))
    return
  end

  local pid = app:pid()
  local windows = windowFilter:getWindows()
  log(string.format("check app=%q pid=%d all-space-windows=%d", appName, pid, #windows))

  for _, window in ipairs(windows) do
    local windowApp = window:application()
    if windowApp and windowApp:pid() == pid then
      log(string.format(
        "skip has-window app=%q id=%s title=%q",
        appName,
        tostring(window:id()),
        window:title() or ""
      ))
      return
    end
  end

  log(string.format("quit app=%q", appName))
  app:kill()
end

windowFilter:subscribe(hs.window.filter.windowDestroyed, function(window, appName)
  local app = window and window:application() or hs.application.get(appName)
  log(string.format("window-destroyed app=%q id=%s", appName, tostring(window and window:id())))
  quitIfNoWindows(appName, app)
end)
