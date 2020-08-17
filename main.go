package main

// import "github.com/elek/flekszible/cmd"
import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/hashicorp/go-getter"
	"github.com/sirupsen/logrus"

	"github.com/elek/flekszible/api/processor"
	"github.com/elek/flekszible/pkg"
	"github.com/elek/flekszible/pkg/operator"
	"github.com/urfave/cli"
)

var version string
var commit string
var date string

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true, ForceQuote: false})
	var inputDir, outputDir, imageOverride, namespaceOverride string
	var minikube bool

	app := cli.NewApp()
	app.Name = "Flekszible"
	app.Usage = "Generate kubernetes sources files"
	app.Description = "Kubernetes resource file generator"
	app.Version = fmt.Sprintf("%s (%s, %s)", version, commit, date)
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "source, s",
			Value:       "",
			Usage:       "Source directory to read the Flekszible definition (default: current dir)",
			Destination: &inputDir,
		},
		cli.StringFlag{
			Name:        "destination, d",
			Value:       "",
			Usage:       "Destination directory to generate the k8s resource (default: current dir)",
			Destination: &outputDir,
		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "generate",
			Aliases: []string{"g"},
			Usage:   "Generate kubernetes resources files",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "namespace, n",
					Value: "",
					Usage: "kubernetes namespace override. With empty value (\"\") the current namespace will " +
						"be used. With exact value the namespace will be forced (set even if the namespace is not added to the to the k8s resources",
					Destination: &namespaceOverride,
				},
				cli.StringFlag{
					Name:        "image, i",
					Usage:       "docker image name override",
					Destination: &imageOverride,
				},

				cli.StringSliceFlag{
					Name:  "transformations, t",
					Usage: "manually defined transformations",
				},

				cli.StringSliceFlag{
					Name:  "import, p",
					Usage: "Define additional import paths",
				},
				cli.BoolFlag{
					Name:  "print",
					Usage: "Print out the result to the standard output (same as to use - as the destination'-d -')",
				},
				cli.BoolFlag{
					Name:        "minikube",
					Usage:       "Enable minikube specific defaults (eg. daemonset to statefulset conversion)",
					Destination: &minikube,
				},
			},
			Action: func(c *cli.Context) error {
				outputDir := findOutputDir(&outputDir, c.Args().Get(0))
				if c.Bool("print") {
					outputDir = "-"
				}
				context := processor.CreateRenderContext("k8s", findInputDir(&inputDir, c.Args().Get(0)), outputDir)
				context.ImageOverride = imageOverride
				context.Namespace = namespaceOverride
				return pkg.Run(context, minikube, c.StringSlice("import"), c.StringSlice("transformations"))
			},
		},
		{
			Name:      "import",
			Usage:     "Import multiple kubernetes resource from file or stdin and generate output to current dir or stdout.",
			ArgsUsage: "[file or -] [output dir]",
			Flags: []cli.Flag{
				cli.StringSliceFlag{
					Name:  "transformations, t",
					Usage: "manually defined transformations",
				},
				cli.BoolFlag{
					Name:  "print",
					Usage: "Print out the result to the standard output (same as to use - as the destination'-d -')",
				},
			},
			Action: func(c *cli.Context) error {
				dir, err := os.Getwd()
				if err != nil {
					return err
				}
				if c.Bool("print") {
					dir = "-"
				}
				err = pkg.Import(c.Args().Get(0), c.StringSlice("transformations"), dir)
				if err != nil {
					return err
				}
				return nil
			},
		},
		{
			Name:      "list",
			Usage:     "List managed kubernetes resources files.",
			ArgsUsage: "[flekszible_dir]",
			Action: func(c *cli.Context) error {
				context := processor.CreateRenderContext("k8s", findInputDir(&inputDir, c.Args().Get(0)), findOutputDir(&outputDir, c.Args().Get(0)))
				pkg.ListResources(context)
				return nil
			},
		},
		{
			Name:      "tree",
			Usage:     "List managed resources files and registered transformations.",
			ArgsUsage: "[flekszible_dir]",
			Flags: []cli.Flag{
				cli.StringSliceFlag{
					Name:  "transformations, t",
					Usage: "manually defined transformations",
				},
			},
			Action: func(c *cli.Context) error {
				context := processor.CreateRenderContext("k8s", findInputDir(&inputDir, c.Args().Get(0)), findOutputDir(&outputDir, c.Args().Get(0)))
				err := context.Init()
				if err != nil {
					return err
				}
				err = context.AddAdHocTransformations(c.StringSlice("transformations"))
				if err != nil {
					return err
				}
				pkg.Tree(context)
				return nil
			},
		},
		{
			Name:  "clean",
			Usage: "Delete yaml files from the destination directories",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "all, a",
					Usage: "Delete all yaml files from destination directories",
				},
			},
			Action: func(c *cli.Context) error {
				context := processor.CreateRenderContext("k8s", findInputDir(&inputDir, c.Args().Get(0)), findOutputDir(&outputDir, c.Args().Get(0)))
				return pkg.Cleanup(context, c.Bool("all"))
			},
		},
		appCommands(&inputDir, &outputDir),
		sourceCommands(&inputDir, &outputDir),
		transformationCommands(&inputDir, &outputDir),
		admissionCommands(),
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func admissionCommands() cli.Command {
	return cli.Command{
		Name:      "operator",
		Usage:     "Start Kubernetes operator (Mutating webhook endpoint)",
		ArgsUsage: "Directory of the Flekszible definition",
		Action: func(c *cli.Context) error {
			return operator.StartServer(c.Args().First())
		},
	}
}
func appCommands(inputDir *string, outputDir *string) cli.Command {
	return cli.Command{
		Name:  "app",
		Usage: "Manage importable dirs/applications.",
		Subcommands: []cli.Command{
			{
				Name:      "list",
				Usage:     "List all the referenced/imported flekszible directory.",
				ArgsUsage: "[flekszible_dir]",
				Action: func(c *cli.Context) error {
					context := processor.CreateRenderContext("k8s", findInputDir(inputDir, c.Args().Get(0)), findOutputDir(outputDir, c.Args().Get(0)))
					pkg.ListApp(context)
					return nil
				},
			},
			{
				Name:      "search",
				Usage:     "Search for available importable apps/dirs from the active sources",
				ArgsUsage: "[flekszible_dir]",
				Action: func(c *cli.Context) error {
					context := processor.CreateRenderContext("k8s", findInputDir(inputDir, c.Args().Get(0)), findOutputDir(outputDir, c.Args().Get(0)))
					pkg.SearchComponent(context)
					return nil
				},
			},
			{
				Name:      "add",
				Usage:     "Add (import) new app to the flekszible.yaml.",
				ArgsUsage: "[flekszible_dir]",
				Action: func(c *cli.Context) error {
					context := processor.CreateRenderContext("k8s", findInputDir(inputDir, ""), findOutputDir(outputDir, c.Args().Get(0)))
					pkg.AddApp(context, *inputDir, c.Args().First())
					return nil
				},
			},
		},
	}
}

