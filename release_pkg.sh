#!/bin/bash
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
    echo "[*] Create package: $pkg_name"
    zip -j "$release_dir/$pkg_name" $files_to_include_all
}

make_zip() {
    if [ -n $1 ]; then
        create_pkg $1
    else
        echo "[!] Should have release target."
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
    local target_archs=("${ARCHS[@]}")
    if [ -n "$1" ]; then
        target_archs=("$1")
    fi

    if [ ! -d "./_release/" ]; then
        mkdir -p ./_release/
    else
        rm -rf ./_release/
        mkdir -p ./_release/
    fi
    for arch in "${target_archs[@]}"; do
        echo -e "\033[34m [★] Compile target =>\033[43;34m ["$arch"]. \033[0m"
        case "$arch" in
            "windows")
                build_windows $arch
                make_zip $arch
                echo -e "\033[33m [v] Compile target => ["$arch"] Ok. \033[0m"
                ;;
            "x64linux")
                build_x64linux $arch
                make_zip $arch
                echo -e "\033[33m [v] Compile target => ["$arch"] Ok. \033[0m"
                ;;
            "arm64linux")
                # sudo apt install gcc-arm-linux-gnueabi -y
                build_arm64linux $arch
                make_zip $arch
                echo -e "\033[33m [v] Compile target => ["$arch"] Ok. \033[0m"
                ;;
            "arm32linux")
                # sudo apt install gcc-arm-linux-gnueabi -y
                build_arm32linux $arch
                make_zip $arch
                echo -e "\033[33m [v] Compile target => ["$arch"] Ok. \033[0m"
                ;;
            "riscv64linux")
                # sudo apt install g++-riscv64-linux-gnu gcc-riscv64-linux-gnu -y
                build_riscv64_linux $arch
                make_zip $arch
                echo -e "\033[33m [v] Compile target => ["$arch"] Ok. \033[0m"
                ;;
            *)
                echo "[!] Unknown architecture: $arch"
                ;;
        esac
    done
}

calculate_and_save_md5() {
    if [ $# -ne 1 ]; then
        echo "[*] Usage: $0 <file_path>"
        exit 1
    fi
    local file_path="$1"
    local md5_hash
    if [ ! -f "$file_path" ]; then
        echo "[!] File not found: $file_path"
        return 1
    fi
    md5_hash=$(md5sum "$file_path" | awk '{print $1}')
    echo -n "$md5_hash" > md5.sum
}

gen_changelog() {
    echo -e "[*] Version Change log:"
    log=$(git log --oneline --pretty=format:" \033[0;31m[*]\033[0m%s\n" $(git describe --abbrev=0 --tags).. | cat)
    echo -e $log
}

upload_to_file_server(){
    BASIC_AUTH="rhilex-file-server-admin:rhilex-file-server-admin_secret"
    VERSION="$(git describe --tags $(git rev-list --tags --max-count=1))"
    cd _build/_release/
    UPLOAD_URL="http://112.5.155.64:10120/release/${VERSION}/"
    ZIP_FILES=$(find . -maxdepth 1 -type f -name "rhilex*.zip")
    if [ -z "$ZIP_FILES" ]; then
        echo "[!] No .zip files found in the current directory."
        exit 1
    fi
    for FILE in $ZIP_FILES; do
        upload_path=$(echo "$FILE" | sed 's|^./||')
        echo "[*] Uploading [$FILE] to [${UPLOAD_URL}${upload_path}]"
        curl -T "$FILE" "${UPLOAD_URL}${upload_path}" --user $BASIC_AUTH
        if [ $? -eq 0 ]; then
            echo "[v] Upload $FILE successfully."
        else
            echo "[x] Upload $FILE failure."
        fi
    done
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
            echo -e "\033[32m [v] $dep has been installed. \033[0m"
        fi
    done
}

build_project() {
    check_cmd
    init_env
    cp -r $(ls | egrep -v '^_build$') ./_build/
    cd ./_build/
    cross_compile $2
    gen_changelog
}

display_help() {
    cat << EOF
Usage: $0 <command> [target_arch]

Commands:
  upload   Uploads the build to the file server.
  build    Builds the project. You can specify a target architecture as the second argument.
           Available target architectures are:
             - arm32linux
             - arm64linux
             - riscv64linux
             - x64linux
             - windows
           If no target architecture is specified, all architectures will be built.
  help     Displays this help message.

EOF
}

main() {
    case "$1" in
        upload)
            upload_to_file_server
        ;;
        build)
            build_project $1 $2
        ;;
        help)
            display_help
        ;;
        *)
            echo "Error: Unknown command '$1'"
            display_help
            exit 1
        ;;
    esac
}

main "$@"