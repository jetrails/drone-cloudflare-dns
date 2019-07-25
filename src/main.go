package main

import "fmt"
import "os"
import "strings"
import "github.com/caarlos0/env"
import "github.com/cloudflare/cloudflare-go"

type Input struct {
	Debug           bool          `env:"PLUGIN_DEBUG" envDefault:false`
	ApiToken        string        `env:"PLUGIN_API_TOKEN,required"`
	ZoneIdentifier  string        `env:"PLUGIN_ZONE_IDENTIFIER,required"`
	Action          string        `env:"PLUGIN_ACTION,required"`
	RecordType      string        `env:"PLUGIN_RECORD_TYPE,required"`
	RecordName      string        `env:"PLUGIN_RECORD_NAME,required"`
	RecordContent   string        `env:"PLUGIN_RECORD_CONTENT"`
	RecordProxied   bool          `env:"PLUGIN_RECORD_PROXIED" envDefault:true`
	RecordTTL       int           `env:"PLUGIN_RECORD_TTL" envDefault:1`
	RecordPriority  int           `env:"PLUGIN_RECORD_PRIORITY" envDefault:1`
}

func handleCloudflareError ( instance error ) {
	if instance != nil {
		fmt.Printf ( "error: %s\n", instance.Error () )
		os.Exit ( 3 )
	}
}

func main () {
	// Parse environmental variables
	input := Input {}
	if error := env.Parse ( &input ); error != nil {
		fmt.Printf ( "%+v\n", error )
		os.Exit ( 1 )
	}
	// Normalize input
	input.Action = strings.ToLower ( input.Action )
	input.RecordType = strings.ToUpper ( input.RecordType )
	// Semantic checks
	if input.Action != "set" && input.Action != "unset" {
		fmt.Printf ("error: PLUGIN_ACTION must equal one of [set,unset]\n")
		os.Exit ( 2 )
	}
	if input.Action == "set" && input.RecordContent == "" {
		fmt.Printf ("error: PLUGIN_RECORD_CONTENT must be set with 'set' action\n")
		os.Exit ( 2 )
	}
	// Create client for cloudflare api communication
	api, error := cloudflare.NewWithAPIToken ( input.ApiToken )
	handleCloudflareError ( error )
	// Print input if debug mode is on
	if input.Debug {
		fmt.Printf ( "Input: %+v\n", input )
	}
	// Initialize all records that will be used
	searchRecord := cloudflare.DNSRecord {
		Type: input.RecordType,
		Name: input.RecordName,
	}
	createRecord := cloudflare.DNSRecord {
		Type: input.RecordType,
		Name: input.RecordName,
		Content: input.RecordContent,
		Proxied: input.RecordProxied,
		TTL: input.RecordTTL,
		Priority: input.RecordPriority,
	}
	// Check to see if record was found
	found, error := api.DNSRecords ( input.ZoneIdentifier, searchRecord )
	handleCloudflareError ( error )
	// Switch based on action
	if input.Action == "unset" && len ( found ) > 0 {
		id := found [ 0 ].ID
		error := api.DeleteDNSRecord ( input.ZoneIdentifier, id )
		handleCloudflareError ( error )
	} else if input.Action == "set" && len ( found ) > 0 {
		updateRecord := found [ 0 ]
		id := updateRecord.ID
		updateRecord.Content = input.RecordContent
		if os.Getenv ("PLUGIN_RECORD_PROXIED") != "" {
			updateRecord.Proxied = input.RecordProxied
		}
		if os.Getenv ("PLUGIN_RECORD_TTL") != "" {
			updateRecord.TTL = input.RecordTTL
		}
		if os.Getenv ("PLUGIN_RECORD_PRIORITY") != "" {
			updateRecord.Priority = input.RecordPriority
		}
		error := api.UpdateDNSRecord ( input.ZoneIdentifier, id, updateRecord )
		handleCloudflareError ( error )
	} else if input.Action == "set" {
		res, error := api.CreateDNSRecord ( input.ZoneIdentifier, createRecord )
		handleCloudflareError ( error )
		if input.Debug {
			fmt.Printf ( "CreateDNSRecord: %+v\n", res )
		}
	} else {
		fmt.Printf ("error: could not find record to unset\n")
		os.Exit ( 4 )
	}
}
