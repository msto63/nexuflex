// generate.go
/**
 * Nexuflex Shared - Protocol Buffer Code Generation
 *
 * This file contains the generation directives for Protocol Buffers.
 * It is used to generate the Go code from the .proto files.
 *
 * @author msto63
 * @version 1.0.0
 * @date 2025-03-12
 */

package proto

//go:generate protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative nexuflex.proto
