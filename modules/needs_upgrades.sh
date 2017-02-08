#!/bin/bash
#gen:module a package:string,installed:string,available:string

DISTRO=$(grep -oP '^ID="?.*"?' /etc/os-release | cut -d"=" -f2 | tr -d '" ')

debian() {
    aptList=$(apt-get --just-print upgrade 2>&1 | grep '^Inst')
    packages=$(echo "$aptList" | perl -ne 'if (/Inst\s([\w,\-,\d,\.,~,:,\+]+)\s\[([\w,\-,\d,\.,~,:,\+]+)\]\s\(([\w,\-,\d,\.,~,:,\+]+)\)? /i) {print "{\"package\": \"$1\", \"installed\": \"$2\", \"available\": \"$3\"},"}')
    echo '['${packages%?}']'
}

redhat() {
    dnfCmd="$(which dnf 2>/dev/null)"
    yumCmd="$(which yum 2>/dev/null)"

    if [ -n "$dnfCmd" ]; then
        runRedhatList $dnfCmd
    elif [ -n "$yumCmd" ]; then
        runRedhatList $yumCmd
    else
        echo '[]'
    fi
}

runRedhatList() {
    $1 makecache >/dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo '[]'
        return
    fi

    packages=$($1 list updates \
        | sed -r '1,/^Up(graded|dated) Packages/ d' \
        | tr -s ' ' \
        | awk '{print "{\"package\": \""$1"\", \"installed\": \"\", \"available\": \""$2"\"},"}')

    echo '['${packages%?}']'
}

case "$DISTRO" in
debian|ubuntu) debian ;;
fedora|rhel|centos|ol) redhat ;;
esac
