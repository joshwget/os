#!/bin/bash
set -ex

cat > /etc/respawn.conf << EOF
/sbin/getty 115200 tty6
/sbin/getty 115200 tty5
/sbin/getty 115200 tty4
/sbin/getty 115200 tty3
/sbin/getty 115200 tty2
/sbin/getty 115200 tty1
EOF

for i in ttyS{0..4} ttyAMA0; do
    if grep -q 'console='$i /proc/cmdline; then
        echo '/sbin/getty 115200' $i >> /etc/respawn.conf
    fi
done

export TERM=linux
echo AA
sleep 3
echo AA
#exec /sbin/agetty --autologin rancher 15200 ttyS0
exec respawn -f /etc/respawn.conf
echo AA
