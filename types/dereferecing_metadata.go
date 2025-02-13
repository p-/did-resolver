package types

import (
	didTypes "github.com/cheqd/cheqd-node/x/did/types"
)

type DereferencingMetadata ResolutionMetadata

type DidDereferencing struct {
	Context               string                   `json:"@context,omitempty" example:"https://w3id.org/did-resolution/v1"`
	DereferencingMetadata DereferencingMetadata    `json:"dereferencingMetadata"`
	ContentStream         ContentStreamI           `json:"contentStream"`
	Metadata              ResolutionDidDocMetadata `json:"contentMetadata"`
}

func NewDereferencingMetadata(did string, contentType ContentType, resolutionError string) DereferencingMetadata {
	return DereferencingMetadata(NewResolutionMetadata(did, contentType, resolutionError))
}

func (d DidDereferencing) GetContentType() string {
	return string(d.DereferencingMetadata.ContentType)
}

func (d DidDereferencing) GetBytes() []byte {
	if d.ContentStream == nil {
		return []byte{}
	}
	return d.ContentStream.GetBytes()
}

type DereferencedDidVersionsList struct {
	Versions []*didTypes.Metadata `json:"versions,omitempty"`
}

func NewDereferencedDidVersionsList(versions []*didTypes.Metadata) *DereferencedDidVersionsList {
	return &DereferencedDidVersionsList{
		Versions: versions,
	}
}

func (e *DereferencedDidVersionsList) AddContext(newProtocol string) {}
func (e *DereferencedDidVersionsList) RemoveContext()                {}
func (e *DereferencedDidVersionsList) GetBytes() []byte              { return []byte{} }
