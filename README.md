# rApp for RAN slicing control
The rApp for scaleing RAN slice on demand.

## Development
```
go run main.go
```

## Usage
### Unregister rApp
```
curl -k -X DELETE https://<catalogue service IP>:<catalogue service port>/services/ranslice-scale-control
```

### Support
* The [api doc](https://docs.o-ran-sc.org/projects/o-ran-sc-nonrtric-plt-rappcatalogue/en/f-release/rac-api.html) for rapp catalogue service