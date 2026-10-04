#!/bin/bash

source ./etc/git-config.sh
source ./etc/shell-config.sh

echo "Configuring shell settings..."
bash_path=$(which bash)
if ! grep -qxF "$bash_path" /etc/shells; then
	echo "$bash_path" | sudo tee -a /etc/shells >/dev/null
fi
chsh -s "$bash_path"

echo "Configuring macOS settings..."
sudo launchctl disable system/com.apple.assistantd
sudo spctl --master-disable
sudo mdutil -a -i off
sudo mdutil -a -E

echo "Configuring Finder and Dock settings..."
chflags nohidden ~/Library

echo "Configuring power settings..."
sudo pmset -a sleep 0
sudo pmset -a powernap 0
sudo pmset -a ttyskeepawake 0
sudo pmset -a womp 0
sudo pmset -a disksleep 0
sudo pmset -a displaysleep 5
sudo pmset -b lessbright 0

echo "Configuring battery settings..."
sudo batt limit 70
sudo batt lower-limit-delta 10
sudo batt prevent-system-sleep enable
sudo batt magsafe-led enable

echo "Disabling Software Update notifications..."
defaults write com.apple.SoftwareUpdate AutomaticCheckEnabled -bool false
sudo defaults write /Library/Preferences/com.apple.SoftwareUpdate AutomaticCheckEnabled -bool false
sudo defaults write /Library/Preferences/com.apple.SoftwareUpdate AutomaticDownload -bool false
sudo defaults write /Library/Preferences/com.apple.SoftwareUpdate ConfigDataInstall -bool false
sudo defaults write /Library/Preferences/com.apple.SoftwareUpdate CriticalUpdateInstall -bool false
sudo defaults write /Library/Preferences/com.apple.SoftwareUpdate AutomaticallyInstallMacOSUpdates -bool false
sudo softwareupdate --schedule off

defaults write NSGlobalDomain NSQuitAlwaysKeepsWindows -bool false
defaults write com.apple.LaunchServices LSQuarantine -bool false

defaults write NSGlobalDomain AppleSymbolicHotKeysEnabled -bool false
defaults write NSGlobalDomain NSDocumentSaveNewDocumentsToCloud -bool false
defaults write com.apple.applicationaccess allowPasswordAutoFill -bool false
defaults write com.apple.applicationaccess safariAllowAutoFill -bool false
defaults write com.apple.Safari AutoFillPasswords -bool false

sudo defaults write /Library/Preferences/com.apple.iokit.AmbientLightSensor "Automatic Display Enabled" -bool false

disabled_symbolic_hotkeys=(
	7 8 9 10 11 12 13
	15 16 17 18 19 20 21 22 23 24 25 26
	27 28 29 30 31 32 33 34 35 36 37
	52 53 54 55 56 57 59 60 61 64 65 80 82
	118 119 120 121 122 123 124 125 126 127 128 129 130 131 132 133
	159 160 162 163 175 190
	215 216 217 218 219 222
	225 226 227 228 229 230 231 232
	260
)
for hotkey in "${disabled_symbolic_hotkeys[@]}"; do
	defaults write com.apple.symbolichotkeys AppleSymbolicHotKeys -dict-add "$hotkey" '<dict><key>enabled</key><false/></dict>'
done

disabled_text_services=(
	"com.apple.services.addToiTunesAsSpokenTrack - Add to Music as a Spoken Track - runWorkflowAsService"
	"com.apple.ChineseTextConverterService - Convert Text to Full Width - convertTextToFullWidth"
	"com.apple.ChineseTextConverterService - Convert Text to Half Width - convertTextToHalfWidth"
	"com.apple.ChineseTextConverterService - Convert Text from Traditional to Simplified Chinese - convertTextToSimplifiedChinese"
	"com.apple.ChineseTextConverterService - Convert Text from Simplified to Traditional Chinese - convertTextToTraditionalChinese"
	"com.apple.Stickies - Make Sticky - makeStickyFromTextService"
	"com.apple.TextEdit - New TextEdit Window Containing Selection - openSelection"
	"com.apple.finder - Finder/Open - open"
	"com.apple.finder - Finder/Reveal - reveal"
	"com.apple.finder - Finder/Show Info - showInfo"
	"com.apple.Terminal - Open man Page in Terminal - openManPage"
	"com.apple.Terminal - Search man Page Index in Terminal - searchManPages"
	"org.hammerspoon.Hammerspoon - Send to Hammerspoon - processDockIconDraggedText"
	"org.localsend.localsendApp - Send to LocalSend - handleSendTextService"
	"com.apple.services.showMap - Show Map - runWorkflowAsService"
)
for service in "${disabled_text_services[@]}"; do
	defaults write pbs NSServicesStatus -dict-add "$service" '<dict><key>presentation_modes</key><dict><key>ContextMenu</key><false/><key>ServicesMenu</key><false/></dict><key>enabled_context_menu</key><false/><key>enabled_services_menu</key><false/><key>key_equivalent</key><string></string></dict>'
done

defaults write com.apple.dock tilesize -int 64
defaults write com.apple.dock largesize -int 81
defaults write com.apple.dock magnification -bool true

