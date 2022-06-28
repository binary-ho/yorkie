/*
 * Copyright 2021 The Yorkie Authors. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"context"
	"errors"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"

	"github.com/yorkie-team/yorkie/admin"
	"github.com/yorkie-team/yorkie/pkg/document/key"
)

func newHistoryCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "history [project name] [document key]",
		Short: "Show the history of a document",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				return errors.New("project name and document key are required")
			}

			// TODO(hackerwins): use adminAddr from env or addr flag.
			cli, err := admin.Dial("localhost:11103")
			if err != nil {
				return err
			}
			defer func() {
				_ = cli.Close()
			}()

			ctx := context.Background()
			changes, err := cli.ListChangeSummaries(ctx, args[0], key.Key(args[1]))
			if err != nil {
				return err
			}

			tw := table.NewWriter()
			tw.Style().Options.DrawBorder = false
			tw.Style().Options.SeparateColumns = false
			tw.Style().Options.SeparateFooter = false
			tw.Style().Options.SeparateHeader = false
			tw.Style().Options.SeparateRows = false
			tw.AppendHeader(table.Row{
				"SEQ",
				"MESSAGE",
				"SNAPSHOT",
			})
			for _, change := range changes {
				tw.AppendRow(table.Row{
					change.ID.ServerSeq(),
					change.Message,
					change.Snapshot,
				})
			}
			cmd.Printf("%s\n", tw.Render())
			return nil
		},
	}
}

func init() {
	rootCmd.AddCommand(newHistoryCmd())
}
