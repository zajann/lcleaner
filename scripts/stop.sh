#!/bin/sh

# Source function library.
. /etc/rc.d/init.d/functions

# Source .bash_profile of root
. $HOME/.bash_profile

RETVAL=0
PROG="lcleaner"

# Stop daemon
echo -n $"Stopping ${PROG}: "
killproc ${PROG}
RETVAL=$?
echo
exit ${RETVAL}

