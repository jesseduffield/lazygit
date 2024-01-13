package m

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// MorrowindScript lexer.
var MorrowindScript = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "MorrowindScript",
		Aliases:   []string{"morrowind", "mwscript"},
		Filenames: []string{},
		MimeTypes: []string{},
	},
	morrowindScriptRules,
))

func morrowindScriptRules() Rules {
	return Rules{
		"root": {
			{`\s+`, Text, nil},
			{`;.*$`, Comment, nil},
			{`(["'])(?:(?=(\\?))\2.)*?\1`, LiteralString, nil},
			{`[0-9]+`, LiteralNumberInteger, nil},
			{`[0-9]+\.[0-9]*(?!\.)`, LiteralNumberFloat, nil},
			Include("keywords"),
			Include("types"),
			Include("builtins"),
			Include("punct"),
			Include("operators"),
			{`\n`, Text, nil},
			{`\S+\s+`, Text, nil},
			{`[a-zA-Z0-9_]\w*`, Name, nil},
		},
		"keywords": {
			{`(?i)(begin|if|else|elseif|endif|while|endwhile|return|to)\b`, Keyword, nil},
			{`(?i)(end)\b`, Keyword, nil},
			{`(?i)(end)\w+.*$`, Text, nil},
			{`[\w+]->[\w+]`, Operator, nil},
		},
		"builtins": {
			{`(?i)(Activate|AddItem|AddSoulGem|AddSpell|AddToLevCreature|AddToLevItem|AddTopic|AIActivate|AIEscort|AIEscortCell|AIFollow|AiFollowCell|AITravel|AIWander|BecomeWerewolf|Cast|ChangeWeather|Choice|ClearForceJump|ClearForceMoveJump|ClearForceRun|ClearForceSneak|ClearInfoActor|Disable|DisableLevitation|DisablePlayerControls|DisablePlayerFighting|DisablePlayerJumping|DisablePlayerLooking|DisablePlayerMagic|DisablePlayerViewSwitch|DisableTeleporting|DisableVanityMode|DontSaveObject|Drop|Enable|EnableBirthMenu|EnableClassMenu|EnableInventoryMenu|EnableLevelUpMenu|EnableLevitation|EnableMagicMenu|EnableMapMenu|EnableNameMenu|EnablePlayerControls|EnablePlayerFighting|EnablePlayerJumping|EnablePlayerLooking|EnablePlayerMagic|EnablePlayerViewSwitch|EnableRaceMenu|EnableRest|EnableStatsMenu|EnableTeleporting|EnableVanityMode|Equip|ExplodeSpell|Face|FadeIn|FadeOut|FadeTo|Fall|ForceGreeting|ForceJump|ForceRun|ForceSneak|Flee|GotoJail|HurtCollidingActor|HurtStandingActor|Journal|Lock|LoopGroup|LowerRank|MenuTest|MessageBox|ModAcrobatics|ModAgility|ModAlarm|ModAlchemy|ModAlteration|ModArmorBonus|ModArmorer|ModAthletics|ModAttackBonus|ModAxe|ModBlock|ModBluntWeapon|ModCastPenalty|ModChameleon|ModConjuration|ModCurrentFatigue|ModCurrentHealth|ModCurrentMagicka|ModDefendBonus|ModDestruction|ModDisposition|ModEnchant|ModEndurance|ModFactionReaction|ModFatigue|ModFight|ModFlee|ModFlying|ModHandToHand|ModHealth|ModHeavyArmor|ModIllusion|ModIntelligence|ModInvisible|ModLightArmor|ModLongBlade|ModLuck|ModMagicka|ModMarksman|ModMediumArmor|ModMercantile|ModMysticism|ModParalysis|ModPCCrimeLevel|ModPCFacRep|ModPersonality|ModRegion|ModReputation|ModResistBlight|ModResistCorprus|ModResistDisease|ModResistFire|ModResistFrost|ModResistMagicka|ModResistNormalWeapons|ModResistParalysis|ModResistPoison|ModResistShock|ModRestoration|ModScale|ModSecurity|ModShortBlade|ModSilence|ModSneak|ModSpear|ModSpeechcraft|ModSpeed|ModStrength|ModSuperJump|ModSwimSpeed|ModUnarmored|ModWaterBreathing|ModWaterLevel|ModWaterWalking|ModWillpower|Move|MoveWorld|PayFine|PayFineThief|PCClearExpelled|PCExpell|PCForce1stPerson|PCForce3rdPerson|PCJoinFaction|PCLowerRank|PCRaiseRank|PlaceAtMe|PlaceAtPC|PlaceItem|PlaceItemCell|PlayBink|PlayGroup|PlayLoopSound3D|PlayLoopSound3DVP|PlaySound|PlaySound3D|PlaySound3DVP|PlaySoundVP|Position|PositionCell|RaiseRank|RemoveEffects|RemoveFromLevCreature|RemoveFromLevItem|RemoveItem|RemoveSoulgem|RemoveSpell|RemoveSpellEffects|ResetActors|Resurrect|Rotate|RotateWorld|Say|StartScript|[S|s]et|SetAcrobatics|SetAgility|SetAlarm|SetAlchemy|SetAlteration|SetAngle|SetArmorBonus|SetArmorer|SetAthletics|SetAtStart|SetAttackBonus|SetAxe|SetBlock|SetBluntWeapon|SetCastPenalty|SetChameleon|SetConjuration|SetDelete|SetDefendBonus|SetDestruction|SetDisposition|SetEnchant|SetEndurance|SetFactionReaction|SetFatigue|SetFight|SetFlee|SetFlying|SetHandToHand|SetHealth|SetHeavyArmor|SetIllusion|SetIntelligence|SetInvisible|SetJournalIndex|SetLightArmor|SetLevel|SetLongBlade|SetLuck|SetMagicka|SetMarksman|SetMediumArmor|SetMercantile|SetMysticism|SetParalysis|SetPCCCrimeLevel|SetPCFacRep|SetPersonality|SetPos|SetReputation|SetResistBlight|SetResistCorprus|SetResistDisease|SetResistFire|SetResistFrost|SetResistMagicka|SetResistNormalWeapons|SetResistParalysis|SetResistPoison|SetResistShock|SetRestoration|SetScale|SetSecurity|SetShortBlade|SetSilence|SetSneak|SetSpear|SetSpeechcraft|SetSpeed|SetStrength|SetSuperJump|SetSwimSpeed|SetUnarmored|SetWaterBreathing|SetWaterlevel|SetWaterWalking|SetWerewolfAcrobatics|SetWillpower|ShowMap|ShowRestMenu|SkipAnim|StartCombat|StopCombat|StopScript|StopSound|StreamMusic|TurnMoonRed|TurnMoonWhite|UndoWerewolf|Unlock|WakeUpPC|CenterOnCell|CenterOnExterior|FillMap|FixMe|ToggleAI|ToggleCollision|ToggleFogOfWar|ToggleGodMode|ToggleMenus|ToggleSky|ToggleWorld|ToggleVanityMode|CellChanged|GetAcrobatics|GetAgility|GetAIPackageDone|GetAlarm|GetAlchemy|GetAlteration|GetAngle|GetArmorBonus|GetArmorer|GetAthletics|GetAttackBonus|GetAttacked|GetArmorType,|GetAxe|GetBlightDisease|GetBlock|GetBluntWeapon|GetButtonPressed|GetCastPenalty|GetChameleon|GetCollidingActor|GetCollidingPC|GetCommonDisease|GetConjuration|GetCurrentAIPackage|GetCurrentTime|GetCurrentWeather|GetDeadCount|GetDefendBonus|GetDestruction|GetDetected|GetDisabled|GetDisposition|GetDistance|GetEffect|GetEnchant|GetEndurance|GetFatigue|GetFight|GetFlee|GetFlying|GetForceJump|GetForceRun|GetForceSneak|GetHandToHand|GetHealth|GetHealthGetRatio|GetHeavyArmor|GetIllusion|GetIntelligence|GetInterior|GetInvisible|GetItemCount|GetJournalIndex|GetLightArmor|GetLineOfSight|GetLOS|GetLevel|GetLocked|GetLongBlade|GetLuck|GetMagicka|GetMarksman|GetMasserPhase|GetSecundaPhase|GetMediumArmor|GetMercantile|GetMysticism|GetParalysis|GetPCCell|GetPCCrimeLevel|GetPCinJail|GetPCJumping|GetPCRank|GetPCRunning|GetPCSleep|GetPCSneaking|GetPCTraveling|GetPersonality|GetPlayerControlsDisabled|GetPlayerFightingDisabled|GetPlayerJumpingDisabled|GetPlayerLookingDisabled|GetPlayerMagicDisabled|GetPos|GetRace|GetReputation|GetResistBlight|GetResistCorprus|GetResistDisease|GetResistFire|GetResistFrost|GetResistMagicka|GetResistNormalWeapons|GetResistParalysis|GetResistPoison|GetResistShock|GetRestoration|GetScale|GetSecondsPassed|GetSecurity|GetShortBlade|GetSilence|GetSneak|GetSoundPlaying|GetSpear|GetSpeechcraft|GetSpeed|GetSpell|GetSpellEffects|GetSpellReadied|GetSquareRoot|GetStandingActor|GetStandingPC|GetStrength|GetSuperJump|GetSwimSpeed|GetTarget|GetUnarmored|GetVanityModeDisabled|GetWaterBreathing|GetWaterLevel|GetWaterWalking|GetWeaponDrawn|GetWeaponType|GetWerewolfKills|GetWillpower|GetWindSpeed|HasItemEquipped|HasSoulgem|HitAttemptOnMe|HitOnMe|IsWerewolf|MenuMode|OnActivate|OnDeath|OnKnockout|OnMurder|PCExpelled|PCGet3rdPerson|PCKnownWerewolf|Random|RepairedOnMe|SameFaction|SayDone|ScriptRunning|AllowWereWolfForceGreeting|Companion|MinimumProfit|NoFlee|NoHello|NoIdle|NoLore|OnPCAdd|OnPCDrop|OnPCEquip|OnPCHitMe|OnPCRepair|PCSkipEquip|OnPCSoulGemUse|StayOutside|CrimeGoldDiscount|CrimeGoldTurnIn|Day|DaysPassed|GameHour|Month|NPCVoiceDistance|PCRace|PCWerewolf|PCVampire|TimeScale|VampClan|Year)\b`, NameBuiltin, nil},
			{`(?i)(sEffectWaterBreathing|sEffectSwiftSwim|sEffectWaterWalking|sEffectShield|sEffectFireShield|sEffectLightningShield|sEffectFrostShield|sEffectBurden|sEffectFeather|sEffectJump|sEffectLevitate|sEffectSlowFall|sEffectLock|sEffectOpen|sEffectFireDamage|sEffectShockDamage|sEffectFrostDamage|sEffectDrainAttribute|sEffectDrainHealth|sEffectDrainSpellpoints|sEffectDrainFatigue|sEffectDrainSkill|sEffectDamageAttribute|sEffectDamageHealth|sEffectDamageMagicka|sEffectDamageFatigue|sEffectDamageSkill|sEffectPoison|sEffectWeaknessToFire|sEffectWeaknessToFrost|sEffectWeaknessToShock|sEffectWeaknessToMagicka|sEffectWeaknessToCommonDisease|sEffectWeaknessToBlightDisease|sEffectWeaknessToCorprusDisease|sEffectWeaknessToPoison|sEffectWeaknessToNormalWeapons|sEffectDisintegrateWeapon|sEffectDisintegrateArmor|sEffectInvisibility|sEffectChameleon|sEffectLight|sEffectSanctuary|sEffectNightEye|sEffectCharm|sEffectParalyze|sEffectSilence|sEffectBlind|sEffectSound|sEffectCalmHumanoid|sEffectCalmCreature|sEffectFrenzyHumanoid|sEffectFrenzyCreature|sEffectDemoralizeHumanoid|sEffectDemoralizeCreature|sEffectRallyHumanoid|sEffectRallyCreature|sEffectDispel|sEffectSoultrap|sEffectTelekinesis|sEffectMark|sEffectRecall|sEffectDivineIntervention|sEffectAlmsiviIntervention|sEffectDetectAnimal|sEffectDetectEnchantment|sEffectDetectKey|sEffectSpellAbsorption|sEffectReflect|sEffectCureCommonDisease|sEffectCureBlightDisease|sEffectCureCorprusDisease|sEffectCurePoison|sEffectCureParalyzation|sEffectRestoreAttribute|sEffectRestoreHealth|sEffectRestoreSpellPoints|sEffectRestoreFatigue|sEffectRestoreSkill|sEffectFortifyAttribute|sEffectFortifyHealth|sEffectFortifySpellpoints|sEffectFortifyFatigue|sEffectFortifySkill|sEffectFortifyMagickaMultiplier|sEffectAbsorbAttribute|sEffectAbsorbHealth|sEffectAbsorbSpellPoints|sEffectAbsorbFatigue|sEffectAbsorbSkill|sEffectResistFire|sEffectResistFrost|sEffectResistShock|sEffectResistMagicka|sEffectResistCommonDisease|sEffectResistBlightDisease|sEffectResistCorprusDisease|sEffectResistPoison|sEffectResistNormalWeapons|sEffectResistParalysis|sEffectRemoveCurse|sEffectTurnUndead|sEffectSummonScamp|sEffectSummonClannfear|sEffectSummonDaedroth|sEffectSummonDremora|sEffectSummonAncestralGhost|sEffectSummonSkeletalMinion|sEffectSummonLeastBonewalker|sEffectSummonGreaterBonewalker|sEffectSummonBonelord|sEffectSummonWingedTwilight|sEffectSummonHunger|sEffectSummonGoldensaint|sEffectSummonFlameAtronach|sEffectSummonFrostAtronach|sEffectSummonStormAtronach|sEffectFortifyAttackBonus|sEffectCommandCreatures|sEffectCommandHumanoids|sEffectBoundDagger|sEffectBoundLongsword|sEffectBoundMace|sEffectBoundBattleAxe|sEffectBoundSpear|sEffectBoundLongbow|sEffectExtraSpell|sEffectBoundCuirass|sEffectBoundHelm|sEffectBoundBoots|sEffectBoundShield|sEffectBoundGloves|sEffectCorpus|sEffectVampirism|sEffectSummonCenturionSphere|sEffectSunDamage|sEffectStuntedMagicka)`, NameBuiltin, nil},
		},
		"types": {
			{`(?i)(short|long|float)\b`, KeywordType, nil},
		},
		"punct": {
			{`[()]`, Punctuation, nil},
		},
		"operators": {
			{`[#=,./%+\-?]`, Operator, nil},
			{`(==|<=|<|>=|>|!=)`, Operator, nil},
		},
	}
}
