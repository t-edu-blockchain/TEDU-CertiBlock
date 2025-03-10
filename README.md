Install wsl2 - (Windows)

Install Wails: https://wails.io/docs/gettingstarted/installation

Run: 
Khởi động wsl trên máy:

Install Go for Linux:
```
wget https://go.dev/dl/go1.24.1.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.24.1.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```
Install Fabric
Get the install script:
```sh
curl -sSLO https://raw.githubusercontent.com/hyperledger/fabric/main/scripts/install-fabric.sh && chmod +x install-fabric.sh
```
Install bin and docker:
```sh
./install-fabric.sh docker binary
```

```
cd CertiBlock/network
./network.sh up createChannel -ca -c mychannel -s couchdb
./network.sh deployCC -ccn certicontract -ccp ../application/backend/chaincode -ccv 1 -ccl go

cd application/nodeapp
wails dev -tags webkit2_41


cd application/backend/api 
go run . ```
