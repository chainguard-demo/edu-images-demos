/*
Copyright 2025 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

/*
Package main demonstrates how to use the Chainguard API client.

This example shows how to:
1. Create and authenticate with the Chainguard API
2. List repositories with optional filtering
3. Apply image customizations using a build.yaml file
4. List and monitor build reports

The example uses the Chainguard SDK to interact with the API and
demonstrates proper patterns for authentication, error handling,
and resource management.
*/
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"chainguard.dev/sdk/auth"
	"chainguard.dev/sdk/proto/platform"
	commonv1 "chainguard.dev/sdk/proto/platform/common/v1"
	iamv1 "chainguard.dev/sdk/proto/platform/iam/v1"
	registryv1 "chainguard.dev/sdk/proto/platform/registry/v1"

	"gopkg.in/yaml.v2"
)

const (
	// API configuration
	defaultAPIURL    = "https://console-api.enforce.dev"
	tokenEnvVariable = "TOK"

	// Group and repository settings
	defaultGroupName = "ORGANIZATION"
	demoRepoName     = "CUSTOM-IMAGE-NAME"

	// File paths
	buildConfigFile = "build.yaml"
)

// listRepositories lists repositories in a group with optional name filtering.
//
// This function demonstrates how to use the Registry client to list repositories
// and how to use filters to narrow down the results. The Chainguard API uses
// UIDP (Unique IDentifier Path) for hierarchical resource navigation, and this
// function shows the proper pattern for filtering by both group ID and name.
//
// Parameters:
//   - ctx: The context for the API call, which can include timeouts, cancellation signals, etc.
//   - clients: The platform clients object that contains all API clients
//   - groupID: The ID of the group containing the repositories to list
//   - repoName: Optional repository name to filter by (empty string lists all repositories)
//
// Returns:
//   - A slice of repository objects matching the filter criteria
//   - An error if the API call fails
func listRepositories(ctx context.Context, clients platform.Clients, groupID string, repoName string) ([]*registryv1.Repo, error) {
	// Create the base filter with the group ID using the UIDPFilter
	// UIDPFilter.DescendantsOf finds all resources that are descendants of the specified ID
	// This is how we navigate the hierarchy of resources in the Chainguard API
	repoFilter := &registryv1.RepoFilter{
		Uidp: &commonv1.UIDPFilter{
			DescendantsOf: groupID,
		},
	}

	// Add name filter if provided
	// This is an optional additional filter that can be applied
	// When combined with the UIDP filter, it will find repositories with the given name
	// that are also descendants of the specified group
	if repoName != "" {
		repoFilter.Name = repoName
	}

	// Make the API call through the Registry client
	// The Registry client provides access to repository-related operations
	// ListRepos returns a RepoList object that contains the Items field with the actual repositories
	repoResp, err := clients.Registry().Registry().ListRepos(ctx, repoFilter)
	if err != nil {
		return nil, fmt.Errorf("failed to list repositories: %w", err)
	}

	return repoResp.Items, nil
}

// listBuildReports lists build reports for a repository.
//
// Build reports provide information about image builds, including status,
// timestamps, and digests of the resulting images. This function shows how to
// query build reports for a specific repository.
//
// Parameters:
//   - ctx: The context for the API call
//   - clients: The platform clients object
//   - repoID: The ID of the repository to list build reports for
//
// Returns:
//   - A slice of build report objects for the repository
//   - An error if the API call fails
func listBuildReports(ctx context.Context, clients platform.Clients, repoID string) ([]*registryv1.BuildReport, error) {
	// Create the filter for the repository using the UIDPFilter
	// Similar to repositories, build reports are hierarchical resources
	// Build reports are children of repositories in the resource hierarchy
	buildReportFilter := &registryv1.BuildReportFilter{
		Uidp: &commonv1.UIDPFilter{
			DescendantsOf: repoID,
		},
	}

	// Make the API call through the Registry client
	// ListBuildReports returns information about image builds for the specified repository
	buildReportResp, err := clients.Registry().Registry().ListBuildReports(ctx, buildReportFilter)
	if err != nil {
		return nil, fmt.Errorf("failed to list build reports: %w", err)
	}

	// Return the Reports field from the response, which contains the actual build reports
	return buildReportResp.Reports, nil
}

// printBuildReports prints build reports in a formatted way.
//
// This helper function displays build reports in a user-friendly format,
// showing the start time, result status, and image digest for each report.
//
// Parameters:
//   - reports: A slice of build report objects to display
func printBuildReports(reports []*registryv1.BuildReport) {
	if len(reports) == 0 {
		fmt.Println("No build reports found")
	} else {
		for _, report := range reports {
			// Format the timestamp in a human-readable format
			startTime := report.StartedAt.AsTime().Format(time.RFC1123)

			// Get the build result status (Success, Failure, Unknown, etc.)
			result := report.Result.String()

			// Get the image digest (if available)
			digest := report.Digest

			// Print the formatted report line
			fmt.Printf("- Started: %s, Result: %s, Digest: %s\n", startTime, result, digest)
		}
	}
}

