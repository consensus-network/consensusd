# consensusminer

`consensusminer` is a CPU-based miner for `consensusd`.

## Requirements

Go 1.19 or later.

## Build from Source

* Install Go according to the installation instructions here:
  http://golang.org/doc/install

* Ensure Go was installed properly and is a supported version:

```bash
go version
```

* Run the following commands to obtain and install `consensusd`
  including all dependencies:

```bash
git clone https://github.com/consensus-network/consensusd
cd consensusd/cmd/consensusminer
go install .
```

* `consensusminer` should now be installed in `$(go env GOPATH)/bin`.
  If you did not already add the bin directory to your system path
  during Go installation, you are encouraged to do so now.
  
## Usage

The full `consensusminer` configuration options can be seen with:

```bash
consensusminer --help
```

But the minimum configuration needed to run it is:

```bash
consensusminer --miningaddr=<YOUR_MINING_ADDRESS>
```
