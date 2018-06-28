package json

import (
	"bytes"
	"fmt"

	jsonpatchapplier "github.com/evanphx/json-patch"
	"github.com/mattbaird/jsonpatch"
)

// NoiseOperation struct
type NoiseOperation struct {
	Patch []jsonpatch.JsonPatchOperation
}

func (nd NoiseOperation) ContainsNoise() bool {
	return len(nd.Patch) > 0
}

// Remove noise from primary and candidate documents
func (nd *NoiseOperation) Remove(primary, candidate []byte) ([]byte, []byte, error) {

	primaryWithoutNoise := primary
	candidateWithoutNoise := candidate

	if nd.ContainsNoise() {
		patch, err := jsonpatchapplier.DecodePatch(nd.materializePatchOperations())

		if err != nil {
			return nil, nil, err
		}

		primaryWithoutNoise, err = patch.Apply(primary)

		if err != nil {
			return nil, nil, err
		}

		candidateWithoutNoise, err = patch.Apply(candidate)

		if err != nil {
			return nil, nil, err
		}
	}
	return primaryWithoutNoise, candidateWithoutNoise, nil
}

func (nd NoiseOperation) materializePatchOperations() []byte {
	var b bytes.Buffer
	b.Write([]byte("["))
	i := 0
	for _, operation := range nd.Patch {
		patchOp, _ := operation.MarshalJSON()
		b.Write(patchOp)
		if i != len(nd.Patch)-1 {
			b.Write([]byte(","))
		}
		i++
	}
	b.Write([]byte("]"))

	return b.Bytes()
}

// Detect Noise between documents
func (nd *NoiseOperation) Detect(primary, secondary []byte) error {

	patch, err := jsonpatch.CreatePatch(primary, secondary)

	if err != nil {
		return err
	}

	newPatch, err := validatePatchToContainOnlyReplaceAndChangeToRemove(patch)

	if err != nil {
		return err
	}

	nd.Patch = newPatch
	return nil

}

func validatePatchToContainOnlyReplaceAndChangeToRemove(patch []jsonpatch.JsonPatchOperation) ([]jsonpatch.JsonPatchOperation, error) {

	var removePatches []jsonpatch.JsonPatchOperation

	for _, operation := range patch {
		if operation.Operation != "replace" {
			return nil, fmt.Errorf("Primary and Secondary payload contains other changes apart from replacingvalues %s", operation.Json())
		}

		removePatches = append(removePatches, jsonpatch.NewPatch("remove", operation.Path, nil))

	}

	return removePatches, nil
}
