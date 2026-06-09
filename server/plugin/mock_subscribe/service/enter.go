package service

var Service = new(service)

type service struct {
	Merchant  merchant
	Contract  contract
	Deduct    deduct
	Callback  callback
	Signature signature
	XMLCodec  xmlCodec
}
