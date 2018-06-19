package main

import (
	"errors"
	"fmt"
	"github.com/quynhdang-vt/vt-engine-tools/model"
	"github.com/spf13/cobra"
)

func main() {
	var outputFile string
	var idTag string
	var cmdMerge = &cobra.Command{
		Use:   "merge-vlf [vlf_file1] [vlf_file2] [vlf_file3] ...",
		Short: "merge-vlf vlf_file1, vlf_file2, vlf_file3, ...",
		Long: "merge-vlf:  merging VLF transcriptions as output per channels of a multi-channel audio file into a single engine output. " +
			" Each file will be tagged as `id:$idTag_N` where N is the order of the file in the arguments, and" +
			" idTag, default `channel`,  can be configured per the --idTag parameter." +
			"  The output filename, default `output.json`, can be configured with the --output parameter.",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("requires at least one file")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			transcriptionOutput := model.NewEngineOutput()
			for i, v := range args {
				vlf, err := model.ParseVLF(v)
				if err == nil {
					transcriptionOutput.AddVLF(fmt.Sprintf("%s_%d", idTag, i), vlf)
				}
			}
			transcriptionOutput.Sort()
			// output to outputFile
			transcriptionOutput.Write(outputFile)
		},
	}
	var rootCmd = &cobra.Command{Use: "vt-engine-tools"}
	cmdMerge.Flags().StringVarP(&outputFile, "output", "o", "output.json", "Output filename")
	cmdMerge.Flags().StringVarP(&idTag, "idTag", "t", "channel", "Tag prefix for utterance id")

	rootCmd.AddCommand(cmdMerge)
	rootCmd.Execute()
}
