#!/bin/sh

# Block ifup until DAD completion
# Copyright (c) 2016-2018 Kaarle Ritvanen

has_flag() {
	ip address show dev $IFACE | grep -q " $1 "
}

while has_flag tentative && ! has_flag dadfailed; do
	sleep 0.2
done
