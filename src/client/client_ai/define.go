package client_ai

type UserState int

const (
	USER_STATE_READY UserState = iota
	USER_STATE_IDLE
	USER_STATE_MOVE
	USER_STATE_ATTACK
)

type AttributeId int

const (
	Attribute_CollectionLv      AttributeId = 1
	Attribute_CombatLv          AttributeId = 2
	Attribute_FarmingLv         AttributeId = 3
	Attribute_CollectionExp     AttributeId = 4
	Attribute_CombatExp         AttributeId = 5
	Attribute_FarmingExp        AttributeId = 6
	Attribute_MoveSpd           AttributeId = 7
	Attribute_HomeExtraMoveSpd  AttributeId = 8
	Attribute_GearAvailableLv   AttributeId = 9
	Attribute_HP                AttributeId = 10
	Attribute_MaxHP             AttributeId = 11
	Attribute_HpRecovery        AttributeId = 12
	Attribute_CombatAtt         AttributeId = 13
	Attribute_CombatAttSpd      AttributeId = 14
	Attribute_CombatDef         AttributeId = 15
	Attribute_CombatHit         AttributeId = 16
	Attribute_CombatDodge       AttributeId = 17
	Attribute_CombatVulnerable  AttributeId = 18
	Attribute_CombatCritRate    AttributeId = 19
	Attribute_CombatCritDmg     AttributeId = 20
	Attribute_CombatDmgBonus    AttributeId = 21
	Attribute_TreeAtt           AttributeId = 22
	Attribute_AxeSpd            AttributeId = 23
	Attribute_TreeCritRate      AttributeId = 24
	Attribute_TreeCritDmg       AttributeId = 25
	Attribute_TreeDmgBonus      AttributeId = 26
	Attribute_ExtraWoodRate     AttributeId = 27
	Attribute_TreeAvailableLv   AttributeId = 28
	Attribute_OreAtt            AttributeId = 29
	Attribute_PickaxeSpd        AttributeId = 30
	Attribute_OreCritRate       AttributeId = 31
	Attribute_OreCritDmg        AttributeId = 32
	Attribute_OreDmgBonus       AttributeId = 33
	Attribute_ExtraOreRate      AttributeId = 34
	Attribute_OreAvailableLv    AttributeId = 35
	Attribute_HoeingEffect      AttributeId = 36
	Attribute_HoeSpd            AttributeId = 37
	Attribute_WateringEffect    AttributeId = 38
	Attribute_ExtraWateringRate AttributeId = 39
	Attribute_ExtraHarvestRate  AttributeId = 40
	Attribute_HarvestSeedRate   AttributeId = 41
	Attribute_CropAvailableLv   AttributeId = 42
)
