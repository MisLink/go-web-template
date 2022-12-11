//go:build tools

package tools

import (
	_ "entgo.io/ent/cmd/ent"
	_ "github.com/google/wire/cmd/wire"
	_ "github.com/twitchtv/twirp/protoc-gen-twirp"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)
