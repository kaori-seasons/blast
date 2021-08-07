package http_api

type ErrorCode int32

const (
	OK ErrorCode = 0

	// input is invalid. e.g. format error, auth failed
	// this is offen caused by client
	// code range is [100, 1000)
	INVALID_PARAM  ErrorCode = 100
	AUTHORIZE_FAIL ErrorCode = 101

	// query fail. e.g. sql error, execution timeout
	// this is offen caused by client, sometimes by server
	// code range is [1000, 2000)
	INVALID_SQL        ErrorCode = 1000
	SLOW_QUERY_TIMEOUT ErrorCode = 1001
	METHOD_NOT_SUPPORT ErrorCode = 1002

	// upserts fail. e.g. required field missed
	// this is offen caused by client
	// code range is [2000, 3000)
	REQUIRED_FIELD_MISSED ErrorCode = 2000
	FIELD_TYPE_RRROR      ErrorCode = 2001
	TABLE_NOT_EXIST       ErrorCode = 2002
	FIELD_NOT_EXIST       ErrorCode = 2003
	//TOO_MANY_DATA         ErrorCode = 2004
	UPSERTS_FAIL_IN_ES    ErrorCode = 2005
	UPSERTS_FAIL_IN_DORIS ErrorCode = 2006

	// schema validate failed
	// code range is [3000, 4000)
	SCHEMA_VALIDATE_FAILED ErrorCode = 3000

	// rate limiter control
	// code range is [4000, 4100)
	EXCEED_RATE_LIMIT ErrorCode = 4000

	// meta-store
	ENV_DB_TBL_NOT_FOUND ErrorCode = 5000
	ENV_DB_TBL_EXIST     ErrorCode = 5001
	STORAGE_NOT_MATCH    ErrorCode = 5002
	COLUMN_NOT_MATCH     ErrorCode = 5003
	COLUMN_EXIST         ErrorCode = 5004
	COLUMN_NOT_EXIST     ErrorCode = 5005
	TABLE_TX_RUNNING     ErrorCode = 5006
	STORAGE_NOT_SUPPORT  ErrorCode = 5007
	INDEX_EXIST          ErrorCode = 5008
	INDEX_NOT_EXIST      ErrorCode = 5009
	ENV_DB_EXIST         ErrorCode = 5010
	ENV_DB_NOT_FOUND     ErrorCode = 5011
	TAG_EXIST            ErrorCode = 5012
	TAG_NOT_EXIST        ErrorCode = 5013

	// scheduler [5200, 5300)
	SCHEDULER_JOB_NOT_FOUND ErrorCode = 5200
	SCHEDULER_JOB_EXIST     ErrorCode = 5201

	// kafka-cluster
	// code range is [6000, 6500)
	// CLUSTER_NOT_FOUND     ErrorCode = 6000
	// REP_MORE_THAN_BROKERS ErrorCode = 6001
	// TOPIC_EXIST           ErrorCode = 6002

	// privilege
	USER_GROUP_EXIST     ErrorCode = 6500
	USER_GROUP_NOT_FOUND ErrorCode = 6501
	RESOURCE_EXIST       ErrorCode = 6502
	RESOURCE_NOT_FOUND   ErrorCode = 6503
	DEPARTMENT_EXIST     ErrorCode = 6504
	DEPARTMENT_NOT_FOUND ErrorCode = 6505
	ROLE_NOT_FOUND       ErrorCode = 6506
	PRIVILEGE_EXIST      ErrorCode = 6507
	PRIVILEGE_NOT_FOUND  ErrorCode = 6508

	// internal error. e.g. db connection broken
	// this is offen caused by server
	// code range is [9000, 10000)
	INTERNAL_ERROR      ErrorCode = 9000
	API_NOT_ONLINE      ErrorCode = 9001
	FEATURE_NOT_SUPPORT ErrorCode = 9002

	// inner error. for debug and testing...
	// code range is [200000000, -)
)

type ErrorWrapper struct {
	Code  ErrorCode
	Error error
}

func NewErrorWrapper(c ErrorCode, ctx error) *ErrorWrapper {
	return &ErrorWrapper{c, ctx}
}

