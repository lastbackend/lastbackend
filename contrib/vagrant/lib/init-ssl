#!/usr/bin/env bash
set -e

# define location of openssl binary manually since running this
# script under Vagrant fails on some systems without it
OPENSSL=$(which openssl)

function usage {
    echo "USAGE: $0 <output-dir> <cert-base-name> <CN> [SAN,SAN,SAN]"
    echo "  example: $0 ./ssl/ worker kube-worker IP.1=127.0.0.1,IP.2=10.0.0.1"
}

if [ -z "$1" ] || [ -z "$2" ] || [ -z "$3" ]; then
    usage
    exit 1
fi

OUTDIR="$1"
CERTBASE="$2"
CN="$3"
SANS="$4"

if [ ! -d $OUTDIR ]; then
    echo "ERROR: output directory does not exist:  $OUTDIR"
    exit 1
fi

OUTFILE="$OUTDIR/$CN.tar"

if [ -f "$OUTFILE" ];then
    exit 0
fi

CNF_TEMPLATE="
[req]
req_extensions = v3_req
distinguished_name = req_distinguished_name
[req_distinguished_name]
[ v3_req ]
basicConstraints = CA:FALSE
keyUsage = nonRepudiation, digitalSignature, keyEncipherment
subjectAltName = @alt_names
[alt_names]
DNS.1 = kubernetes
DNS.2 = kubernetes.default
DNS.3 = kubernetes.default.svc
DNS.4 = kubernetes.default.svc.cluster.local
"
echo "Generating SSL artifacts in $OUTDIR"


CONFIGFILE="$OUTDIR/$CERTBASE-req.cnf"
CAFILE="$OUTDIR/ca.pem"
CAKEYFILE="$OUTDIR/ca-key.pem"
KEYFILE="$OUTDIR/$CERTBASE-key.pem"
CSRFILE="$OUTDIR/$CERTBASE.csr"
PEMFILE="$OUTDIR/$CERTBASE.pem"

CONTENTS="${CAFILE} ${KEYFILE} ${PEMFILE}"


# Add SANs to openssl config
echo "$CNF_TEMPLATE$(echo $SANS | tr ',' '\n')" > "$CONFIGFILE"

$OPENSSL genrsa -out "$KEYFILE" 2048
$OPENSSL req -new -key "$KEYFILE" -out "$CSRFILE" -subj "/CN=$CN" -config "$CONFIGFILE"
$OPENSSL x509 -req -in "$CSRFILE" -CA "$CAFILE" -CAkey "$CAKEYFILE" -CAcreateserial -out "$PEMFILE" -days 365 -extensions v3_req -extfile "$CONFIGFILE"

tar -cf $OUTFILE -C $OUTDIR $(for  f in $CONTENTS;do printf "$(basename $f) ";done)

echo "Bundled SSL artifacts into $OUTFILE"
echo "$CONTENTS"
