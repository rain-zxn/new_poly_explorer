module new_poly_explorer

go 1.16

require github.com/beego/beego/v2 v2.0.1

require (
	github.com/cosmos/cosmos-sdk v0.0.0-00010101000000-000000000000
	github.com/ethereum/go-ethereum v1.9.13
	github.com/go-sql-driver/mysql v1.6.0
	github.com/joeqian10/neo-gogogo v1.3.0
	github.com/ontio/ontology v1.13.2
	github.com/shopspring/decimal v1.2.0
	github.com/smartystreets/goconvey v1.6.4
	gorm.io/driver/mysql v1.1.1
	gorm.io/gorm v1.21.11
)

replace (
	github.com/cosmos/cosmos-sdk => github.com/Switcheo/cosmos-sdk v0.39.2-0.20200814061308-474a0dbbe4ba
	github.com/ethereum/go-ethereum => github.com/ethereum/go-ethereum v1.9.13
	github.com/joeqian10/neo-gogogo => github.com/blockchain-develop/neo-gogogo v0.0.0-20200824102609-fddf06a45f66
)
