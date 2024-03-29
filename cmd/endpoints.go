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
	"os"
	"strings"

	"emperror.dev/errors"
	"github.com/dave/jennifer/jen"
	"github.com/nnnewb/jk/internal/domain"
	"github.com/nnnewb/jk/internal/gen/endpoints"
	"github.com/nnnewb/jk/internal/utils"
	"github.com/spf13/cobra"
)

// endpointsCmd represents the endpoints command
var endpointsCmd = &cobra.Command{
	Use:   "endpoints",
	Short: "generate endpoint code",
	Long:  "generate endpoint code",
	Run: func(cmd *cobra.Command, args []string) {
		pkgPath, service, err := parse(cmd)
		cobra.CheckErr(err)
		err = genEndpoint(service, pkgPath, "endpoints.go")
		cobra.CheckErr(err)
	},
}

func init() {
	generateCmd.AddCommand(endpointsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// endpointsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// endpointsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// genEndpoint 生成服务端点代码。
func genEndpoint(service *domain.Service, pkg, filename string) error {
	f := jen.NewFilePath(pkg)
	f.HeaderComment(fmt.Sprintf("Code generated by jk %s; DO NOT EDIT.", strings.Join(os.Args[1:], " ")))
	utils.InitializeFileCommon(f)
	err := endpoints.GenerateEndpoints(f, service.Interface)
	if err != nil {
		return errors.Wrap(err, "generate endpoints for service failed")
	}

	err = f.Save(filename)
	if err != nil {
		return errors.Wrap(err, "render generated service layer code failed")
	}

	return nil
}
