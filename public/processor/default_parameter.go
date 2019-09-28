package processor

var TriggerParameter = ProcessorParameter{
	Name:        "trigger",
	Description: "A yaml struct to define the condition when the rule, should be applied.",
}

var PathParameter = ProcessorParameter{
	Name:        "path",
	Description: "A string array to define the path of the transformation in the kubernetes yaml.",
}

var TriggerDoc = `

#### Trigger parameter

Trigger can define a filter to apply the transformations only on a subset of the k8s manifests.

For example:

'''
- type: Add
  trigger:
  metadata:
    name: datanode
  path:
  - metadata
  - labels
  value:
  flokkr.github.io/monitoring: false
'''

This definition will apply only to the k8s resources where the value of 'metadata.name' is 'datanode'.

You can use multiple values in the trigger. All the key nodes will be collected and should be the same in the target resource.
`

var PathDoc = `

#### Path parameter

Path is a string array where each element represents a new level in the kubernetes manifest.

For example the '["spec","spec", "spec", "containers"]' array address a list in the kubernetes manifest files.

Array elements are indexed with a number from zero (eg. ["foo","barr",0]) _except_ if the elements in the array are maps with an existing _name_ key. In this case the index is this name.

For example with this standard kubernetes manifest:

'''
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
'''

The path of the KEY environment variable is ''[spec, template, spec, containers, nginx, env, KEY ]' and not '[ spec, template, spec, containers, 0, env, 0]'

For matching, path segments are used as regular expressions. Therefore the following path matches for both the init and main containers:

'''yaml
path:
- spec
- template
- spec
- (initC|c)ontainers
- .*
'''

Matching works only if the Yaml file already has the specified path. But for kubernetes resources a few standard paths are pre-defined'
`
