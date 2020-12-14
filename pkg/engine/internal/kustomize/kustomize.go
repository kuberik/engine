package kustomize

import (
	"encoding/json"
	"fmt"

	"sigs.k8s.io/kustomize/api/filesys"
	"sigs.k8s.io/kustomize/api/krusty"
	"sigs.k8s.io/kustomize/api/resmap"
	"sigs.k8s.io/kustomize/api/types"
)

type KustomizeLayer struct {
	fs            filesys.FileSystem
	id            int
	Kustomization types.Kustomization
	parent        *KustomizeLayer
}

const (
	layerResourceObjectsFile       = "resources.yaml"
	layerPatchesStrategicMergeFile = "patchesStrategicMerge.yaml"
)

func (kl *KustomizeLayer) fileAbsPath(name string) string {
	return filesys.RootedPath(fmt.Sprint(kl.id), name)
}

func (kl *KustomizeLayer) writeFile(name string, contents []byte) {
	kl.fs.WriteFile(kl.fileAbsPath(name), contents)
}

func (kl *KustomizeLayer) appendFile(name string, contents []byte) {
	oldContents, _ := kl.fs.ReadFile(kl.fileAbsPath(name))
	kl.writeFile(name, append(oldContents, contents...))
}

func NewKustomizeLayerRoot() KustomizeLayer {
	kl := KustomizeLayer{}
	kl.fs = filesys.MakeFsInMemory()

	return kl
}

func (kl *KustomizeLayer) writeKustomizationFiles() {
	if kl.parent != nil {
		kl.parent.writeKustomizationFiles()
		kl.Kustomization.Resources = append(kl.Kustomization.Resources, fmt.Sprintf("../%d", kl.parent.id))
	}

	if kl.fs.Exists(kl.fileAbsPath(layerResourceObjectsFile)) {
		kl.Kustomization.Resources = append(kl.Kustomization.Resources, layerResourceObjectsFile)
	}

	if kl.fs.Exists(kl.fileAbsPath(layerPatchesStrategicMergeFile)) {
		kl.Kustomization.PatchesStrategicMerge = append(kl.Kustomization.PatchesStrategicMerge, layerPatchesStrategicMergeFile)
	}

	kl.Kustomization.FixKustomizationPreMarshalling()
	kustomizeContents, _ := json.Marshal(kl.Kustomization)
	kl.writeFile("kustomization.yaml", kustomizeContents)
}

func (kl *KustomizeLayer) Run() (resmap.ResMap, error) {
	kl.writeKustomizationFiles()
	kustomizer := krusty.MakeKustomizer(kl.fs, krusty.MakeDefaultOptions())
	return kustomizer.Run(filesys.RootedPath(fmt.Sprint(kl.id)))
}

func (kl *KustomizeLayer) AddObjectRaw(raw []byte) {
	kl.appendFile(layerResourceObjectsFile, []byte(fmt.Sprintf("---\n%s\n", raw)))
}
func (kl *KustomizeLayer) AddObject(object interface{}) {
	marshaled, _ := json.Marshal(object)
	kl.AddObjectRaw(marshaled)
}

func (kl *KustomizeLayer) AddLayer() KustomizeLayer {
	nl := KustomizeLayer{}
	nl.fs = kl.fs
	nl.id = kl.id + 1
	nl.parent = kl

	return nl
}
