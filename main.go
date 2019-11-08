package main

// import "github.com/elek/flekszible/cmd"
import (
	"fmt"
	"github.com/hashicorp/go-getter"
	"log"
	"os"
	"path"

	"github.com/elek/flekszible/api/processor"
	"github.com/elek/flekszible/pkg"
	"github.com/urfave/cli"
)

var version string
var commit string
var date string

func main() {

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
			Usage:   "Genrate kubernetes resources files",
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
				outputDir := findOutputDir(&outputDir)
				if c.Bool("print") {
					outputDir = "-"
				}
				context := processor.CreateRenderContext("k8s", findInputDir(&inputDir), outputDir)
				context.ImageOverride = imageOverride
				context.Namespace = namespaceOverride
				pkg.Run(context, minikube, c.StringSlice("import"), c.StringSlice("transformations"))
				return nil
			},
		},
		{
			Name:  "list",
			Usage: "List managed kubernetes resources files.",
			Action: func(c *cli.Context) error {
				context := processor.CreateRenderContext("k8s", findInputDir(&inputDir), findOutputDir(&outputDir))
				pkg.ListResources(context)
				return nil
			},
		},
		appCommands(&inputDir, &outputDir),
		sourceCommands(&inputDir, &outputDir),
		transformationCommands(&inputDir, &outputDir),
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func appCommands(inputDir *string, outputDir *string) cli.Command {
	return cli.Command{
		Name:  "app",
		Usage: "Manage importable dirs/applications.",
		Subcommands: []cli.Command{
			{
				Name:  "list",
				Usage: "List all the referenced/imported flekszible directory.",
				Action: func(c *cli.Context) error {
					context := processor.CreateRenderContext("k8s", findInputDir(inputDir), findOutputDir(outputDir))
					pkg.ListApp(context)
					return nil
				},
			},
			{
				Name:  "search",
				Usage: "Search for available importable apps/dirs from the active sources",
				Action: func(c *cli.Context) error {
					context := processor.CreateRenderContext("k8s", findInputDir(inputDir), findOutputDir(outputDir))
					pkg.SearchComponent(context)
					return nil
				},
			},
			{
				Name:  "add",
				Usage: "Add (import) new app to the flekszible.yaml.",
				Action: func(c *cli.Context) error {
					context := processor.CreateRenderContext("k8s", findInputDir(inputDir), findOutputDir(outputDir))
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
				Name:  "list",
				Usage: "List registered sources (directories / repositories where other directories are imported from)",
				Action: func(c *cli.Context) error {
					context := processor.CreateRenderContext("k8s", findInputDir(inputDir), findOutputDir(outputDir))
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
				Name:  "add",
				Usage: "Register source to your Flekszible descriptor file",
				Action: func(c *cli.Context) error {
					context := processor.CreateRenderContext("k8s", findInputDir(inputDir), findOutputDir(outputDir))
					return pkg.AddSource(context, findInputDir(inputDir), c.Args().First())
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
				Name:  "search",
				Usage: "List available transformation types.",
				Action: func(c *cli.Context) error {
					context := processor.CreateRenderContext("k8s", findInputDir(inputDir), findOutputDir(outputDir))
					pkg.ListProcessor(context)
					return nil
				},
			},
			{
				Name:  "show",
				Usage: "Show details information / help about a specific transformation type",
				Action: func(c *cli.Context) error {
					context := processor.CreateRenderContext("k8s", findInputDir(inputDir), findOutputDir(outputDir))
					pkg.ShowProcessor(context, c.Args().First())
					return nil
				},
			},
		},
	}
}

func findInputDir(argument *string) string {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
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

func findOutputDir(argument *string) string {
	if *argument != "" {
		return *argument
	}
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return pwd
}
