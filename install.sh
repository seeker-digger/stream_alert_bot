#!/bin/bash

red='\033[0;31m'
green='\033[0;32m'
yellow='\033[0;33m'
plain='\033[0m'

# Check root
if [[ "$EUID" -ne 0 ]]; then
  echo -e "${red}""Fatal error: ${plain}""Please, run this script as root!"
  exit 1
fi

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

arch() {
    case "$(uname -m)" in
    x86_64 | x64 | amd64) echo 'amd64' ;;
    #i*86 | x86) echo '386' ;;
    #armv8* | armv8 | arm64 | aarch64) echo 'arm64' ;;
    #armv7* | armv7 | arm) echo 'armv7' ;;
    #armv6* | armv6) echo 'armv6' ;;
    #armv5* | armv5) echo 'armv5' ;;
    #s390x) echo 's390x' ;;
    *) echo -e "${green}Unsupported CPU architecture! ${plain}" && rm -f install.sh && exit 1 ;;
    esac
}

echo "Arch: $(arch)"

install_base() {
    case "${release}" in
    ubuntu | debian | armbian)
        apt-get update && apt-get install -y -q wget curl tar jq
        ;;
    centos | rhel | almalinux | rocky | ol)
        yum -y update && yum install -y -q wget curl tar jq
        ;;
    fedora | amzn | virtuozzo)
        dnf -y update && dnf install -y -q wget curl tar jq
        ;;
    arch | manjaro | parch)
        pacman -Syu && pacman -Syu --noconfirm wget curl tar jq
        ;;
    opensuse-tumbleweed | opensuse-leap)
        zypper refresh && zypper -q install -y wget curl tar jq
        ;;
    alpine)
        apk update && apk add wget curl tar jq
        ;;
    *)
        apt-get update && apt-get install -y -q wget curl tar jq
        ;;
    esac
}
#REMOVE WHEN BE PUBLIC!!
init_private_repo() {
  echo -ne "${yellow}""Enter the bearer key: ""${plain}"
  read -r key
}

#CHANGE WHEN BE PUBLIC!!
install_app() {
    cd /usr/local || exit 1

    release_json=$(curl -s -H "Authorization: Bearer ${key}" \
                        -H "Accept: application/vnd.github+json" \
                        "https://api.github.com/repos/seeker-digger/stream_alert_bot/releases")

    if [[ -z "$release_json" ]]; then
        echo "${red}Fatal error:${plain} Failed to fetch releases from GitHub"
        exit 1
    fi

    tag_version=$(echo "$release_json" | jq -r '.[0].tag_name')
    asset_url=$(echo "$release_json" | jq -r '.[0].assets[] | select(.name=="stream_alert_bot-linux_amd64") | .url')

    if [[ -z "$tag_version" || -z "$asset_url" ]]; then
        echo "${red}Fatal error:${plain} Could not determine latest release or asset"
        exit 1
    fi

    echo "Got latest version: ${tag_version}, beginning installation..."

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

    wget --header="Authorization: Bearer ${key}" \
         --header="Accept: application/octet-stream" \
         -O bot "$asset_url"

    if [[ $? -ne 0 ]]; then
        echo "${red}Fatal error:${plain} Failed to download release asset!"
        exit 1
    fi

    curl -s -H "Authorization: Bearer ${key}" \
         -H "Accept: application/vnd.github.v3.raw" \
         -L -o stream-alert-bot.service \
         "https://api.github.com/repos/seeker-digger/stream_alert_bot/contents/stream-alert-bot.service?ref=master"

    if [[ $? -ne 0 ]]; then
        echo "${red}Fatal error:${plain} Failed to download service file!"
        rm -f bot
        exit 1
    fi

    chmod +x bot

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
init_private_repo
install_app

