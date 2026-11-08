#!/bin/sh

if [ "$1" != "0" ] && [ "$1" != "remove" ] && [ "$1" != "purge" ]; then
    exit 0
fi
if ! command -V systemctl >/dev/null 2>&1; then
    exit 0
fi

name=simple-fileserver

echo "Clean up ${name} service"
systemctl unmask "${name}.service" || true
systemctl disable --now "${name}.service"
