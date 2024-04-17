#!/bin/sh
set -e

ASN_URL=https://github.com/P3TERX/GeoLite.mmdb/releases/latest/download/GeoLite2-ASN.mmdb
COUNTRY_URL=https://github.com/P3TERX/GeoLite.mmdb/releases/latest/download/GeoLite2-Country.mmdb
CITY_URL=https://github.com/P3TERX/GeoLite.mmdb/releases/latest/download/GeoLite2-City.mmdb


destination="db/base_mmdb"
wget -P $destination -N $ASN_URL
wget -P $destination -N $COUNTRY_URL
wget -P $destination -N $CITY_URL
