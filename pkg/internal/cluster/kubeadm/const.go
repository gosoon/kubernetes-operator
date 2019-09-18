package kubeadm

// APIServerPort is the expected default APIServerPort on the control plane node(s)
// https://kubernetes.io/docs/reference/access-authn-authz/controlling-access/#api-server-ports-and-ips
const APIServerPort = 6443

// Token defines a dummy, well known token for automating TLS bootstrap process
const Token = "abcdef.0123456789abcdef"

// ObjectName is the name every generated object will have
// I.E. `metadata:\nname: config`
const ObjectName = "config"
