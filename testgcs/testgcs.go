package main

import (
	"net/http"
	"os"

	"code.google.com/p/goauth2/appengine/serviceaccount"
	"golang.org/x/net/context"
	"google.golang.org/cloud"
	"google.golang.org/cloud/storage"

	"appengine"
)

func init() {
	http.HandleFunc("/", handleFunc)
}

func handleFunc(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	c.Infof("env: %v", os.Environ())
	client, err := serviceaccount.NewClient(c, "https://www.googleapis.com/auth/devstorage.read_write")
	c.Infof("client=%+v err=%+v", client, err)
	cc := cloud.WithContext(context.Background(), appengine.AppID(c), client)
	bucket, err := storage.BucketInfo(cc, appengine.AppID(c)+".appspot.com")
	c.Infof("bucket=%+v err=%+v", bucket, err)
}
