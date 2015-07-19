export GOPATH=$PWD
go get code.google.com/p/goauth2
go get google.golang.org/cloud
go get google.golang.org/api/storage
go get go get google.golang.org/appengine # one of these other packages requires urlfetch from here
goapp deploy testgcs/testgcs.yaml
