package main

import (
	"fmt"
	"log"

	"github.com/osbuild/images/pkg/distrofactory"
	"github.com/osbuild/images/pkg/imagefilter"
	testrepos "github.com/osbuild/images/test/data/repositories"
)

func main() {
	fac := distrofactory.NewDefault()
	repos, err := testrepos.New()
	if err != nil {
		log.Fatalf("Failed to create repos: %v", err)
	}

	filter, err := imagefilter.New(fac, repos)
	if err != nil {
		log.Fatalf("Failed to create filter: %v", err)
	}

	// Test filtering by guest-image alias
	results, err := filter.Filter("type:guest-image")
	if err != nil {
		log.Fatalf("Failed to filter: %v", err)
	}

	fmt.Printf("Found %d results for type:guest-image\n", len(results))
	for _, result := range results {
		fmt.Printf("- %s/%s/%s\n", result.Distro.Name(), result.Arch.Name(), result.ImgType.Name())
		fmt.Printf("  Aliases: %v\n", result.ImgType.Aliases())
	}

	// Test filtering by qcow2 for comparison
	results2, err := filter.Filter("type:qcow2")
	if err != nil {
		log.Fatalf("Failed to filter: %v", err)
	}

	fmt.Printf("\nFound %d results for type:qcow2\n", len(results2))
	for _, result := range results2 {
		fmt.Printf("- %s/%s/%s\n", result.Distro.Name(), result.Arch.Name(), result.ImgType.Name())
		fmt.Printf("  Aliases: %v\n", result.ImgType.Aliases())
	}

	// Verify that guest-image and qcow2 return the same results for Fedora
	guestImageFedoraCount := 0
	qcow2FedoraCount := 0

	for _, result := range results {
		if result.Distro.Name() == "fedora-42" || result.Distro.Name() == "fedora-43" {
			guestImageFedoraCount++
		}
	}

	for _, result := range results2 {
		if result.Distro.Name() == "fedora-42" || result.Distro.Name() == "fedora-43" {
			qcow2FedoraCount++
		}
	}

	fmt.Printf("\nFedora results: guest-image=%d, qcow2=%d\n", guestImageFedoraCount, qcow2FedoraCount)
	if guestImageFedoraCount > 0 && guestImageFedoraCount == qcow2FedoraCount {
		fmt.Println("✅ guest-image alias works correctly!")
	} else {
		fmt.Println("❌ guest-image alias not working as expected")
	}
}
