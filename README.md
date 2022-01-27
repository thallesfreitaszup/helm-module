## Description
A helm lib that downloads the contents from the source, do the helm template
and return an array of manifests
## How to use
Instantiate  client of your preference(current we support the go getter client)
and pass as an argument of the **New** method with the source url of manifests,  the destination path  and the options.
See the example below:

    pwd, _ := os.Getwd()
	getter := getter.Client{
		Ctx:  context.TODO(),
		Pwd:  pwd,
		Src:  "git::git@gitlab.com:thalleslmf/event-receiver.git//event-receiver?ref=main",
		Dst:  filepath.Join(os.TempDir(), "helm"+strconv.Itoa(int(z.FastRand()))),
		Mode: getter.ClientModeAny,
	}
	h := helm.New(getter.Src, &getter, helm.Options{}, getter.Dst)


After that call the **Render** method that will return the desired manifests:

    ```
        manifests, err := h.Render()

    ``
## Options
In the helm options you can pass a SSHKEY if you want to download from private repositories, this key must not have a passphrase.For more explanation see ```https://www.ssh.com/academy/ssh/passphrase```
It's also possible to configure a cache passing a struct that fits with the cache interface.For a cache implementation use https://github.com/ZupIT/charlescd-k8s-config-cache/.
See example
``` 
helm.Options{ Cache: someCache, SSHKey: 'some-key' }
	
```
