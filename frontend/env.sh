#!/bin/sh

## Helper script to run the app in Docker container

cp env-template.js env.js

if [ ! -z "$API_BASE_URL" ]
then
    echo "Setting API Base URL to $API_BASE_URL"

    sed -i "s~__API_BASE_URL__~$API_BASE_URL~g" env.js
else
    echo "API Base Url env var is required"
    exit 1
fi

mv env.js /usr/share/nginx/html/env.js
