---
title: Transformations
menu: main
---

# Add

Extends existing k8s resources with additional elements.

## Parameters

| Name    | Type     | Value 
|---------|----------|-------
| path    | []string | Path in the resource file to 
| value   | yaml     | Any yaml fragment which should be added to the node which is selected by the `Path`

## Example

```
- type: Add
  path: 
    - metadata
    - annotations
  value: 
    felkszible: generated
```

## Supported add methods

| Type of the destination node (selected by `Path`) | Type of the `Value` | Supported
|---------------------------------------------------|---------------------|------------
| Map                                               | Map                 | Yes
| Array                                             | Array               | Yes
| Array                                             | Map                 | Yes


# Image

Replaces the docker image definition everywhere

| Name    | Type     | Value 
|---------|----------|-------
| image   | string   | Full name of the required docker image. Such as 'elek/ozone:trunk'


Note: This transformations could also added with the `--image` CLI argument.

# Namespace

Similar to the image namespace also can be changed with simple transformation:

### namespace

Use explicit namespace

#### Parameters

+-----------+---------+----------------------------------------------------------------------------------------------------------------------------------+
| name      | default | description                                                                                                                      |
+-----------+---------+----------------------------------------------------------------------------------------------------------------------------------+
| namespace |         | The namespace to use in the k8s resources. If empty, the current namespace will be used (from ~/.kube/config or $KUBECONFIG)     |
| force     | false   | If false (default) only the existing namespace attributes will be changed. If yes, namespace will be added to all the resources. |
+-----------+---------+----------------------------------------------------------------------------------------------------------------------------------+


Note: This transformations could also added with the '--namespace' CLI argument.

Example):

```yaml
- type: Namespace
  namespace: myns
```

# Change

Change is a simple replacement like `sed`. You can apply a regular expression based replacement to a specific _string_ value:

Example: 

```
- type: Change
  trigger:
    metadata:
      name: namenode
  path:
  - spec
  - serviceName
  pattern: (.*)
  replacement: prefix-$1
```

| Name     | Type     | Value 
|----------|----------|-------
| pattern  | string   | Regular expression to replace a string value (see: https://github.com/google/re2/wiki/Syntax)
| replacement | string   | Replacement value.

# Prefix

Add a specific prefix for all of the names.


| Name     | Type     | Value 
|----------|----------|-------
| prefix   | string   | Prefix which will be added to all the names.

Example (`transformations/set.yaml`):

```yaml
- type: Namespace
  namespace: myns
```

# Pipe

Transform content with external shell command.

#### Parameters

| Name      | Type     | Value 
|-----------|----------|-------
| command   | string   | External program which transforms standard input to output.
| args      | []string | List of the arguments of the command.



Pipe executes a specific command to transform a k8s resources to a new resources.

The original manifest will be sent to the stdin of the process and the stdout will be processed as a the stdout of the file.

Example:

```
- type: Pipe
  command: sed
  args: ["s/nginx/qwe/g"]
```


# ConfigHash

Add a kubernetes annotation with the hash of the used configmap. With 
this approach you can force to re-create the k8s resources in case of config change. 
In case of configmap change the annotation value will be different and the resource
will be recreated.

As of now it supports only one configmap per resource and only the top-level
resource will be annotated (in case of statefulset this is the statefulset 
not the pod).

Example (`transformations/config.yaml`):

```yaml
- type: ConfigHash
```

# PublishStatefulSet

Creates additional NodeType service for StatefulSet internal services.


# DaemonToStatefulset

Converts daemonset to statefulset.

Useful for minikube based environments where you may not have enough node to run a daemonset based cluster.

# Composite

You can create additional transformations with grouping existing transformations. For example the following definition register a new transformation type:

```yaml
name: flokkr.github.io/prometheus
---
  - type: Add
    path:
      - spec
      - template
      - metadata
      - annotations
    value:
      prometheus.io/scrape: "true"
      prometheus.io/port: "28942"
  - type: Add
    path:
      - spec
      - template
      - spec
      - containers
      - "*"
      - env
    value:
      - name: "PROMETHEUSJMX_ENABLED"
        value: "true"
      - name:  "PROMETHEUSJMX_AGENTOPTS"
        value: "port=28942"

```

Put the previous file to the `definitions/prometheus.yaml`. By default it won't be applied to any k8s resources but from now you can use the `flokkr.github.io/prometheus` type in your transformations:

For example in `transformations/monitor.yaml` you can write:

```yaml
- type: flokkr.github.io/prometheus
```
