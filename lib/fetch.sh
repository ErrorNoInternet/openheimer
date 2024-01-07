#!/usr/bin/env bash
set -e

PREFIX=$(dirname "$0")
curl -Lo "$PREFIX/npcap-sdk.zip" https://npcap.com/dist/npcap-sdk-1.13.zip
unzip -o "$PREFIX/npcap-sdk.zip" "Lib/*" -d "$PREFIX/npcap"
rm "$PREFIX/npcap-sdk.zip"
