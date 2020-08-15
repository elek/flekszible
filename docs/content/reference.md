---
title: Reference
menu: main
---

# The main concept


 * **Import and source** Every directory (which contains Kubernetes resource files) can be used as the source of can be used as a base of the flekszible transformations. With using the optional `Flekszible` descriptor you can *import* resources from multiple directories (with optional transformations) to the destination dir. *Source* path also can be defined: dirs can be reported from local dir or remote git repositories.   

 * **Transformations and definitions**: Transformations can be applied to any Kubernetes resources to modify them. Transformations has a type (like `add`, `replace`, `image`,...). New composite transformations can be created with the help of existing transformation types.

 * **Generators**: arbitrary files (like shell scripts, keytab / keystore definitions) can be converted to Kubernetes resources files (eg. configmap can be generated from simple config files, secrets can be generated with shell scripts)

# Directory structure

Flekszible directories can contain:

 * kubernetes resource files
 * transformations
 * definition (transformation templates)
 * raw configuration (to convert them to configmaps) 

## Simple os dir

A simple OS dir can work as a source directory for flekszible __without any descriptor__. In this case all the Kubernetes resource files will be imported and used from the dir.

## `Flekszible/flekszible.yaml` descriptor

In case of a `Flekszible` file does exist, it will be parsed and all the imports/sources/transformations will be used from the Flekszible file and from the directories relative to the Flekszible file (eg `./transformations`, `./configmaps`, ...). In this case __none__ of the other files in the directory will be used as source. Best to use this structure to generate the final k8s resources.

The structure of a flekszible directory is the following:

 * `Flekszible`: configuration file. Could include other directories (all the resources + transformations + definitions will be added from that directory.)
 * `transformations/*.yaml`: will be applied to all the resources according to the specified rules
 * `definitions/*.yaml`: composite definitions which could be used in `transformations.yaml`. Won't be applied by default.
 * `configmaps/*_*.*`: all the files from here will be imported as configmaps. The first part of the filename (before the first `_`) will be used as the name of the configmap, the remaining part is the key inside the configmap.
 * `resources/*.yaml`: used as Kubernetes resources
 * `*.yaml`: used as Kubernetes resources
 
 **Note**: to avoid circular dependencies, destination directories are never read for additional resources. For example if the current directory is the destination directory, `*.yaml` files won't be used from the current directory (use `./resources`)
  
# Imports

## Simple import

You can import other directory structures with adding references to the `flekszible.yaml`/`Flekszible`

For example

```
import:
  - path: ../../hadoop
```

All the transformations + definitions + k8s resources will be added and applied (see the previous section about the directories). Note: the transformations from the imported directory will be applied only to the imported resources.


## Import to subdirectory

The imported resources could be generated to a subdirectory:

```
import;
  - path: ../../hadoop
  - path: ../../prometheus
    destination: monitoring
```

With this approach the prometheus related resources will be saved to the `monitoring` subdirectory of the destination path.

## Import with transformations

Transformations also can be applied to the imported resources:

```
import:
  - path: ../../hadoop
    transformations:
       - type: Image
         image: elek/ozone
```


The hadoop resources are imported here and the image reference is changed during the import (only for the imported resources).

## Import from external source

Imported path is checked in the following location (in this order):

 1. The directory which is defined with the `FLEKSZIBLE_PATH` environment variable
 2. In the current directory (or more preciously: relative to the current directory).
 3. In any local remote source

Remote sources can be defined with the `source` tag:

```
source:
   - url: github.com/flokkr/k8s
import:
   path: ozone
```

