#!/bin/bash

red='\033[0;31m'
green='\033[0;32m'
plain='\033[0m'

# Check root
if [[ "$EUID" -ne 0 ]]; then
  echo -e "${red}Fatal error: ${plain}Please, run this script as root!"
  exit 1
fi

# Detect OS
if [[ -f /etc/os-release ]]; then
    source /etc/os-release
    release=$ID
elif [[ -f /usr/lib/os-release ]]; then
    source /usr/lib/os-release
    release=$ID
else
    echo "Failed to check the system OS, please contact the author!" >&2
    exit 1
fi
echo "The OS release is: $release"

# Detect architecture
arch() {
    case "$(uname -m)" in
    x86_64 | x64 | amd64) echo 'amd64' ;;
    *) echo -e "${red}Unsupported CPU architecture! ${plain}" && exit 1 ;;
    esac
}
echo "Arch: $(arch)"

# Install required tools
install_base() {
    case "${release}" in
    ubuntu | debian | armbian)
        apt-get update && apt-get install -y wget curl tar jq
        ;;
    centos | rhel | almalinux | rocky | ol)
        yum -y update && yum install -y wget curl tar jq
        ;;
    fedora | amzn | virtuozzo)
        dnf -y update && dnf install -y wget curl tar jq
        ;;
    arch | manjaro | parch)
        pacman -Syu --noconfirm wget curl tar jq
        ;;
    opensuse-tumbleweed | opensuse-leap)
        zypper refresh && zypper install -y wget curl tar jq
        ;;
    alpine)
        apk update && apk add wget curl tar jq
        ;;
    *)
        apt-get update && apt-get install -y wget curl tar jq
        ;;
    esac
}

# Install the bot
install_app() {
    cd /usr/local || exit 1

    # Get latest release (including pre-release)
    release_json=$(curl -s -H "Accept: application/vnd.github+json" \
                        "https://api.github.com/repos/seeker-digger/stream_alert_bot/releases/latest")

    if [[ -z "$release_json" ]]; then
        echo "${red}Fatal error:${plain} Failed to fetch releases from GitHub"
        exit 1
    fi

    tag_version=$(echo "$release_json" | jq -r '.tag_name')
    asset_url=$(echo "$release_json" | jq -r '.assets[0].browser_download_url')

    if [[ -z "$tag_version" || -z "$asset_url" ]]; then
        echo "${red}Fatal error:${plain} Could not determine latest release or asset"
        exit 1
    fi

    echo "Got latest version: ${tag_version}, beginning installation..."

    # Stop existing service
    if [[ -d /usr/local/stream-alert-bot/ ]]; then
        if [[ $release == "alpine" ]]; then
            rc-service stream-alert-bot stop
        else
            systemctl stop stream-alert-bot
        fi
        rm -rf /usr/local/stream-alert-bot/bot
    fi

    mkdir -p /usr/local/stream-alert-bot/
    cd /usr/local/stream-alert-bot/ || exit 1

    # Download binary
    wget -O bot "$asset_url"
    if [[ $? -ne 0 ]]; then
        echo "${red}Fatal error:${plain} Failed to download release asset!"
        exit 1
    fi

    # Download service file
    curl -s -L -o stream-alert-bot.service \
         "https://raw.githubusercontent.com/seeker-digger/stream_alert_bot/master/stream-alert-bot.service"
    if [[ $? -ne 0 ]]; then
        echo "${red}Fatal error:${plain} Failed to download service file!"
        rm -f bot
        exit 1
    fi

    chmod +x bot

    # Install service
    if [[ $release == "alpine" ]]; then
        cp stream-alert-bot.service /etc/init.d/stream_alert-bot
        chmod +x /etc/init.d/stream_alert-bot
        rc-update add stream_alert-bot default
        rc-service stream_alert-bot start
    else
        cp stream-alert-bot.service /etc/systemd/system/stream-alert-bot.service
        systemctl daemon-reload
        systemctl enable stream-alert-bot
        systemctl start stream-alert-bot
    fi

    echo -e "${green}Stream Alert Bot ${tag_version} installed successfully!${plain}"
}

install_base
install_app
