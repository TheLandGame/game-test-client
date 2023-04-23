package message

import (
	"game-message-core/proto"

	"github.com/spf13/cast"
)

func ToProtoNftData(nft NFT) *proto.NftData {
	return &proto.NftData{
		Network:    nft.Network,
		TokenId:    nft.TokenId,
		IsMelandAi: nft.IsMelandAI,
		Metadata:   GetNftMetaData(nft),
	}
}

func GetNftMetaData(nft NFT) *proto.NftMetadata {
	return &proto.NftMetadata{
		Name:            nft.Metadata.Name,
		Description:     nft.Metadata.Description,
		Image:           *nft.Metadata.Image,
		BackGroundColor: *nft.Metadata.BackgroundColor,
		Type:            GetNftTraitType(nft),
		Attributes:      GetNftEntityAttributes(nft),
		TraitData:       GetNftTraitData(nft),
		CreateSec:       cast.ToInt64(nft.CreatedAt),
	}
}

func GetNftTraitData(nft NFT) []*proto.NftTraitData {
	attributes := []*proto.NftTraitData{}
	for _, na := range nft.Metadata.Attributes {
		// 属性
		// _, exist := configData.ConfigMgr().EntityAttributeByName(na.TraitType)
		// if !exist {
		attr := &proto.NftTraitData{
			TraitType: na.TraitType,
			Value:     na.Value,
		}
		if na.DisplayType != nil {
			attr.DisplayType = *na.DisplayType
		}
		attributes = append(attributes, attr)
		// }
	}
	return attributes
}

func GetNftTypeTrait(nft NFT) *NFTAttribute {
	for _, na := range nft.Metadata.Attributes {
		if na.TraitType == string(NFTTraitTypesType) {
			return &na
		}
	}
	return nil
}

func GetNftTraitType(nft NFT) (traitType proto.NftTraitType) {
	if !nft.IsMelandAI {
		return proto.NftTraitType_Third
	}

	typeAttribute := GetNftTypeTrait(nft)
	if typeAttribute == nil {
		return traitType
	}

	switch typeAttribute.Value {
	case string(NFTTraitTypeHandsArmor): // "Hands Armor" 手部装备
		traitType = proto.NftTraitType_HandsArmor

	case string(NFTTraitTypeChestArmor): // "Chest Armor" 胸部装备
		traitType = proto.NftTraitType_ChestArmor

	case string(NFTTraitTypeHeadArmor): // "Head Armor" 头部装备
		traitType = proto.NftTraitType_HeadArmor

	case string(NFTTraitTypeLegsArmor): // "Legs Armor" 腿部装备
		traitType = proto.NftTraitType_LegsArmor

	case string(NFTTraitTypeFeetArmor): // "Feet Armor" 脚部装备
		traitType = proto.NftTraitType_FeetArmor

	case string(NFTTraitTypeSword): // "Sword" 剑
		traitType = proto.NftTraitType_Sword

	case string(NFTTraitTypeBow): // "Bow"  弓
		traitType = proto.NftTraitType_Bow

	case string(NFTTraitTypeDagger): // "Dagger" 匕首
		traitType = proto.NftTraitType_Dagger

	case string(NFTTraitTypeSpear): // "Spear"枪
		traitType = proto.NftTraitType_Spear

	case string(NFTTraitTypeConsumable): // "Consumable" 消耗品
		traitType = proto.NftTraitType_Consumable

	case string(NFTTraitTypeMaterial): // "Material" 材料
		traitType = proto.NftTraitType_Material

	case string(NFTTraitTypeMysteryBox): // "MysteryBox" 神秘宝箱
		traitType = proto.NftTraitType_MysteryBox

	case string(NFTTraitTypePlaceable): // "Placeable" 可放置
		traitType = proto.NftTraitType_Placeable

	case string(NFTTraitTypeWearable): // "Wearable" 可穿戴
		traitType = proto.NftTraitType_Wearable
	}

	return traitType
}

func GetNftEntityAttributes(nft NFT) []*proto.AttributeData {
	attributes := []*proto.AttributeData{}
	// for _, na := range nft.Metadata.Attributes {
	// 	// 属性
	// 	attrSetting, exist := configData.ConfigMgr().EntityAttributeByName(na.TraitType)
	// 	if exist {
	// 		attr := &proto.AttributeData{
	// 			Id:    int32(attrSetting.Id),
	// 			Value: cast.ToInt32(na.Value),
	// 		}
	// 		if na.DisplayType != nil {
	// 			attr.DisplayType = GetNftEntityAttributeDisType(*na.DisplayType)
	// 		}
	// 		attributes = append(attributes, attr)
	// 		continue
	// 	}
	// }
	return attributes
}

func GetNftEntityAttributeDisType(displayType string) proto.AttributeDisplayType {
	if displayType == "boost_percentage" {
		return proto.AttributeDisplayType_BoostPercentage
	}
	return proto.AttributeDisplayType_Number
}

func EquipmentPosition(traitType proto.NftTraitType) (position proto.AvatarPosition) {
	switch traitType {
	case proto.NftTraitType_HandsArmor: // "Hands Armor" 手部装备
		position = proto.AvatarPosition_AvatarPositionHand

	case proto.NftTraitType_ChestArmor: // "Chest Armor" 胸部装备
		position = proto.AvatarPosition_AvatarPositionCoat

	case proto.NftTraitType_HeadArmor: // "Head Armor" 头部装备
		position = proto.AvatarPosition_AvatarPositionHead

	case proto.NftTraitType_LegsArmor: // "Legs Armor" 腿部装备
		position = proto.AvatarPosition_AvatarPositionPant

	case proto.NftTraitType_FeetArmor: // "Feet Armor" 脚部装备
		position = proto.AvatarPosition_AvatarPositionShoe

	case proto.NftTraitType_Sword: // "Sword" 剑
		position = proto.AvatarPosition_AvatarPositionWeapon

	case proto.NftTraitType_Bow: // "Bow"  弓
		position = proto.AvatarPosition_AvatarPositionWeapon

	case proto.NftTraitType_Dagger: // "Dagger" 匕首
		position = proto.AvatarPosition_AvatarPositionWeapon

	case proto.NftTraitType_Spear: // "Spear"枪
		position = proto.AvatarPosition_AvatarPositionWeapon
	}
	return position
}

func GetConsumableData(traitData []*proto.NftTraitData) *proto.NFTConsumableInfo {
	data := &proto.NFTConsumableInfo{}
	// for _, na := range traitData {
	// 	if na.TraitType == string(NFTTraitTypesQuality) {
	// 		data.Quality = na.Value
	// 		continue
	// 	}

	// 	t, err := configData.ConsumableKeyToPbType(na.TraitType)
	// 	if err == nil {
	// 		data.ConsumableType = t
	// 		data.Value = cast.ToInt32(na.Value)
	// 		break
	// 	}
	// }
	return data
}

func GetNftUseLevel(traitData []*proto.NftTraitData) (needLv int32) {
	for _, na := range traitData {
		if na.TraitType == string(NFTTraitTypesRequiresLevel) {
			needLv = cast.ToInt32(na.Value)
			break
		}
	}
	return needLv
}
