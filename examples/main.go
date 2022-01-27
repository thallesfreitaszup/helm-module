package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/ristretto/z"
	"github.com/hashicorp/go-getter"
	"github.com/thallesfreitaszup/helm-module/helm"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	pwd, _ := os.Getwd()
	getter := getter.Client{
		Ctx:  context.TODO(),
		Pwd:  pwd,
		Src:  "git::git@gitlab.com:thalleslmf/event-receiver.git",
		Dst:  filepath.Join(os.TempDir(), "helm"+strconv.Itoa(int(z.FastRand()))),
		Mode: getter.ClientModeAny,
	}

	gitOptions := helm.Git{Branch: "main", Path: "/event-receiver"}
	h := helm.New(getter.Src, &getter, helm.Options{Git: gitOptions}, getter.Dst)
	manifests, err := h.Render()
	if err != nil {
		log.Fatal(err)
	}
	manifestBytes, _ := json.Marshal(manifests)
	fmt.Println(string(manifestBytes))
}
