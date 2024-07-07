package consts

const (
	DATA_FMT_MICRO_SECONDS = "2006-01-02 15:04:05.000000"
)

// relate to message
const (
	SPACE                 = " "
	UINT32_MAX            = 4294967295
	INT64_MAX       int64 = 9_223_372_036_854_775_807
	RESEVERD_PREFIX       = "__toPIIK"
)

// relate to persistent
const (
	// the Queue that waiting to be persistent
	PERSISTENT_BUF_QUEUE = "__toPIIK_persistent_buf_queue"
)

// relate to Topic
const (
	CONSUMER_OFFSET_PREFIX = "__toPIIK_consumer_"
)

const (
	EMAIL_PATTERN = `^.+:\d{4,5}`
)