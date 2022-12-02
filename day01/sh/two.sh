#!/bin/sh

awk '{ e += $1 } $1 == ""{ print e; e=0 } END{ print e }' | sort -rn | head -3 | awk '{ t += $1 } END{ print t }'
