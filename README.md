# aur-autoupdater

Job to auto update version of AUR package if newer version is available.

# How to run locally
Set up environment variables.
```shell
export SSH_KEY=$(cat /home/njkevlani/.ssh/aur_key)
export SSH_KEY_PASSWORD=$(cat /home/njkevlani/.ssh/aur_key.password)
export GH_TOKEN_FOR_AUR_AUTO_UPDATE=your_token
```

Run code.
```shell
make go-run
```
