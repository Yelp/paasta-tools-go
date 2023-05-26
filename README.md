# paasta-tools-go
Go library for interacting with [PaaSTA](https://github.com/Yelp/paasta)

# Troubleshooting
### How to clone
The project must be clonned via `https` as it is required to perform SSO, this means no push via ssh would work.

### Private libraries
You may need to export the following env variable

    export GOPRIVATE=*github.yelpcorp.com

otherwise doing `go mod tidy` will have problems finding the monk library.