The remote repositories are downloaded with [go-getter](https://github.com/hashicorp/go-getter) and the downloaded directory is stored in the `.cache` directory (relative to the input dir). 


Local sources can be defined with the `path` tag:

```
source:
  - url: git::https://github.com/flokkr/docker-ozone.git
  - url: git::https://github.com/elek/docker-byteman
  - url: git::https://github.com/elek/docker-java-async-profiler
  - path: ../flekszible
import:
    - path: ozone
```

Note: During local and remote imports, the Flekszible files can be defined in a `flekszible` subdirectory of the source repository or directory (to make it easier to organize them).

## Automatic/implicit import

Normally resource dirs are imported based on the `import` tag definitions. But external sources might have defintions which can be usefull to be imported automatically.

For example if a remote `Flekszible` dir defines only `definitions`, it might look like this:

```
.
├── Dockerfile
├── flekszible
│   └── byteman-helpers
│       ├── definitions
│       │   └── byteman.yaml
│       └── flekszible.yaml
└── LICENSE
```

To use this repository, it should be added to the source ** and import:

```
source:
  - url: git::https://github.com/elek/docker-byteman
import:
    - path: byteman-helpers
```

To make it easier to share simple definitions, the `_global` directory is always auto imported.

Let's rename the name of flekszible subdir to `_global`:

```
.
├── Dockerfile
├── flekszible
│   └── _global
│       ├── definitions
│       │   └── byteman.yaml
│       └── flekszible.yaml
└── LICENSE
```

Now, the content of the `_global` is auto-imported, therefore the following two definitions are equivalent:

```
source:
  - url: git::https://github.com/elek/docker-byteman
```

```
source:
  - url: git::https://github.com/elek/docker-byteman
import:
    - path: _global
```

# Definitions and Transformations


The heart of flekszible is modifying Kubernetes resources files. There are multiple way to modify existing YAML files and they are defined with different types of transformation definitions. (For example you can `Add` additional fragments or `Replace` existing one).

With the help of the transformation definitions you can instantiate a transformations in the `Flekszible` file (under `transformations` key) or in the `transformation` subdirectory.

For example: 

```
- type: Add
  path:
  - spec
  - template
  - spec
  - containers
  - "datanode"
  - env
  value:
  - name: KEY1
    value: VALUE1
  - name: KEY2
    value: VALUE2
```

## Scope of transformations

It's very important that the transformation is activated only in the subtree where it's used.

Imagine the following structure (hdfs-ha _imports_ zookeeper and hdfs resources).

```             +---------+
       import   |         |      import
      +---------+ hdfs-ha +------------+
      |         |         |            |
      |         +---------+            |
      v                                |
+-----+----+                   +-------v-----+
|          |                   |             |
|   hdfs   |                   |  zookeeper  |
|          |                   |             |
+----------+                   +-------------+

```

In this case a transformation which is defined in the `zookeeper/transformations` directory or in the `zookeeper/Flekszible` file are applied only to the zookeeper resources.

To use global transformations you have two options.

  1.) You can define a composite transformation (`./definitions`) definition which may or may not be activated.
  
  2.) You can set the `scope` of the transformation to `global`.
  
 Let's check examples for both of these cases:

If you have doubts, check with `flekszible tree` subcommand.

## Cli transformations

Some cli commands (such as `flekszible generate`) accept an additional `-t` arguments which can define additional, ad-hoc transformations. For example this script executes an ad-hoc test:

```
flekszible generate --print \
   -t namefilter:include=test-runner \
   -t run:args="/opt/test.sh -c custom -p parameters " | kubectl apply -f -
```

(Note: `namefilter` transformation keeps only the matched resource)