func sourceCommands(inputDir, outputDir *string) cli.Command {
	return cli.Command{
		Name:  "source",
		Usage: "Manage sources of the importable applications.",
		Subcommands: []cli.Command{
			{
				Name:      "list",
				Usage:     "List registered sources (directories / repositories where other directories are imported from)",
				ArgsUsage: "[flekszible_dir]",
				Action: func(c *cli.Context) error {
					context := processor.CreateRenderContext("k8s", findInputDir(inputDir, c.Args().Get(0)), findOutputDir(outputDir, c.Args().Get(0)))
					pkg.ListSources(context)
					return nil
				},
			},
			{
				Name:  "search",
				Usage: "Search for the available source repositories (github repositories with flekszible:topic).",
				Action: func(c *cli.Context) error {
					pkg.SearchSource()
					return nil
				},
			},
			{
				Name:      "add",
				Usage:     "Register source to your Flekszible descriptor file",
				ArgsUsage: "[flekszible_dir]",
				Action: func(c *cli.Context) error {
					context := processor.CreateRenderContext("k8s", findInputDir(inputDir, c.Args().Get(0)), findOutputDir(outputDir, c.Args().Get(0)))
					return pkg.AddSource(context, findInputDir(inputDir, ""), c.Args().First())
				},
			},
		},
	}
}

func transformationCommands(inputDir, outputDir *string) cli.Command {
	return cli.Command{
		Name:    "transformation",
		Aliases: []string{"trafo", "definitions"},
		Usage:   "Show available transformation types/definitions",
		Subcommands: []cli.Command{
			{
				Name:      "list",
				Aliases:   []string{"search"},
				Usage:     "List available transformation types.",
				ArgsUsage: "[flekszible_dir]",
				Action: func(c *cli.Context) error {
					context := processor.CreateRenderContext("k8s", findInputDir(inputDir, c.Args().Get(0)), findOutputDir(outputDir, c.Args().Get(0)))
					return pkg.ListProcessor(context)
				},
			},
			{
				Name:      "info",
				Aliases:   []string{"show"},
				Usage:     "Show details information / help about a specific transformation type",
				ArgsUsage: "[flekszible_dir]",
				Action: func(c *cli.Context) error {
					context := processor.CreateRenderContext("k8s", findInputDir(inputDir, ""), findOutputDir(outputDir, c.Args().Get(0)))
					return pkg.ShowProcessor(context, c.Args().First())
				},
			},
		},
	}
}

func findInputDir(argument *string, inputDirFromArg string) string {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	if len(inputDirFromArg) > 0 {
		return inputDirFromArg
	}
	if *argument != "" {
		//input dir is specificed
		if _, err := os.Stat(*argument); err == nil {
			return *argument
		} else {
			//no such file, assuming it's an url
			workDir := path.Join(pwd, ".flekszible-source")
			err := getter.Get(workDir, *argument)
			if err != nil {
				panic(err)
			}
			return workDir
		}
		return *argument
	}

	return pwd
}

func findOutputDir(argument *string, outputDirFromArg string) string {
	if len(outputDirFromArg) > 0 {
		return outputDirFromArg
	}
	if *argument != "" {
		return *argument
	}
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return pwd
}
