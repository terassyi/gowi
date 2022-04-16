package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/terassyi/gowi/decoder"
	"github.com/terassyi/gowi/runtime"
	"github.com/terassyi/gowi/runtime/debugger"
	"github.com/terassyi/gowi/runtime/instance"
	"github.com/terassyi/gowi/runtime/value"
	"github.com/terassyi/gowi/types"
	"github.com/terassyi/gowi/validator"
)

var execCommand = &cobra.Command{
	Use:   "exec",
	Short: "execute WASM binary file",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		file := args[0]
		d, err := decoder.New(file)
		if err != nil {
			log.Fatalln(err)
		}
		mod, err := d.Decode()
		if err != nil {
			log.Fatalln(err)
		}
		v, err := validator.New(mod)
		if err != nil {
			log.Fatalln(err)
		}
		if _, err := v.Validate(); err != nil {
			log.Fatalln(err)
		}
		inst, err := instance.New(mod)
		if err != nil {
			log.Fatalln(err)
		}
		listExports, err := cmd.Flags().GetBool("list-all-exports")
		if err != nil {
			log.Fatalln(err)
		}
		if listExports {
			fmt.Printf("WASM fule: %s\n\n", file)
			fmt.Println("List all exported functions")
			for _, exp := range mod.Exports {
				ext, err := inst.GetExport(exp.Name)
				if err != nil {
					log.Fatalln(err)
				}
				f := instance.GetExternVal[*instance.Function](ext)
				fmt.Printf("\t%s(%s) -> (%s)\n", exp.Name, f.Type.Params, f.Type.Returns)
			}
			fmt.Println()
			return
		}
		invoke, err := cmd.Flags().GetString("invoke")
		if err != nil {
			log.Fatalln(err)
		}
		if invoke != "" {
			args, err := cmd.Flags().GetStringSlice("args")
			if err != nil {
				log.Fatalln(err)
			}
			ext, err := inst.GetExport(invoke)
			if err != nil {
				log.Fatalln(err)
			}
			f := instance.GetExternVal[*instance.Function](ext)
			locals, err := parseArgs(f.Type.Params, args)
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Println(locals)
			debugLevel, err := cmd.Flags().GetInt("debug")
			if err != nil {
				log.Fatalln(err)
			}
			runner, err := runtime.New(mod, nil, debugger.DebugLevel(debugLevel))
			if err != nil {
				log.Fatalln(err)
			}
			results, err := runner.Invoke(invoke, locals)
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Println(results)
		}
	},
}

func parseArgs(params types.ResultType, args []string) ([]value.Value, error) {
	values := make([]value.Value, 0, len(params))
	if len(params) != len(args) {
		return nil, fmt.Errorf("The number of given arguments is not matched. required=%d actual=%d", len(params), len(args))
	}
	for i, arg := range args {
		typ := params[i]
		v, err := value.FromString(arg, typ)
		if err != nil {
			return nil, err
		}
		values = append(values, v)
	}
	return values, nil
}
