package cmd

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/beringresearch/macpine/host"
	qemu "github.com/beringresearch/macpine/qemu"
	"github.com/beringresearch/macpine/utils"
	"github.com/spf13/cobra"
)

// publishCmd stops an Alpine instance
var publishCmd = &cobra.Command{
	Use:   "publish NAME",
	Short: "Publish an Alpine VM.",
	Run:   publish,
}

func publish(cmd *cobra.Command, args []string) {

	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	if len(args) == 0 {
		log.Fatal("missing name - please provide VM name")
		return
	}

	machineConfig := qemu.MachineConfig{
		Alias: args[0],
	}
	machineConfig.Location = filepath.Join(userHomeDir, ".macpine", machineConfig.Alias)

	err = host.Stop(machineConfig)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Second)

	fileInfo, err := ioutil.ReadDir(machineConfig.Location)
	if err != nil {
		log.Fatal(err)
	}

	files := []string{}
	for _, f := range fileInfo {
		files = append(files, filepath.Join(machineConfig.Location, f.Name()))
	}

	out, err := os.Create(machineConfig.Alias + ".tar.gz")
	if err != nil {
		log.Fatalln("Error writing archive:", err)
	}
	defer out.Close()

	// Create the archive and write the output to the "out" Writer
	err = utils.Compress(files, out)
	if err != nil {
		log.Fatalln("error creating archive:", err)
	}

	err = host.Start(machineConfig)
	if err != nil {
		log.Fatal(err)
	}

}