package radiusop

import (
	"context"
	"time"

	"github.com/asmexie/gopub/common"
	"github.com/asmexie/gopub/dbutils"
)

// RadiusCheck ...
type RadiusCheck struct {
	UserName   string `db:"username"`
	Attribute  string `db:"attribute"`
	Op         string `db:"op"`
	Value      string `db:"value"`
	ExpireTime time.Time
}

var _db *dbutils.Database

// OpenDB ...
func OpenDB(ctx context.Context, connStr string) error {
	if _db == nil {
		_db = dbutils.NewDatabase(ctx, nil)
	}
	return _db.Open(connStr)
}

// CloseDB ...
func CloseDB() {
	if _db != nil {
		_db.Close()
	}
}

// NewRadiusCheck ...
func NewRadiusCheck(user, pwd string) RadiusCheck {
	return RadiusCheck{
		UserName:  user,
		Value:     pwd,
		Op:        ":=",
		Attribute: "Cleartext-Password",
	}
}

// AddRadiusUser ...
func AddRadiusUser(rc RadiusCheck) error {
	ds := _db.NewSession()
	_, err := ds.InsertInto(`radcheck`).Columns(ds.DBInsColumns(rc)...).Record(rc).Exec()
	common.LogError(err)
	return err
}

// DelRadiusUser ...
func DelRadiusUser(userName string) error {
	ds := _db.NewSession()
	_, err := ds.DeleteFrom(`radcheck`).Where("username=?", userName).Exec()
	common.LogError(err)
	return err
}
