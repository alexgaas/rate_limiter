package bucket_quoter

type BucketQuoterStat struct {
	MsgPassed        int64
	BucketUnderflows int64
	TokensUsed       int64
	UsecWaited       int64
	AggregateInflow  int64
}
