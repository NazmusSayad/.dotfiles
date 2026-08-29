local ignoredApps = {
  ["Finder"] = true,
  ["Hammerspoon"] = true,
  ["Google Chrome"] = true,
  ["Code"] = true,
  ["Code - Insiders"] = true,
  ["OiPer"] = true
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

local fnFPressed = false
fnFTap = hs.eventtap.new({
  hs.eventtap.event.types.keyDown,
  hs.eventtap.event.types.keyUp,
  hs.eventtap.event.types.tapDisabledByTimeout,
  hs.eventtap.event.types.tapDisabledByUserInput
}, function(event)
  local eventType = event:getType()

  if eventType == hs.eventtap.event.types.tapDisabledByTimeout
    or eventType == hs.eventtap.event.types.tapDisabledByUserInput then
    fnFPressed = false
    fnFTap:start()
    return false
  end

  if event:getKeyCode() ~= hs.keycodes.map.f then return false end

  if eventType == hs.eventtap.event.types.keyDown and event:getFlags().fn then
    if not fnFPressed then
      fnFPressed = true
      hs.eventtap.keyStroke({ "ctrl", "cmd" }, "f", 0)
    end
    return true
  end

  if eventType == hs.eventtap.event.types.keyUp and fnFPressed then
    fnFPressed = false
    return true
  end

  return false
end)

fnFTap:start()

hs.hotkey.bind({ "alt" }, "space", function()
  hs.eventtap.keyStroke({ "alt", "shift" }, "f", 0)
end)
