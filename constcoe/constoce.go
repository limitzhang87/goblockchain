package constcoe

const (
	Difficult = 12
	InitCoin  = 1000 // This line is new

	LHKey         = "lh"
	OgPrevHashKey = "ogPrevHash"

	TransactionPoolFile = "./tmp/transaction_pool.data"
	BCPatch             = "./tmp/blocks"
	BCFile              = "./tmp/blocks/MANIFEST"

	ChecksumLength = 4
	NetworkVersion = byte(0x00)
	Wallets        = "./tmp/wallets/"
	WalletsRefList = "./tmp/ref_list"
)
