/*
	Copyright (C) 2022  ikafly144

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/
package mc

import (
	_ "embed"
	"encoding/json"
	"strings"
)

type CommandName string

const (
	CommandNameAdvancement     CommandName = "advancement"
	CommandNameAttribute       CommandName = "attribute"
	CommandNameBan             CommandName = "ban"
	CommandNameBanIp           CommandName = "ban-ip"
	CommandNameBanlist         CommandName = "banlist"
	CommandNameBossbar         CommandName = "bossbar"
	CommandNameClear           CommandName = "clear"
	CommandNameClone           CommandName = "clone"
	CommandNameData            CommandName = "data"
	CommandNameDatapack        CommandName = "datapack"
	CommandNameDebug           CommandName = "debug"
	CommandNameDefaultgamemode CommandName = "defaultgamemode"
	CommandNameDeop            CommandName = "deop"
	CommandNameDifficulty      CommandName = "difficulty"
	CommandNameEffect          CommandName = "effect"
	CommandNameEnchant         CommandName = "enchant"
	CommandNameExecute         CommandName = "execute"
	CommandNameExperience      CommandName = "experience"
	CommandNameFill            CommandName = "fill"
	CommandNameFillbiome       CommandName = "fillbiome"
	CommandNameForceload       CommandName = "forceload"
	CommandNameFunction        CommandName = "function"
	CommandNameGamemode        CommandName = "gamemode"
	CommandNameGamerule        CommandName = "gamerule"
	CommandNameGive            CommandName = "give"
	CommandNameHelp            CommandName = "help"
	CommandNameItem            CommandName = "item"
	CommandNameJfr             CommandName = "jfr"
	CommandNameKick            CommandName = "kick"
	CommandNameKill            CommandName = "kill"
	CommandNameList            CommandName = "list"
	CommandNameLocate          CommandName = "locate"
	CommandNameLoot            CommandName = "loot"
	CommandNameMe              CommandName = "me"
	CommandNameMsg             CommandName = "msg"
	CommandNameOp              CommandName = "op"
	CommandNamePardon          CommandName = "pardon"
	CommandNamePardonIp        CommandName = "pardon-ip"
	CommandNameParticle        CommandName = "particle"
	CommandNamePerf            CommandName = "perf"
	CommandNamePlace           CommandName = "place"
	CommandNamePlaysound       CommandName = "playsound"
	CommandNamePublish         CommandName = "publish"
	CommandNameRecipe          CommandName = "recipe"
	CommandNameReload          CommandName = "reload"
	CommandNameSaveAll         CommandName = "save-all"
	CommandNameSaveOff         CommandName = "save-off"
	CommandNameSaveOn          CommandName = "save-on"
	CommandNameSay             CommandName = "say"
	CommandNameSchedule        CommandName = "schedule"
	CommandNameScoreboard      CommandName = "scoreboard"
	CommandNameSeed            CommandName = "seed"
	CommandNameSetblock        CommandName = "setblock"
	CommandNameSetidletimeout  CommandName = "setidletimeout"
	CommandNameSetworldspawn   CommandName = "setworldspawn"
	CommandNameSpawnpoint      CommandName = "spawnpoint"
	CommandNameSpectate        CommandName = "spectate"
	CommandNameSpreadplayers   CommandName = "spreadplayers"
	CommandNameStop            CommandName = "stop"
	CommandNameStopsound       CommandName = "stopsound"
	CommandNameSummon          CommandName = "summon"
	CommandNameTag             CommandName = "tag"
	CommandNameTeam            CommandName = "team"
	CommandNameTeammsg         CommandName = "teammsg"
	CommandNameTeleport        CommandName = "teleport"
	CommandNameTell            CommandName = "tell"
	CommandNameTellraw         CommandName = "tellraw"
	CommandNameTime            CommandName = "time"
	CommandNameTitle           CommandName = "title"
	CommandNameTm              CommandName = "tm"
	CommandNameTp              CommandName = "tp"
	CommandNameTrigger         CommandName = "trigger"
	CommandNameW               CommandName = "w"
	CommandNameWeather         CommandName = "weather"
	CommandNameWhitelist       CommandName = "whitelist"
	CommandNameWorldborder     CommandName = "worldborder"
	CommandNameXp              CommandName = "xp"
)

type CommandType string

const (
	CommandTypeRoot     CommandType = "root"
	CommandTypeLiteral  CommandType = "literal"
	CommandTypeArgument CommandType = "argument"
)

type CommandChildren map[string]Command

type CommandParser string

const (
	CommandParserBrigadierBool              CommandParser = "brigadier:bool"
	CommandParserBrigadierDouble            CommandParser = "brigadier:double"
	CommandParserBrigadierFloat             CommandParser = "brigadier:float"
	CommandParserBrigadierInteger           CommandParser = "brigadier:integer"
	CommandParserBrigadierString            CommandParser = "brigadier:string"
	CommandParserMinecraftAngle             CommandParser = "minecraft:angle"
	CommandParserMinecraftBlockPos          CommandParser = "minecraft:block_pos"
	CommandParserMinecraftBlockPredicate    CommandParser = "minecraft:block_predicate"
	CommandParserMinecraftBlockState        CommandParser = "minecraft:block_state"
	CommandParserMinecraftColor             CommandParser = "minecraft:color"
	CommandParserMinecraftColumnPos         CommandParser = "minecraft:column_pos"
	CommandParserMinecraftComponent         CommandParser = "minecraft:component"
	CommandParserMinecraftDimension         CommandParser = "minecraft:dimension"
	CommandParserMinecraftEnchantment       CommandParser = "minecraft:enchantment"
	CommandParserMinecraftEntity            CommandParser = "minecraft:entity"
	CommandParserMinecraftEntityAnchor      CommandParser = "minecraft:entity_anchor"
	CommandParserMinecraftFunction          CommandParser = "minecraft:function"
	CommandParserMinecraftGameProfile       CommandParser = "minecraft:game_profile"
	CommandParserMinecraftGamemode          CommandParser = "minecraft:gamemode"
	CommandParserMinecraftIntRange          CommandParser = "minecraft:int_range"
	CommandParserMinecraftItemPredicate     CommandParser = "minecraft:item_predicate"
	CommandParserMinecraftItemSlot          CommandParser = "minecraft:item_slot"
	CommandParserMinecraftItemStack         CommandParser = "minecraft:item_stack"
	CommandParserMinecraftMessage           CommandParser = "minecraft:message"
	CommandParserMinecraftNBTCompoundTag    CommandParser = "minecraft:nbt_compound_tag"
	CommandParserMinecraftNBTPath           CommandParser = "minecraft:nbt_path"
	CommandParserMinecraftNBTTag            CommandParser = "minecraft:nbt_tag"
	CommandParserMinecraftObjective         CommandParser = "minecraft:objective"
	CommandParserMinecraftObjectiveCriteria CommandParser = "minecraft:objective_criteria"
	CommandParserMinecraftOperation         CommandParser = "minecraft:operation"
	CommandParserMinecraftParticle          CommandParser = "minecraft:particle"
	CommandParserMinecraftResource          CommandParser = "minecraft:resource"
	CommandParserMinecraftResourceKey       CommandParser = "minecraft:resource_key"
	CommandParserMinecraftResourceLocation  CommandParser = "minecraft:resource_location"
	CommandParserMinecraftResourceOrTag     CommandParser = "minecraft:resource_or_tag"
	CommandParserMinecraftResourceOrTagKey  CommandParser = "minecraft:resource_or_tag_key"
	CommandParserMinecraftRotation          CommandParser = "minecraft:rotation"
	CommandParserMinecraftScoreHolder       CommandParser = "minecraft:score_holder"
	CommandParserMinecraftScoreboardSlot    CommandParser = "minecraft:scoreboard_slot"
	CommandParserMinecraftSwizzle           CommandParser = "minecraft:swizzle"
	CommandParserMinecraftTeam              CommandParser = "minecraft:team"
	CommandParserMinecraftTemplateMirror    CommandParser = "minecraft:template_mirror"
	CommandParserMinecraftTime              CommandParser = "minecraft:time"
	CommandParserMinecraftUUID              CommandParser = "minecraft:uuid"
	CommandParserMinecraftVec2              CommandParser = "minecraft:vec2"
	CommandParserMinecraftVec3              CommandParser = "minecraft:vec3"
)

type CommandPropertiesType string

const (
	CommandPropertiesTypePlayers  CommandPropertiesType = "players"
	CommandPropertiesTypeEntities CommandPropertiesType = "entities"
	CommandPropertiesTypeGreedy   CommandPropertiesType = "greedy"
	CommandPropertiesTypePhrase   CommandPropertiesType = "phrase"
	CommandPropertiesTypeWord     CommandPropertiesType = "word"
)

type CommandPropertiesAmount string

const (
	CommandPropertiesAmountMultiple CommandPropertiesAmount = "multiple"
	CommandPropertiesAmountSingle   CommandPropertiesAmount = "single"
)

type CommandPropertiesRegistry string

const (
	CommandPropertiesRegistryMinecraftAttribute                 CommandPropertiesRegistry = "minecraft:attribute"
	CommandPropertiesRegistryMinecraftEnchantment               CommandPropertiesRegistry = "minecraft:enchantment"
	CommandPropertiesRegistryMinecraftEntityType                CommandPropertiesRegistry = "minecraft:entity_type"
	CommandPropertiesRegistryMinecraftMobEffect                 CommandPropertiesRegistry = "minecraft:mob_effect"
	CommandPropertiesRegistryMinecraftPointOfInterestType       CommandPropertiesRegistry = "minecraft:point_of_interest_type"
	CommandPropertiesRegistryMinecraftWorldgenBiome             CommandPropertiesRegistry = "minecraft:worldgen/biome"
	CommandPropertiesRegistryMinecraftWorldgenConfiguredFeature CommandPropertiesRegistry = "minecraft:worldgen/configured_feature"
	CommandPropertiesRegistryMinecraftWorldgenStructure         CommandPropertiesRegistry = "minecraft:worldgen/structure"
	CommandPropertiesRegistryMinecraftWorldgenTemplatePool      CommandPropertiesRegistry = "minecraft:worldgen/template_pool"
)

type CommandProperties struct {
	Type     *CommandPropertiesType     `json:"type,omitempty"`
	Amount   *CommandPropertiesAmount   `json:"amount,omitempty"`
	Registry *CommandPropertiesRegistry `json:"registry,omitempty"`
	Min      *float32                   `json:"min,omitempty"`
	Max      *float32                   `json:"max,omitempty"`
}

type Command struct {
	Type       CommandType        `json:"type"`
	Children   *CommandChildren   `json:"children,omitempty"`
	Executable *bool              `json:"executable,omitempty"`
	Parser     *CommandParser     `json:"parser,omitempty"`
	Properties *CommandProperties `json:"properties,omitempty"`
	Redirect   *[]CommandName     `json:"redirect,omitempty"`
}

type CommandNameSpace map[CommandName]Command

type CommandRoot struct {
	Type     CommandType      `json:"type"`
	Commands CommandNameSpace `json:"children"`
}

//go:embed commands.json
var commandsByte []byte

var commands CommandRoot

func init() {
	json.Unmarshal(commandsByte, &commands)
}

func Test(str string) (check bool, err error) {
	args := strings.Split(str, " ")
	var c Command
	for i, v := range args {
		switch i {
		case 0:
			s := CommandName(v)
			var ok bool
			c, ok = commands.Commands[s]
			if ok {
				if *c.Executable {
					return true, nil
				}
			}
		default:
			cm := *c.Children
			var ok bool
			c, ok = cm[v]
			if ok {
				if *c.Executable {
					return true, nil
				}
			}
		}
	}
	return
}
