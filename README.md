# docker-registry-platformsh

A Docker Registry implementation for Platform.sh

## Required Build Variables

```sh
platform variable:create --level project --name=env:VARIABLE --value=VALUE --sensitive=true --visible-runtime=false --yes --no-wait
```

- PKI_SERVER_PRIVATE_KEY - the private key used by the token authentication server

## Example

```sh
cd ~/src
git clone https;//github.com/OpenVPN/easy-rsa.git
cd easy-rsa/easyrsa3
./easyrsa init-pki
./easyrsa build-ca
./easyrsa build-server-full token-auth nopass
cd ~/src/docker-registry-platformsh
platform variable:create --level project --name=env:PKI_SERVER_PRIVATE_KEY --sensitive=true --value="$(< $HOME/src/easy-rsa/easyrsa3/pki/private/token-auth.key)" --yes --no-wait
cat $HOME/src/easy-rsa/easyrsa3/pki/{ca.crt,auth/server.crt} > registry/bundle.crt
cp $HOME/src/easy-rsa/easyrsa3/pki/issued/token-auth.crt auth/server.crt
```

```sh
echo -n "PASSWORD" | docker login --password-stdin --username USERNAME registry.URL
```
