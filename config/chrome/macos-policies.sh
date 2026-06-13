#!/bin/bash
# macOS Chrome Enterprise policies
# Equivalent of Chrome Policy.reg for Windows
# Run this script to apply policies. Chrome will read them on next launch.

set -euo pipefail

POLICY_DOMAIN="com.google.Chrome"

apply_bool() {
  defaults write "$POLICY_DOMAIN" "$1" -bool "$2"
}

apply_int() {
  defaults write "$POLICY_DOMAIN" "$1" -int "$2"
}

echo "Applying Chrome policies..."

# Disable background apps
apply_bool "BackgroundModeEnabled" false

# Disable guest mode
apply_bool "BrowserGuestModeEnabled" false

# Disable browser experiments (labs icon)
apply_bool "BrowserLabsEnabled" false

# Disable URL-keyed anonymized data collection
apply_bool "UrlKeyedAnonymizedDataCollectionEnabled" false

# Disable cloud reporting
apply_bool "CloudReportingEnabled" false

# Disable spell check service
apply_bool "SpellCheckServiceEnabled" false

# Disable metrics reporting
apply_bool "MetricsReportingEnabled" false

# Disable promotional tabs
apply_bool "PromotionalTabsEnabled" false

# Disable promotions
apply_bool "PromotionsEnabled" false

# Disable reporting (all require CloudReportingEnabled to take effect)
apply_bool "ReportExtensionsAndPluginsData" false
apply_bool "ReportMachineIDData" false
apply_bool "ReportPolicyData" false
apply_bool "ReportUserIDData" false
apply_bool "ReportVersionData" false

# Disable password manager
apply_bool "PasswordManagerEnabled" false

# Disable password leak detection
apply_bool "PasswordLeakDetectionEnabled" false

# Disable autofill
apply_bool "AutofillAddressEnabled" false
apply_bool "AutofillCreditCardEnabled" false

# Enable AI settings (1 = Allow)
apply_int "AIModeSettings" 1
apply_int "GeminiSettings" 1
apply_int "GeminiActOnWebSettings" 1
apply_int "GenAILocalFoundationalModelSettings" 1

echo "Done. Restart Chrome for policies to take effect."
echo "Verify at chrome://policy"
