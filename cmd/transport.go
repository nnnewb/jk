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
	"log"
	"os"
	"strings"

	"emperror.dev/errors"
	"github.com/dave/jennifer/jen"
	"github.com/nnnewb/jk/internal/domain"
	stdcli "github.com/nnnewb/jk/internal/gen/http/client/go/std"
	"github.com/nnnewb/jk/internal/gen/http/client/typescript/fetch"
	"github.com/nnnewb/jk/internal/gen/http/doc"
	"github.com/nnnewb/jk/internal/gen/http/server/go/gin"
	stdsvr "github.com/nnnewb/jk/internal/gen/http/server/go/std"
	"github.com/nnnewb/jk/internal/utils"
	"github.com/spf13/cobra"
)

// transportCmd represents the transport command
var transportCmd = &cobra.Command{
	Use:   "transport",
	Short: "generate transport layer code",
	Long:  "generate transport layer code",
	Run: func(cmd *cobra.Command, args []string) {
		var allErrors error
		proto, err := cmd.Flags().GetString("protocol")
		allErrors = errors.Combine(allErrors, err)
		lang, err := cmd.Flags().GetString("language")
		allErrors = errors.Combine(allErrors, err)
		server, err := cmd.Flags().GetBool("server")
		allErrors = errors.Combine(allErrors, err)
		client, err := cmd.Flags().GetBool("client")
		allErrors = errors.Combine(allErrors, err)
		framework, err := cmd.Flags().GetString("framework")
		allErrors = errors.Combine(allErrors, err)
		swagger, err := cmd.Flags().GetBool("swagger")
		allErrors = errors.Combine(allErrors, err)
		embedSwagger, err := cmd.Flags().GetBool("embed-swagger")
		allErrors = errors.Combine(allErrors, err)
		cobra.CheckErr(allErrors)

		pkgPath, service, err := parse(cmd)
		cobra.CheckErr(err)

		switch proto {
		case "http":
			if server {
				switch lang {
				case "go":
					switch framework {
					case "http":
						err = genHTTPServer(service, pkgPath, "transport_http_server.go", embedSwagger)
						cobra.CheckErr(err)
					case "gin":
						err = genGinServer(service, pkgPath, "transport_gin_server.go", embedSwagger)
						cobra.CheckErr(err)
					default:
						cobra.CheckErr(fmt.Errorf("protocol %s server code generation does not support framework %s (%s)", proto, framework, lang))
					}
				default:
					cobra.CheckErr(fmt.Errorf("protocol %s server code generation does not support language %s", proto, lang))
				}
				if swagger {
					err = genSwagger(service, "swagger.json")
					cobra.CheckErr(err)
				}
			} else if client {
				switch lang {
				case "go":
					switch framework {
					case "http":
						err = genHTTPClient(service, pkgPath, "transport_http_client.go")
						cobra.CheckErr(err)
					default:
						cobra.CheckErr(fmt.Errorf("protocol %s server code generation does not support framework %s (%s)", proto, framework, lang))
					}
				case "ts":
					switch framework {
					case "fetch":
						err = genTypeScriptClient(service, "client.ts")
						cobra.CheckErr(err)
					default:
						cobra.CheckErr(fmt.Errorf("protocol %s server code generation does not support framework %s (%s)", proto, framework, lang))
					}
				default:
					cobra.CheckErr(fmt.Errorf("protocol %s client code generation does not support language %s", proto, lang))
				}
			}
		}
	},
}

