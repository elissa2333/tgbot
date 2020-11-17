package telegram

import (
	"fmt"

	"github.com/elissa2333/tgbot/utils"
)

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

// SetPassportDataErrors 通知用户他们提供的某些Telegram Passport元素包含错误。在错误得到纠正之前，用户将无法重新提交其护照（必须更改返回错误的字段的内容）。成功返回True。
//如果用户提交的数据出于任何原因不满足您的服务要求的标准，请使用此选项。例如，如果生日似乎无效，提交的文档模糊，扫描显示篡改的证据等。请在错误消息中提供一些详细信息，以确保用户知道如何解决问题。
// https://core.telegram.org/bots/api#setpassportdataerrors
func (a API) SetPassportDataErrors(userID int64, errors []PassportElementError) (bool, error) {
	var err error
	var cache []map[string]interface{}
	for k, v := range errors {
		m := map[string]interface{}{}
		switch {
		default:
			return false, fmt.Errorf("%d all fiald is nil", k)
		case v.PassportElementErrorDataField != nil:
			m, err = utils.StructToMap(v.PassportElementErrorDataField)
		case v.PassportElementErrorFrontSide != nil:
			m, err = utils.StructToMap(v.PassportElementErrorFrontSide)
		case v.PassportElementErrorReverseSide != nil:
			m, err = utils.StructToMap(v.PassportElementErrorReverseSide)
		case v.PassportElementErrorSelfie != nil:
			m, err = utils.StructToMap(v.PassportElementErrorSelfie)
		case v.PassportElementErrorFile != nil:
			m, err = utils.StructToMap(v.PassportElementErrorFile)
		case v.PassportElementErrorFiles != nil:
			m, err = utils.StructToMap(v.PassportElementErrorFiles)
		case v.PassportElementErrorTranslationFile != nil:
			m, err = utils.StructToMap(v.PassportElementErrorTranslationFile)
		case v.PassportElementErrorTranslationFiles != nil:
			m, err = utils.StructToMap(v.PassportElementErrorTranslationFiles)
		case v.PassportElementErrorUnspecified != nil:
			m, err = utils.StructToMap(v.PassportElementErrorUnspecified)
		}
		if err != nil {
			return false, err
		}
		cache = append(cache, m)
	}

	var result bool
	err = a.handleOptional("/setPassportDataErrors", map[string]interface{}{"user_id": userID, "errors": cache}, nil, &result)
	return result, err
}

const (
	// PassportElementErrorTypeAtPersonalDetails 个人资料
	PassportElementErrorTypeAtPersonalDetails = "personal_details"
	// PassportElementErrorTypeAtPassport 护照
	PassportElementErrorTypeAtPassport = "passport"
	// PassportElementErrorTypeAtDriverLicense 驾驶执照
	PassportElementErrorTypeAtDriverLicense = "driver_license"
	// PassportElementErrorTypeAtIdentityCard 身份证
	PassportElementErrorTypeAtIdentityCard = "identity_card"
	// PassportElementErrorTypeAtInternalPassport 内部护照
	PassportElementErrorTypeAtInternalPassport = "internal_passport"
	// PassportElementErrorTypeAtAddress 地址
	PassportElementErrorTypeAtAddress = "address"
	// PassportElementErrorTypeAtUtilityBill 水电费
	PassportElementErrorTypeAtUtilityBill = "utility_bill"
	// PassportElementErrorTypeAtBankStatement 银行对帐单
	PassportElementErrorTypeAtBankStatement = "bank_statement"
	// PassportElementErrorTypeAtRentalAgreement 租赁协议
	PassportElementErrorTypeAtRentalAgreement = "rental_agreement"
	// PassportElementErrorTypeAtPassportRegistration 护照登记
	PassportElementErrorTypeAtPassportRegistration = "passport_registration"
	// PassportElementErrorTypeAtTemporaryRegistration 临时注册
	PassportElementErrorTypeAtTemporaryRegistration = "temporary_registration"
)

// PassportElementError 该对象表示已提交的“电报护照”元素中的错误，应由用户解决。应该是以下之一：
//https://core.telegram.org/bots/api#passportelementerror
type PassportElementError struct {
	*PassportElementErrorDataField
	*PassportElementErrorFrontSide
	*PassportElementErrorReverseSide
	*PassportElementErrorSelfie
	*PassportElementErrorFile
	*PassportElementErrorFiles
	*PassportElementErrorTranslationFile
	*PassportElementErrorTranslationFiles
	*PassportElementErrorUnspecified
}

// PassportElementErrorDataFieldSource PassportElementErrorDataField source 类型
const PassportElementErrorDataFieldSource = "data"

// PassportElementErrorDataField 代表用户提供的数据字段之一中的问题。当字段的值更改时，该错误被视为已解决。
// https://core.telegram.org/bots/api#passportelementerrordatafield
type PassportElementErrorDataField struct {
	Source    string `json:"source"`     // 错误源，必须 data
	Type      string `json:"type"`       // 用户电报护照中出现错误的部分，必须是 “personal_details”, “passport”, “driver_license”, “identity_card”, “internal_passport”, “address” 之一
	FieldName string `json:"field_name"` // 出现错误的数据字段名称
	DataHash  string `json:"data_hash"`  // Base64编码的数据哈希
	Message   string `json:"message"`    // 错误信息
}

// PassportElementErrorFrontSideSource PassportElementErrorFrontSide source 类型
const PassportElementErrorFrontSideSource = "front_side"

