package ent

//go:generate go run entgo.io/ent/cmd/ent generate ./schema --feature sql/lock,sql/modifier,sql/execquery,sql/upsert,schema/snapshot,privacy,entql --idtype int64
