# Custom Transformers
Individual custom transformers can be composed together from any number of external repositories and executed as a
single process using the `compose` and `execute` commands. This is accomplished by
generating a Go plugin which allows the `vulcanizedb` binary to link to the external transformers, so long as they
abide by the [event transformer interface](../libraries/shared/factories/event/transformer.go) or the
[storage transformer interface](../libraries/shared/factories/storage/transformer.go).

## Writing custom transformers
For help with writing different types of custom transformers please see below:

Storage Transformers: transform data derived from contract storage tries
   * [Guide](../libraries/shared/factories/storage/README.md)
   * [Example](../libraries/shared/factories/storage/EXAMPLE.md)

Event Transformers: transform data derived from Ethereum log events
   * [Guide](../libraries/shared/factories/event/README.md)
   * [Example 1](https://github.com/vulcanize/ens_transformers/tree/master/transformers/registar)
   * [Example 2](https://github.com/vulcanize/ens_transformers/tree/master/transformers/registry)
   * [Example 3](https://github.com/vulcanize/ens_transformers/tree/master/transformers/resolver)

Contract Transformers: transform data derived from Ethereum log events provided only a contract address
   * [Example 1](https://github.com/vulcanize/account_transformers)
   * [Example 2](https://github.com/vulcanize/ens_transformers/tree/master/transformers/domain_records)

## Preparing custom transformers to work as part of a plugin
To plug in an external transformer we need to:

1. Create a package that exports a variable `EventTransformerInitializer`, `StorageTransformerInitializer`, or `ContractTransformerInitializer`
that are of type [event.TransformerInitializer](../libraries/shared/factories/event/transformer.go#L31)
or [storage.TransformerInitializer](../libraries/shared/factories/storage/transformer.go#L33),
or [ContractTransformerInitializer](../libraries/shared/transformer/contract_transformer.go#L31), respectively
2. Design the transformers to work in the context of their [event](../libraries/shared/watcher/event_watcher.go),
[storage](../libraries/shared/watcher/storage_watcher.go),
or [contract](../libraries/shared/watcher/contract_watcher.go) watcher `Execute` methods
3. Create db migrations to run against vulcanizeDB so that we can store the transformer output
    * Do not `goose fix` the transformer migrations, this is to ensure they are always ran after the core vulcanizedb migrations which are kept in their fixed form
    * Specify migration locations for each transformer in the config with the `exporter.transformer.migrations` fields
    * If the base vDB migrations occupy this path as well, they need to be in their `goose fix`ed form
    as they are [here](../../staging/db/migrations)

To update a plugin repository with changes to the core vulcanizedb repository, use your dependency manager to install the desired version of vDB.

## Building and Running Custom Transformers
### Commands
* The `compose` and `execute` commands require Go 1.11+ and use [Go plugins](https://golang
.org/pkg/plugin/) which only work on Unix-based systems.

* There is an ongoing [conflict](https://github.com/golang/go/issues/20481) between Go plugins and the use of vendored
dependencies which imposes certain limitations on how the plugins are built.

* Separate `compose` and `execute` commands allow pre-building and linking to the pre-built .so file. A couple of things
need to be considered:
    * It is necessary that the .so file was built with the same exact dependencies that are present in the execution
    environment, i.e. we need to `compose` and `execute` the plugin .so file with the same exact version of vulcanizeDB.
    * The plugin migrations are run during the plugin's composition. As such, if `execute` is used to run a prebuilt .so
    in a different environment than the one it was composed in, then the database structure will need to be loaded 
    into the environment's Postgres database. This can either be done by manually loading the plugin's schema into 
    Postgres, or by manually running the plugin's migrations.
     
* The `compose` command assumes you are in the vulcanizdb directory located at your system's 
`$GOPATH`, and that the plugin dependencies are present at their `$GOPATH` directories.

* The `execute` command does not require the plugin transformer dependencies be located in their `$GOPATH` directories,
instead it expects a .so file (of the name specified in the config file) to be in
`$GOPATH/src/github.com/makerdao/vulcanizedb/plugins/` and, as noted above, also expects the plugin db migrations to
 have already been ran against the database.

 * Usage:
     * compose: `./vulcanizedb compose --config=environments/config_name.toml`

     * execute: `./vulcanizedb execute --config=environments/config_name.toml`

### Flags
The `execute` command can be passed optional flags to specify the operation of the watchers:

- `--recheck-headers`/`-r` - specifies whether to re-check headers for events after the header has already been queried for watched logs.
Can be useful for redundancy if you suspect that your node is not always returning all desired logs on every query.
Argument is expected to be a boolean: e.g. `-r=true`.
Defaults to `false`.

### Configuration
A .toml config file is specified when executing the commands.
The config provides information for composing a set of transformers from external repositories:

```toml
[database]
    name     = "vulcanize_public"
    hostname = "localhost"
    user     = "vulcanize"
    password = "vulcanize"
    port     = 5432

[client]
    ipcPath  = "/Users/user/Library/Ethereum/geth.ipc"

[exporter]
    home     = "github.com/makerdao/vulcanizedb"
    name     = "exampleTransformerExporter"
    save     = false
    transformerNames = [
        "transformer1",
        "transformer2",
        "transformer3",
        "transformer4",
    ]
    [exporter.transformer1]
        path = "path/to/transformer1"
        type = "eth_event"
        repository = "github.com/account/repo"
        migrations = "db/migrations"
        rank = "0"
    [exporter.transformer2]
        path = "path/to/transformer2"
        type = "eth_contract"
        repository = "github.com/account/repo"
        migrations = "db/migrations"
        rank = "0"
    [exporter.transformer3]
        path = "path/to/transformer3"
        type = "eth_event"
        repository = "github.com/account/repo"
        migrations = "db/migrations"
        rank = "0"
    [exporter.transformer4]
        path = "path/to/transformer4"
        type = "eth_storage"
        repository = "github.com/account2/repo2"
        migrations = "to/db/migrations"
        rank = "1"
```
- `home` is the name of the package you are building the plugin for, in most cases this is github.com/makerdao/vulcanizedb
- `name` is the name used for the plugin files (.so and .go)   
- `save` indicates whether or not the user wants to save the .go file instead of removing it after .so compilation. Sometimes useful for debugging/trouble-shooting purposes.
- `transformerNames` is the list of the names of the transformers we are composing together, so we know how to access their submaps in the exporter map
- `exporter.<transformerName>`s are the sub-mappings containing config info for the transformers
    - `repository` is the path for the repository which contains the transformer and its `TransformerInitializer`
    - `path` is the relative path from `repository` to the transformer's `TransformerInitializer` directory (initializer package).
        - Transformer repositories need to be cloned into the user's $GOPATH (`go get`)
    - `type` is the type of the transformer; indicating which type of watcher it works with (for now, there are only two options: `eth_event` and `eth_storage`)
        - `eth_storage` indicates the transformer works with the [storage watcher](../libraries/shared/watcher/storage_watcher.go)
         that fetches state and storage diffs from an ETH node (instead of, for example, from IPFS)
        - `eth_event` indicates the transformer works with the [event watcher](../libraries/shared/watcher/event_watcher.go)
         that fetches event logs from an ETH node
        - `eth_contract` indicates the transformer works with the [contract watcher](../libraries/shared/watcher/contract_watcher.go)
        that is made to work with [contract_watcher pkg](../pkg/contract_watcher)
        based transformers which work with vDB to watch events provided only a contract address ([example1](https://github.com/vulcanize/account_transformers/tree/master/transformers/account/light), [example2](https://github.com/vulcanize/ens_transformers/tree/working/transformers/domain_records))
    - `migrations` is the relative path from `repository` to the db migrations directory for the transformer
    - `rank` determines the order that migrations are ran, with lower ranked migrations running first
        - this is to help isolate any potential conflicts between transformer migrations
        - start at "0" 
        - use strings
        - don't leave gaps
        - transformers with identical migrations/migration paths should share the same rank
- Note: If any of the imported transformers need additional config variables those need to be included as well   

This information is used to write and build a Go plugin which exports the configured transformers.
These transformers are loaded onto their specified watchers and executed.

Transformers of different types can be run together in the same command using a single config file or in separate instances using different config files   

The general structure of a plugin .go file, and what we would see built with the above config is shown below

```go
package main

import (
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
    "github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
    interface1 "github.com/makerdao/vulcanizedb/libraries/shared/transformer"
	transformer1 "github.com/account/repo/path/to/transformer1"
	transformer2 "github.com/account/repo/path/to/transformer2"
	transformer3 "github.com/account/repo/path/to/transformer3"
	transformer4 "github.com/account2/repo2/path/to/transformer4"
)

type exporter string

var Exporter exporter

func (e exporter) Export() []interface1.EventTransformerInitializer, []interface1.StorageTransformerInitializer, []interface1.ContractTransformerInitializer {
	return []interface1.TransformerInitializer{
            transformer1.TransformerInitializer,
            transformer3.TransformerInitializer,
        },     []interface1.StorageTransformerInitializer{
            transformer4.StorageTransformerInitializer,
        },     []interface1.ContractTransformerInitializer{
            transformer2.TransformerInitializer,
        }
}
```
