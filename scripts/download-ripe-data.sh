#!/bin/sh

url="https://ftp.ripe.net/ripe/dbase/split/ripe.db.inetnum.gz"
ipV6Url="https://ftp.ripe.net/ripe/dbase/split/ripe.db.inet6num.gz"
apnicV4Url=https://ftp.apnic.net/apnic/whois/apnic.db.inetnum.gz

destination="db"
force="false"

while [ $# -gt 0 ]; do
    if [ "$1" = "--force" ]; then
        force="true"
    fi
    shift
done

if [ ! -f "$destination/ripe.db.inetnum" ] || [ "$force" = "true" ]; then
    echo "Downloading file from $url..."
    wget -q "$url" -P "$destination"
    echo "Extracting file..."
    gzip -d "$destination/ripe.db.inetnum.gz"
    echo "Extraction complete."
    rm -rf "$destination/ripe.db.inetnum.gz"
else
    echo "File already exists in $destination. Use --force to download again."
fi

if [ ! -f "$destination/ripe.db.inet6num" ] || [ "$force" = "true" ]; then
    echo "Downloading IPv6 file from $ipV6Url..."
    wget -q "$ipV6Url" -P "$destination"
    echo "Extracting file..."
    gzip -d "$destination/ripe.db.inet6num.gz"
    echo "Extraction complete."
    rm -rf "$destination/ripe.db.inet6num.gz"
else
    echo "IPv6 file already exists in $destination. Use --force to download again."
fi



if [ ! -f "$destination/apnic.db.inetnum" ] || [ "$force" = "true" ]; then
    echo "Downloading file from $apnicV4Url..."
    wget -q "$apnicV4Url" -P "$destination"
    echo "Extracting file..."
    gzip -d "$destination/apnic.db.inetnum.gz"
    echo "Extraction complete."
    rm -rf "$destination/apnic.db.inetnum.gz"
else
    echo "File already exists in $destination. Use --force to download again."
fi
