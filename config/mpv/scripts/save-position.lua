local current_path = nil

mp.register_event("file-loaded", function()
    current_path = mp.get_property("path")
end)

mp.add_periodic_timer(5, function()
    if current_path then
        mp.command("write-watch-later-config")
    end
end)

mp.add_hook("on_unload", 50, function()
    mp.command("write-watch-later-config")
end)

mp.register_event("end-file", function(event)
    if event.reason == "eof" and current_path then
        mp.commandv("delete-watch-later-config", current_path)
    end
    current_path = nil
end)
