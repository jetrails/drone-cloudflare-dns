# Drone — Cloudflare DNS
> Drone plugin to control simple DNS records via Cloudflare's API

![](https://img.shields.io/badge/License-MIT-lightgray.svg?style=for-the-badge)
![](https://img.shields.io/docker/stars/jetrails/drone-cloudflare-dns.svg?style=for-the-badge&colorB=9f9f9f)
![](https://img.shields.io/docker/pulls/jetrails/drone-cloudflare-dns.svg?style=for-the-badge&colorB=9f9f9f)

## About

Our Drone plugin enables the ability for your pipeline to interface with Cloudflare's API to create/update/delete DNS records. This plugin is written in Go and it uses the [cloudflare-go](https://github.com/cloudflare/cloudflare-go) package to communicate with Cloudflare's API. For information on Cloudflare's API please refer to their [documentation](https://api.cloudflare.com/#dns-records-for-a-zone-properties) page.

## Cloudflare Token

The API token that is used to authenticate with Cloudflare's API can be created in Cloudflare's dashboard. It is recommended to create an API token that includes only the zone resource you want to manipulate and give edit permissions only to the DNS resource.

## Build

Develop locally by running the plugin with the following commands. Also please note that you should specify environmental variables for the plugin either via the inline method (`FOO=bar go run src/main`) or via the `export FOO=bar` method before the `go run` command.

```shell
$ go mod download
$ go run cmd/main/main.go
```

## Docker

Drone plugins work off of docker images. The following commands will go over building, pushing, and running the docker image for this plugin.

###### Build Docker Image:

```shell
$ docker build -t jetrails/drone-cloudflare-dns .
```

###### Run Docker Container:

You can then replicate the command that Drone will use to launch the plugin by running:

```shell
$ docker run --rm \
	-e PLUGIN_API_TOKEN="u4C7ev06GMS8_vWBTpjqtVReT3I7FwGpW7MG44ZD" \
	-e PLUGIN_ZONE_IDENTIFIER="eJzrjE44s6Ki67x1tSDJzI8LdXxM3nj7" \
	-e PLUGIN_ACTION="set" \
	-e PLUGIN_RECORD_TYPE="cname" \
	-e PLUGIN_RECORD_NAME="example.com" \
	-e PLUGIN_RECORD_CONTENT="google.com" \
	-e PLUGIN_RECORD_PROXIED="true" \
	-e PLUGIN_RECORD_TTL="3600" \
	-v $(pwd):/drone/src \
	-w /drone/src \
	jetrails/drone-cloudflare-dns
```

###### Push Docker Image:

Finally, push this image to our Docker Hub [repository](https://hub.docker.com/r/jetrails/drone-cloudflare-dns) (assuming you have permission):

```shell
$ docker push jetrails/drone-cloudflare-dns
```

## Usage

There are two values for `action`, set and unset. If set is choosen and a record does not already exist, then the record is created. If set is choosen and the record already exists, then the record is updated. It is recommended that the `record_name` value be always set to the FQDN. Please refer to the table below with all possible settings that can be passed to the plugin:

|       Name      |    Required   | Default | Case-Sensitive |        Type       |
|:---------------:|:-------------:|:-------:|:--------------:|:-----------------:|
|    api_token    |      Yes      |    -    |       Yes      |       STRING      |
| zone_identifier |      Yes      |    -    |       Yes      |       STRING      |
|      action     |      Yes      |    -    |       No       |  ENUM[set,unset]  |
|   record_type   |      Yes      |    -    |       No       | ENUM[a,cname,...] |
|   record_name   |      Yes      |    -    |       Yes      |       STRING      |
|  record_content | action == set |    -    |       Yes      |       STRING      |
|  record_proxied |       No      |   true  |       N/A      |       STRING      |
|    record_ttl   |       No      |    1    |       N/A      |        INT        |
| record_priority |       No      |    1    |       N/A      |        INT        |
|      debug      |       No      |  false  |       N/A      |       STRING      |

## Examples

```yaml
kind: pipeline
name: default

steps:
-   name: cloudflare
    image: jetrails/drone-cloudflare-dns
    settings:
        api_token:
            from_secret: cloudflare_token
        zone_identifier:
            from_secret: cloudflare_zone_identifier
        action: set
        record_type: a
        record_name: example.com
        record_content: 127.0.0.1
        record_proxied: false
```

```yaml
kind: pipeline
name: default

steps:
-   name: cloudflare
    image: jetrails/drone-cloudflare-dns
    settings:
        api_token:
            from_secret: cloudflare_token
        zone_identifier:
            from_secret: cloudflare_zone_identifier
        action: unset
        record_type: a
        record_name: example.com
```

## Feature Requests / Issues

Feel free to open an issue for any feature requests and issues that you may come across. For furthur inquery, please contact [development@jetrails.com](mailto://development@jetrails.com).
