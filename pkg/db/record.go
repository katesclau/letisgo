package db

type RecordInterface interface {
	Record() any
}

type Record struct {
	Pk      string `json:"pk" dynamodbav:"pk"`
	Sk      string `json:"sk" dynamodbav:"sk"`
	Type    string `json:"typ" dynamodbav:"type"`
	Version int    `json:"v" dynamodbav:"version"`
}