// applyCustomization applies a custom package overlay to a repository using the specified yaml file.
//
// This function demonstrates the pattern for updating a repository with a custom
// overlay that defines the packages to include in the image. The overlay is applied
// to the repository, which triggers a new build.
//
// The build.yaml file should have the following format:
//
//	contents:
//	  packages:
//	    - package1
//	    - package2
//
// Parameters:
//   - ctx: The context for the API call
//   - clients: The platform clients object
//   - repoID: The ID of the repository to update
//   - configFile: The path to the YAML file containing the package configuration
//
// Returns:
//   - A string indicating the result of the operation
//   - An error if the operation fails
func applyCustomization(ctx context.Context, clients platform.Clients, repoID string, configFile string) (string, error) {
	// Read the build config file from disk
	yamlData, err := os.ReadFile(configFile)
	if err != nil {
		return "", fmt.Errorf("failed to read config file %s: %w", configFile, err)
	}

	// Define structs that match the expected format of the build.yaml file
	// These structs are used to parse the YAML into Go objects
	type ImageContents struct {
		Packages []string `yaml:"packages"`
	}

	type ImageConfiguration struct {
		Contents ImageContents `yaml:"contents"`
	}

	// Parse the YAML into the ImageConfiguration struct
	var config ImageConfiguration
	if err := yaml.Unmarshal(yamlData, &config); err != nil {
		return "", fmt.Errorf("failed to parse YAML: %w", err)
	}

	// Find the repository by ID
	// We use ListRepos with a filter on the ID since there's no direct GetRepo method
	repoFilter := &registryv1.RepoFilter{
		Id: repoID,
	}
	repoResp, err := clients.Registry().Registry().ListRepos(ctx, repoFilter)
	if err != nil {
		return "", fmt.Errorf("failed to find repository: %w", err)
	}

	// Ensure the repository was found
	if len(repoResp.Items) == 0 {
		return "", fmt.Errorf("repository with ID %s not found", repoID)
	}

	// Get the first (and should be only) repository from the response
	repo := repoResp.Items[0]

	// Create a CustomOverlay object with the packages from the config
	// The CustomOverlay is what defines the packages to include in the image
	newOverlay := &registryv1.CustomOverlay{
		Contents: &registryv1.ImageContents{
			Packages: config.Contents.Packages,
		},
	}

	// Assign the overlay to the repository
	// This updates the repository's CustomOverlay field
	repo.CustomOverlay = newOverlay

	// Update the repository using the UpdateRepo method
	// This sends the modified repository back to the API
	// The API will then trigger a new build with the updated configuration
	_, err = clients.Registry().Registry().UpdateRepo(ctx, repo)
	if err != nil {
		return "", fmt.Errorf("failed to update repository: %w", err)
	}

	// Return a success message
	// Note: The build is asynchronous, so we don't have a build report ID yet
	// You can use listBuildReports to monitor for new build reports after this call
	return "Repository updated successfully", nil
}

// createClient creates a new Chainguard API client using a token from the environment.
//
// This function demonstrates how to create and authenticate a client for the
// Chainguard API. It shows the pattern for:
// 1. Getting an authentication token from the environment
// 2. Creating credentials from the token
// 3. Initializing the platform clients
//
// The returned clients object provides access to all Chainguard API services.
//
// Parameters:
//   - ctx: The context for the API client
//
// Returns:
//   - A platform.Clients object for making authenticated API calls
//   - An error if client creation fails
func createClient(ctx context.Context) (platform.Clients, error) {
	// Get token from environment variable
	// This shows how to use environment variables for secure token handling
	// Common sources include:
	// - TOK environment variable (from chainctl auth token)
	// - User-provided token
	// - Service account credentials
	token := os.Getenv(tokenEnvVariable)
	if token == "" {
		return nil, fmt.Errorf("no token found in environment variable %s", tokenEnvVariable)
	}

	// Create credentials from the token
	// The auth.NewFromToken function creates a credential provider from a token
	// - The first argument is the context
	// - The second argument is the token with "Bearer " prefix
	// - The third argument indicates whether to disable verification (false = verify)
	cred := auth.NewFromToken(ctx, fmt.Sprintf("Bearer %s", token), false)

	// Create the platform client
	// platform.NewPlatformClients initializes clients for all Chainguard API services
	// - The first argument is the context
	// - The second argument is the API URL
	// - The third argument is the credential provider
	clients, err := platform.NewPlatformClients(ctx, defaultAPIURL, cred)
	if err != nil {
		return nil, fmt.Errorf("failed to create platform clients: %w", err)
	}

	// Return the clients object, which provides access to all API services
	return clients, nil
}

