# Initial Setup

## Create devspace.yaml
The first step after installing DevSpace is to create a `devspace.yaml` configuration file. This provides the DevSpace CLI with all the information it needs to help you develop your application. DevSpace uses a mix of declarative and imperative instructions. We'll focus on the declarative instructions first and then use devspace to peek under the hood and see the imperative instructions that are automatically added.

Create a file titled `devspace.yaml` and save the following content to it:
```yaml
version: v2beta1
```

This version is the schema version. In order to support upgrades between versions a schema version is used. This allows DevSpace to automatically "upgrade" to the latest schema and maintain backwards compatibility as the schema evolves. This is also the absolute bare minimum configuration needed to run the DevSpace CLI.

## `devspace print`
The first DevSpace command we'll run is `devspace print`. This is mostly used to show the complete configuration and help understand what DevSpace will do when you run other commands. Go ahead and run it.

```sh
$ devspace print

-------------------

Vars:
  
               NAME             |                       VALUE                        
  ------------------------------+----------------------------------------------------
    DEVSPACE_CONTEXT            | kind-kind                                          
    DEVSPACE_EXECUTABLE         | /opt/homebrew/bin/devspace                         
    DEVSPACE_KUBECTL_EXECUTABLE | kubectl                                            
    DEVSPACE_NAME               | 01-initial                                         
    DEVSPACE_NAMESPACE          | devspace                                               
    DEVSPACE_PROFILE            |                                                    
    DEVSPACE_PROFILES           |                                                    
    DEVSPACE_RANDOM             | QGPjoO                                             
    DEVSPACE_TIMESTAMP          | 1671577341                                         
    DEVSPACE_TMPDIR             | /var/folders/3k/f2wz_n393392ky9z0999d4bm0000gn/T/  
    DEVSPACE_USER_HOME          | /Users/russellcentanni                             
    DEVSPACE_VERSION            | 6.2.3                                       
    devspace.context            | kind-kind                                          
    devspace.namespace          | devspace                                               
  

-------------------

Loaded path: /Users/russellcentanni/Projects/devspace-example-go/01-initial/devspace.yaml

-------------------

version: v2beta1
name: 01-initial
```

The top section shows [built-in variables](https://www.devspace.sh/docs/configuration/variables#built-in-variables) that are available to use in your configuration. We'll use some of those as we add more configuration.

The bottom section shows the complete configuration at runtime. This will include dynamic values that are determined at runtime. For example, the `name` property has the value `01-initial`, which is the current working directory. Go ahead and override this value by adding it to `devspace.yaml` and run `devspace print` again. You'll see that it now uses the value from your `devspace.yaml` instead of the default value. Using this knowledge, if DevSpace does something by default that you'd like to change, you can use `devspace print` to find out what configuration to override.

You may have noticed that DevSpace `Vars` section has `DEVSPACE_NAMESPACE` and `devspace.namespace` variables, as well as `DEVSPACE_CONTEXT` and `devspace.context`. These are derived from your kube config, and can be changed by running `devspace use context` and `devspace use namespace`. When run without passing a name you'll be able to select from existing options. If you add a name, such as `devspace use context kind-kind`, it will update DevSpace to use that context. In the case of `devspace use namespace`, you can add a name that doesn't exist, and DevSpace will create the namespace for you, as we'll see in the next section.

## `devspace dev`
Running `devspace print` only shows the configuration after subsituting variables and applying profiles (more on this later!). To do something slightly more interesting, run:

```sh
$ devspace use namespace russ
$ devspace dev
info Using namespace 'russ'
info Using kube context 'kind-kind'
info Created namespace: russ
```

We can see that DevSpace has created the namespace `russ`. This doesn't seem like much, but that's only because it's the only thing we've configured DevSpace to do. Next we'll make a deployment from which we can develop our golang based web application.