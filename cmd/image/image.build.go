package cmd

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/c12s/cockpit/clients"
	"github.com/c12s/cockpit/constants"
	"github.com/c12s/cockpit/model"
	"github.com/c12s/cockpit/utils"
	"github.com/spf13/cobra"
)

var (
	buildName          string
	buildTag           string
	buildContext       string
	buildContainerfile string
	buildRegistry      string
	buildArgs          []string
	buildWait          bool
	buildTimeout       time.Duration
	buildOutput        string

	buildQcow2Name     string
	buildQcow2ImageRef string
	buildQcow2Arch     string
	buildQcow2RootFS   string
)

var errBuildTimeout = errors.New("build did not finish before the deadline")

var ImageBuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build an OS image in Singularity",
	Long: "Build a bootc OS image from a Containerfile context through Supernova. " +
		"Supernova builds, pushes and signs the image, then returns the digest " +
		"that `cockpit image update` pins. The build goes through the gateway " +
		"(authenticated, same entry point as every other command).",
	Run: executeImageBuild,
}

var ImageBuildQCOW2Cmd = &cobra.Command{
	Use:   "qcow2",
	Short: "Build a qcow2 provisioning disk from an existing OS image",
	Long: "Turn an already-built bootc image into a qcow2 disk, upload it to the " +
		"Singularity artifact server, and print the artifact name and checksum.",
	Run: executeImageBuildQCOW2,
}

var ImageBuildStatusCmd = &cobra.Command{
	Use:   "status <job-id>",
	Short: "Show one build job",
	Args:  cobra.ExactArgs(1),
	Run:   executeImageBuildStatus,
}

var ImageBuildListCmd = &cobra.Command{
	Use:   "list",
	Short: "List build jobs",
	Run:   executeImageBuildList,
}

func executeImageBuild(cmd *cobra.Command, args []string) {
	if buildName == "" {
		buildFail(2, "--name is required (repo path, e.g. os/edge-base)")
	}
	argsMap, err := parseBuildArgs(buildArgs)
	if err != nil {
		buildFail(2, "%v", err)
	}
	tgz, err := utils.TarGzDir(buildContext)
	if err != nil {
		buildFail(2, "packing context %s: %v", buildContext, err)
	}

	req := model.BuildImageRequest{
		Spec: model.BuildSpec{
			Name:          buildName,
			Tag:           buildTag,
			Containerfile: buildContainerfile,
			Registry:      buildRegistry,
			BuildArgs:     argsMap,
		},
		Context: base64.StdEncoding.EncodeToString(tgz),
	}
	var job model.BuildJob
	if err := sendBuildPost(clients.BuildURL("core", "v1", "BuildImage"), req, &job); err != nil {
		buildFail(1, "starting build: %v", err)
	}
	fmt.Printf("Build job %s started (%s:%s, %d bytes context).\n", job.JobID, buildName, buildTag, len(tgz))
	if !buildWait {
		fmt.Printf("Track with: cockpit image build status %s\n", job.JobID)
		return
	}
	waitAndPrintBuild(job.JobID)
}

func executeImageBuildQCOW2(cmd *cobra.Command, args []string) {
	if buildQcow2Name == "" {
		buildFail(2, "--name is required (artifact base name, e.g. edge-base-v3.18)")
	}
	if buildQcow2ImageRef == "" {
		buildFail(2, "--image-ref is required (the already-built bootc image)")
	}
	req := model.BuildQCOW2Request{
		Name:       buildQcow2Name,
		ImageRef:   buildQcow2ImageRef,
		TargetArch: buildQcow2Arch,
		RootFS:     buildQcow2RootFS,
	}
	var job model.BuildJob
	if err := sendBuildPost(clients.BuildURL("core", "v1", "BuildQCOW2"), req, &job); err != nil {
		buildFail(1, "starting qcow2 build: %v", err)
	}
	fmt.Printf("Build job %s started (qcow2 %s from %s).\n", job.JobID, buildQcow2Name, buildQcow2ImageRef)
	if !buildWait {
		fmt.Printf("Track with: cockpit image build status %s\n", job.JobID)
		return
	}
	waitAndPrintBuild(job.JobID)
}

func executeImageBuildStatus(cmd *cobra.Command, args []string) {
	job, err := getBuildJob(args[0])
	if err != nil {
		buildFail(1, "fetching job %s: %v", args[0], err)
	}
	printBuildJob(job)
	if job.Status == "failed" {
		os.Exit(1)
	}
}

func executeImageBuildList(cmd *cobra.Command, args []string) {
	var resp model.BuildJobsResponse
	if err := sendBuildGet(clients.BuildURL("core", "v1", "ListJobs"), &resp); err != nil {
		buildFail(1, "listing jobs: %v", err)
	}
	printBuildJobs(resp.Jobs)
}

// waitAndPrintBuild polls the job until it reaches a terminal state.
func waitAndPrintBuild(jobID string) {
	deadline := time.Now().Add(buildTimeout)
	for {
		job, err := getBuildJob(jobID)
		if err != nil {
			buildFail(1, "polling job %s: %v", jobID, err)
		}
		if job.Status == "succeeded" || job.Status == "failed" {
			printBuildJob(job)
			if job.Status == "failed" {
				os.Exit(1)
			}
			if job.Type == "image" && job.ImageRef != "" && job.Digest != "" {
				fmt.Printf("\nNext: cockpit image update --node-id <node-id> --image-ref %s --digest %s\n", job.ImageRef, job.Digest)
			}
			return
		}
		if time.Now().After(deadline) {
			printBuildJob(job)
			buildFail(3, "%v (job %s still %s)", errBuildTimeout, jobID, job.Status)
		}
		time.Sleep(2 * time.Second)
	}
}

