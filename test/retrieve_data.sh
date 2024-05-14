#! /bin/bash

if [ $# -ne 1 ]; then
    echo "Needed only one argument, the config path."
    exit 0
fi 

config=$1
if [ ! -f "${config}" ]; then
    echo "Invalid config path ${config}"
    exit 0
fi

while IFS="" read -r p || [ -n "$p" ]
do
    line=${p}
    IFS='|' read -r -a parts <<< "$line"

    auth=${parts[0]}
    server=${parts[1]}
    force_download=${parts[2]}
    name=${parts[3]}
    preview=${parts[4]}
    output=${parts[5]}

    json="{ \"server\": ${server}, \"force_download\": ${force_download}, \"name\": ${name}, \"preview_pattern\": ${preview} }"

    http_status_code=$(curl -s -w "%{http_code}" -X POST -d "${json}" https://acsm.domain.com/live-timing -H "Content-Type:application/json" -H "Auth: ${auth}" -o "${output}.tmp")

    if [ "$http_status_code" == "200" ]; then
        mv "${output}.tmp" "${output}"
    fi

done < "${config}"