// String convert ErrorCode to human readable string
func (ec ErrorCode) String(err error) string {
	switch ec {
	case OK:
		return "success"
	case INVALID_PARAM:
		return "input param is invalid: (" + err.Error() + ")"
	case AUTHORIZE_FAIL:
		return "authorize failed: (" + err.Error() + ")"
	case INVALID_SQL:
		return "sql is invalid: (" + err.Error() + ")"
	case SLOW_QUERY_TIMEOUT:
		return "query timeout"
	case METHOD_NOT_SUPPORT:
		return "unsupported method: (only SELECT is supported now)"
	case REQUIRED_FIELD_MISSED:
		return "required field missed: (" + err.Error() + ")"
	case FIELD_TYPE_RRROR:
		return "field type error: (" + err.Error() + ")"
	case TABLE_NOT_EXIST:
		return "table not exist: (" + err.Error() + ")"
	case FIELD_NOT_EXIST:
		return "field not exist: (" + err.Error() + ")"
	//case TOO_MANY_DATA:
	//	return "too many data: (" + err.Error() + ")"
	case UPSERTS_FAIL_IN_ES:
		return "upserts failed against ES"
	case UPSERTS_FAIL_IN_DORIS:
		return "upserts failed against Doris"
	case SCHEMA_VALIDATE_FAILED:
		return "schema validate failed: (" + err.Error() + ")"
	case EXCEED_RATE_LIMIT:
		return "exceed rate limit"
	case ENV_DB_TBL_NOT_FOUND:
		return "env.db.table not found. (" + err.Error() + ")"
	case ENV_DB_TBL_EXIST:
		return "env.db.table exist. (" + err.Error() + ")"
	case ENV_DB_EXIST:
		return "env.db exist. (" + err.Error() + ")"
	case ENV_DB_NOT_FOUND:
		return "env.db not found. (" + err.Error() + ")"
	case STORAGE_NOT_MATCH:
		return "storage.name not match. (" + err.Error() + ")"
	case COLUMN_NOT_MATCH:
		return "column not match. (" + err.Error() + ")"
	case COLUMN_EXIST:
		return "column field exist. (" + err.Error() + ")"
	case COLUMN_NOT_EXIST:
		return "column field is not exist. (" + err.Error() + ")"
	case TABLE_TX_RUNNING:
		return "table is runing tx."
	case STORAGE_NOT_SUPPORT:
		return "storage is not support this feature. (" + err.Error() + ")"
	case INDEX_EXIST:
		return "index exist. (" + err.Error() + ")"
	case INDEX_NOT_EXIST:
		return "index is not exist. (" + err.Error() + ")"
	case SCHEDULER_JOB_NOT_FOUND:
		return "job is not exist. (" + err.Error() + ")"
	case SCHEDULER_JOB_EXIST:
		return "job is exist. (" + err.Error() + ")"
	case USER_GROUP_EXIST:
		return "user group exit. (" + err.Error() + ")"
	case USER_GROUP_NOT_FOUND:
		return "user group is not exist. (" + err.Error() + ")"
	case RESOURCE_EXIST:
		return "resource exist. (" + err.Error() + ")"
	case RESOURCE_NOT_FOUND:
		return "resource is not exist. (" + err.Error() + ")"
	case DEPARTMENT_EXIST:
		return "department exist. (" + err.Error() + ")"
	case DEPARTMENT_NOT_FOUND:
		return "department is not exist. (" + err.Error() + ")"
	case ROLE_NOT_FOUND:
		return "role is not exist. (" + err.Error() + ")"
	case PRIVILEGE_EXIST:
		return "privilege is exist. (" + err.Error() + ")"
	case PRIVILEGE_NOT_FOUND:
		return "privilege is not exist. (" + err.Error() + ")"
	case INTERNAL_ERROR:
		return "internal error: (" + err.Error() + ")"
	case API_NOT_ONLINE:
		return "api is offline. error: (" + err.Error() + ")"
	case FEATURE_NOT_SUPPORT:
		return "feature is not supported. error: (" + err.Error() + ")"
	case TAG_EXIST:
		return "tag is exist, error: (" + err.Error() + ")"
	case TAG_NOT_EXIST:
		return "tag is not exist, error: (" + err.Error() + ")"
	default:
		return "unknow errcode..."
	}
}
