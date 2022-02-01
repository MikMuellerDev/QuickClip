#!/bin/bash
if [ ! -f "/app/QuickClip/config/config.json" ]; then
    echo -e "[PRELOAD] The config file does not exist, creating."
    touch /app/QuickClip/config/config.json
fi

if [ ! -f "/app/QuickClip/config/clipboard.json" ]; then
    echo -e "[PRELOAD] The clipboard file does not exist, creating."
    touch /app/QuickClip/config/clipboard.json
    echo "{\"Clips\":[{\"Name\":\"Default Document\",\"Id\":\"default\",\"Content\":\"Write something.\",\"Description\":\"Test Document's description\",\"Restricted\":false,\"Refresh\":false,\"RefreshInterval\":0,\"ReadOnly\":false}]}" > /app/QuickClip/config/clipboard.json
fi

cd ./bin || exit 1
./QuickClip
