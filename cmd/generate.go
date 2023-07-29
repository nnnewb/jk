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
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"strings"

	"emperror.dev/errors"
	"github.com/nnnewb/battery/maps"
	"github.com/nnnewb/jk/internal/domain"
	"github.com/nnnewb/jk/internal/utils"
	"github.com/spf13/cobra"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "generate code",
	Long:  `generate code`,
}

func init() {
	rootCmd.AddCommand(generateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateCmd.PersistentFlags().String("foo", "", "A help for foo")
	generateCmd.PersistentFlags().StringP("typename", "t", "", "name of service interface")
	generateCmd.PersistentFlags().StringP("package", "p", "", "package name of service interface")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// parse 解析命令行参数，返回包路径、服务对象和错误信息。
func parse(cmd *cobra.Command) (string, *domain.Service, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", nil, errors.WithStack(err)
	}

	pkgPath, err := cmd.Flags().GetString("package")
	if err != nil {
		return "", nil, errors.WithStack(err)
	}

	if pkgPath == "" {
		pkgPath, err = utils.ResolveFullPackagePath(wd, wd)
		if err != nil {
			return "", nil, errors.Errorf("can not resolve full path of current working directory %s, error %+v", wd, err)
		}
	}

	typeName, err := cmd.Flags().GetString("typename")
	if err != nil {
		return "", nil, errors.WithStack(err)
	}

	fileSet := token.NewFileSet()
	parsedPackages, err := parser.ParseDir(fileSet, wd, nil, parser.ParseComments)
	if err != nil {
		return "", nil, errors.WithStack(err)
	}

	var astPkg *ast.Package
	for name, p := range parsedPackages {
		if strings.Contains(name, "_test") {
			continue
		}
		astPkg = p
	}

	if astPkg == nil || astPkg.Files == nil {
		return "", nil, errors.Errorf("no valid package found in path %s, test package skipped", wd)
	}

	info := &types.Info{}
	typeCheckerConfig := types.Config{Importer: importer.ForCompiler(fileSet, "source", nil)}
	pkg, err := typeCheckerConfig.Check(pkgPath, fileSet, maps.Values(astPkg.Files), info)
	if err != nil {
		return "", nil, errors.WithStack(err)
	}

	service, err := domain.ParseInterfaceData(pkg, astPkg, typeName)
	if err != nil {
		return "", nil, errors.WithStack(err)
	}

	return pkgPath, service, nil
}
