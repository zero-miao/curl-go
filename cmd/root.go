package cmd

import "github.com/spf13/cobra"

// ==========================
// @Author : zero-miao
// @Date   : 2019-08-22 13:55
// @File   : root.go
// @Project: curl-go/entry
// ==========================

var rootCmd *cobra.Command

func Entry() error {
	rootCmd = manual()
	autoCmd := auto()
	rootCmd.AddCommand(autoCmd)
	return rootCmd.Execute()
}