// PassportElementErrorFrontSide 代表文件正面有问题。当文档正面的文件更改时，该错误被视为已解决
// https://core.telegram.org/bots/api#passportelementerrorfrontside
type PassportElementErrorFrontSide struct {
	Source   string `json:"source"`    // 错误源，必须为front_side
	Type     string `json:"type"`      // 用户电报护照中出现问题的部分，“passport”, “driver_license”, “identity_card”, “internal_passport” 之一
	FileHash string `json:"file_hash"` // 带有文档正面的文件的Base64编码的哈希
	Message  string `json:"message"`   // 错误信息
}

// PassportElementErrorReverseSideSource PassportElementErrorReverseSide source 类型
const PassportElementErrorReverseSideSource = "reverse_side"

// PassportElementErrorReverseSide 代表文档背面出现问题。当文档背面的文件发生更改时，该错误被视为已解决
// https://core.telegram.org/bots/api#passportelementerrorreverseside
type PassportElementErrorReverseSide struct {
	Source   string `json:"source"`    // 错误源，必须为reverse_side
	Type     string `json:"type"`      // 用户电报护照中出现问题的部分，“ driver_license”，“ identity_card”之一
	FileHash string `json:"file_hash"` // 文件的Base64编码的哈希值以及文档的背面
	Message  string `json:"message"`   // 错误信息
}

// PassportElementErrorSelfieSource PassportElementErrorSelfie source 类型
const PassportElementErrorSelfieSource = "selfie"

// PassportElementErrorSelfie 代表带有文档的自拍照的问题。带有自拍照的文件更改时，该错误被视为已解决
// https://core.telegram.org/bots/api#passportelementerrorselfie
type PassportElementErrorSelfie struct {
	Source   string `json:"source"`    // 错误来源，必须是 selfie
	Type     string `json:"type"`      // 用户电报护照中出现问题的部分，“passport”, “driver_license”, “identity_card”, “internal_passport”之一
	FileHash string `json:"file_hash"` // 文件的Base64编码的哈希值以及文档的背面
	Message  string `json:"message"`   // 错误信息
}

// PassportElementErrorFileSource PassportElementErrorFile source 类型
const PassportElementErrorFileSource = "file"

// PassportElementErrorFile 代表文件扫描有问题。当包含文档扫描的文件更改时，该错误被视为已解决
// https://core.telegram.org/bots/api#passportelementerrorfile
type PassportElementErrorFile struct {
	Source   string `json:"source"`    // 错误源，必须是文件
	Type     string `json:"type"`      // 用户电报护照中出现问题的部分，其中一个是 “utility_bill”, “bank_statement”, “rental_agreement”, “passport_registration”, “temporary_registration”
	FileHash string `json:"file_hash"` // 文件的Base64编码的哈希值以及文档的背面
	Message  string `json:"message"`   // 错误信息
}

// PassportElementErrorFilesSource PassportElementErrorFiles source 类型
const PassportElementErrorFilesSource = " files"

// PassportElementErrorFiles 代表扫描列表问题。当包含扫描的文件列表更改时，该错误被视为已解决
// https://core.telegram.org/bots/api#passportelementerrorfiles
type PassportElementErrorFiles struct {
	Source     string   `json:"source"`      // 错误源，必须是文件
	Type       string   `json:"type"`        // 用户电报护照中出现问题的部分，其中一个是 “utility_bill”, “bank_statement”, “rental_agreement”, “passport_registration”, “temporary_registration”
	FileHashes []string `json:"file_hashes"` // base64编码的文件哈希列表
	Message    string   `json:"message"`     // 错误信息
}

// PassportElementErrorTranslationFileSource PassportElementErrorTranslationFile source 类型
const PassportElementErrorTranslationFileSource = "translation_file"

// PassportElementErrorTranslationFile 代表构成文档翻译的文件之一存在问题。文件更改后，该错误被视为已解决
// https://core.telegram.org/bots/api#passportelementerrortranslationfile
type PassportElementErrorTranslationFile struct {
	Source   string `json:"source"`    // 错误源，必须是translation_file
	Type     string `json:"type"`      // 出现问题的用户电报护照的元素类型，“passport”, “driver_license”, “identity_card”, “internal_passport”, “utility_bill”, “bank_statement”, “rental_agreement”, “passport_registration”, “temporary_registration” 之一
	FileHash string `json:"file_hash"` // 文件的Base64编码的哈希值以及文档的背面
	Message  string `json:"message"`   // 错误信息
}

// PassportElementErrorTranslationFilesSource PassportElementErrorTranslationFiles source 类型
const PassportElementErrorTranslationFilesSource = "translation_files"

// PassportElementErrorTranslationFiles 代表文档翻译版本存在问题。当包含文档翻译的文件更改时，该错误被视为已解决
// https://core.telegram.org/bots/api#passportelementerrortranslationfiles
type PassportElementErrorTranslationFiles struct {
	Source     string   `json:"source"`      // 错误源，必须是translation_files
	Type       string   `json:"type"`        // 出现问题的用户电报护照的元素类型，“passport”, “driver_license”, “identity_card”, “internal_passport”, “utility_bill”, “bank_statement”, “rental_agreement”, “passport_registration”, “temporary_registration” 之一
	FileHashes []string `json:"file_hashes"` // base64编码的文件哈希列表
	Message    string   `json:"message"`     // 错误信息
}

// PassportElementErrorUnspecifiedSource PassportElementErrorUnspecified source 类型
const PassportElementErrorUnspecifiedSource = "unspecified"

// PassportElementErrorUnspecified 代表未指定位置的问题。添加新数据后，该错误被视为已解决
// https://core.telegram.org/bots/api#passportelementerrorunspecified
type PassportElementErrorUnspecified struct {
	Source      string `json:"source"`       // 错误源，必须未指定
	Type        string `json:"type"`         // 有问题的用户电报护照元素的类型
	ElementHash string `json:"element_hash"` // Base64编码的元素哈希
	Message     string `json:"message"`      // 错误信息
}
