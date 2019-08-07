package types

// Bank module event types
var (
	EventTypeTransfer     = "transfer"
	EventTypeIssueAsset   = "issueAsset"
	EventTypeIssueFiat    = "issueFiat"
	EventTypeRedeemAsset  = "redeemAsset"
	EventTypeRedeemFiat   = "redeemFiat"
	EventTypeSendAsset    = "sendAsset"
	EventTypeSendFiat     = "sendFiat"
	EventTypeExecuteOrder = "executeOrder"
	EventTypeReleaseAsset = "releaseAsset"
	
	AttributeKeyRecipient = "recipient"
	AttributeKeySender    = "sender"
	
	AttributeValueCategory = ModuleName
)
