#!/bin/sh

cleanup() {
    # Cleanup s6 event folder
    # On start and on shutdown in case container has been killed
    rm -rf "$(find /app/ggit/docker/s6/ -name 'event')"
}

create_volume_subfolder() {
    # only change ownership if needed, if using an nfs mount this could be expensive
    if [ "$USER:$USER" != "$(stat /data -c '%U:%G')" ]
    then
        # Modify the owner of /data dir, make $USER(git) user have permission to create sub-dir in /data.
        chown -R "$USER:$USER" /data
    fi

    # Create VOLUME subfolder
    for f in /data/data /data/repositories /data/conf /data/log /data/ssh /data/data/attachments; do
        if ! test -d $f; then
            mkdir -p $f
        fi
    done
}

setids() {
    export USER=git
    PUID=${PUID:-1000}
    PGID=${PGID:-1000}
    groupmod -o -g "$PGID" $USER
    usermod -o -u "$PUID" $USER
}

setids
cleanup
create_volume_subfolder

# Cron Daemon (crond)
CROND=$(echo "$RUN_CROND" | tr '[:upper:]' '[:lower:]')
if [ "$CROND" = "true" ] || [ "$CROND" = "1" ]; then
    echo "init:crond  | Cron Daemon (crond) will be run as requested by s6" 1>&2
    rm -f /app/ggit/docker/s6/crond/down
    /bin/sh /app/ggit/docker/runtime/backup-init.sh "${PUID}"
else
    # Tell s6 not to run the crond service
    touch /app/ggit/docker/s6/crond/down
fi

# Exec CMD or S6 by default if nothing present
if [ $# -gt 0 ];then
    exec "$@"
else
    exec /usr/bin/s6-svscan /app/ggit/docker/s6/
fi