// confirmAction handles user confirmation for potentially destructive actions.
//
// This utility function demonstrates a pattern for getting user confirmation
// before performing potentially destructive operations. It supports both
// interactive confirmation and automatic confirmation via a flag.
//
// Parameters:
//   - assumeYes: If true, skips the prompt and automatically confirms
//   - prompt: The message to display to the user
//
// Returns:
//   - true if the action is confirmed, false otherwise
func confirmAction(assumeYes bool, prompt string) bool {
	// If assumeYes is true, skip the prompt and return true
	// This provides a non-interactive mode for automation
	if assumeYes {
		fmt.Println("Automatically confirming with --yes flag")
		return true
	}

	// Display a prompt to the user and wait for input
	fmt.Printf("%s (y/n): ", prompt)
	var answer string
	if _, err := fmt.Scanln(&answer); err != nil {
		log.Printf("Failed to read input: %v", err)
		return false
	}

	// Return true only if the user typed 'y' or 'Y'
	return strings.ToLower(answer) == "y"
}

// main is the entry point for the example.
//
// This function demonstrates the complete workflow for using the Chainguard API:
// 1. Creating and authenticating an API client
// 2. Listing repositories with optional filtering
// 3. Applying image customizations using a build.yaml file
// 4. Monitoring build reports
//
// The example accepts a --yes flag to automatically confirm actions,
// which is useful for automation.
func main() {
	// Parse command line flags
	// The --yes flag allows for non-interactive operation
	assumeYes := flag.Bool("yes", false, "Assume yes for confirmation prompts")
	flag.Parse()

	// Create a background context for the API calls
	// You can also use contexts with timeouts, deadlines, or cancellation
	ctx := context.Background()

	// SECTION 1: CREATE API CLIENT
	// Create the authenticated API client using the createClient function
	clients, err := createClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create API client: %v", err)
	}
	// Ensure the client is closed when the program exits
	defer clients.Close()

	// SECTION 2: LIST REPOSITORIES
	// First, get the group ID by name
	groupFilter := &iamv1.GroupFilter{
		Name: defaultGroupName,
	}
	groupResp, err := clients.IAM().Groups().List(ctx, groupFilter)
	if err != nil {
		log.Fatalf("Failed to list groups: %v", err)
	}

	// Find the group ID from the response
	var groupID string
	if len(groupResp.Items) > 0 {
		for _, group := range groupResp.Items {
			if group.Name == defaultGroupName {
				groupID = group.Id
				break
			}
		}
	} else {
		log.Fatalf("Group %s not found", defaultGroupName)
	}

	// Display the group information
	fmt.Printf("Group: %s (ID: %s)\n", defaultGroupName, groupID)

	// List all repositories in the group using the listRepositories function
	// This demonstrates listing repositories without any name filtering
	fmt.Printf("\nAll repositories in %s:\n", defaultGroupName)
	allRepos, err := listRepositories(ctx, clients, groupID, "")
	if err != nil {
		log.Fatalf("Failed to list repositories: %v", err)
	}

	// Display all repositories
	for _, repo := range allRepos {
		fmt.Printf("- %s\n", repo.Name)
	}

	// Get a specific repository by name using the same function with a name filter
	// This demonstrates filtering repositories by name
	demoRepos, err := listRepositories(ctx, clients, groupID, demoRepoName)
	if err != nil {
		log.Fatalf("Failed to find %s repository: %v", demoRepoName, err)
	}

	// Ensure the repository was found
	if len(demoRepos) == 0 {
		log.Fatalf("Repository %s not found", demoRepoName)
	}

	// Get the repository object
	demoRepo := demoRepos[0]
	fmt.Printf("\nRepository: %s (ID: %s)\n", demoRepo.Name, demoRepo.Id)

	// SECTION 3: LIST BUILD REPORTS
	// List existing build reports for the repository using the listBuildReports function
	// This demonstrates querying build history for a repository
	fmt.Printf("\nBuild Reports for %s repository:\n", demoRepoName)
	reports, err := listBuildReports(ctx, clients, demoRepo.Id)
	if err != nil {
		log.Fatalf("Failed to list build reports: %v", err)
	}

	// Display the build reports in a formatted way
	printBuildReports(reports)

	// SECTION 4: APPLY CUSTOMIZATION
	// Prepare to apply customization using the build.yaml file
	fmt.Printf("\nAbout to apply customization using configuration file: %s\n", buildConfigFile)

	// Ask for confirmation before proceeding with the customization
	// This demonstrates the pattern for getting user confirmation
	if !confirmAction(*assumeYes, fmt.Sprintf("Are you sure you want to update repository %s?", demoRepo.Name)) {
		fmt.Println("Customization cancelled.")
		return
	}

	// Apply the customization using the applyCustomization function
	// This demonstrates the pattern for updating a repository with a custom overlay
	fmt.Println("Applying customization...")
	result, err := applyCustomization(ctx, clients, demoRepo.Id, buildConfigFile)
	if err != nil {
		log.Fatalf("Failed to apply customization: %v", err)
	}
	fmt.Printf("Customization result: %s\n", result)

	// SECTION 5: MONITOR BUILD PROGRESS
	// List updated build reports to see the new build
	// This demonstrates monitoring build progress after submission
	fmt.Printf("\nUpdated Build Reports for %s repository:\n", demoRepoName)
	updatedReports, err := listBuildReports(ctx, clients, demoRepo.Id)
	if err != nil {
		log.Fatalf("Failed to list updated build reports: %v", err)
	}
	printBuildReports(updatedReports)
}
