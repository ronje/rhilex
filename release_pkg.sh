#! /bin/bash
set -e
APP=rhilex
RESPOSITORY="https://github.com/hootrhino"
ARCHS=("arm32linux" "arm64linux" "riscv64linux" "x64linux" "windows")

create_pkg() {
    local target=$1
    local version="$(git describe --tags $(git rev-list --tags --max-count=1))"
    local release_dir="_release"
    local pkg_name="${APP}-$target-$version.zip"
    local common_files="./config/*.ini ./md5.sum"
    local files_to_include_all="./${APP} $common_files"
    local files_to_include_win="./${APP}.exe $common_files"

    if [[ "$target" != "windows" ]]; then
        files_to_include_all="$files_to_include_all ./script/*.sh"
        mv ./${APP}-$target ./${APP}
        chmod +x ./${APP}
        calculate_and_save_md5 ./${APP}
    else
        files_to_include_all="$files_to_include_win"
        mv ./${APP}-$target.exe ./${APP}.exe
        calculate_and_save_md5 ./${APP}.exe
    fi
    echo "Create package: $pkg_name"
    zip -j "$release_dir/$pkg_name" $files_to_include_all
}

make_zip() {
    if [ -n $1 ]; then
        create_pkg $1
    else
        echo "Should have release target."
        exit 1
    fi

}

build_windows() {
    make windows
}
build_x64linux() {
    make x64linux
}

build_arm64linux() {
    make arm64
}

build_arm32linux() {
    make arm32
}

build_riscv64_linux() {
    make riscv64
}

cross_compile() {
    if [ ! -d "./_release/" ]; then
        mkdir -p ./_release/
    else
        rm -rf ./_release/
        mkdir -p ./_release/
    fi
    for arch in "${ARCHS[@]}"; do
        echo -e "\033[34m [★] Compile target =>\033[43;34m ["$arch"]. \033[0m"
        if [[ "${arch}" == "windows" ]]; then
            build_windows $arch
            make_zip $arch
            echo -e "\033[33m [√] Compile target => ["$arch"] Ok. \033[0m"
        fi
        if [[ "${arch}" == "x86linux" ]]; then
            build_x86linux $arch
            make_zip $arch
            echo -e "\033[33m [√] Compile target => ["$arch"] Ok. \033[0m"
        fi
        if [[ "${arch}" == "x64linux" ]]; then
            build_x64linux $arch
            make_zip $arch
            echo -e "\033[33m [√] Compile target => ["$arch"] Ok. \033[0m"

        fi
        if [[ "${arch}" == "arm64linux" ]]; then
            # sudo apt install gcc-arm-linux-gnueabi -y
            build_arm64linux $arch
            make_zip $arch
            echo -e "\033[33m [√] Compile target => ["$arch"] Ok. \033[0m"

        fi
        if [[ "${arch}" == "arm32linux" ]]; then
            # sudo apt install gcc-arm-linux-gnueabi -y
            build_arm32linux $arch
            make_zip $arch
            echo -e "\033[33m [√] Compile target => ["$arch"] Ok. \033[0m"
        fi
        if [[ "${arch}" == "riscv64linux" ]]; then
            # sudo apt install g++-riscv64-linux-gnu gcc-riscv64-linux-gnu -y
            build_riscv64_linux $arch
            make_zip $arch
            echo -e "\033[33m [√] Compile target => ["$arch"] Ok. \033[0m"
        fi
    done
}

