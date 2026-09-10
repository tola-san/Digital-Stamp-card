package migrations

import (
	"embed"
	"io/fs"
)

/*
Files keeps database migrations in the API binary so production deploys
do not depend on local files being present at runtime.
*/



//go:embed *.up.sql
var Files embed.FS

var Up = []string{
	"000001_initial.up.sql",
	"000002_seed_default_reward.up.sql",
	"000003_sessions_and_relationships.up.sql",
}

func ReadUp(name string) ([]byte, error) {
	return fs.ReadFile(Files, name)
}
