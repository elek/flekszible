---
title: Reference
menu: main
---

# Path

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

# Trigger

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


# Directory structure

Flekszible directories can contain:

 * kubernetes resource files
 * transformations
 * definition (transformation templates)
 * raw configuration (to convert them to configmaps) 

## Simple os dir

A simple OS dir can work as a source directory for flekszible __without any descriptor__. In this case all the kubernetes resource files will be imported and used from the dir.

## Simplified descriptor (Flekszible)

In case of a `Flekszible` file does exist, it will be parsed and all the imports/sources/transformations will be used based on the Flekszible file. In this case __none__ of the other files in the directory will be used. Best to use this structure to generate the final k8s resources.

## Fully featured descriptor (flekszible.yaml)

If `flekszible.yaml` does exist in the directory it's handled as a full flekszible definition and all the resources/transformations/configs will be used from the directory based on the following naming convention:

 * `*.yaml`: used as k8s resources
 * `transformations/*.yaml`: will be applied to all the resources according to the specified rules
 * `definitions/*.yaml`: composit definitions which could be used in `transformations.yaml`. Won't be applied by default.
 * `configmaps/*_*.*`: all the files from here will be imported as configmaps. The first part of the filename (before the first `_`) will be used as the name of the configmap, the remaining part is the key inside the configmap.
 * `flekszible.yaml`: configuration file. Could include other directories (all the resources + transformations + definitions will be added from that directory.)

In the output directory all the `yaml` files (except the `flekszible.yaml` configuration file) are considered to be a k8s resource. One file could contain multiple resources.

All the yaml files from the `transformations` subdirectory are parsed as transformation definitions. Each file should contain an array of objects. The object should have a `type` definition (see the available transformations below) 

Example: `./transformations/label.yaml`

```
- type: Add
  path: 
    - metadata
    - annotations
  value: 
    felkszible: generated
```

All the yaml files from the `definitions` directory will be parsed as composit transformation type. You can define multiple transformation and name it. It may be used form other transformation files.

## Summary

You have the following options to read a directory.

| descriptor file   | k8s resource files to load      | transformations               | definitions    | configs  
|-------------------|---------------------------------|-------------------------------|----------------|-------------
| None              | all the yaml files              | None                          | None           | None
| `Flekszible`      | None                            | desc(1)                       | None           | None            
| `flekszible.yaml` | *.yaml                          | ./transformations/* + desc(1) | ./definitions  | ./configmaps/*

(1) Only from the descriptor file

# Imports

## Simple import

You can import other directory structures with adding references to the `flekszible.yaml`/`Flekszible`

For example

```
import:
  - path: ../../hadoop
```

All the transformations + definitions + k8s resources will be added and applied (see the previos section about the directories). Note: the transformations from the imported directory will be applied only to the imported resources.


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
 2. In the current directory (or more  preciously: relative to the current directory).
 3. In any remote source

Remote sources can be defined with the `source` tag:

```
source:
   - url: github.com/flokkr/k8s
import:
   path: ozone
```

The remote repositories are downloaded with [go-getter](https://github.com/hashicorp/go-getter) and the downloaded directory is stored in the `.cache` directory (relative to the input dir). 


# Definitions and Transformations

The heart of flekszible is modifying kubernetes resources files. There are multiple way to modify existing YAML files and they are defined with different types of transformation definitions. (For example you can `Add` additional fragments or `Replace` existing one).

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

It's very important that the transformation is activated only the subtree ver it's used.

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

In this case a transformation which is defined in the zookeeper directory is applied only the zookeeper resources.

To use global transformations you have two options.

  1.) You can define a composite transformation definition which may or may not be activated.
  
  2.) You can set the `scope` of the transformation to true.
  
 Let's check examples for both of these cases:
  
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

With using `scope: globale` parameters, you can ask flekszible to apply the transformations to _all_ the resources files not just the current subtree:

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