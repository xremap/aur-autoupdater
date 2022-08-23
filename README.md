# aur-autoupdater

Job to auto update version of AUR package if newer version is available.

# How to run locally
```shell
SSH_KEY=$(cat /home/njkevlani/.ssh/aur_key) SSH_KEY_PASSWORD=$(cat /home/njkevlani/.ssh/aur_key.password) make go-run
```
