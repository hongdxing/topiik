package consts

import "os"

const (
	DATA_FMT_MICRO_SECONDS = "2006-01-02 15:04:05.000000"
)

// relate to message
const (
	SPACE                 = " "
	UINT32_MAX            = 4294967295
	INT64_MAX       int64 = 9_223_372_036_854_775_807
	INT64_MIN       int64 = -9_223_372_036_854_775_808
	RESEVERD_PREFIX       = "__toPIIK"
)

// relate to persistence
const (
	// the Queue that waiting to be persistence
	PERSISTENT_BUF_QUEUE = "__toPIIK_persistence_buf_queue"
)

// relate to Topic
const (
	CONSUMER_OFFSET_PREFIX = "__toPIIK_consumer_"
)

// relate to Response
const (
	RES_INVALID_CMD = "INVALID_CMD"
)

const (
	HOST_PATTERN = `^.+:\d{4,5}`
)

const (
	SLASH    = string(os.PathSeparator)
	DATA_DIR = "data"
)

const (
	VOTE_ACCEPTED = "A"
	VOTE_REJECTED = "R"
)

const (
	CLUSTER_ID_LEN = 10
	NODE_ID_LEN    = 10
	PTN_ID_LEN     = 10
)

const SLOTS = 256
