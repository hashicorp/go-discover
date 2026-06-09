// Copyright IBM Corp. 2017, 2026
// SPDX-License-Identifier: MPL-2.0

// Package gce provides node discovery for Google Cloud.
package gce

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	compute "google.golang.org/api/compute/v1"
)

type Provider struct {
	userAgent string
}

func (p *Provider) SetUserAgent(s string) {
	p.userAgent = s
}

func (p *Provider) Help() string {
	return `Google Cloud:

    provider:         "gce"
    project_name:     The name of the project. discovered if not set
    tag_value:        The tag value to filter on.
    label_key:        The label key to filter on. Required if label_value is set.
    label_value:      The label value to filter on. Required if label_key is set.

    Tag and label filters can be used independently or combined.
    When combined, only instances matching both filters are returned.

    zone_pattern:     A RE2 regular expression for filtering zones, e.g. us-west1-.*, or us-(?west|east).*
    credentials_file: The path to the credentials file. See below for more details

    The credentials for a GCE Service Account are required and are searched in
    the following locations:

    1. Use credentials from "credentials_file", if provided.
    2. Use JSON file from GOOGLE_APPLICATION_CREDENTIALS environment variable.
    3. Use JSON file in a location known to the gcloud command-line tool.
      On Windows, this is %APPDATA%/gcloud/application_default_credentials.json.
      On other systems, $HOME/.config/gcloud/application_default_credentials.json.
    4. On Google Compute Engine, use credentials from the metadata
      server. In this final case any provided scopes are ignored.
`
}

func (p *Provider) Addrs(args map[string]string, l *log.Logger) ([]string, error) {
	if args["provider"] != "gce" {
		return nil, fmt.Errorf("discover-gce: invalid provider %s", args["provider"])
	}

	if l == nil {
		l = log.New(io.Discard, "", 0)
	}

	project := args["project_name"]
	zone := args["zone_pattern"]
	creds := args["credentials_file"]
	tagValue := args["tag_value"]
	labelKey := args["label_key"]
	labelValue := args["label_value"]

	filter, err := buildFilter(tagValue, labelKey, labelValue)
	if err != nil {
		return nil, err
	}
	if filter == "" {
		l.Printf("[INFO] discover-gce: No tag or label filter configured")
		return nil, nil
	}

	// determine the project name
	if project == "" {
		l.Println("[INFO] discover-gce: Looking up project name")
		p, err := lookupProject()
		if err != nil {
			return nil, fmt.Errorf("discover-gce: %s", err)
		}
		project = p
	}
	l.Printf("[INFO] discover-gce: Project name is %q", project)

	// create an authenticated client
	if creds != "" {
		l.Printf("[INFO] discover-gce: Loading credentials from %s", creds)
	}
	client, err := client(creds)
	if err != nil {
		return nil, fmt.Errorf("discover-gce: %s", err)
	}
	svc, err := compute.New(client)
	if err != nil {
		return nil, fmt.Errorf("discover-gce: %s", err)
	}
	if p.userAgent != "" {
		svc.UserAgent = p.userAgent
	}

	// lookup the project zones to look in
	if zone != "" {
		l.Printf("[INFO] discover-gce: Looking up zones matching %s", zone)
	} else {
		l.Printf("[INFO] discover-gce: Looking up all zones")
	}
	zones, err := lookupZones(svc, project, zone)
	if err != nil {
		return nil, fmt.Errorf("discover-gce: %s", err)
	}
	l.Printf("[INFO] discover-gce: Found zones %v", zones)

	// lookup the instance addresses across all zones
	var addrs []string
	for _, zone := range zones {
		a, err := lookupAddrsByFilter(svc, project, zone, filter)
		if err != nil {
			return nil, fmt.Errorf("discover-gce: %s", err)
		}
		l.Printf("[INFO] discover-gce: Zone %q has matches: %v", zone, a)
		addrs = append(addrs, a...)
	}
	return addrs, nil
}

func buildFilter(tagValue, labelKey, labelValue string) (string, error) {
	if tagValue == "" && labelKey == "" {
		return "", fmt.Errorf("discover-gce: tag_value or label_key must be provided")
	}

	if (labelKey != "" && labelValue == "") || (labelKey == "" && labelValue != "") {
		return "", fmt.Errorf("discover-gce: label_key and label_value must both be set or both be empty")
	}

	if tagValue != "" && labelKey != "" {
		return fmt.Sprintf("(tags.items:\"%s\") AND (labels.%s=\"%s\")", tagValue, labelKey, labelValue), nil
	}
	if tagValue != "" {
		return fmt.Sprintf("tags.items:\"%s\"", tagValue), nil
	}
	if labelKey != "" {
		return fmt.Sprintf("labels.%s=\"%s\"", labelKey, labelValue), nil
	}

	return "", nil
}

// client returns an authenticated HTTP client for use with GCE.
func client(path string) (*http.Client, error) {
	if path == "" {
		return google.DefaultClient(oauth2.NoContext, compute.ComputeScope)
	}

	key, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	jwtConfig, err := google.JWTConfigFromJSON(key, compute.ComputeScope)
	if err != nil {
		return nil, err
	}

	return jwtConfig.Client(oauth2.NoContext), nil
}

// lookupProject retrieves the project name from the metadata of the current node.
func lookupProject() (string, error) {
	req, err := http.NewRequest("GET", "http://metadata.google.internal/computeMetadata/v1/project/project-id", nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Metadata-Flavor", "Google")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("discover-gce: invalid status code %d when fetching project id", resp.StatusCode)
	}

	project, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(project), nil
}

// lookupZones retrieves the zones of the project and filters them by pattern.
func lookupZones(svc *compute.Service, project, pattern string) ([]string, error) {
	call := svc.Zones.List(project)
	if pattern != "" {
		call = call.Filter("name eq " + pattern)
	}

	var zones []string
	f := func(page *compute.ZoneList) error {
		for _, v := range page.Items {
			zones = append(zones, v.Name)
		}
		return nil
	}

	if err := call.Pages(oauth2.NoContext, f); err != nil {
		return nil, err
	}
	return zones, nil
}

// lookupAddrsByFilter retrieves the private ip addresses of all instances in a given
// project and zone which match the provided filter string.
func lookupAddrsByFilter(svc *compute.Service, project, zone, filter string) ([]string, error) {
	var addrs []string
	f := func(page *compute.InstanceList) error {
		for _, v := range page.Items {
			if len(v.NetworkInterfaces) == 0 || v.NetworkInterfaces[0].NetworkIP == "" {
				continue
			}
			addrs = append(addrs, v.NetworkInterfaces[0].NetworkIP)
		}
		return nil
	}

	call := svc.Instances.List(project, zone).Filter(filter)
	if err := call.Pages(oauth2.NoContext, f); err != nil {
		return nil, err
	}
	return addrs, nil
}
