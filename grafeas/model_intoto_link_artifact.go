/*
 * grafeas.proto
 *
 * No description provided (generated by Swagger Codegen https://github.com/swagger-api/swagger-codegen)
 *
 * API version: version not set
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package grafeas

type IntotoLinkArtifact struct {
	ResourceUri string `json:"resource_uri,omitempty"`
	Hashes *LinkArtifactHashes `json:"hashes,omitempty"`
}