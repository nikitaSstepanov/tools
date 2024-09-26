package migrate

type Dialect string

const (
	Postgres   Dialect = "postgres"
	Mysql      Dialect = "mysql"
	Sqlite3    Dialect = "sqlite3"
	Mssql      Dialect = "mssql"
	Redshift   Dialect = "redshift"
	TiDB       Dialect = "tidb"
	ClickHouse Dialect = "clickhouse"
)

func diaToStr(d Dialect) string {
	switch d {

	case Postgres:
		return "postgres"

	case Mysql:
		return "mysql"

	case Sqlite3:
		return "sqlite3"

	case Mssql:
		return "mssql"

	case Redshift:
		return "redshift"

	case TiDB:
		return "tidb"

	case ClickHouse:
		return "clickhouse"

	default:
		return "postgres"
		
	}
}
