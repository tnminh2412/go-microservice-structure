#!/bin/bash
/usr/bin/wget --no-verbose -Y off --tries=1 --spider http://127.0.0.1:$app_port || exit 1
