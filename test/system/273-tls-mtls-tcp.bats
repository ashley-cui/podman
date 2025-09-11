#
# Tests that spot check connectivity for each of the supported remote transports,
# unix, tcp, tls, mtls

load helpers
load helpers.systemd
load helpers.network
load helpers.tls

SERVICE_NAME="podman-service-$(random_string)"


function setup() {
  basic_setup
  mkdir -p $PODMAN_TMPDIR/certs
}

function teardown() {
  # Ignore exit status: this is just a backup stop in case tests failed
  run systemctl stop "$SERVICE_NAME"
  rm -f $PODMAN_TMPDIR/myunix.sock
  rm -rf $PODMAN_TMPDIR/certs
  basic_teardown
}

@test "unix" {
  skip_if_remote "testing unix service only works locally"
  URL=unix:$PODMAN_TMPDIR/myunix.sock

  systemd-run --unit=$SERVICE_NAME ${PODMAN%%-remote*} system service $URL --time=0
  wait_for_file $PODMAN_TMPDIR/myunix.sock

  run_podman --url="$URL" info --format '{{.Host.RemoteSocket.Path}}'
  is "$output" "$URL" "RemoteSocket.Path using unix:"
  # Streaming command works
  run_podman --url="$URL" run --rm -i $IMAGE /bin/sh -c 'echo -n foo; sleep 0.1; echo -n bar; sleep 0.1; echo -n baz'
  is "$output" foobarbaz

}

@test "tcp" {
  skip_if_remote "testing tcp service only works locally"
  port=$(random_free_port)
  URL=tcp://127.0.0.1:$port

  systemd-run --unit=$SERVICE_NAME ${PODMAN%%-remote*} system service $URL --time=0
  wait_for_port 127.0.0.1 $port

  # Flag works
  run_podman --url="$URL" info --format '{{.Host.RemoteSocket.Path}}'
  is "$output" "$URL" "RemoteSocket.Path using unix:"
  # Streaming command works
  run_podman --url="$URL" run --rm -i $IMAGE /bin/sh -c 'echo -n foo; sleep 0.1; echo -n bar; sleep 0.1; echo -n baz'
  is "$output" foobarbaz

}

@test "tls" {
  skip_if_remote "testing tls service only works locally"
  port=$(random_free_port)
  URL=tcp://127.0.0.1:$port

  CA_CRT=$PODMAN_TMPDIR/certs/ca.crt.pem
  CA_KEY=$PODMAN_TMPDIR/certs/ca.key.pem
  SERVER_CRT=$PODMAN_TMPDIR/certs/server.crt.pem
  SERVER_KEY=$PODMAN_TMPDIR/certs/server.key.pem

  cert-pair "ca" \
    "${CA_KEY}" "${CA_CRT}" \
		-addext basicConstraints=critical,CA:TRUE,pathlen:1

  signed-cert-pair "localhost" \
    "${SERVER_KEY}" "${SERVER_CRT}" \
    "${CA_KEY}" "${CA_CRT}" \
		-addext "subjectAltName=DNS:localhost,IP:127.0.0.1"

  systemd-run --unit=$SERVICE_NAME ${PODMAN%%-remote*} system service $URL --time=0 \
    --tls-key="${SERVER_KEY}" \
    --tls-cert="${SERVER_CRT}"
  wait_for_port 127.0.0.1 $port

  # Flags work
  run_podman \
    --url="$URL" \
    --tls-ca="${CA_CRT}" \
    info --format '{{.Host.RemoteSocket.Path}}'
  is "$output" "$URL" "RemoteSocket.Path using unix:"
  # Streaming command works
  run_podman --url="$URL" \
    --tls-ca="${CA_CRT}" \
    run --rm -i $IMAGE /bin/sh -c 'echo -n foo; sleep 0.1; echo -n bar; sleep 0.1; echo -n baz'
  is "$output" foobarbaz

}