func init() {
	generateCmd.AddCommand(transportCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// transportCmd.PersistentFlags().String("foo", "", "A help for foo")
	generateCmd.PersistentFlags().StringP("protocol", "P", "http", "transport layer protocol")
	generateCmd.PersistentFlags().StringP("language", "l", "", "programming language")
	generateCmd.PersistentFlags().StringP("framework", "f", "", "server/client framework")
	generateCmd.PersistentFlags().BoolP("swagger", "S", false, "generate swagger document, works with http protocol")
	generateCmd.PersistentFlags().Bool("embed-swagger", false, "embed swagger ui into server code")
	generateCmd.PersistentFlags().BoolP("server", "s", false, "generate server code")
	generateCmd.PersistentFlags().BoolP("client", "c", false, "generate client code")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// transportCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// genHTTPClient 生成HTTP客户端代码。
func genHTTPClient(service *domain.Service, pkg, filename string) error {
	f := jen.NewFilePath(pkg)
	f.HeaderComment(fmt.Sprintf("Code generated by jk %s; DO NOT EDIT.", strings.Join(os.Args[1:], " ")))
	utils.InitializeFileCommon(f)
	stdcli.GenerateHTTPTransportClient(f, service)

	err := f.Save(filename)
	if err != nil {
		return errors.Wrap(err, "render generated transport layer code failed")
	}

	return nil
}

// genHTTPServer 生成HTTP服务器代码。
func genHTTPServer(service *domain.Service, pkg, filename string, embedSwaggerUI bool) error {
	f := jen.NewFilePath(pkg)
	f.HeaderComment(fmt.Sprintf("Code generated by jk %s; DO NOT EDIT.", strings.Join(os.Args[1:], " ")))
	utils.InitializeFileCommon(f)
	stdsvr.GenerateHTTPTransportServer(f, service)

	if embedSwaggerUI {
		stdsvr.GenerateEmbedSwaggerJSON(f, service)
	}

	err := f.Save(filename)
	if err != nil {
		return errors.Wrap(err, "render generated transport layer code failed")
	}

	return nil
}

// genGinServer 生成 gin 服务器代码。
func genGinServer(service *domain.Service, pkg, filename string, embedSwaggerUI bool) error {
	f := jen.NewFilePath(pkg)
	f.HeaderComment(fmt.Sprintf("Code generated by jk %s; DO NOT EDIT.", strings.Join(os.Args[1:], " ")))
	utils.InitializeFileCommon(f)
	err := gin.GenerateGin(f, service)
	if err != nil {
		return errors.Wrap(err, "generate gin server code failed")
	}

	if embedSwaggerUI {
		gin.GenerateGinEmbedSwaggerUI(f, service)
	}

	err = f.Save(filename)
	if err != nil {
		return errors.Wrap(err, "render generated transport layer (gin) code failed")
	}

	return nil
}

// genSwagger 生成Swagger文档。
func genSwagger(service *domain.Service, filename string) error {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return errors.Wrap(err, "open file failed")
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("close file failed, error %+v", err)
		}
	}(file)

	err = doc.GenerateSwagger(file, service)
	if err != nil {
		return errors.Wrap(err, "generate swagger failed")
	}
	return nil
}

// genTypeScriptClient 生成 typescript 客户端代码。
func genTypeScriptClient(service *domain.Service, filename string) error {
	_, err := os.Stat("tsconfig.json")
	if os.IsNotExist(err) {
		file, err := os.OpenFile("tsconfig.json", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
		if err != nil {
			return err
		}
		defer func(file *os.File) {
			_ = file.Close()
		}(file)

		_, err = fmt.Fprintf(file, `{
  "compilerOptions": {
    "target": "ES2015",
    "module": "ES6",
    "sourceMap": true,
    "declaration": true,
    "declarationMap": true,
    "preserveConstEnums": true,
    "lib": [
      "DOM",
      "ES2015"
    ]
  },
  "files": [
    "%s"
  ]
}`, filename)
		if err != nil {
			return err
		}
	}

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return errors.Wrap(err, "open file failed")
	}

	_, err = fmt.Fprintf(file, "// Code generated by jk %s; DO NOT EDIT.\n", strings.Join(os.Args[1:], " "))
	if err != nil {
		return errors.Wrap(err, "write typescript api client failed")
	}

	err = fetch.GenerateTypeScriptClient(file, service)
	if err != nil {
		return errors.Wrap(err, "generate typescript api client failed")
	}

	return nil
}
