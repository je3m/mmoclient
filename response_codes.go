package main

import "fmt"

type ArtifactsResponseCode int

const (
	// general
	CodeInvalidPayload  ArtifactsResponseCode = 422
	CodeTooManyRequests ArtifactsResponseCode = 429
	CodeNotFound        ArtifactsResponseCode = 404
	CodeFatalError      ArtifactsResponseCode = 500

	// account error codes
	CodeTokenInvalid           ArtifactsResponseCode = 452
	CodeTokenExpired           ArtifactsResponseCode = 453
	CodeTokenMissing           ArtifactsResponseCode = 454
	CodeTokenGenerationFail    ArtifactsResponseCode = 455
	CodeUsernameAlreadyUsed    ArtifactsResponseCode = 456
	CodeEmailAlreadyUsed       ArtifactsResponseCode = 457
	CodeSamePassword           ArtifactsResponseCode = 458
	CodeCurrentPasswordInvalid ArtifactsResponseCode = 459

	// character error codes
	CodeCharacterNotEnoughHp            ArtifactsResponseCode = 483
	CodeCharacterMaximumUtilitesEquiped ArtifactsResponseCode = 484
	CodeCharacterItemAlreadyEquiped     ArtifactsResponseCode = 485
	CodeCharacterLocked                 ArtifactsResponseCode = 486
	CodeCharacterNotThisTask            ArtifactsResponseCode = 474
	CodeCharacterTooManyItemsTask       ArtifactsResponseCode = 475
	CodeCharacterNoTask                 ArtifactsResponseCode = 487
	CodeCharacterTaskNotCompleted       ArtifactsResponseCode = 488
	CodeCharacterAlreadyTask            ArtifactsResponseCode = 489
	CodeCharacterAlreadyMap             ArtifactsResponseCode = 490
	CodeCharacterSlotEquipmentError     ArtifactsResponseCode = 491
	CodeCharacterGoldInsufficient       ArtifactsResponseCode = 492
	CodeCharacterNotSkillLevelRequired  ArtifactsResponseCode = 493
	CodeCharacterNameAlreadyUsed        ArtifactsResponseCode = 494
	CodeMaxCharactersReached            ArtifactsResponseCode = 495
	CodeCharacterNotLevelRequired       ArtifactsResponseCode = 496
	CodeCharacterInventoryFull          ArtifactsResponseCode = 497
	CodeCharacterNotFound               ArtifactsResponseCode = 498
	CodeCharacterInCooldown             ArtifactsResponseCode = 499

	// item error codes
	CodeItemInsufficientQuantity ArtifactsResponseCode = 471
	CodeItemInvalidEquipment     ArtifactsResponseCode = 472
	CodeItemRecyclingInvalidItem ArtifactsResponseCode = 473
	CodeItemInvalidConsumable    ArtifactsResponseCode = 476
	CodeMissingItem              ArtifactsResponseCode = 478

	// grand exchange error codes
	CodeGeMaxQuantity           ArtifactsResponseCode = 479
	CodeGeNotInStock            ArtifactsResponseCode = 480
	CodeGeNotThePrice           ArtifactsResponseCode = 482
	CodeGeTransactionInProgress ArtifactsResponseCode = 436
	CodeGeNoOrders              ArtifactsResponseCode = 431
	CodeGeMaxOrders             ArtifactsResponseCode = 433
	CodeGeTooManyItems          ArtifactsResponseCode = 434
	CodeGeSameAccount           ArtifactsResponseCode = 435
	CodeGeInvalidItem           ArtifactsResponseCode = 437
	CodeGeNotYourOrder          ArtifactsResponseCode = 438

	// bank error codes
	CodeBankInsufficientGold      ArtifactsResponseCode = 460
	CodeBankTransactionInProgress ArtifactsResponseCode = 461
	CodeBankFull                  ArtifactsResponseCode = 462

	// maps error codes
	CodeMapNotFound        ArtifactsResponseCode = 597
	CodeMapContentNotFound ArtifactsResponseCode = 598
)

type ResponseCodeError struct {
	code ArtifactsResponseCode
}

