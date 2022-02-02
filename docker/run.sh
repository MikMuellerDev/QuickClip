#!/bin/bash

### Config files ###
if [ ! -f "/app/quickclip/config/config.json" ]; then
    echo -e "[PRELOAD] The config file does not exist, creating."
    touch /app/quickclip/config/config.json
fi

if [ ! -f "/app/quickclip/config/clipboard.json" ]; then
    echo -e "[PRELOAD] The clipboard file does not exist, creating."
    touch /app/quickclip/config/clipboard.json
    echo "{\"Clips\":[{\"Name\":\"Default Document\",\"Id\":\"default\",\"Content\":\"Write something.\",\"Description\":\"Test Document's description\",\"Restricted\":false,\"Refresh\":false,\"RefreshInterval\":0,\"ReadOnly\":false}]}" > /app/quickclip/config/clipboard.json
fi

### Log files ###
mkdir -p /app/quickclip/log/
if [ ! -f "/app/quickclip/log/application.log" ]; then
    echo -e "[PRELOAD] The application.log file does not exist, creating."
    touch /app/quickclip/log/application.log
fi

if [ ! -f "/app/quickclip/log/error.log" ]; then
    echo -e "[PRELOAD] The error.log file does not exist, creating."
    touch /app/quickclip/log/error.log
fi

cd ./bin || exit 1
./QuickClip
