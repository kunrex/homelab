#!/bin/bash
echo "start entrypoint.sh"

set -e

# Folder for sshd. No Change.
mkdir -p /var/run/sshd

# These folders are mounted at runtime, so change their ownership to avoid permission issues
chown -R $WORKSPACE_USER:$WORKSPACE_USER /home/$WORKSPACE_USER/.ssh
chown -R $WORKSPACE_USER:$WORKSPACE_USER /home/$WORKSPACE_USER/.config

exec "$@"
