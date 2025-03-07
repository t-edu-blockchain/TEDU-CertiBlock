#Install Fabric
Get the install script:
```sh
curl -sSLO https://raw.githubusercontent.com/hyperledger/fabric/main/scripts/install-fabric.sh && chmod +x install-fabric.sh
```

#Install bin and docker:
```sh
./install-fabric.sh docker binary
```

#Install Wails: https://wails.io/docs/gettingstarted/installation

Run: 
Khởi động wsl trên máy:
```
./network.sh up createChannel -ca -c mychannel -s couchdb
./network.sh deployCC -ccn certicontract -ccp ../application/backend/chaincode -ccv 1 -ccl go

cd application/nodeapp
wails dev -tags webkit2_41


cd application/backend/api 
go run . ```
