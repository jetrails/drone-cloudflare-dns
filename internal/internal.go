package plugin

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/caarlos0/env/v7"
	"github.com/cloudflare/cloudflare-go"
)

type Input struct {
	Debug          bool    `env:"PLUGIN_DEBUG" envDefault:false`
	ApiToken       string  `env:"PLUGIN_API_TOKEN,required"`
	ZoneIdentifier string  `env:"PLUGIN_ZONE_IDENTIFIER,required"`
	Action         string  `env:"PLUGIN_ACTION,required"`
	RecordType     string  `env:"PLUGIN_RECORD_TYPE,required"`
	RecordName     string  `env:"PLUGIN_RECORD_NAME,required"`
	RecordContent  string  `env:"PLUGIN_RECORD_CONTENT"`
	RecordProxied  *bool   `env:"PLUGIN_RECORD_PROXIED" envDefault:true`
	RecordTTL      int     `env:"PLUGIN_RECORD_TTL" envDefault:1`
	RecordPriority *uint16 `env:"PLUGIN_RECORD_PRIORITY" envDefault:1`
}

func handleCloudflareError(err error) {
	if err != nil {
		fmt.Printf("err: %s\n", err.Error())
		os.Exit(3)
	}
}

func Run() {
	// Parse environmental variables
	input := Input{}
	if err := env.Parse(&input); err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
	// Normalize input
	input.Action = strings.ToLower(input.Action)
	input.RecordType = strings.ToUpper(input.RecordType)
	// Semantic checks
	if input.Action != "set" && input.Action != "unset" {
		fmt.Printf("err: PLUGIN_ACTION must equal one of [set,unset]\n")
		os.Exit(2)
	}
	if input.Action == "set" && input.RecordContent == "" {
		fmt.Printf("err: PLUGIN_RECORD_CONTENT must be set with 'set' action\n")
		os.Exit(2)
	}
	// Create client for cloudflare api communication
	api, err := cloudflare.NewWithAPIToken(input.ApiToken)
	handleCloudflareError(err)
	// Print input if debug mode is on
	if input.Debug {
		fmt.Printf("Input: %+v\n", input)
	}
	// Initialize all records that will be used
	searchRecord := cloudflare.ListDNSRecordsParams{
		Type: input.RecordType,
		Name: input.RecordName,
	}
	createRecord := cloudflare.CreateDNSRecordParams{
		Type:     input.RecordType,
		Name:     input.RecordName,
		Content:  input.RecordContent,
		Proxied:  input.RecordProxied,
		TTL:      input.RecordTTL,
		Priority: input.RecordPriority,
	}
	// Check to see if record was found
	found, _, err := api.ListDNSRecords(context.Background(), cloudflare.ZoneIdentifier(input.ZoneIdentifier), searchRecord)
	handleCloudflareError(err)
	// Switch based on action
	if input.Action == "unset" && len(found) > 0 {
		id := found[0].ID
		err := api.DeleteDNSRecord(context.Background(), cloudflare.ZoneIdentifier(input.ZoneIdentifier), id)
		handleCloudflareError(err)
	} else if input.Action == "set" && len(found) > 0 {
		foundRecord := found[0]
		updateRecord := cloudflare.UpdateDNSRecordParams{
			ID:      foundRecord.ID,
			Type:    foundRecord.Type,
			Name:    foundRecord.Name,
			Content: input.RecordContent,
		}
		if os.Getenv("PLUGIN_RECORD_PROXIED") != "" {
			updateRecord.Proxied = input.RecordProxied
		}
		if os.Getenv("PLUGIN_RECORD_TTL") != "" {
			updateRecord.TTL = input.RecordTTL
		}
		if os.Getenv("PLUGIN_RECORD_PRIORITY") != "" {
			updateRecord.Priority = input.RecordPriority
		}
		err := api.UpdateDNSRecord(context.Background(), cloudflare.ZoneIdentifier(input.ZoneIdentifier), updateRecord)
		handleCloudflareError(err)
	} else if input.Action == "set" {
		res, err := api.CreateDNSRecord(context.Background(), cloudflare.ZoneIdentifier(input.ZoneIdentifier), createRecord)
		handleCloudflareError(err)
		if input.Debug {
			fmt.Printf("CreateDNSRecord: %+v\n", res)
		}
	} else {
		fmt.Printf("err: could not find record to unset\n")
		os.Exit(4)
	}
}
