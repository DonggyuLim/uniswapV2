// /*
// Copyright © 2022 NAME HERE <EMAIL ADDRESS>
// */
package cmd

// import (
// 	"fmt"

// 	"github.com/DonggyuLim/db"
// 	"github.com/DonggyuLim/uniswap/erc20"
// 	"github.com/spf13/cobra"
// )

// // go run main.go deploy --name="USDT" --symbol="USDT" --account="0x123" --totalbalance=100000
// // deployCmd represents the deploy command
// var deployCmd = &cobra.Command{
// 	Use:     "deploy",
// 	Short:   "Contract erc20 Deploy",
// 	Example: `--name="name" --symbol="symbol" --account="account" --totalbalance=100000`,
// 	Run: func(cmd *cobra.Command, args []string) {
// 		name, _ := cmd.Flags().GetString("name")

// 		symbol, _ := cmd.Flags().GetString("symbol")

// 		account, _ := cmd.Flags().GetString("account")

// 		totalBalance, _ := cmd.Flags().GetInt("totalbalance")

// 		token := erc20.Initalize(name, symbol, account, int64(totalBalance))

// 		fmt.Println("token Deploy!", token)
// 		db.NewDB(name).Add("name", []byte(token.Name()))
// 		db.NewDB(name).Add("symbol", []byte(token.Symbol()))
// 		db.NewDB(name).Add("totalbalance", []byte(fmt.Sprintf("%v", token.TotalSupply())))
// 		defer db.Close()
// 	},
// }

// func init() {
// 	rootCmd.AddCommand(deployCmd)
// 	var name string
// 	deployCmd.Flags().StringVarP(&name, "name", "n", "", "Name require()")
// 	var symbol string
// 	deployCmd.Flags().StringVarP(&symbol, "symbol", "s", "", "Symbol is require")
// 	// Here you will define your flags and configuration settings.
// 	var totalBalance int
// 	deployCmd.Flags().IntVarP(&totalBalance, "totalbalance", "t", 0, "totalBalance is totalSupply")

// 	var contractAddress string
// 	deployCmd.Flags().StringVarP(&contractAddress, "account", "a", "", "account is needed")
// 	deployCmd.MarkFlagRequired("name")
// 	deployCmd.MarkFlagRequired("symbol")
// 	deployCmd.MarkFlagRequired("totalbalance")
// 	deployCmd.MarkFlagRequired("account")
// }
//
