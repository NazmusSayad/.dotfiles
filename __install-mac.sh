#!/bin/bash

sudo spctl --master-disable
sudo mdutil -a -i off
sudo mdutil -a -E

sudo launchctl disable system/com.apple.assistantd

defaults write com.apple.dock tilesize -int 64
defaults write com.apple.dock autohide -bool true
defaults write com.apple.dock magnification -bool false
defaults write com.apple.dock autohide-delay -float 0
defaults write com.apple.dock autohide-time-modifier -float 0
defaults write com.apple.dock show-recents -bool false

defaults write com.apple.finder QuitMenuItem -bool true
defaults write com.apple.finder CreateDesktop -bool false

defaults write com.apple.finder ShowPathbar -bool true
defaults write com.apple.finder ShowStatusBar -bool true
defaults write com.apple.finder ShowRecentTags -bool false
defaults write com.apple.finder _FXSortFoldersFirst -bool true
defaults write com.apple.finder FXEnableExtensionChangeWarning -bool false
defaults write NSGlobalDomain AppleShowAllExtensions -bool true

defaults write com.apple.AppleMultitouchTrackpad Clicking -bool true
defaults write NSGlobalDomain com.apple.mouse.tapBehavior -int 1

defaults write com.apple.LaunchServices LSQuarantine -bool false

defaults write com.apple.desktopservices DSDontWriteUSBStores -bool true
defaults write com.apple.desktopservices DSDontWriteNetworkStores -bool true

defaults write NSGlobalDomain NSAutomaticSpellingCorrectionEnabled -bool false
defaults write NSGlobalDomain NSAutomaticInlinePredictionEnabled -bool false
defaults write NSGlobalDomain NSAutomaticPeriodSubstitutionEnabled -bool false
defaults write NSGlobalDomain NSAutomaticQuoteSubstitutionEnabled -bool false
defaults write NSGlobalDomain NSAutomaticDashSubstitutionEnabled -bool false
defaults write NSGlobalDomain NSAutomaticCapitalizationEnabled -bool false

defaults write NSGlobalDomain AppleLiveTextEnabled -bool false
defaults write NSGlobalDomain NSQuitAlwaysKeepsWindows -bool false
