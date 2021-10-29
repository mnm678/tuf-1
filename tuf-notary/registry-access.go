package tufnotary

import (
	"context"
	"io/ioutil"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"

    "oras.land/oras-go/pkg/content"
    "oras.land/oras-go/pkg/oras"
)

func UploadTUFMetadata(registry string, repository string, name string, reference string) ocispec.Descriptor{
	ref := registry + "/" + repository + ":" + name
	fileName := repository + "/staged/" + name + ".json"
	contents, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	fileContent := []byte(contents)
	mediaType := "application/vnd.cncf.notary.tuf+json"

	ctx := context.Background()

	// TODO: add reference once it's supported in oras-go: https://github.com/oras-project/oras-go/pull/35

	memoryStore := content.NewMemory()
    desc, err := memoryStore.Add(fileName, mediaType, fileContent)
	if err != nil {
		return err
	}

	manifest, manifestDesc, config, configDesc, err := content.GenerateManifestAndConfig(nil, nil, desc)
	if err != nil {
		return err
	}

	memoryStore.Set(configDesc, config)
	err = memoryStore.StoreManifest(ref, manifestDesc, manifest)
	if err != nil {
		return err
	}

	reg, err := content.NewRegistry(content.RegistryOptions{PlainHTTP: true})
	fmt.Println(reg)

	//pushContents := []ocispec.Descriptor{desc}
	//desc, err = oras.Push(ctx, resolver, ref, memoryStore, pushContents)
	desc, err = oras.Copy(ctx, memoryStore, ref, reg, "")

	if err != nil {
		return err
	}

	return desc
}
