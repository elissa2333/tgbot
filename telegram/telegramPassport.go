package telegram

// pack

/*Telegram Passport
Telegram Passport是用于需要个人识别的服务的统一授权方法。用户可以一次上传他们的文档，然后立即与需要真实世界ID（金融，ICO等）的服务共享数据。有关详细信息，请参见手册。
https://core.telegram.org/bots/api#telegram-passport*/

// PassportData 用户与机器人共享的Telegram Passport数据的信息
// https://core.telegram.org/bots/api#passportdata
type PassportData struct {
	Data        []EncryptedPassportElement `json:"data"`        // 数组，其中包含与机器人共享的有关文档和其他Telegram Passport元素的信息
	Credentials *EncryptedCredentials      `json:"credentials"` // 解密数据所需的加密凭据
}

// PassportFile 上传到Telegram Passport的文件。当前，所有Telegram Passport文件在解密时均为JPEG格式，并且不超过10MB
// https://core.telegram.org/bots/api#passportfile
type PassportFile struct {
	FileID       string `json:"file_id"`        // 该文件的标识符，可用于下载或重复使用该文件
	FileUniqueID string `json:"file_unique_id"` // 此文件的唯一标识符，随着时间的推移，对于不同的漫游器，应该是相同的。不能用于下载或重复使用文件。
	FileSize     int64  `json:"file_size"`      // 文件大小
	FileDate     int64  `json:"file_date"`      // 文件上传的Unix时间
}

// EncryptedPassportElement 用户与机器人共享的文档或其他Telegram Passport元素的信息
// https://core.telegram.org/bots/api#encryptedpassportelement
type EncryptedPassportElement struct {
	Type        string         `json:"type"`         // 元素类型。“personal_details”, “passport”, “driver_license”, “identity_card”, “internal_passport”, “address”, “utility_bill”, “bank_statement”, “rental_agreement”, “passport_registration”, “temporary_registration”, “phone_number”, “email”
	Data        string         `json:"data"`         // 可选的。用户提供的Base64编码的加密电报Passport元素数据，可用于“ personal_details”，“ passport”，“ driver_license”，“ identity_card”，“ internal_passport”和“ address”类型。可以使用随附的EncryptedCredentials解密和验证。
	PhoneNumber string         `json:"phone_number"` // 可选的。用户的验证电话号码，仅适用于“ phone_number”类型
	Email       string         `json:"email"`        // 可选的。用户的验证电子邮件地址，仅适用于“电子邮件”类型
	Files       []PassportFile `json:"files"`        // 可选的。带有用户提供的文档的加密文件数组，可用于“ utility_bill”，“ bank_statement”，“ rental_agreement”，“ passport_registration”和“ temporary_registration”类型。可以使用随附的EncryptedCredentials解密和验证文件。
	FrontSide   *PassportFile  `json:"front_side"`   // 可选的。用户提供的带有文档正面的加密文件。适用于“护照”，“驾驶执照”，“身份证”和“内部护照”。可以使用随附的EncryptedCredentials解密和验证文件。
	ReverseSide *PassportFile  `json:"reverse_side"` // 可选的。由用户提供的带有文档反面的加密文件。适用于“ driver_license”和“ identity_card”。可以使用随附的EncryptedCredentials解密和验证文件。
	Selfie      *PassportFile  `json:"selfie"`       // 可选的。用户提供的带有文件持有人自拍照的加密文件；适用于“护照”，“驾驶执照”，“身份证”和“内部护照”。可以使用随附的EncryptedCredentials解密和验证文件。
	Translation []PassportFile `json:"translation"`  // 可选的。带有用户提供的文档翻译版本的加密文件数组。如果要求提供“passport”, “driver_license”, “identity_card”, “internal_passport”, “utility_bill”, “bank_statement”, “rental_agreement”, “passport_registration” and “temporary_registration” 类型，则可用。可以使用随附的EncryptedCredentials解密和验证文件。
	Hash        string         `json:"hash"`         // 在PassportElementErrorUnspecified中使用的Base64编码的元素哈希
}

// EncryptedCredentials 包含解密和认证EncryptedPassportElement所需的数据。有关数据解密和身份验证过程的完整说明，请参阅电报护照文档。
// https://core.telegram.org/bots/api#encryptedcredentials
type EncryptedCredentials struct {
	Data   string `json:"data"`   // Base64编码的加密JSON序列化数据，具有唯一的用户有效载荷，EncryptedPassportElement解密和身份验证所需的数据哈希和秘密
	Hash   string `json:"hash"`   // 用于数据认证的Base64编码的数据哈希
	Secret string `json:"secret"` // Base64编码的机密，已用漫游器的公共RSA密钥加密，是数据解密所必需的
}
