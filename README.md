# OpenTofu Mixin for Porter

This is a OpenTofu mixin for [Porter](https://porter.sh).

[![porter/tofu-mixin](https://github.com/getporter/tofu-mixin/actions/workflows/tofu-mixin.yml/badge.svg?branch=main)](https://github.com/getporter/tofu-mixin/actions/workflows/tofu-mixin.yml)

## Install via Porter

This will install the latest mixin release via the Porter CLI.

```bash
porter mixin install tofu
```

## Build from source

Following commands build the OpenTofu mixin.

```bash
git clone https://github.com/getporter/tofu-mixin.git
cd tofu-mixin
# Learn about Mage in our CONTRIBUTING.md
go run mage.go EnsureMage
mage build
```

Then, to install the resulting mixin into PORTER_HOME, execute
`mage install`

## Mixin Configuration

```yaml
mixins:
- tofu:
    clientVersion: 1.0.3
    workingDir: myinfra
    initFile: providers.tf
```

### clientVersion

The OpenTofu client version can be specified via the `clientVersion` configuration when declaring this mixin.

### workingDir

The `workingDir` configuration setting is the relative path to your OpenTofu files. Defaults to "opentofu".

### initFile

OpenTofu providers are installed into the bundle during porter build.
We recommend that you put your provider declarations into a single file, e.g. "opentofu/providers.tf".
Then use `initFile` to specify the relative path to this file within workingDir.
This will dramatically improve Docker image layer caching and performance when building, publishing and installing the bundle.
> Note: this approach isn't suitable when using OpenTofu modules as those need to be "initilized" as well but aren't specified in the `initFile`. You shouldn't specifiy an `initFile` in this situation.

### User Agent Opt Out

When you declare the mixin, you can disable the mixin from customizing the azure user agent string

```yaml
mixins:
- tofu:
    userAgentOptOut: true
```

By default, the OpenTofu mixin adds the porter and mixin version to the user agent string used by the azure provider.
We use this to understand which version of porter and the mixin are being used by a bundle, and assist with troubleshooting.
Below is an example of what the user agent string looks like:

```bash
AZURE_HTTP_USER_AGENT="getporter/porter/v1.0.0 getporter/tofu/v1.2.3"
```

You can add your own custom strings to the user agent string by editing your [template Dockerfile] and setting the AZURE_HTTP_USER_AGENT environment variable.

[template Dockerfile]: https://getporter.org/bundle/custom-dockerfile/

## OpenTofu state

### Let Porter do the heavy lifting

The simplest way to use this mixin with Porter is to let Porter track the OpenTofu [state](https://opentofu.org/docs/language/state/) as actions are executed.  This can be done through the state section. Each time the bundle is executed, the output will capture the updated state file and inject it into the next action.

```yaml
state:
  - name: tfstate
    path: opentofu/opentofu.tfstate
  - name: tfvars
    path: opentofu/opentofu.tfvars.json
```

The specified path inside the installer (`/cnab/app/opentofu/opentofu.tfstate`) should be where OpenTofu will be looking to read/write its state.  For a full example bundle using this approach, see the [basic-tofu-example](examples/basic-tofu-example).

### Remote Backends

Alternatively, state can be managed by a remote backend.  When doing so, each action step needs to supply the remote backend config via `backendConfig`.  In the step examples below, the configuration has key/value pairs according to the [Azurerm](https://opentofu.org/docs/language/settings/backends/azurerm/) backend.

## OpenTofu variables file

By default the mixin will create a default [`opentofu.tfvars.json`](https://opentofu.org/docs/language/values/variables/)
file from the `vars` block during during the install step.

To use this file, a `tfvars` file parameter and output must be added to persist it for subsequent steps.

This can be disabled by setting `disableVarFile` to `true` during install.

Here is an example setup using the tfvar file:

```yaml
parameters:
  - name: tfvars
    type: file
    # This designates the path within the installer to place the parameter value
    path: /cnab/app/opentofu/opentofu.tfvars.json
    # Here we tell Porter that the value for this parameter should come from the 'tfvars' output
    source:
      output: tfvars
  - name: foo
    type: string
    applyTo:
      - install 
  - name: baz
    type: string
    default: blaz
    applyTo:
      - install 

outputs:
  - name: tfvars
    type: file
    # This designates the path within the installer to read the output from
    path: /cnab/app/opentofu/opentofu.tfvars.json
    
install:
  - tofu:
      description: "Install Azure Key Vault"
      vars:
        foo: bar
        baz: biz
      outputs:
      - name: vault_uri
upgrade: # No var block required
  - tofu:
      description: "Install Azure Key Vault"
      outputs:
      - name: vault_uri
uninstall: # No var block required
  - tofu:
      description: "Install Azure Key Vault"
      outputs:
      - name: vault_uri
```

and with var file disabled

```yaml
parameters:
  - name: foo
    type: string
    applyTo:
      - install 
  - name: baz
    type: string
    default: blaz
    applyTo:
      - install 

install:
  - tofu:
      description: "Install Azure Key Vault"
      disableVarFile: true
      vars:
        foo: bar
        baz: biz
      outputs:
      - name: vault_uri
uninstall: # Var block required
  - tofu:
      description: "Install Azure Key Vault"
      vars:
        foo: bar
        baz: biz
```

## Examples

### Install

```yaml
install:
  - tofu:
      description: "Install Azure Key Vault"
      backendConfig:
        key: "mybundle.tfstate"
        storage_account_name: "mystorageacct"
        container_name: "mycontainer"
        access_key: "myaccesskey"
      outputs:
      - name: vault_uri
```

### Upgrade

```yaml
upgrade:
  - tofu:
      description: "Upgrade Azure Key Vault"
      backendConfig:
        key: "mybundle.tfstate"
        storage_account_name: "mystorageacct"
        container_name: "mycontainer"
        access_key: "myaccesskey"
      outputs:
      - name: vault_uri
```

### Invoke

An invoke step is used for any custom action (not one of `install`, `upgrade` or `uninstall`).

By default, the command given to `tofu` will be the step name.  Here it is `show`,
resulting in `tofu show` with the provided configuration.

```yaml
show:
  - tofu:
      description: "Invoke 'tofu show'"
      backendConfig:
        key: "mybundle.tfstate"
        storage_account_name: "mystorageacct"
        container_name: "mycontainer"
        access_key: "myaccesskey"
```

Or, if the step name does not match the intended OpenTofu command, the command
can be supplied via the `arguments:` section, like so:

```yaml
printVersion:
  - tofu:
      description: "Invoke 'opentofu version'"
      arguments:
        - version
```

### Uninstall

```yaml
uninstall:
  - tofu:
      description: "Uninstall Azure Key Vault"
      backendConfig:
        key: "mybundle.tfstate"
        storage_account_name: "mystorageacct"
        container_name: "mycontainer"
        access_key: "myaccesskey"
```

See further examples in the [Examples](examples) directory

## Step Outputs

As seen above, outputs can be declared for a step.  All that is needed is the name of the output.

For each output listed, `opentofu output <output name>` is invoked to fetch the output value
from the state file for use by Porter. Outputs can be saved to the filesystem so that subsequent
steps can use the file by specifying the `destinationFile` field. This is particularly useful
when your OpenTofu module creates a Kubernetes cluster. In the example below, the module
creates a cluster, and then writes the kubeconfig to /root/.kube/config so that the rest of the
bundle can immediately use the cluster.

```yaml
install:
  - tofu:
      description: "Create a Kubernetes cluster"
      outputs:
      - name: kubeconfig
        destinationFile: /root/.kube/config
```

See the Porter [Outputs documentation](https://porter.sh/wiring/#outputs) on how to wire up
outputs for use in a bundle.
