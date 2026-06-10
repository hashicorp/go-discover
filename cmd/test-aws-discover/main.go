// Test program for AWS discover with dual-stack endpoint support.
// Use this to verify the fix works in regions with and without dual-stack endpoints.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	discover "github.com/hashicorp/go-discover"
	_ "github.com/hashicorp/go-discover/provider/aws"
)

func main() {
	region := flag.String("region", "", "AWS region (required)")
	tagKey := flag.String("tag-key", "", "Tag key to search for (required)")
	tagValue := flag.String("tag-value", "", "Tag value to search for (required)")
	addrType := flag.String("addr-type", "private_v4", "Address type: private_v4, public_v4, public_v6")
	flag.Parse()

	if *region == "" || *tagKey == "" || *tagValue == "" {
		fmt.Fprintf(os.Stderr, "Usage: %s -region <region> -tag-key <key> -tag-value <value> [-addr-type <type>]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nRequired flags:\n")
		fmt.Fprintf(os.Stderr, "  -region     AWS region (e.g., us-east-1, me-central-1)\n")
		fmt.Fprintf(os.Stderr, "  -tag-key    Tag key to search for\n")
		fmt.Fprintf(os.Stderr, "  -tag-value  Tag value to search for\n")
		fmt.Fprintf(os.Stderr, "\nOptional flags:\n")
		fmt.Fprintf(os.Stderr, "  -addr-type  Address type: private_v4 (default), public_v4, public_v6\n")
		fmt.Fprintf(os.Stderr, "\nEnvironment variables:\n")
		fmt.Fprintf(os.Stderr, "  AWS_USE_DUALSTACK_ENDPOINT  Set to 'false' to disable dual-stack endpoints\n")
		fmt.Fprintf(os.Stderr, "  AWS_REGION                  Default region (overridden by -region flag)\n")
		fmt.Fprintf(os.Stderr, "  AWS_ACCESS_KEY_ID           AWS access key (optional, uses default chain)\n")
		fmt.Fprintf(os.Stderr, "  AWS_SECRET_ACCESS_KEY       AWS secret key (optional, uses default chain)\n")
		os.Exit(1)
	}

	// Show current configuration
	dualStackEnv := os.Getenv("AWS_USE_DUALSTACK_ENDPOINT")
	if dualStackEnv == "" {
		dualStackEnv = "(not set - dual-stack enabled by default)"
	}
	fmt.Printf("Configuration:\n")
	fmt.Printf("  Region:                       %s\n", *region)
	fmt.Printf("  Tag:                          %s=%s\n", *tagKey, *tagValue)
	fmt.Printf("  Address type:                 %s\n", *addrType)
	fmt.Printf("  AWS_USE_DUALSTACK_ENDPOINT:   %s\n", dualStackEnv)
	fmt.Println()

	// Build discover config
	cfg := discover.Config{
		"provider":  "aws",
		"region":    *region,
		"tag_key":   *tagKey,
		"tag_value": *tagValue,
		"addr_type": *addrType,
	}

	// Create discoverer
	d, err := discover.New()
	if err != nil {
		log.Fatalf("Failed to create discoverer: %v", err)
	}

	// Create logger that writes to stderr
	l := log.New(os.Stderr, "", log.LstdFlags)

	fmt.Printf("Discovering instances...\n\n")

	// Perform discovery
	addrs, err := d.Addrs(cfg.String(), l)
	if err != nil {
		log.Fatalf("Discovery failed: %v", err)
	}

	// Print results
	if len(addrs) == 0 {
		fmt.Printf("No instances found with tag %s=%s in region %s\n", *tagKey, *tagValue, *region)
	} else {
		fmt.Printf("Found %d instance(s):\n", len(addrs))
		for i, addr := range addrs {
			fmt.Printf("  %d. %s\n", i+1, addr)
		}
	}
}
