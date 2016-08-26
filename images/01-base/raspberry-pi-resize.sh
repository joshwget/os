#!/bin/bash
set -ex

fdisk /dev/mmcblk0 <<EOF || true
d

n
p



w
EOF

partx --update /dev/mmcblk0
partx --update /dev/mmcblk0p2

resize2fs /dev/mmcblk0p2
