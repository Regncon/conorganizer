#!/usr/bin/env bash
set -euo pipefail

username="$1"

useradd -m -s /bin/bash "$username"
usermod -aG adm,sudo,www-data,users,deploy "$username"

password="$(pwgen -s 20 1)"
echo "${username}:${password}" | chpasswd
passwd -e "$username" >/dev/null

home_dir="/home/${username}"
ln -sfn /opt/conorganizer "${home_dir}/conorganizer"
ln -sfn /opt/meetupJan2026 "${home_dir}/meetupJan2026"
chown -h "$username:$username" "${home_dir}/conorganizer" "${home_dir}/meetupJan2026"

echo "$password"
