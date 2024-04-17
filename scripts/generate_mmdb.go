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

	rows, err := sqlite_db.Query("SELECT * FROM ip_data ORDER BY first_ip_int, subnet ASC")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var ipData IPData
		err := rows.Scan(&ipData.ID, &ipData.FirstIP, &ipData.LastIP, &ipData.FirstIPInt, &ipData.LastIPInt, &ipData.IPVersion, &ipData.Subnet, &ipData.NetworkPrefix, &ipData.Netname, &ipData.Country, &ipData.Description, &ipData.MntBy)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("ID: %d, FirstIP: %s, LastIP: %s, FirstIPInt: %s, LastIPInt: %s, IPVersion: %d, Subnet: %d, NetworkPrefix: %s, Netname: %s, Country: %s, Description: %s, MntBy: %s\n", ipData.ID, ipData.FirstIP, ipData.LastIP, ipData.FirstIPInt, ipData.LastIPInt, ipData.IPVersion, ipData.Subnet, ipData.NetworkPrefix, ipData.Netname, ipData.Country, ipData.Description, ipData.MntBy)
	}

	// writer, err := mmdbwriter.Load("../db/base_mmdb/GeoLite2-ASN.mmdb", mmdbwriter.Options{})
	writer, err := mmdbwriter.New(mmdbwriter.Options{}) //("../db/base_mmdb/GeoLite2-ASN_2.mmdb", mmdbwriter.Options{})
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

	_, network, err := net.ParseCIDR("194.243.217.64/31")

	if err != nil {
		log.Fatal(err)
	}

	record := mmdbtype.Map{}

	// record["autonomous_system_number"] = mmdbtype.Uint32(3269)

	ip := net.ParseIP("194.243.217.64")

	var asn_record map[string]any
	if err := asn_db.Lookup(ip, &asn_record); err != nil {
		log.Fatal(err)
	}

	fmt.Println(asn_record)
	asn_number, ok := asn_record["autonomous_system_number"].(uint64)
	if ok {
		record["asn_number"] = mmdbtype.Uint32(asn_number)
		// fmt.Printf("asn_number:, Type %T", asn_number, asn_number)
	}

	asn_name, ok := asn_record["autonomous_system_organization"].(string)
	if ok {
		record["asn_name"] = mmdbtype.String(asn_name)
		fmt.Println("asn_name:", asn_name)
	}

	var city_record map[string]interface{}
	if err := city_db.Lookup(ip, &city_record); err != nil {
		log.Fatal(err)
	}

	if city, ok := city_record["city"].(map[string]interface{}); ok {
		if names, ok := city["names"].(map[string]interface{}); ok {
			if cityName, ok := names["en"].(string); ok {
				fmt.Println("cityName:", cityName)
				record["city_name"] = mmdbtype.String(cityName)

			}
		}
	}

	if country, ok := city_record["country"].(map[string]interface{}); ok {
		if names, ok := country["names"].(map[string]interface{}); ok {
			if countryName, ok := names["en"].(string); ok {
				record["country_name"] = mmdbtype.String(countryName)
				fmt.Println("countryName:", countryName)
			}
		}

		if iso_code, ok := country["iso_code"].(string); ok {
			record["iso_code"] = mmdbtype.String(iso_code)
		}
	}

	if continent, ok := city_record["continent"].(map[string]interface{}); ok {
		if names, ok := continent["names"].(map[string]interface{}); ok {
			if continentName, ok := names["en"].(string); ok {
				record["continent_name"] = mmdbtype.String(continentName)
				fmt.Println("continentName:", continentName)
			}
		}

		if code, ok := continent["code"].(string); ok {
			record["continent_code"] = mmdbtype.String(code)
			fmt.Println("code:", code)
		}
	}

	if location, ok := city_record["location"].(map[string]interface{}); ok {
		if latitude, ok := location["latitude"].(string); ok {
			record["latitude"] = mmdbtype.String(latitude)
			fmt.Println("latitude:", latitude)
		}
		if longitude, ok := location["longitude"].(string); ok {
			record["longitude"] = mmdbtype.String(longitude)
			fmt.Println("longitude:", longitude)
		}

		if time_zone, ok := location["time_zone"].(string); ok {
			record["time_zone"] = mmdbtype.String(time_zone)
			fmt.Println("time_zone:", time_zone)

		}

		if accuracy_radius, ok := location["accuracy_radius"].(string); ok {
			record["accuracy_radius"] = mmdbtype.String(accuracy_radius)
			fmt.Println("accuracy_radius:", accuracy_radius)
		}

	}

	// record["autonomous_system_organization"] = mmdbtype.String("Cloudflare_ Mia khalifa")
	record["mnt-by"] = mmdbtype.String("mnt-by")

	err = writer.Insert(network, record)
	if err != nil {
		log.Fatal(err)
	}

	fh, err := os.Create("../output/mixed.mmdb")
	if err != nil {
		log.Fatal(err)
	}

	_, err = writer.WriteTo(fh)
	if err != nil {
		log.Fatal(err)
	}

	// writer.InsertRange(startIpBytes,endIpBytes,)
	// log.Println(writer.)

	// writer.InsertRange("")
	log.Println("Loaded database")

}
