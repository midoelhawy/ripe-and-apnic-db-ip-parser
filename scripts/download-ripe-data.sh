#!/bin/sh

url="https://ftp.ripe.net/ripe/dbase/split/ripe.db.inetnum.gz"
ipV6Url="https://ftp.ripe.net/ripe/dbase/split/ripe.db.inet6num.gz"
apnicV4Url=https://ftp.apnic.net/apnic/whois/apnic.db.inetnum.gz
apnicV6Url=https://ftp.apnic.net/apnic/whois/apnic.db.inet6num.gz
afrinicUrl=https://ftp.afrinic.net/dbase/afrinic.db.gz
lacnicUrl=https://ftp.lacnic.net/lacnic/dbase/lacnic.db.gz

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


if [ ! -f "$destination/apnic.db.inet6num" ] || [ "$force" = "true" ]; then
    echo "Downloading IPv6 file from $apnicV6Url..."
    wget -q "$apnicV6Url" -P "$destination"
    echo "Extracting file..."
    gzip -d "$destination/apnic.db.inet6num.gz"
    echo "Extraction complete."
    rm -rf "$destination/apnic.db.inet6num.gz"
    echo "Extraction complete."
else
    echo "IPv6 file already exists in $destination. Use --force to download again."
fi



if [ ! -f "$destination/afrinic.db" ] || [ "$force" = "true" ]; then
    echo "Downloading IPv file from $afrinicUrl..."
    wget -q "$afrinicUrl" -P "$destination"
    echo "Extracting file..."
    gzip -d "$destination/afrinic.db.gz"
    echo "Extraction complete."
    rm -rf "$afrinicUrl/afrinic.db.gz"
    echo "Extraction complete."
else
    echo "IP Africa file already exists in $destination. Use --force to download again."
fi


if [ ! -f "$destination/lacnic.db" ] || [ "$force" = "true" ]; then
    echo "Downloading IPv file from $afrinicUrl..."
    wget -q "$lacnicUrl" -P "$destination"
    echo "Extracting file..."
    gzip -d "$destination/lacnic.db.gz"
    echo "Extraction complete."
    rm -rf "$lacnicUrl/lacnic.db.gz"
    echo "Extraction complete."
else
    echo "IP lacnic file already exists in $destination. Use --force to download again."
fi
