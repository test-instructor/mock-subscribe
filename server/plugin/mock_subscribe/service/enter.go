package service

var Service = new(service)

type service struct {
	Merchant       merchant
	Contract       contract
	Deduct         deduct
	Callback       callback
	CallbackRecord callbackRecord
	DeductCallback deductCallback
	Signature      signature
	XMLCodec       xmlCodec
}