calculate_and_save_md5() {
    if [ $# -ne 1 ]; then
        echo "Usage: $0 <file_path>"
        exit 1
    fi
    local file_path="$1"
    local md5_hash
    if [ ! -f "$file_path" ]; then
        echo "File not found: $file_path"
        return 1
    fi
    md5_hash=$(md5sum "$file_path" | awk '{print $1}')
    echo -n "$md5_hash" > md5.sum
}

fetch_dashboard() {
    local owner="hootrhino"
    local repo="rhilex-web"
    if [ -f "www.zip" ]; then
        echo "[!] www.zip already exists. No need to download."
        exit 0
    fi
    local tag=$(curl -s "https://api.github.com/repos/$owner/$repo/releases/latest" | jq -r .tag_name)
    local zip_url=$(curl -s "https://api.github.com/repos/$owner/$repo/releases/latest" | jq -r '.assets[] | select(.name == "www.zip") | .browser_download_url')
    if [ -z "$zip_url" ]; then
        echo "[x] Error: www.zip not found in the release assets."
        exit 1
    fi
    curl -L -o www.zip "$zip_url"
    echo "[√] Download complete. Tag: $tag"
    unzip -o www.zip -d /plugin/apiserver/server/www/
    echo "[√] Extraction complete. www.zip contents have been overwritten to /plugin/apiserver/server/www/."
}

gen_changelog() {
    echo -e "[.]Version Change log:"
    log=$(git log --oneline --pretty=format:" \033[0;31m[*]\033[0m%s\n" $(git describe --abbrev=0 --tags).. | cat)
    echo -e $log
}
gen_release_versions(){
    cd _release
    local json="{\"code\": 200, \"data\": {\"software_versions\": ["
    local version="$(git describe --tags $(git rev-list --tags --max-count=1))"
    declare -A platforms
    platforms["windows"]="Windows:x86_64"
    platforms["ubuntu_arm32"]="Ubuntu:arm32"
    platforms["ubuntu_arm64"]="Ubuntu:arm64"
    platforms["ubuntu_x86_64"]="Ubuntu:x86_64"
    platforms["ubuntu_riscv64"]="Ubuntu:riscv64"
    platforms["debian_arm32"]="Debian:arm32"
    platforms["debian_arm64"]="Debian:arm64"
    platforms["debian_x86_64"]="Debian:x86_64"
    platforms["debian_riscv64"]="Debian:riscv64"
    platforms["busybox_riscv64"]="Busybox:riscv64"
    platforms["openwrt_arm32"]="OpenWRT:arm32"
    platforms["openwrt_arm64"]="OpenWRT:arm64"
    platforms["openwrt_mips"]="OpenWRT:mips"
    platforms["openwrt_riscv64"]="OpenWRT:riscv64"
    json+="{\"version\": \"$version\", \"platforms\": {"
    for platform in "windows" "ubuntu" "debian" "openwrt" "busybox"; do
        json+="\"$platform\": ["
        case $platform in
            windows)
                json+="{\"name\": \"${platforms[windows]}\", \"url\": \"/share/rhilex-windows-$version.zip\"}"
            ;;
            ubuntu)
                json+="{\"name\": \"${platforms[ubuntu_arm32]}\", \"url\": \"/share/rhilex-arm32linux-$version.zip\"},"
                json+="{\"name\": \"${platforms[ubuntu_arm64]}\", \"url\": \"/share/rhilex-arm64linux-$version.zip\"},"
                json+="{\"name\": \"${platforms[ubuntu_x86_64]}\", \"url\": \"/share/rhilex-x64linux-$version.zip\"},"
                json+="{\"name\": \"${platforms[ubuntu_riscv64]}\", \"url\": \"/share/rhilex-riscv64linux-$version.zip\"}"
            ;;
            debian)
                json+="{\"name\": \"${platforms[debian_arm32]}\", \"url\": \"/share/rhilex-arm32linux-$version.zip\"},"
                json+="{\"name\": \"${platforms[debian_arm64]}\", \"url\": \"/share/rhilex-arm64linux-$version.zip\"},"
                json+="{\"name\": \"${platforms[debian_x86_64]}\", \"url\": \"/share/rhilex-x64linux-$version.zip\"},"
                json+="{\"name\": \"${platforms[debian_riscv64]}\", \"url\": \"/share/rhilex-riscv64linux-$version.zip\"}"
            ;;
            openwrt)
                json+="{\"name\": \"${platforms[openwrt_arm32]}\", \"url\": \"/share/rhilex-arm32linux-$version.zip\"},"
                json+="{\"name\": \"${platforms[openwrt_arm64]}\", \"url\": \"/share/rhilex-arm64linux-$version.zip\"},"
                json+="{\"name\": \"${platforms[openwrt_mips]}\", \"url\": \"/share/rhilex-mipslinux-$version.zip\"},"
                json+="{\"name\": \"${platforms[openwrt_riscv64]}\", \"url\": \"/share/rhilex-riscv64linux-$version.zip\"}"
            ;;
            busybox)
                json+="{\"name\": \"${platforms[busybox_riscv64]}\", \"url\": \"/share/rhilex-riscv64linux-$version.zip\"}"
            ;;
        esac

        json+="],"
    done
    json=${json%,}
    json+="}}]}"
    json+=",\"msg\": \"Success\"}"
    echo "$json" > "${version}.json"
    cd ..

}
init_env() {
    if [ ! -d "./_build/" ]; then
        mkdir -p ./_build/
    else
        rm -rf ./_build/
        mkdir -p ./_build/
    fi

}

check_cmd() {
    DEPS=("bash" "git" "jq" "gcc" "make" "x86_64-w64-mingw32-gcc" "aarch64-linux-gnu-gcc" "arm-linux-gnueabi-gcc")
    for dep in ${DEPS[@]}; do
        echo -e "\033[34m [*] Check Env: $dep. \033[0m"
        if ! [ -x "$(command -v $dep)" ]; then
            echo -e "\033[31m |x| Error: $dep is not installed. \033[0m"
            exit 1
        else
            echo -e "\033[32m [√] $dep has been installed. \033[0m"
        fi
    done

}
main(){
    check_cmd
    init_env
    cp -r $(ls | egrep -v '^_build$') ./_build/
    cd ./_build/
    # fetch_dashboard
    cross_compile
    gen_release_versions
    gen_changelog
    find . -mindepth 1 -not -path "./_release/*" -not -name "_release" -exec rm -rf {} +
}
#
#-----------------------------------
#
main