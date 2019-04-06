# docker-registry-platformsh

A Docker Registry implementation for Platform.sh

## Required Environment Variables

```sh
platform variable:create --level project --name=env:VARIABLE --value=VALUE --sensitive=true --yes --no-wait
```

- env:PKI_SERVER_PRIVATE_KEY - the private key used by the token authentication server
- users - a JSON dictionary of user: password pairs

## Example

```sh
cd ~/src
git clone https://github.com/OpenVPN/easy-rsa.git
git clone https://github.com/demosdemon/docker-registry-platformsh

cd ~/src/easy-rsa/easyrsa3
./easyrsa init-pki
./easyrsa build-ca
./easyrsa build-server-full token-auth nopass

cd ~/src/docker-registry-platformsh
platform project:create --title "Docker Registry" --region=us-2.platform.sh --yes > .platform/project
platform project:set-remote $(< .platform/project)
platform variable:create --level project --name=env:PKI_SERVER_PRIVATE_KEY --sensitive=true --value="$(< $HOME/src/easy-rsa/easyrsa3/pki/private/token-auth.key)" --yes --no-wait
platform var:create --level=environment --name=users --value="{\"brandon\": \"$PW\"}" --json=true --sensitive=true --no-wait --yes
cat $HOME/src/easy-rsa/easyrsa3/pki/{ca.crt,auth/server.crt} > registry/bundle.crt
cp $HOME/src/easy-rsa/easyrsa3/pki/issued/token-auth.crt auth/server.crt
# configure auth/config.yml
# (see https://github.com/cesanta/docker_auth/tree/master/examples)
git add .
git commit -m 'customize'
platform push --force-with-lease --set-upstream --wait
```

```sh
echo -n "$PW" | docker login --password-stdin --username brandon registry.docker.apimxdx.com
```
