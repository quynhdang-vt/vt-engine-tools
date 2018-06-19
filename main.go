package main

import (
	"errors"
	"fmt"
	"github.com/quynhdang-vt/vt-engine-tools/model"
	"github.com/spf13/cobra"
)

func main() {
	var outputFile string
	var cmdMerge = &cobra.Command{
		Use:   "merge-vlf [vlf_file1] [vlf_file2] [vlf_file3] ...",
		Short: "merge-vlf vlf_file1, vlf_file2, vlf_file3, ...",
		Long: "merge-vlf VLF_files for merging -- the files are assumed to be channels of an audio file.  Each file will be assigned a Speaker_# where # is the order that the file appeared in the arguments.  Specify -output for the resulting output file",
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
					transcriptionOutput.AddVLF(fmt.Sprintf("Speaker_%d", i), vlf)
				}
			}
			transcriptionOutput.Sort()
			// output to outputFile
			transcriptionOutput.Write(outputFile)
		},
	}
	var rootCmd = &cobra.Command{Use: "vt-speaker"}
	cmdMerge.Flags().StringVarP(&outputFile, "output", "o", "output.json", "Output filename")

	rootCmd.AddCommand(cmdMerge)
	rootCmd.Execute()
}
