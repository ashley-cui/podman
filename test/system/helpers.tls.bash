# -*- bash -*-
#
# BATS helpers for tls functionality
#
function cert-pair {
  cn=$1 key=$2 cert=$3
  shift 3
  openssl req -x509 \
		-nodes \
		-newkey rsa:4096 -keyout "${key}" \
		-out "${cert}" \
		-days 1 \
		-subj "/C=??/ST=System/L=Test/O=Containers/OU=Podman/CN=${cn}" \
    -quiet \
    "$@"
}

function signed-cert-pair {
  cn=$1 key=$2 cert=$3 ca_key=$4 ca_cert=$5
  shift 5
  cert-pair "${cn}" \
    "${key}" "${cert}" \
		-CAkey "${ca_key}" -CA "${ca_cert}" \
    "$@"
}

function tls-certs {
    mkdir -p $PODMAN_TMPDIR/certs
  # CA
  cert-pair "ca" \
    "${CA_KEY}" "${CA_CRT}" \
		-addext basicConstraints=critical,CA:TRUE,pathlen:1
  # Client, signed by CA
  signed-cert-pair "client" \
    "${CLIENT_KEY}" "${CLIENT_CRT}" \
    "${CA_KEY}" "${CA_CRT}" \

  # Server, signed by CA, valid for localhost, 127.0.0.1
	# NOTE: Go refuses certs without SAN's
  signed-cert-pair "localhost" \
    "${SERVER_KEY}" "${SERVER_CRT}" \
    "${CA_KEY}" "${CA_CRT}" \
		-addext "subjectAltName=DNS:localhost,IP:127.0.0.1"

  # Bogus, self-signed
  cert-pair "bogus" \
    "${BOGUS_KEY}" "${BOGUS_CRT}" \
		-addext basicConstraints=critical,CA:TRUE,pathlen:1
}
