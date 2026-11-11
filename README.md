# cockpit — c12s operator CLI

Cobra-based command-line client for the c12s control plane. Everything an
operator does — register/login, claim nodes, manage labels/schemas/configs,
declare OS image updates and rollouts, manage VM fleets — goes through this
binary via the lunar-gateway HTTP API.

## Role in c12s

| Command group | Backend service |
|---|---|
| `register`, `login` | Apollo (identity; token stored in `token.txt`) |
| `list/claim/delete nodes`, labels, metrics | Magnetar |
| `create/put/delete schema`, `validate` | Quasar |
| config groups and standalone configs (`put`, `list`, `place`, `diff`) | Kuiper (+ Gravity for delivery) |
| `image update/rollout/rollback/status` | Gravity (`DisseminateOSImage`), Magnetar (desired image, status) |
| `image build`, `image build qcow2`, `image build status/list` | Supernova (`BuildImage`, `BuildQCOW2`, `GetJob`, `ListJobs`) |
| `image list` | Supernova (`ListImages`; registry + artifact server) |
| `vms set/list/status` | Magnetar (VM desired state / status on the host's node record) |

## Building and running

1. The control plane must be up (`cd tools && ./start.sh` in the c12s
   workspace; gateway on `localhost:5555`).

2. Build the CLI:

    ```sh
    go build -o ~/.cache/c12s/cockpit .
    ```

3. Run it **from the repository directory** — it reads two local files:

    | File | Purpose |
    |---|---|
    | `.env` | `CONFIG_PATH=../lunar-gateway/config/config.yml` — the gateway routes and ports |
    | `token.txt` | bearer token written by `cockpit login` (mode 0600) |

4. Log in (the lab user is `e2e` / `testpass123`):

    ```sh
    cockpit login --username 'e2e' --password 'testpass123'
    # or: C12S_PASSWORD=testpass123 cockpit login --username 'e2e'
    ```

## Features

- **User Management**: Register and login to the system.
- **Node Management**: List, query, and claim nodes; delete records (with the `node.delete` permission). Lists show each node's liveness (`online`/`stale`/`offline`/`unknown`) from its last heartbeat.
- **Label Management**: Add, update, and delete labels for nodes.
- **Schema Management**: Create, retrieve, validate, and delete schemas.
- **Configuration Management**: Manage both config groups and standalone configs, including placing, listing, and diffing configurations.
- **OS Image Management**: Build images and qcow2 disks through Supernova, then apply digest-pinned updates, canary-gated rollouts, and rollback for `os_update`-capable nodes.
- **VM Management**: Declare the desired VM fleet on a `run_vms` host and inspect observed status.

## Table of Contents
- [Getting Started](#getting-started)
- [Command Reference](#command-reference)
  - [User Management](#user-management)
  - [Node Management](#node-management)
  - [OS Image Management (bootc)](#os-image-management-bootc)
  - [VM Management](#vm-management)
  - [Relationship Management](#relationship-management)
  - [Label Management](#label-management)
  - [Schema Management](#schema-management)
  - [Config Group Management](#config-group-management)
  - [Standalone Config Management](#standalone-config-management)
  - [Node Metrics Management](#node-metrics-management)
- [Bugs, disabled and unfinished features](#bugs-disabled-and-unfinished-features)
- [Contributing](#contributing)
- [License](#license)

## Getting Started

### Prerequisites
Ensure you have the following installed:
- [Go](https://golang.org/doc/install) (version 1.22 or higher)
- A running control plane (the [c12s/tools](https://github.com/c12s/tools) repository provides `./install.sh` + `./start.sh` for the full compose stack)

### Installation

Clone this repository next to the other c12s repos (the workspace layout is
`c12s/tools`, `c12s/cockpit`, `c12s/lunar-gateway`, ... — cockpit reads the
gateway route config via a relative path):

```sh
git clone https://github.com/c12s/cockpit.git
cd cockpit
go build -o ~/.cache/c12s/cockpit .
```

Run the binary from the cockpit directory (it needs `.env` and writes
`token.txt` there), or use the helper from the lab experiments:

```sh
cd /path/to/c12s/cockpit && ~/.cache/c12s/cockpit login --username e2e --password testpass123
```

## Command Reference

### User Management

#### Register
Register a new user.
- **Command**: cockpit register
- **Options**:
  - --email: Email address of the user.
  - --name: Name of the user.
  - --org: Organization of the user.
  - --surname: Surname of the user.
  - --username: Username of the user.
- **Example**:

    ```sh
    cockpit register --email 'user@gmail.com' --name 'name' --org 'c12s' --surname 'surname' --username 'user'
    ```

#### Login
Login with an existing user.
- **Command**: cockpit login
- **Options**:
  - --username: Username of the user.
  - --password, -p: Password (optional; also read from `C12S_PASSWORD`,
    otherwise prompted).
- **Example**:

    ```sh
    cockpit login --username 'user'
    C12S_PASSWORD=secret cockpit login --username 'user'
    ```

### Node Management

#### List Nodes
List all nodes. The table shows a `Liveness` column derived from the node's
last heartbeat (`online` within 90 s, `stale` within 5 min, `offline` after
that, `unknown` for records that never heartbeated; thresholds are Magnetar's
`NODE_STALE_AFTER`/`NODE_OFFLINE_AFTER`).
- **Command**: cockpit list nodes
- **Options**:
  - --query: Query to filter nodes.
  - --details: show liveness age and every label per node.
- **Example**:

    ```sh
    cockpit list nodes
    cockpit list nodes --details
    cockpit list nodes --query 'memory-totalGB > 2'
    ```

#### Claim Nodes
Claim nodes based on specific criteria.
- **Command**: cockpit claim nodes
- **Options**:
  - --org: Organization.
  - --query: Query to filter nodes.
- **Example**:

    ```sh
    cockpit claim nodes --org 'c12s' --query 'memory-totalGB > 2'
    ```

#### Delete a Node
Remove a node record completely: the record, its label index, liveness state
and desired OS image. Org-owned nodes are never removed automatically (the
control-plane GC only reaps stale *pool* records), so this is the explicit
path for them. Requires the `node.delete` permission on the node — create a
policy first, like for `os.update`:

```yaml
objectScope: { id: "<node-id>", kind: "node" }
permission:  { condition: { expression: "" }, kind: "ALLOW", name: "node.delete" }
subjectScope: { id: "e2e", kind: "user" }
```

- **Command**: cockpit delete node
- **Options**:
  - --node-id: Node ID (required).
- **Example**:

    ```sh
    cockpit delete node --node-id <id>
    ```

Deleting a host does not stop the VMs it declared: they keep running and keep
their own node records.

#### List Allocated Nodes
List nodes allocated to an organization.
- **Command**: cockpit list nodes allocated
- **Options**:
  - --org: Organization.
  - --query: Query to filter nodes.
- **Example**:

    ```sh
    cockpit list nodes allocated --org 'c12s'
    cockpit list nodes allocated --org 'c12s' --query 'memory-totalGB > 2'
    ```

### OS Image Management (bootc)

Nodes with the `os_update` capability run bootc; images are pinned by digest.
Image refs may be short (`os/edge-base:v4.1`): the node completes them with
`IMAGE_REGISTRY` from its image's `node.env`, so the control plane stores no
registry address. A ref that names a host (`host:port/…`) overrides it.
Updating a node needs `os.update` for the caller (`create policies`): on the
node's org (`objectScope: {kind: org, id: <org>}`, covers every node claimed
into it) or on the node itself (`{kind: node, id: <node-id>}`, also works for
pool nodes). Apollo mints the permissions on every request, so a new policy
applies immediately.

#### Build an image
```bash
cockpit image build --name os/edge-base --tag v3.18 \
  --context tools/bootc --containerfile Containerfile.bootc \
  --build-arg OS_VERSION=v3.18 [--wait] [--timeout 30m]
```
Sends the build context to Supernova through the gateway, which builds, pushes
and signs the image in Singularity and returns the digest. Without `--wait` it
prints the job id; `cockpit image build status <job-id>` shows progress and
`cockpit image build list` shows all jobs. The build needs an `image.build`
policy on the image repo (`image.build|image|os/edge-base`), created with
`create policies` and baked into the token at the next `login`.

#### Build a qcow2 provisioning disk
```bash
cockpit image build qcow2 --name edge-base-v3.18 \
  --image-ref 192.168.122.1:5001/os/edge-base:v3.18 [--target-arch x86_64] [--rootfs ext4]
```
Turns an already-built bootc image into a qcow2 disk, uploads it to the
Singularity artifact server and prints the artifact name and checksum. Use the
artifact URL + checksum in `cockpit vms set`.

#### Update one node
```bash
cockpit image update --node-id <id> --image-ref 192.168.122.1:5001/os/edge-base:v3.10 --digest sha256:...
```
The desired image is recorded in Magnetar first, then the node is notified.
The node stages the image, reboots, and greenboot health-gates it; a failing
image is rolled back automatically and reported as `rolled_back`. A node that
is offline does not acknowledge, but the update is still recorded and the node
applies it when it next pulls (every 30 s once it is back).

#### Canary rollout
```bash
cockpit image rollout --org c12s --canary <node-id> \
  --nodes <id1,id2> | --query 'labelKey = value' \
  --image-ref <ref> --digest sha256:... [--soak 2m] [--timeout 15m] [--poll 10s] [--no-wait]
```
Updates the canary first; it must boot the digest, report `healthy`, and stay
healthy for `--soak`. Only then are the targets (os_update-capable org nodes,
canary excluded) updated and awaited. A canary rollback, leaving healthy
during the soak, or a timeout stops the rollout with exit code 1 before any
target is touched.

#### Apply a fleet spec (declarative)
```bash
cockpit image apply --spec fleet.yaml --dry-run   # the plan
cockpit image apply --spec fleet.yaml             # make it so
```
```yaml
apiVersion: c12s/os.v1
kind: HostImageSpec
metadata: {org: c12s}
spec:
  images:
    stable: {ref: os/edge-base:v4.1}   # short ref: nodes pull from their IMAGE_REGISTRY; digest optional
  nodes:
    - {id: <node-id>, name: host-1, image: stable}
    - id: <vm-host-id>
      name: node-1
      image: stable
      vms:                                  # the host's complete VM set
        - {name: vm-01, disk: edge-base-v4.0.qcow2, vcpus: 1, memoryMb: 1024, image: stable}
  rollout: {soak: 2m, timeout: 15m}
```
Only the listed nodes are managed; removing a node from the file stops
managing it (nothing is changed on it). `image` is optional on nodes and VMs
(absent: OS not managed). `vms:` absent leaves a host's VMs alone, `vms: []`
deletes them. A VM's `disk` (an artifact name, or a URL with `digest`) is used
only to create it; its node id is derived from the host id and the VM name
(UUIDv5), so it is known before the VM exists.

Order: VM sets → wait for new VMs to enroll and claim them (label
`vm-host = <host-id>`) → OS for plain nodes and VMs → OS for VM hosts. Per
image group the first node is the canary (gate + soak, as in `rollout`). A
rollback or timeout stops the run (exit 1); re-running continues. `--force`
re-sends updates whose desired image already matches. Every node and VM with
an `image` needs `os.update`: one policy on the org (`kind: org`) covers all
its nodes, or one per node (`--dry-run` lists the ids).

#### Status
```bash
cockpit image status --org c12s
```
Per node: phase (`healthy`, `switching`, `degraded`, `rolled_back`,
`unmanaged`), version, BOOTED and DESIRED as `tag (short digest)` (DESIRED is
the image recorded in Magnetar; `-` = none declared yet) and the last
attempt. `--wide` prints full digests.

#### List images
```bash
cockpit image list
```
Read-only view of the registry and the artifact server, through the gateway's
`ListImages` route: every image tag with the digest it points to (the
value `image update --digest` takes) and every qcow2 artifact with its sha256.
It uses the same token as the other commands; `--output json` prints the raw
response.

#### Rollback
```bash
cockpit image rollback --org <org> --node-id <id>
```
Declares the node's previous deployment (the bootc rollback entry it
reports) as its desired image: the node switches back through the normal
update path (greenboot gate, attempt record) and stays managed. Refused when
that deployment is the image greenboot rejected, during an update, or when
the node reports no previous deployment.

#### Declare the desired VM set
```bash
cockpit vms set --host <host-node-id> --spec '{"vms":[{"name":"vm-01","qcow2Url":"edge-base-v4.0.qcow2","digest":"sha256:...","vcpus":2,"memoryMb":2048}]}'
```
`--spec` accepts inline JSON or a path to a YAML/JSON file. Each VM needs
`name`, `qcow2Url` and a `sha256:<64 hex>` digest; `vcpus`/`memoryMb` are
optional. `qcow2Url` may be a bare artifact name: the host downloads it from
`DISK_REGISTRY_URL` (image default `http://192.168.122.1:8081`). New VMs boot
from the downloaded disk and enroll as c12s nodes with id
UUIDv5(host id, VM name) and labels `vm-host`/`vm-name`. `image apply` does
the same declaratively, together with the OS.

#### Inspect
```bash
cockpit vms list --host <host-node-id>     # desired set (in Magnetar) + each VM's node id
cockpit vms status --host <host-node-id> [--org <org>]   # observed state; with --org also each VM's booted OS
```

### Relationship Management

#### Create Relations
Create relations between entities.
- **Command**: cockpit create relations
- **Options**:
  - --ids: IDs of the entities.
  - --kinds: Kinds of the entities.
- **Example**:

    ```sh
    cockpit create relations --ids 'c12s|dev' --kinds 'org|namespace'
    ```

### Label Management

#### Add Label
Add a label to a node.
- **Command**: cockpit put label
- **Options**:
  - --org: Organization of the node.
  - --node-id: ID of the node.
  - --key: Key of the label.
  - --value: Value of the label.
- **Example**:

    ```sh
    cockpit put label --org 'c12s' --node-id 'nodeId' --key 'newlabel' --value '25.00'
    ```

#### Delete Label
Delete a label from a node.
- **Command**: cockpit delete label
- **Options**:
  - --org: Organization of the node.
  - --nodeId: ID of the node.
  - --key: Key of the label.
- **Example**:

    ```sh
    cockpit delete label --org 'c12s' --node-id 'nodeId' --key 'newlabel'
    ```

### Schema Management

#### Create Schema
Create a new schema.
- **Command**: cockpit create schema
- **Options**:
  - --org: Organization.
  - --namespace: Namespace.
  - --schema-name: Name of the schema.
  - --version: Version of the schema.
  - --path: Path to the schema file.
- **Example**:

    ```sh
    cockpit create schema --org 'c12s' --namespace 'default' --schema-name 'schema' --version 'v1.0.0' --path 'request/schema/create-schema.yaml'
    ```

#### Get Schema Version
Retrieve schema version details.
- **Command**: cockpit get schema version
- **Options**:
  - --org: Organization.
  - --namespace: Namespace.
  - --schema: Name of the schema.
- **Example**:

    ```sh
    cockpit get schema version --org 'c12s' --namespace 'default' --schema-name 'schema'
    ```

#### Get Schema
Retrieve a specific schema.
- **Command**: cockpit get schema
- **Options**:
  - --org: Organization.
  - --namespace: Namespace.
  - --schema: Name of the schema.
  - --version: Version of the schema.
- **Example**:

    ```sh
    cockpit get schema --org 'c12s' --namespace 'default' --schema-name 'schema' --version 'v1.0.0'
    ```

#### Validate Schema
Validate a schema.
- **Command**: cockpit validate schema
- **Options**:
  - --org: Organization.
  - --namespace: Namespace.
  - --schema: Name of the schema.
  - --version: Version of the schema.
  - --path: Path to the validation file.
- **Example**:

    ```sh
    cockpit validate schema --org 'c12s' --namespace 'default' --schema-name 'schema' --version 'v1.0.0' --path 'request/schema/validate-schema.yaml'
    ```

#### Delete Schema
Delete a schema.
- **Command**: cockpit delete schema
- **Options**:
  - --org: Organization.
  - --namespace: Namespace.
  - --schema: Name of the schema.
  - --version: Version of the schema.
- **Example**:

    ```sh
    cockpit delete schema --org 'c12s' --namespace 'default' --schema-name 'schema' --version 'v1.0.0'
    ```

### Config Group Management

#### Add Config Group
Add a configuration group.
- **Command**: cockpit put config group
- **Options**:
  - --path: Path to the config group file.
- **Example**:

    ```sh
    cockpit put config group --path 'request/config-group/create-config-group.yaml'
    ```

#### Get Config Group
Retrieve a configuration group.
- **Command**: cockpit get config group
- **Options**:
  - --org: Organization.
  - --namespace: Namespace.
  - --name: Name of the config group.
  - --version: Version of the config group.
  - --output: Output format (json, yaml).
- **Example**:

    ```sh
    cockpit get config group --org 'c12s' --namespace 'default' --name 'app_config' --version 'v1.0.1'
    ```

#### List Config Groups
List all configuration groups.
- **Command**: cockpit list config group
- **Options**:
  - --org: Organization.
  - --namespace: Namespace.
  - --output: Output format (json, yaml).
- **Example**:

    ```sh
    cockpit list config group --org 'c12s' --namespace 'default'
    ```

#### Diff Config Groups
Compare differences between configuration groups.
- **Command**: cockpit diff config group
- **Options**:
  - --org: Organization.
  - --namespace: Namespace.
  - --names: Names of the config groups.
  - --versions: Versions of the config groups.
  - --output: Output format (json, yaml).
- **Example**:

    ```sh
    cockpit diff config group --org 'c12s' --namespace 'default' --names 'app_config|app_config' --versions 'v1.0.0|v1.0.1'
    cockpit diff config group --org 'c12s' --namespace 'default' --names 'app_config' --versions 'v1.0.0|v1.0.1'
    cockpit diff config group --org 'c12s' --namespace 'default' --names 'app_config|app_config' --versions 'v1.0.0'
    ```

#### Place Config Group
Place a configuration group.
- **Command**: cockpit place config group
- **Options**:
  - --path: Path to the placement file.
- **Example**:

    ```sh
    cockpit place config group --path 'request/config-group/create-config-group-placements.yaml'
    ```

#### List Config Group Placements
List all placements of a configuration group.
- **Command**: cockpit list config group placements
- **Options**:
  - --org: Organization.
  - --namespace: Namespace.
  - --name: Name of the config group.
  - --version: Version of the config group.
- **Example**:

    ```sh
    cockpit list config group placements --org 'c12s' --namespace 'default' --name 'app_config' --version 'v1.0.0'
    ```

#### Delete Config Group
Delete a configuration group.
- **Command**: cockpit delete config group
- **Options**:
  - --org: Organization.
  - --namespace: Namespace.
  - --name: Name of the config group.
  - --version: Version of the config group.
  - --output: Output format (json, yaml).
- **Example**:

    ```sh
    cockpit delete config group --org 'c12s' --namespace 'default' --name 'app_config' --version 'v1.0.1'
    ```

### Standalone Config Management

#### Add Standalone Config
Add a standalone configuration.
- **Command**: cockpit put standalone config
- **Options**:
  - --path: Path to the config file.
- **Example**:

    ```sh
    cockpit put standalone config --path 'request/standalone-config/create-standalone-config.json'
    ```

#### Get Standalone Config
Retrieve a standalone configuration.
- **Command**: cockpit get standalone config
- **Options**:
  - --org: Organization.
  - --namespace: Namespace.
  - --name: Name of the config.
  - --version: Version of the config.
  - --output: Output format (json, yaml).
- **Example**:

    ```sh
    cockpit get standalone config --org 'c12s' --namespace 'default' --name 'db_config' --version 'v1.0.1'
    ```

#### List Standalone Configs
List all standalone configurations.
- **Command**: cockpit list standalone config
- **Options**:
  - --org: Organization.
  - --namespace: Namespace.
  - --output: Output format (json, yaml).
- **Example**:

    ```sh
    cockpit list standalone config --org 'c12s' --namespace 'default'
    ```

#### Diff Standalone Configs
Compare differences between standalone configurations.
- **Command**: cockpit diff standalone config
- **Options**:
  - --org: Organization.
  - --namespace: Namespace.
  - --names: Names of the configs.
  - --versions: Versions of the configs.
  - --output: Output format (json, yaml).
- **Example**:

    ```sh
    cockpit diff standalone config --org 'c12s' --namespace 'default' --names 'db_config|db_config' --versions 'v1.0.1|v1.0.0'
    cockpit diff standalone config --org 'c12s' --namespace 'default' --names 'db_config' --versions 'v1.0.1|v1.0.0'
    cockpit diff standalone config --org 'c12s' --namespace 'default' --names 'db_config|db_config' --versions 'v1.0.1'
    ```

#### Place Standalone Config
Place a standalone configuration.
- **Command**: cockpit place standalone config
- **Options**:
  - --path: Path to the placement file.
- **Example**:

    ```sh
    cockpit place standalone config --path 'request/standalone-config/create-standalone-config-placements.yaml'
    ```

#### List Standalone Config Placements
List all placements of a standalone configuration.
- **Command**: cockpit list standalone config placements
- **Options**:
  - --org: Organization.
  - --namespace: Namespace.
  - --name: Name of the config.
  - --version: Version of the config.
- **Example**:

    ```sh
    cockpit list standalone config placements --org 'c12s' --namespace 'default' --name 'db_config' --version 'v1.0.0'
    ```

#### Delete Standalone Config
Delete a standalone configuration.
- **Command**: cockpit delete standalone config
- **Options**:
  - --org: Organization.
  - --namespace: Namespace.
  - --name: Name of the config.
  - --version: Version of the config.
  - --output: Output format (json, yaml).
- **Example**:

    ```sh
    cockpit delete standalone config --org 'c12s' --namespace 'default' --name 'db_config' --version 'v1.0.1'
    ```

### Node Metrics Management

#### Get Node Metrics
Retrieve metrics for a specific node.
- **Command**: cockpit get node metrics
- **Options**:
  - --node-id: Node ID (required).
  - --all: Display all metrics (optional).
  - --sort: Sort metrics by 'cpu', 'memory', 'disk', 'network receive', 'network transmit' or 'bandwidth'.
- **Example**:

    ```sh
    cockpit get node metrics --node-id 'nodeID'
    cockpit get node metrics --node-id 'nodeID' --all-services --sort 'memory'
    ```

## Status

Stable for all documented commands; the `image` and `vms` command groups
were added in Phases 3–4 and are exercised end-to-end by the lab
experiments (`tools/experiments/`).

## Bugs, disabled and unfinished features

### Confirmed bugs

| Where | Effect |
|---|---|
| `cmd/image/image.status.go`, `cmd/vm/vm.list.go`, `cmd/vm/vm.status.go` | These read commands send a JSON **body on a GET request**. It works through the current gateway (which forwards bodies), but it violates HTTP semantics and some proxies strip GET bodies (review BUG-3, deferred; a query-param or POST rework is the fix). |
| `image update`/`rollback` accepted:false | "did not acknowledge (node offline?)" is printed for *any* no-ack — including a node that rejected an invalid command silently (review R4). The real state is visible via `image status`. |

### Disabled / unfinished features

- HTTP + token boilerplate is duplicated across commands; only
  `update`/`rollout` share `requestImageUpdate` (review M1).
- `vms` commands do not validate that the host actually has the `run_vms`
  capability — the declaration is stored regardless.

### PoC boundaries (by design, not bugs)

- The CLI talks plain HTTP to the gateway and runs from the repository
  directory (`.env`, `token.txt` are local files).
- `--password` on the command line is visible in the shell history/process
  list; prefer `C12S_PASSWORD`.

## Contributing

Contributions are welcome! Please follow these steps to contribute:

1. Fork the repository.
2. Create a new branch:

    ```sh
    git checkout -b feature/YourFeature
    ```

3. Make your changes.
4. Commit your changes:

    ```sh
    git commit -m 'Add new feature'
    ```

5. Push to the branch:

    ```sh
    git push origin feature/YourFeature
    ```

6. Open a pull request and describe your changes.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more information.
