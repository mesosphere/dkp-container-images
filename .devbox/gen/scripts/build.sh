set -e

if [ -z "$__DEVBOX_SKIP_INIT_HOOK_51aa22d944991e3bf345c85a36cfc4a63d00babbc4e8857aaccdb7e25db6e49e" ]; then
    . "/Users/vedansh.kakkar/dkp-container-images/.devbox/gen/scripts/.hooks.sh"
fi

just --list
