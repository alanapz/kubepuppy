COPY 

/var/lib/rancher/k3s/server/tls/client-ca.crt
/var/lib/rancher/k3s/server/tls/client-ca.key
/var/lib/rancher/k3s/server/tls/server-ca.crt

openssl ecparam -genkey -name prime256v1 -noout -out alan.key

openssl req -new -key alan.key -out alan.csr -subj "/CN=system:admin/O=system:masters"

openssl x509 -req -in alan.csr -CA client-ca.crt -CAkey client-ca.key -CAcreateserial -out alan.crt -days 500

To test:

curl -vvvv --cacert client-ca.crt --cert alan.crt --key alan.key https://apz-dev-u22.vbox:6443/api/v1

kubectl get clusterrolebinding -A -o json | jq ".items[].subjects" | more
