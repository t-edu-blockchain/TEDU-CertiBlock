#!/bin/bash
#
# Copyright IBM Corp All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#

# default to using UET
ORG=${1:-UET}

# Exit on first error, print all commands.
set -e
set -o pipefail

# Where am I?
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." && pwd )"

ORDERER_CA=${DIR}/network/organizations/ordererOrganizations/example/tlsca/tlsca.example-cert.pem
PEER0_UET_CA=${DIR}/network/organizations/peerOrganizations/uet/tlsca/tlsca.uet-cert.pem
PEER0_HUS_CA=${DIR}/network/organizations/peerOrganizations/hus/tlsca/tlsca.hus-cert.pem
PEER0_ORG3_CA=${DIR}/network/organizations/peerOrganizations/org3/tlsca/tlsca.org3-cert.pem
PEER0_ORG4_CA=${DIR}/network/organizations/peerOrganizations/org4/tlsca/tlsca.org4-cert.pem

if [[ ${ORG,,} == "uet" || ${ORG,,} == "digibank" ]]; then
   CORE_PEER_LOCALMSPID=UETMSP
   CORE_PEER_MSPCONFIGPATH=${DIR}/network/organizations/peerOrganizations/uet/users/Admin@uet/msp
   CORE_PEER_ADDRESS=localhost:7051
   CORE_PEER_TLS_ROOTCERT_FILE=${DIR}/network/organizations/peerOrganizations/uet/tlsca/tlsca.uet-cert.pem
elif [[ ${ORG,,} == "hus" || ${ORG,,} == "magnetocorp" ]]; then
   CORE_PEER_LOCALMSPID=HUSMSP
   CORE_PEER_MSPCONFIGPATH=${DIR}/network/organizations/peerOrganizations/hus/users/Admin@hus/msp
   CORE_PEER_ADDRESS=localhost:9051
   CORE_PEER_TLS_ROOTCERT_FILE=${DIR}/network/organizations/peerOrganizations/hus/tlsca/tlsca.hus-cert.pem
else
   echo "Unknown \"$ORG\", please choose UET/Digibank or HUS/Magnetocorp"
   echo "For example to get the environment variables to set up a HUS shell environment run:  ./setOrgEnv.sh HUS"
   echo
   echo "This can be automated to set them as well with:"
   echo
   echo 'export $(./setOrgEnv.sh HUS | xargs)'
   exit 1
fi

# output the variables that need to be set
echo "CORE_PEER_TLS_ENABLED=true"
echo "ORDERER_CA=${ORDERER_CA}"
echo "PEER0_UET_CA=${PEER0_UET_CA}"
echo "PEER0_HUS_CA=${PEER0_HUS_CA}"
echo "PEER0_ORG3_CA=${PEER0_ORG3_CA}"
echo "PEER0_ORG4_CA=${PEER0_ORG4_CA}"

echo "CORE_PEER_MSPCONFIGPATH=${CORE_PEER_MSPCONFIGPATH}"
echo "CORE_PEER_ADDRESS=${CORE_PEER_ADDRESS}"
echo "CORE_PEER_TLS_ROOTCERT_FILE=${CORE_PEER_TLS_ROOTCERT_FILE}"

echo "CORE_PEER_LOCALMSPID=${CORE_PEER_LOCALMSPID}"