func (e ResponseCodeError) Error() string {
	return fmt.Sprintf("%v", e.code)
}

func (me ArtifactsResponseCode) String() string {
	switch me {
	case CodeInvalidPayload:
		return "InvalidPayload"
	case CodeTooManyRequests:
		return "TooManyRequests"
	case CodeNotFound:
		return "NotFound"
	case CodeFatalError:
		return "FatalError"
	case CodeTokenInvalid:
		return "TokenInvalid"
	case CodeTokenExpired:
		return "TokenExpired"
	case CodeTokenMissing:
		return "TokenMissing"
	case CodeTokenGenerationFail:
		return "TokenGenerationFail"
	case CodeUsernameAlreadyUsed:
		return "UsernameAlreadyUsed"
	case CodeEmailAlreadyUsed:
		return "EmailAlreadyUsed"
	case CodeSamePassword:
		return "SamePassword"
	case CodeCurrentPasswordInvalid:
		return "CurrentPasswordInvalid"
	case CodeCharacterNotEnoughHp:
		return "CharacterNotEnoughHp"
	case CodeCharacterMaximumUtilitesEquiped:
		return "CharacterMaximumUtilitesEquiped"
	case CodeCharacterItemAlreadyEquiped:
		return "CharacterItemAlreadyEquiped"
	case CodeCharacterLocked:
		return "CharacterLocked"
	case CodeCharacterNotThisTask:
		return "CharacterNotThisTask"
	case CodeCharacterTooManyItemsTask:
		return "CharacterTooManyItemsTask"
	case CodeCharacterNoTask:
		return "CharacterNoTask"
	case CodeCharacterTaskNotCompleted:
		return "CharacterTaskNotCompleted"
	case CodeCharacterAlreadyTask:
		return "CharacterAlreadyTask"
	case CodeCharacterAlreadyMap:
		return "CharacterAlreadyMap"
	case CodeCharacterSlotEquipmentError:
		return "CharacterSlotEquipmentError"
	case CodeCharacterGoldInsufficient:
		return "CharacterGoldInsufficient"
	case CodeCharacterNotSkillLevelRequired:
		return "CharacterNotSkillLevelRequired"
	case CodeCharacterNameAlreadyUsed:
		return "CharacterNameAlreadyUsed"
	case CodeMaxCharactersReached:
		return "MaxCharactersReached"
	case CodeCharacterNotLevelRequired:
		return "CharacterNotLevelRequired"
	case CodeCharacterInventoryFull:
		return "CharacterInventoryFull"
	case CodeCharacterNotFound:
		return "CharacterNotFound"
	case CodeCharacterInCooldown:
		return "CharacterInCooldown"
	case CodeItemInsufficientQuantity:
		return "ItemInsufficientQuantity"
	case CodeItemInvalidEquipment:
		return "ItemInvalidEquipment"
	case CodeItemRecyclingInvalidItem:
		return "ItemRecyclingInvalidItem"
	case CodeItemInvalidConsumable:
		return "ItemInvalidConsumable"
	case CodeMissingItem:
		return "MissingItem"
	case CodeGeMaxQuantity:
		return "GeMaxQuantity"
	case CodeGeNotInStock:
		return "GeNotInStock"
	case CodeGeNotThePrice:
		return "GeNotThePrice"
	case CodeGeTransactionInProgress:
		return "GeTransactionInProgress"
	case CodeGeNoOrders:
		return "GeNoOrders"
	case CodeGeMaxOrders:
		return "GeMaxOrders"
	case CodeGeTooManyItems:
		return "GeTooManyItems"
	case CodeGeSameAccount:
		return "GeSameAccount"
	case CodeGeInvalidItem:
		return "GeInvalidItem"
	case CodeGeNotYourOrder:
		return "GeNotYourOrder"
	case CodeBankInsufficientGold:
		return "BankInsufficientGold"
	case CodeBankTransactionInProgress:
		return "BankTransactionInProgress"
	case CodeBankFull:
		return "BankFull"
	case CodeMapNotFound:
		return "MapNotFound"
	case CodeMapContentNotFound:
		return "MapContentNotFound"
	default:
		return "BAD_ERROR_CODE"
	}
}
