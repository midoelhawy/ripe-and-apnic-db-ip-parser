package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"net"

	_ "github.com/mattn/go-sqlite3"
	"github.com/maxmind/mmdbwriter"
	"github.com/maxmind/mmdbwriter/mmdbtype"
	"github.com/oschwald/maxminddb-golang"
)

type IPData struct {
	ID            int
	FirstIP       string
	LastIP        string
	FirstIPInt    string
	LastIPInt     string
	IPVersion     int
	Subnet        int
	NetworkPrefix sql.NullString
	Netname       string
	Country       string
	Description   string
	MntBy         sql.NullString
}

func main() {

	sqlite_db, err := sql.Open("sqlite3", "../ripe_data.db")
	if err != nil {
		log.Fatal(err)
	}
	defer sqlite_db.Close()

	writer, err := mmdbwriter.New(mmdbwriter.Options{
		RecordSize:              32,
		IncludeReservedNetworks: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	asn_db, err := maxminddb.Open("../db/base_mmdb/GeoLite2-ASN.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	defer asn_db.Close()

	city_db, err := maxminddb.Open("../db/base_mmdb/GeoLite2-City.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	defer city_db.Close()

	rows, err := sqlite_db.Query("SELECT * FROM ip_data where subnet > 0 ORDER by cast( first_ip_int as unsigned) ASC, subnet ASC")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var counter = 0
	for rows.Next() {
		var ipData IPData
		err := rows.Scan(&ipData.ID, &ipData.FirstIP, &ipData.LastIP, &ipData.FirstIPInt, &ipData.LastIPInt, &ipData.IPVersion, &ipData.Subnet, &ipData.NetworkPrefix, &ipData.Netname, &ipData.Country, &ipData.Description, &ipData.MntBy)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Processing IP; %s", (ipData.FirstIP))

		checkIP := net.ParseIP(ipData.FirstIP)

		if checkIP.IsPrivate() {
			log.Printf("FirstIP is Reserved; %s", (ipData.FirstIP))
			continue
		}

		record, err := buildMMDBRecord(ipData, asn_db, city_db)
		if err != nil {
			log.Fatal(err)
		}

		network, err := getNetworkFromRecord(ipData)
		if err != nil {
			log.Fatal(err)
		}

		err = writer.Insert(network, record)
		if err != nil {
			log.Fatal(err)
		}
		counter++
		// fmt.Printf("ID: %d, FirstIP: %s, LastIP: %s, FirstIPInt: %s, LastIPInt: %s, IPVersion: %d, Subnet: %d, NetworkPrefix: %s, Netname: %s, Country: %s, Description: %s, MntBy: %s\n", ipData.ID, ipData.FirstIP, ipData.LastIP, ipData.FirstIPInt, ipData.LastIPInt, ipData.IPVersion, ipData.Subnet, ipData.NetworkPrefix, ipData.Netname, ipData.Country, ipData.Description, ipData.MntBy)
		fmt.Printf("Counter %d, FirstIP: %s/%d\n", counter, ipData.FirstIP, ipData.Subnet)
	}

	fh, err := os.Create("../output/ASN_COUNTRY_AND_CITY.mmdb")
	if err != nil {
		log.Fatal(err)
	}

	_, err = writer.WriteTo(fh)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Loaded generated in ../output/ASN_COUNTRY_AND_CITY.mmdb")

}

func buildMMDBRecord(ipData IPData, asnDB, cityDB *maxminddb.Reader) (mmdbtype.Map, error) {
	record := mmdbtype.Map{}

	ip := net.ParseIP(ipData.FirstIP)

	var asnRecord map[string]interface{}
	if err := asnDB.Lookup(ip, &asnRecord); err != nil {
		return record, err
	}

	asnNumber, ok := asnRecord["autonomous_system_number"].(uint64)
	if ok {
		record["asn_number"] = mmdbtype.Uint32(asnNumber)
	}

	asnName, ok := asnRecord["autonomous_system_organization"].(string)
	if ok {
		record["asn_name"] = mmdbtype.String(asnName)
	}

	var mntByValue string
	if ipData.MntBy.Valid {
		mntByValue = ipData.MntBy.String
	} else {
		mntByValue = ""
	}
	record["mnt_by"] = mmdbtype.String(mntByValue)
	record["netname"] = mmdbtype.String(ipData.Netname)
	// record["subnet"] = mmdbtype.Uint32(ipData.Subnet)
	// record["first_ip"] = mmdbtype.String(ipData.FirstIP)

	var cityRecord map[string]interface{}
	if err := cityDB.Lookup(ip, &cityRecord); err != nil {
		return nil, err
	}

	if city, ok := cityRecord["city"].(map[string]interface{}); ok {
		if names, ok := city["names"].(map[string]interface{}); ok {
			if cityName, ok := names["en"].(string); ok {
				record["city_name"] = mmdbtype.String(cityName)
			}
		}
	}

	if country, ok := cityRecord["country"].(map[string]interface{}); ok {
		if names, ok := country["names"].(map[string]interface{}); ok {
			if countryName, ok := names["en"].(string); ok {
				record["country_name"] = mmdbtype.String(countryName)
			}
		}

		if isoCode, ok := country["iso_code"].(string); ok {
			record["iso_code"] = mmdbtype.String(isoCode)
		}
	}

	return record, nil
}

func getNetworkFromRecord(ipData IPData) (*net.IPNet, error) {

	_, network, err := net.ParseCIDR(fmt.Sprintf("%s/%d", ipData.FirstIP, ipData.Subnet))
	if err != nil {
		return nil, err
	}

	return network, nil
}