defaults write com.apple.dock autohide -bool true
defaults write com.apple.dock show-recents -bool false
defaults write com.apple.dock mineffect -string scale

defaults write com.apple.dock wvous-tr-corner -int 3
defaults write com.apple.dock wvous-tr-modifier -int 0

defaults write com.apple.dock autohide-delay -float 0
defaults write com.apple.dock autohide-time-modifier -float 0.75
defaults write com.apple.dock expose-animation-duration -float 0.01
defaults write com.apple.dock workspaces-edge-delay -float 0
defaults write NSGlobalDomain NSWindowResizeTime -float 0.001

defaults write com.apple.dock mru-spaces -bool false
defaults write com.apple.spaces spans-displays -bool false

defaults write com.apple.WindowManager StandardHideWidgets -bool true
defaults write com.apple.WindowManager StageManagerHideWidgets -bool true

defaults write com.apple.finder ShowRecentTags -bool false
defaults write com.apple.finder WarnOnEmptyTrash -bool false
defaults write com.apple.finder AppleShowAllFiles -bool true
defaults write com.apple.finder NewWindowTarget -string PfHm
defaults write com.apple.finder FXPreferredViewStyle -string icnv
defaults write com.apple.finder FXDisableAllAnimations -bool true
defaults write com.apple.finder NewWindowTargetPath -string "file://${HOME}/Desktop"
defaults write com.apple.finder _FXSortFoldersFirst -bool true
defaults write com.apple.finder _FXShowPosixPathInTitle -bool true
defaults write com.apple.finder _FXSortFoldersFirstOnDesktop -bool true
defaults write com.apple.finder FXDefaultSearchScope -string SCcf
defaults write com.apple.finder FXEnableExtensionChangeWarning -bool false
defaults write NSGlobalDomain AppleShowAllExtensions -bool true

defaults write com.apple.AppleMultitouchTrackpad Clicking -bool true
defaults write com.apple.AppleMultitouchTrackpad FirstClickThreshold -int 0
defaults write com.apple.AppleMultitouchTrackpad SecondClickThreshold -int 0
defaults write NSGlobalDomain com.apple.mouse.tapBehavior -int 1
defaults write NSGlobalDomain com.apple.trackpad.scaling -float 3.0
defaults write NSGlobalDomain AppleEnableSwipeNavigateWithScrolls -bool false
defaults write com.apple.AppleMultitouchTrackpad TrackpadTwoFingerFromRightEdgeSwipeGesture -int 0
defaults write com.apple.driver.AppleBluetoothMultitouch.trackpad TrackpadTwoFingerFromRightEdgeSwipeGesture -int 0
defaults write com.apple.AppleMultitouchTrackpad TrackpadFiveFingerPinchGesture -int 0
defaults write com.apple.driver.AppleBluetoothMultitouch.trackpad TrackpadFiveFingerPinchGesture -int 0
defaults write com.apple.AppleMultitouchTrackpad TrackpadFourFingerPinchGesture -int 0
defaults write com.apple.driver.AppleBluetoothMultitouch.trackpad TrackpadFourFingerPinchGesture -int 0

defaults write com.apple.desktopservices DSDontWriteUSBStores -bool true
defaults write com.apple.desktopservices DSDontWriteNetworkStores -bool true

defaults write NSGlobalDomain AppleLocale -string "en_US"
defaults write NSGlobalDomain AppleLanguages -array "en-US"
defaults write NSGlobalDomain AppleMeasurementUnits -string "Centimeters"
defaults write NSGlobalDomain AppleMetricUnits -bool true

defaults write NSGlobalDomain NSAutomaticSpellingCorrectionEnabled -bool false
defaults write NSGlobalDomain NSAutomaticInlinePredictionEnabled -bool false
defaults write NSGlobalDomain NSAutomaticPeriodSubstitutionEnabled -bool false
defaults write NSGlobalDomain NSAutomaticQuoteSubstitutionEnabled -bool false
defaults write NSGlobalDomain NSAutomaticDashSubstitutionEnabled -bool false
defaults write NSGlobalDomain NSAutomaticCapitalizationEnabled -bool false

defaults write NSGlobalDomain CGDisableCursorLocationMagnification -bool true
defaults write NSGlobalDomain NSScrollViewRubberbanding -bool false
defaults write NSScrollViewRubberbanding -int 0

defaults write NSGlobalDomain KeyRepeat -int 2
defaults write NSGlobalDomain InitialKeyRepeat -int 25
defaults write NSGlobalDomain ApplePressAndHoldEnabled -bool false
defaults write NSGlobalDomain com.apple.keyboard.fnState -bool true

defaults write NSGlobalDomain AppleLiveTextEnabled -bool false
defaults write NSGlobalDomain NSQuitAlwaysKeepsWindows -bool false
defaults write com.apple.loginwindow TALLogoutSavesState -bool false

defaults write com.apple.Safari WebKitDeveloperExtrasEnabledPreferenceKey -bool true

launchctl disable "gui/$(id -u)/com.google.GoogleUpdater.wake"

echo "Opening Chrome configuration profile..."
open ./config/chrome/com.google.Chrome.mobileconfig
defaults write com.google.Chrome BackgroundModeEnabled -bool false

killall Dock
killall Finder
killall SystemUIServer
