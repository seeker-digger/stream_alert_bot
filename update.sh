#!/bin/bash

red='\033[0;31m'
green='\033[0;32m'
plain='\033[0m'

# --- Root check ---
if [[ "$EUID" -ne 0 ]]; then
  echo -e "${red}Fatal error: ${plain}Please, run this script as root!"
  exit 1
fi

# --- Detect OS ---
if [[ -f /etc/os-release ]]; then
    source /etc/os-release
    release=$ID
elif [[ -f /usr/lib/os-release ]]; then
    source /usr/lib/os-release
    release=$ID
else
    echo "Failed to detect OS!" >&2
    exit 1
fi
echo "The OS release is: $release"

# --- Detect arch ---
arch() {
    case "$(uname -m)" in
    x86_64 | x64 | amd64) echo 'amd64' ;;
    *) echo -e "${red}Unsupported CPU architecture!${plain}" && exit 1 ;;
    esac
}
echo "Arch: $(arch)"

# --- Ensure required tools ---
install_base() {
    case "${release}" in
    ubuntu | debian | armbian)
        apt-get update -y && apt-get install -y wget curl jq
        ;;
    centos | rhel | almalinux | rocky | ol)
        yum install -y wget curl jq
        ;;
    fedora | amzn | virtuozzo)
        dnf install -y wget curl jq
        ;;
    arch | manjaro | parch)
        pacman -Syu --noconfirm wget curl jq
        ;;
    opensuse-tumbleweed | opensuse-leap)
        zypper install -y wget curl jq
        ;;
    alpine)
        apk add wget curl jq
        ;;
    *)
        apt-get install -y wget curl jq
        ;;
    esac
}

# --- Update bot ---
update_app() {
    local install_dir="/usr/local/stream-alert-bot"
    local binary_path="${install_dir}/bot"

    if [[ ! -d "$install_dir" || ! -f "$binary_path" ]]; then
        echo -e "${red}Error:${plain} Stream Alert Bot is not installed!"
        echo "Please run the install script first."
        exit 1
    fi

    echo "Checking for updates..."

    # Get latest release info
    release_json=$(curl -s -H "Accept: application/vnd.github+json" \
                        "https://api.github.com/repos/seeker-digger/stream_alert_bot/releases/latest")

    echo $release_json

    if [[ -z "$release_json" ]]; then
        echo -e "${red}Fatal error:${plain} Failed to fetch release info from GitHub!"
        exit 1
    fi

    latest_version=$(echo "$release_json" | jq -r '.tag_name')
    asset_url=$(echo "$release_json" | jq -r '.assets[0].browser_download_url')
    echo $latest_version
    echo $asset_url
    if [[ -z "$latest_version" || -z "$asset_url" ]]; then
        echo -e "${red}Fatal error:${plain} Could not parse release info!"
        exit 1
    fi

    echo "Latest version:  $latest_version"

    echo "Updating Stream Alert Bot to ${latest_version}..."

    cd "$install_dir" || exit 1

    # Stop service
    if [[ $release == "alpine" ]]; then
        rc-service stream-alert-bot stop
    else
        systemctl stop stream-alert-bot
    fi

    # Download new binary
    wget -O bot.new "$asset_url"
    if [[ $? -ne 0 ]]; then
        echo -e "${red}Fatal error:${plain} Failed to download new release!"
        exit 1
    fi

    chmod +x bot.new
    mv bot.new bot

    # Restart service
    if [[ $release == "alpine" ]]; then
        rc-service stream-alert-bot start
    else
        systemctl start stream-alert-bot
    fi

    echo -e "${green}Stream Alert Bot updated to version ${latest_version}!${plain}"
}

install_base
update_app
