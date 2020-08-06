#!/bin/bash

echo "Install 'lastbackend'"
if [[ "$OSTYPE" == "linux-gnu" || "$OSTYPE" == "linux-musl" ]]; then
    mv  build/linux/lastbackend /usr/bin/lastbackend
elif [[ "$OSTYPE" == "darwin"* ]]; then
    mv  build/darwin/lastbackend /usr/bin/lastbackend
elif [[ "$OSTYPE" == "windows"* ]]; then
    mv  build/windows/lastbackend /usr/bin/lastbackend
fi