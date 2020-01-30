#!/bin/sh
read -p "# Enter the installation directory : " INSDIR

if test -z $INSDIR; then
	echo "ERROR: Invalid installation directory"
	exit 1
fi

set -e

mkdir -p $INSDIR/bin
mkdir -p $INSDIR/etc
mkdir -p $INSDIR/log

rm -rf $INSDIR/bin/*
rm -rf $INSDIR/etc/*

cd ../cmd
go build -o lcleaner main.go
cd ../scripts 

cp -f ../cmd/lcleaner					$INSDIR/bin/lcleaner
cp -f ../configs/lcleaner_config.yml 	$INSDIR/etc/lcleaner_config.yml
cp -f start.sh 							$INSDIR/bin/start.sh
cp -f stop.sh							$INSDIR/bin/stop.sh

rm -rf ../cmd/lcleaner

chmod 755 $INSDIR/bin/lcleaner
chmod 755 $INSDIR/bin/start.sh
chmod 755 $INSDIR/bin/stop.sh
chmod 644 $INSDIR/etc/lcleaner_config.yml

echo "Successfully Installed lcleaner"
