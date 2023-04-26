/*
Copyright © 2023 weak_ptr <weak_ptr@outlook.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/nnnewb/jk/internal/gen"
	"github.com/nnnewb/jk/internal/utils"
	"go/importer"
	"go/token"
	"go/types"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// allCmd represents the all command
var allCmd = &cobra.Command{
	Use:   "all",
	Short: "Generate all we can.",
	Long:  `Generate all we can.`,
	Run: func(cmd *cobra.Command, args []string) {
		pkgPath, err := cmd.Flags().GetString("package")
		cobra.CheckErr(err)
		typeName, err := cmd.Flags().GetString("typename")
		cobra.CheckErr(err)

		if pkgPath == "" {
			wd, err := os.Getwd()
			if err != nil {
				log.Fatalf("get working directory failed, error %+v", err)
			}

			pkgPath, err = utils.ResolveFullPackagePath(wd, wd)
			if err != nil {
				log.Fatalf("can not resolve full path of current working directory %s, error %+v", wd, err)
			}
		}

		fset := token.NewFileSet()
		imp := importer.ForCompiler(fset, "source", nil)
		pkg, err := imp.Import(pkgPath)
		if err != nil {
			log.Fatalf("import package %s failed, error %+v", pkgPath, err)
		}

		result := pkg.Scope().Lookup(typeName)
		if result == nil {
			log.Fatalf("type %s not found in package %s", typeName, pkgPath)
		}

		f := jen.NewFilePath(result.Pkg().Path())
		f.HeaderComment(fmt.Sprintf("Code generated by jk %s; DO NOT EDIT.", strings.Join(os.Args[1:], " ")))
		err = gen.GenerateEndpoints(f, result.Type().(*types.Named))
		if err != nil {
			log.Fatalf("generate endpoints for service failed, error %+v", err)
		}

		err = f.Save("endpoint_gen.go")
		if err != nil {
			log.Fatalf("render generated service layer code failed, error %+v", err)
		}

		f = jen.NewFilePath(result.Pkg().Path())
		f.HeaderComment(fmt.Sprintf("Code generated by jk %s; DO NOT EDIT.", strings.Join(os.Args[1:], " ")))
		err = gen.GenerateTransportLayerHTTP(f, result.Type().(*types.Named))
		if err != nil {
			log.Fatalf("generate endpoint factory for service failed, error %+v", err)
		}

		err = f.Save("transport_gen.go")
		if err != nil {
			log.Fatalf("render generated transport layer code failed, error %+v", err)
		}
	},
}

func init() {
	generateCmd.AddCommand(allCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// allCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// allCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
