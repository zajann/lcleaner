#!/bin/sh

# Source .bash_profile of root
. $HOME/.bash_profile

RETVAL=0
PROG="lcleaner"

# Start daemon
/usr/local/lcleaner/bin/lcleaner -c=/usr/local/lcleaner/etc/lcleaner_config.yml &
RETVAL=$?
exit ${RETVAL}