The syntax of `-t` is `TRANSFORMATION:PARAM1=VALUE,PARAM2=VALUE`, and this can be used only for simplified transformations (for example `Add` transformation requires complex yaml structure as a parameter, therefor it couldn't be used from CLI.)

## Global transformations

Similar tto the CLI `-t` arguments, transformations can be defined with `FLEKSZIBLE_TRANSFORMATION`:

```
export FLEKSZIBLE_TRANSFORMATIONS="ozone/onenode"
```

This example enables `ozone/onenode` transformation until it's defined (`ozone/onenode` is defined to remove the affinity rules, and this approach makes it possible to test a resource set in one-node Kubernetes cluster from the CI).

## Path

Many transformations requires a `Path` definition to address a specific node in the yaml tree.

Path is a string array where each element represents a new level. Array elements are indexed from zero (eg ["foo","barr",0]) __except_ if the elements in the array are maps with _name_ key. In this case the index is this name.

For example:

```
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    felkszible: generated
spec:
  selector:
    matchLabels:
      app: nginx
  replicas: 2
  template:
    metadata:
      labels:
        app: nginx
      annotations: {}
    spec:
      containers:
        - name: nginx
          image: nginx:1.7.9
          ports:
            - containerPort: 80
          env: 
            - name:  KEY
              value: XXX
```

Here the path of the environment variable is `[spec, template, spec, containers, nginx, env, KEY ]` and not `[ spec, template, spec, containers, 0, env, 0]`

For matching, path segments are used as regular expressions. Therefore the following path matches for both the init and main containers:

```yaml
path:
  - spec
  - template
  - spec
  - (initC|c)ontainers
  - .*
```

## Trigger

Most of the processors also use a `trigger` parameter. With trigger you can specify any k8s fragments and only the matched resources will be transformed.

For example:

```yaml
- type: Add
  trigger:
     metadata:
        name: datanode
  path:
    - metadata
    - labels
  value:
     flokkr.github.io/monitoring: false
```

This definition will apply only to the k8s resources where the value of `metadata.name` is `datanode`.

You can use multiple values in the trigger. All the key nodes will be collected and should be the same in the target resource.


## Definition example
 
 From the previous examples you can create a `defintions` directory in the hdfs directory and put a transformation definitions:
 
 In `hdfs/definitions/prometheus.yaml`
 
 ```
name: hdfs/prometheus
description: Enable prometheus monitoring
---
- type: Add
  path:
    - spec
    - template
    - spec
  value:
    replicas: {{.replicas}}
```

This transformation definition is not activated by default (because it's not part of a `transformation` directory).

But you can turn it on if you need it. For example in the `hdfs-ha/Flekszible` file:

```yaml
import:
  - path: ../hdfs
    transformations:
    type: hdfs/prometheus
```

## Using global transformation

The second option is more simple. Imagine that you have a grafana dashboard defined in the `hdfs` application.

You would like to define the grafana statefulset to add your configmap, but the grafana configmap imported from a different location.

With using `scope: global` parameters, you can ask flekszible to apply the transformations to _all_ the resources files not just the current subtree:

```yaml
name: prometheus
description: Enable prometheus monitoring
parameters:
  - name: replicas
    default: 2
    required: false
    type: int
---
- type: Add
  path:
    - spec
    - template
    - spec
  value:
    replicas: {{.replicas}}

```

## Using optional transformation

In some cases the transformation should be optional. For example `grafana` app itself may provide a `grafana/install-dashboard` transformation type. If you would like to create a transformation which is executed only if the used transformation type is used, use the `optional: true ` flag.

For example in the `transformations/grafana-dashboard.yaml` you can define a transformation:

```yaml
- type: grafana/install-dashboard
  scope: global
  optional: true
  configmap: ozone-dashboard
```

If the `grafana/install-dashboard` transformation type is available (which means you used an import to import grafana flekszible definition) this transformation will modify the grafana configmap and register all the dashboards in the `ozone-dashboard` configmap to be available. If grafana is not imported the transformation will be ignored without error.

# Generators

Generators are simple plugins which can 

 * create additional Kubernetes resource files during the import (eg. convert shell scripts / config files to configmaps)
 * can write the destination directory directly (eg. generate helper shell scripts)
 * They are activated for a specific type of subdirectory (eg. to a subdirectory with the name of `configmaps`)

As of now we have three type of generators and will be described here in more details

 ## Configmaps

Configmap generator is wrapping any file to real Kubernetes configmaps. Generator is activated for the `configmaps` subdirectories.

To create a configmap: 

 1. Create a `configmaps` subdirectory
 2. Put any file with the nameconvention: *configmapname*_*keyname*

## Output

Output generator copies resources to the output directories. Create a `output` directory and all the resources will be copied to the destination during the import. Can be used to manage helper shell scripts.

## Secrets

Secret generator is activated for any subdirectory which contains the file `.secretgen`. This file should contain a line `secret: <script>` where `<script>` is the name of a shell script which exists in the output directory. 

Generator will iterate over all the other files in the directory and execute `<script> <file>`. The shell script should output with multiple lines in the format `<key> <secret>` Where `<key>` is the name of the secret inside the Kubernetes secret resource, and `<secret>` is a `base64` encoded link.