@test "mtls" {
  skip_if_remote "testing mtls service only works locally"

  port=$(random_free_port)
  URL=tcp://127.0.0.1:$port

  CA_CRT=$PODMAN_TMPDIR/certs/ca.crt.pem
  CA_KEY=$PODMAN_TMPDIR/certs/ca.key.pem
  SERVER_CRT=$PODMAN_TMPDIR/certs/server.crt.pem
  SERVER_KEY=$PODMAN_TMPDIR/certs/server.key.pem
  CLIENT_CRT=$PODMAN_TMPDIR/certs/client.crt.pem
  CLIENT_KEY=$PODMAN_TMPDIR/certs/client.key.pem

  cert-pair "ca" \
    "${CA_KEY}" "${CA_CRT}" \
		-addext basicConstraints=critical,CA:TRUE,pathlen:1

  signed-cert-pair "client" \
    "${CLIENT_KEY}" "${CLIENT_CRT}" \
    "${CA_KEY}" "${CA_CRT}" \

  signed-cert-pair "localhost" \
    "${SERVER_KEY}" "${SERVER_CRT}" \
    "${CA_KEY}" "${CA_CRT}" \
		-addext "subjectAltName=DNS:localhost,IP:127.0.0.1"

  systemd-run --unit=$SERVICE_NAME ${PODMAN%%-remote*} system service $URL --time=0 \
    --tls-client-ca="${CA_CRT}" \
    --tls-key="${SERVER_KEY}" \
    --tls-cert="${SERVER_CRT}"
  wait_for_port 127.0.0.1 $port

  # Flags work
  run_podman \
    --url="$URL" \
    --tls-key="${CLIENT_KEY}" \
    --tls-cert="${CLIENT_CRT}" \
    --tls-ca="${CA_CRT}" \
    info --format '{{.Host.RemoteSocket.Path}}'
  is "$output" "$URL" "RemoteSocket.Path using unix:"
  # Streaming command works
    run_podman \
    --url="$URL" \
    --tls-key="${CLIENT_KEY}" \
    --tls-cert="${CLIENT_CRT}" \
    --tls-ca="${CA_CRT}" \
    run --rm -i $IMAGE /bin/sh -c 'echo -n foo; sleep 0.1; echo -n bar; sleep 0.1; echo -n baz'
  is "$output" foobarbaz
}

@test "bogus cert should fail" {
  skip_if_remote "testing tls service only works locally"
  port=$(random_free_port)
  URL=tcp://127.0.0.1:$port

  CA_CRT=$PODMAN_TMPDIR/certs/ca.crt.pem
  CA_KEY=$PODMAN_TMPDIR/certs/ca.key.pem
  SERVER_CRT=$PODMAN_TMPDIR/certs/server.crt.pem
  SERVER_KEY=$PODMAN_TMPDIR/certs/server.key.pem
  BOGUS_CRT=$PODMAN_TMPDIR/certs/bogus.crt.pem
  BOGUS_KEY=$PODMAN_TMPDIR/certs/bogus.key.pem


  cert-pair "ca" \
    "${CA_KEY}" "${CA_CRT}" \
		-addext basicConstraints=critical,CA:TRUE,pathlen:1

  signed-cert-pair "localhost" \
    "${SERVER_KEY}" "${SERVER_CRT}" \
    "${CA_KEY}" "${CA_CRT}" \
		-addext "subjectAltName=DNS:localhost,IP:127.0.0.1"

  cert-pair "bogus" \
  "${BOGUS_KEY}" "${BOGUS_CRT}" \
  -addext basicConstraints=critical,CA:TRUE,pathlen:1

  systemd-run --unit=$SERVICE_NAME ${PODMAN%%-remote*} system service $URL --time=0 \
    --tls-client-ca="${CA_CRT}" \
    --tls-key="${SERVER_KEY}" \
    --tls-cert="${SERVER_CRT}"
  wait_for_port 127.0.0.1 $port

  run_podman 125 \
    --url="$URL" \
    --tls-key="${BOGUS_KEY}" \
    --tls-cert="${BOGUS_CRT}" \
    --tls-ca="${CA_CRT}" \
    info --format '{{.Host.RemoteSocket.Path}}'
  is "$output" ".* remote error: tls: certificate required"
}
