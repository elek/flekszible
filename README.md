# fle[ksz]ible

Flekszible is a Kubernetes resource manager.

Features:

 1. Based on k8s resource fragments (yaml) and mixin rules
 2. Can generate final k8s resources
 3. Or can generate helm charts.

### Recipes

## Getting started

Put one kubernetes resource to the directory:

./nginx.yaml:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  selector:
    matchLabels:
      app: nginx
  replicas: 2
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.7.9
        ports:
        - containerPort: 80
```

Create an override file to the `./transformations/mylabels.yaml`
```yaml
- type: Add
  path: 
    - metadata
    - annotations
  value: 
    felkszible: generated
```

And execute the generation:

```bash
./flekszible k8s . ./out
```

The result will be something like this:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    felkszible: generated
  annotations: {}
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
          volumeMounts: []
          env: []
          envFrom: []
      volumes: []
```

As you can see the original k8s source is modified base on the transformation rules.


## Definitions

### Path

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

Here tha path of the environment variable is `[spec, template, spec, containers, nginx, env, KEY ]` and not `[ spec, template, spec, containers, 0, env, 0]`

### Trigger

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


## Directory structure

In the output directory all the yaml files (except the `flekszible.yaml` configuration file) are considered to be a k8s resource. One file could contain multiple resources.

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


## Available transformation types

### Add

Extends existing k8s resources with additional elements.

#### Parameters

| Name    | Type     | Value 
|---------|----------|-------
| path    | []string | Path in the resource file to 
| value   | yaml     | Any yaml fragment which should be added to the node which is selected by the `Path`

#### Example

```
- type: Add
  path: 
    - metadata
    - annotations
  value: 
    felkszible: generated
```

#### Supported add methods

| Type of the destination node (selected by `Path`) | Type of the `Value` | Supported
|---------------------------------------------------|---------------------|------------
| Map                                               | Map                 | Yes
| Array                                             | Array               | Yes


#### ImageSet

Replaces the docker image definition eveywhere

| Name    | Type     | Value 
|---------|----------|-------
| image   | string   | Full name of the required docker image. Such as 'elek/ozone:trunk'


Note: This transformations could also added with the `--image` CLI argument.

#### PublishStatefulSet

Creates additional NodeType service for StatefulSet internal services.


#### DaemonToStatefulset

Converts daemonset to statefulset.

Useful for minikube based environments where you may not have enough node to run a daemonset based cluster.


