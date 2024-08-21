# kubectl-pvcmount

Go binary for temporary mounting of PVCs to a pod. This is useful for debugging purposes or for copying data from a PVC to another location.

## Installation via Homebrew (MacOS/Linux - x86_64/arm64)

```bash
brew install stenstromen/tap/kubectl-pvcmount
```

## Download and Run Binary

* For **MacOS** and **Linux**: Checkout and download the latest binary from [Releases page](https://github.com/Stenstromen/kubectl-pvcmount/releases/latest/)
* For **Windows**: Build the binary yourself.

## Build and Run Binary

```bash
go build
./kubectl-pvcmount
```

## Example Usage

```bash
kubectl-pvcmount pvc --namespace my-namespace --pvc my-pvc

---

Go binary for temporary mounting of PVCs

Usage:
  kubectl-pvcmount [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  pvc         Mount Persistent Volume to a temporary container

Flags:
  -h, --help   help for kubectl-pvcmount

Use "kubectl-pvcmount [command] --help" for more information about a command.
```
