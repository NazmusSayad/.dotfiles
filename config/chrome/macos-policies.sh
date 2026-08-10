#!/bin/bash
# macOS Chrome Enterprise policies
# Equivalent of Chrome Policy.reg for Windows
# Generates com.google.Chrome.mobileconfig from the policy list below and
# installs it as a computer-level configuration profile, so Chrome applies the
# policies as Mandatory (forced) instead of Recommended.

set -euo pipefail

POLICY_DOMAIN="com.google.Chrome"
PROFILE_IDENTIFIER="com.google.Chrome.policies"
PROFILE_PATH="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/com.google.Chrome.mobileconfig"
PROFILE_UUID="64ED36A5-7C62-4701-9CD5-5CA1AAEA2991"
PAYLOAD_UUID="E3559138-8572-465F-96B8-8982D75B9AF2"

POLICY_KEYS=()
POLICY_PAYLOAD=""

apply_bool() {
  POLICY_KEYS+=("$1")
  POLICY_PAYLOAD+="			<key>$1</key>
			<$2/>
"
}

apply_int() {
  POLICY_KEYS+=("$1")
  POLICY_PAYLOAD+="			<key>$1</key>
			<integer>$2</integer>
"
}

echo "Collecting Chrome policies..."

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

# Disable saving passkeys to the password manager
apply_bool "PasswordManagerPasskeysEnabled" false

# Disable passkey creation defaulting to iCloud Keychain (macOS only)
apply_bool "CreatePasskeysInICloudKeychain" false

# Disable password leak detection
apply_bool "PasswordLeakDetectionEnabled" false

# Disable autofill
apply_bool "AutofillAddressEnabled" false
apply_bool "AutofillCreditCardEnabled" false
apply_bool "PaymentMethodQueryEnabled" false

# Enable AI settings (1 = Allow)
apply_int "AIModeSettings" 1
apply_int "GeminiSettings" 1
apply_int "GeminiActOnWebSettings" 1
apply_int "GenAILocalFoundationalModelSettings" 1

echo "Writing $PROFILE_PATH..."
cat >"$PROFILE_PATH" <<EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>PayloadContent</key>
	<array>
		<dict>
			<key>PayloadType</key>
			<string>$POLICY_DOMAIN</string>
			<key>PayloadIdentifier</key>
			<string>$PROFILE_IDENTIFIER.chrome</string>
			<key>PayloadUUID</key>
			<string>$PAYLOAD_UUID</string>
			<key>PayloadVersion</key>
			<integer>1</integer>
			<key>PayloadEnabled</key>
			<true/>
			<key>PayloadDisplayName</key>
			<string>Google Chrome</string>
			<key>PayloadOrganization</key>
			<string>dotfiles</string>
$POLICY_PAYLOAD		</dict>
	</array>
	<key>PayloadType</key>
	<string>Configuration</string>
	<key>PayloadScope</key>
	<string>System</string>
	<key>PayloadIdentifier</key>
	<string>$PROFILE_IDENTIFIER</string>
	<key>PayloadUUID</key>
	<string>$PROFILE_UUID</string>
	<key>PayloadVersion</key>
	<integer>1</integer>
	<key>PayloadDisplayName</key>
	<string>Google Chrome Policies</string>
	<key>PayloadDescription</key>
	<string>Mandatory (forced) Google Chrome policies.</string>
	<key>PayloadOrganization</key>
	<string>dotfiles</string>
	<key>PayloadRemovalDisallowed</key>
	<false/>
</dict>
</plist>
EOF

plutil -lint "$PROFILE_PATH"

if [ "$(id -u)" -ne 0 ]; then
  echo "Profile generated but not installed (needs root)."
  echo "Re-run with sudo to install system-wide:"
  echo "  sudo $0"
  exit 1
fi

if [ -z "${SUDO_USER:-}" ]; then
  echo "Error: run this script through sudo from your own account (SUDO_USER is unset)." >&2
  exit 1
fi

echo "Removing old defaults-based (recommended) Chrome policies..."
for key in "${POLICY_KEYS[@]}"; do
  sudo -u "$SUDO_USER" defaults delete "$POLICY_DOMAIN" "$key" 2>/dev/null || true
  defaults delete "/Library/Preferences/$POLICY_DOMAIN" "$key" 2>/dev/null || true
done

echo "Installing configuration profile at the computer level..."
profiles remove -identifier "$PROFILE_IDENTIFIER" 2>/dev/null || true
if profiles install -type configuration -path "$PROFILE_PATH" 2>/dev/null; then
  echo "Installed with the profiles tool."
else
  # macOS 26 removed profile installs from the profiles tool, so the same
  # forced payload is deployed straight to computer-level managed preferences.
  echo "profiles tool cannot install; deploying the mandatory payload to managed preferences."
  mkdir -p "/Library/Managed Preferences"
  cat >"/Library/Managed Preferences/$POLICY_DOMAIN.plist" <<EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
$POLICY_PAYLOAD</dict>
</plist>
EOF
  plutil -lint "/Library/Managed Preferences/$POLICY_DOMAIN.plist"
  chown root:wheel "/Library/Managed Preferences/$POLICY_DOMAIN.plist"
  chmod 644 "/Library/Managed Preferences/$POLICY_DOMAIN.plist"
  killall cfprefsd 2>/dev/null || true
fi

if [ ! -f "/Library/Managed Preferences/$POLICY_DOMAIN.plist" ]; then
  echo "Error: /Library/Managed Preferences/$POLICY_DOMAIN.plist was not created." >&2
  echo "The profile is not active at the computer level; policies would stay unforced." >&2
  exit 1
fi

if [ -d "/Applications/Google Chrome.app" ]; then
  echo "Restarting Chrome..."
  if pgrep -x "Google Chrome" >/dev/null; then
    sudo -u "$SUDO_USER" osascript -e 'quit app "Google Chrome"'
    sleep 2
  fi
  sudo -u "$SUDO_USER" open -a "Google Chrome"
else
  echo "Google Chrome not found at /Applications/Google Chrome.app; start Chrome manually." >&2
fi

echo "Done."
echo "Verify at chrome://policy - policies must show Source: Platform, Level: Mandatory"
echo "Verify at chrome://settings/payments - controls must be off and locked"
