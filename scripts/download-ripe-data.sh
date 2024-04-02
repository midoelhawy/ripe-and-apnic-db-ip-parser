#!/bin/sh


url="https://ftp.ripe.net/ripe/dbase/split/ripe.db.inetnum.gz"
destination="db"
echo "Downloading file from $url..."
wget -q "$url" -P "$destination"

echo "Extracting file..."
gzip -d "$destination/ripe.db.inetnum.gz"

rm -rf "$destination/ripe.db.inetnum.gz"
echo "Extraction complete."