func getBuildJob(id string) (model.BuildJobStatus, error) {
	var job model.BuildJobStatus
	url := strings.Replace(clients.BuildURL("core", "v1", "GetJob"), "{id}", id, 1)
	err := sendBuildGet(url, &job)
	return job, err
}

func sendBuildPost(url string, body, response interface{}) error {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		return fmt.Errorf("reading token: %w", err)
	}
	return utils.SendHTTPRequest(model.HTTPRequestConfig{
		URL:         url,
		Method:      "POST",
		Token:       token,
		RequestBody: body,
		Response:    response,
		Timeout:     60 * time.Second,
	})
}

func sendBuildGet(url string, response interface{}) error {
	token, err := utils.ReadTokenFromFile()
	if err != nil {
		return fmt.Errorf("reading token: %w", err)
	}
	return utils.SendHTTPRequest(model.HTTPRequestConfig{
		URL:      url,
		Method:   "GET",
		Token:    token,
		Response: response,
		Timeout:  30 * time.Second,
	})
}

func printBuildJob(job model.BuildJobStatus) {
	if buildOutput == "json" {
		out, _ := json.MarshalIndent(job, "", "  ")
		fmt.Println(string(out))
		return
	}
	ref := job.ImageRef
	if ref == "" {
		ref = job.ArtifactName
	}
	sum := job.Digest
	if sum == "" {
		sum = job.ArtifactSum
	}
	fmt.Printf("%-16s %-6s %-10s %-50s %s\n", "JOB", "TYPE", "STATUS", "IMAGE/ARTIFACT", "DIGEST/CHECKSUM")
	fmt.Printf("%-16s %-6s %-10s %-50s %s\n", job.ID, job.Type, job.Status, ref, sum)
	if job.Error != "" {
		fmt.Printf("error: %s\n", job.Error)
	}
	if job.Status == "failed" && job.Log != "" {
		fmt.Printf("log tail:\n%s\n", job.Log)
	}
}

func printBuildJobs(jobs []model.BuildJobStatus) {
	if buildOutput == "json" {
		out, _ := json.MarshalIndent(model.BuildJobsResponse{Jobs: jobs}, "", "  ")
		fmt.Println(string(out))
		return
	}
	fmt.Printf("%-16s %-6s %-10s %-50s %s\n", "JOB", "TYPE", "STATUS", "IMAGE/ARTIFACT", "DIGEST/CHECKSUM")
	for _, job := range jobs {
		ref := job.ImageRef
		if ref == "" {
			ref = job.ArtifactName
		}
		sum := job.Digest
		if sum == "" {
			sum = job.ArtifactSum
		}
		fmt.Printf("%-16s %-6s %-10s %-50s %s\n", job.ID, job.Type, job.Status, ref, sum)
	}
}

// parseBuildArgs turns repeated --build-arg key=value flags into a map.
func parseBuildArgs(raw []string) (map[string]string, error) {
	if len(raw) == 0 {
		return nil, nil
	}
	out := make(map[string]string, len(raw))
	for _, item := range raw {
		k, v, ok := strings.Cut(item, "=")
		if !ok || k == "" {
			return nil, fmt.Errorf("invalid --build-arg %q (want key=value)", item)
		}
		out[k] = v
	}
	return out, nil
}

func buildFail(code int, format string, args ...interface{}) {
	fmt.Printf("Error: "+format+"\n", args...)
	os.Exit(code)
}

func init() {
	ImageCmd.AddCommand(ImageBuildCmd)
	ImageBuildCmd.AddCommand(ImageBuildQCOW2Cmd, ImageBuildStatusCmd, ImageBuildListCmd)

	pf := ImageBuildCmd.PersistentFlags()
	pf.BoolVar(&buildWait, constants.WaitFlag, false, constants.WaitDescription)
	pf.DurationVar(&buildTimeout, constants.TimeoutFlag, 30*time.Minute, constants.TimeoutDescription)
	pf.StringVarP(&buildOutput, constants.OutputFlag, constants.OutputShorthandFlag, "table", "Output format: table or json")

	f := ImageBuildCmd.Flags()
	f.StringVar(&buildName, constants.NameFlag, "", "Repo path under the registry, e.g. os/edge-base (required)")
	f.StringVar(&buildTag, constants.TagFlag, "latest", constants.TagDescription)
	f.StringVar(&buildContext, constants.ContextFlag, ".", constants.ContextDescription)
	f.StringVar(&buildContainerfile, constants.ContainerfileFlag, "Containerfile", constants.ContainerfileDescription)
	f.StringVar(&buildRegistry, constants.RegistryFlag, "192.168.122.1:5001", constants.RegistryDescription)
	f.StringArrayVar(&buildArgs, constants.BuildArgFlag, nil, constants.BuildArgDescription)

	qf := ImageBuildQCOW2Cmd.Flags()
	qf.StringVar(&buildQcow2Name, constants.NameFlag, "", "Artifact base name, e.g. edge-base-v3.18 (required)")
	qf.StringVar(&buildQcow2ImageRef, constants.ImageRefFlag, "", "Already-built bootc image reference (required)")
	qf.StringVar(&buildQcow2Arch, constants.TargetArchFlag, "x86_64", constants.TargetArchDescription)
	qf.StringVar(&buildQcow2RootFS, constants.RootFSFlag, "ext4", constants.RootFSDescription)
}
