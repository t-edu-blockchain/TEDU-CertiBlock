# CertiBlock

A blockchain platform for decentralized, trustworthy issuance and verification of certificates worldwide.

- [CertiBlock](#certiblock)
  - [Run on Local Machine](#run-on-local-machine)
    - [Requirements](#requirements)
    - [Setup Network](#setup-network)
    - [Run Applications](#run-applications)
      - [The Native App for Issuing Certificates](#the-native-app-for-issuing-certificates)
      - [The Web App for Viewing Personal Certificates](#the-web-app-for-viewing-personal-certificates)
      - [The Web App for Verifying Certificates by QR Codes](#the-web-app-for-verifying-certificates-by-qr-codes)
  - [Technologies](#technologies)
  - [Authors and Licensing](#authors-and-licensing)

## Run on Local Machine

### Requirements

- **A Linux-based system.** On Windows, use WSL2.
- **Install Go.** One way to do that is as follows:

    ```sh
    wget https://go.dev/dl/go1.24.1.linux-amd64.tar.gz
    sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.24.1.linux-amd64.tar.gz
    echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
    source ~/.bashrc
    ```

- **Install Wails:** <https://wails.io/docs/gettingstarted/installation>
- **Install Fabric:**
    - Get the install script:

        ```sh
        curl -sSLO https://raw.githubusercontent.com/hyperledger/fabric/main/scripts/install-fabric.sh && chmod +x install-fabric.sh
        ```

    - Install bin and docker components:

        ```sh
        ./install-fabric.sh docker binary
        ```

### Setup Network

In project root:

```sh
cd network
./network.sh up createChannel -ca -c mychannel -s couchdb
./network.sh deployCC -ccn certicontract -ccp ../application/backend/chaincode -ccv 1 -ccl go
```

In case one of the above commands fails, you might want to try
**deleting all Docker containers, volumes and networks.** (this
is dangerous so proceed with caution.) Still in that `network`
directory, run:

```sh
./delete-all-docker-containers-volumes-and-networks.sh
```

### Run Applications

#### The Native App for Issuing Certificates

To be run on each university's private machines.

```sh
cd application/nodeapp
wails dev
# Or, if you're using Ubuntu 24.04 or some new system
wails dev -tags webkit2_41
```

#### The Web App for Viewing Personal Certificates

To be accessed by students alike.

First, in `application/backend/certiblock`, create a new file
named `.env`. Specify the appropriate environment variables
as instructed in `.env.example` (also in that same directory).
The backend makes use of MySQL RDBMS.

```sh
cd application/backend/certiblock
./run

cd ../../../application/frontend/recipient
npm i
npm run dev
```

#### The Web App for Verifying Certificates by QR Codes

TODO

## Technologies

TODO

## Authors and Licensing

TODO
