package main

import (
	. "../utils"
	"fmt"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Group struct {
	id                int
	allegiance        bool
	units             int
	hp                int
	damage            int
	damageType        string
	receiveMultiplier map[string]int
	initiative        int
	target            *Group
}

func (data Group) String() string {
	team := "Evil"
	if data.allegiance {
		team = "Good"
	}
	target := "None"
	if data.target != nil {
		target = strconv.Itoa(data.target.id)
	}
	return fmt.Sprintf("{id: %d, initiative: %d, "+
		"team: %s, units: %d, hp: %d, damage: %d,"+
		"damageType: %s, receiveMultiplier: %v, ep: %d, "+
		"tgt: %s}",
		data.id, data.initiative,
		team, data.units, data.hp, data.damage,
		data.damageType, data.receiveMultiplier,
		data.units*data.damage,
		target,
	)
}

func getDamageMultiplier(g *Group, damageType *string) int {
	if v, ok := g.receiveMultiplier[*damageType]; !ok {
		return 1
	} else {
		return v
	}
}

func main() {
	DEBUG = false

	fileLines := ReadFile("input")
	//fileLines = ReadFile("input_test1")

	groups := parseGroups(fileLines, 0)
	result1 := max(fight(groups))
	fmt.Println("Result1:", result1)

	// Exponential search
	upperBound := 1
	for immuneSurvivors := 0; immuneSurvivors == 0; upperBound = upperBound * 2 {
		groups := parseGroups(fileLines, upperBound)
		IfDebugPrintf("Trying %v\n", upperBound)
		immuneSurvivors, _ = fight(groups)
		IfDebugPrintf("Tried %v: %v\n", upperBound, immuneSurvivors)
	}
	lb := upperBound / 4
	IfDebugPrintln("Bounds found:", lb, "-", upperBound)
	// Binary search (upperBound/4 - upperBound/2)
	r2 := 0
	for lb <= upperBound {
		m := (lb + upperBound) / 2
		boost := m
		groups := parseGroups(fileLines, boost)
		r2, _ = fight(groups)
		IfDebugPrintln("Binary:", lb, upperBound, r2)
		if r2 == 0 {
			lb = m + 1
		} else {
			upperBound = m - 1
		}
	}
	fmt.Println("Result2:", r2)
	// 2364 low
	// 17852 high
}

func fight(groups []*Group) (int, int) {
	for {
		damagePerRound := 0
		// Target selection
		sort.Slice(groups, func(i, j int) bool {
			teamA := groups[i]
			teamB := groups[j]
			epA := teamA.units * teamA.damage
			epB := teamB.units * teamB.damage
			if epA == epB {
				return teamA.initiative > teamB.initiative
			}
			return epA > epB
		})
		IfDebugPrintln("Sorted to choose targets:", groups)

		//decisions := make(map[int]int, len(groups))
		chosen := make(map[int]bool)
		for _, g := range groups {
			if g.units <= 0 {
				continue
			}
			targets := make([]*Group, len(groups))
			copy(targets, groups)
			sort.Slice(targets, func(i, j int) bool {
				dealDamage1 := g.units * g.damage * getDamageMultiplier(targets[i], &g.damageType)
				dealDamage2 := g.units * g.damage * getDamageMultiplier(targets[j], &g.damageType)
				ep1 := targets[i].units * targets[i].damage
				ep2 := targets[j].units * targets[j].damage
				if dealDamage1 == dealDamage2 {
					if ep1 == ep2 {
						return targets[i].initiative > targets[j].initiative
					}
					return ep1 > ep2
				}
				return dealDamage1 > dealDamage2
			})
			//IfDebugPrintln("Sorted to take damage:", targets)
			if DEBUG {
				for _, t := range targets {
					if t.allegiance == g.allegiance {
						continue
					}
					fmt.Println("Group", g.id, "would deal defending group", t.id, g.units*g.damage*getDamageMultiplier(t, &g.damageType), "damage")
				}
			}

			for _, t := range targets {
				if _, alreadyChosen := chosen[t.id]; t.id == g.id || t.allegiance == g.allegiance ||
					t.units <= 0 || alreadyChosen ||
					g.units*g.damage*getDamageMultiplier(t, &g.damageType) == 0 {
					continue
				}
				chosen[t.id] = true
				g.target = t
				break
			}
		}

		// Attack
		sort.Slice(groups, func(i, j int) bool {
			return groups[i].initiative > groups[j].initiative
		})
		IfDebugPrintln("Sorted to attack:", groups)

		for _, g := range groups {
			if g.units <= 0 || g.target == nil {
				continue
			}
			target := g.target
			dealDamage := g.damage * g.units * getDamageMultiplier(target, &g.damageType)
			damagePerRound += dealDamage
			killedUnits := dealDamage / target.hp
			IfDebugPrintln("Group", g.id, "attacks defending group", target.id, "killing", killedUnits, "units")
			target.units = max(0, target.units-killedUnits)
			g.target = nil
		}

		// Endgame condition
		if damagePerRound == 0 { // stale fight
			return 0, 0
		}
		r := map[bool]int{true: 0, false: 0}
		for _, g := range groups {
			r[g.allegiance] = r[g.allegiance] + max(0, g.units)
		}
		if r[true] == 0 || r[false] == 0 {
			return r[true], r[false]
		}
		IfDebugPrintln("Proceeding to next turn")
	}
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func parseGroups(fileLines []string, boost int) []*Group {
	var groups []*Group

	// 3139 units each with 48062 hit points (immune to bludgeoning, slashing, fire; weak to cold)
	//with an attack that does 28 bludgeoning damage at initiative 3
	reLine := regexp.MustCompile(`^(?P<units>\d+) units each with (?P<hp>\d+) hit points (?P<special>\([^)]+\) )?with an attack that does (?P<damage>\d+) (?P<damageType>[a-z]+) damage at initiative (?P<initiative>\d+)$`)
	reImmune := regexp.MustCompile(`immune to ([^;)]+)`)
	reWeak := regexp.MustCompile(`weak to ([^;)]+)`)

	allegiance := true
	id := 1
	for _, line := range fileLines {
		if line == "Immune System:" || line == "Infection:" {
			continue
		} else if line == "" {
			allegiance = false
			continue
		}
		IfDebugPrintln("\n==========\nParsing:", line)

		grep := make(map[string]string)
		match := reLine.FindStringSubmatch(line)
		if match == nil {
			log.Fatal("Couldn't parse: ", line)
		}
		for i, n := range reLine.SubexpNames() {
			grep[n] = match[i]
		}
		group := Group{
			id:                id,
			allegiance:        allegiance,
			damageType:        grep["damageType"],
			receiveMultiplier: make(map[string]int),
		}
		id++
		group.units, _ = strconv.Atoi(grep["units"])
		group.hp, _ = strconv.Atoi(grep["hp"])
		group.damage, _ = strconv.Atoi(grep["damage"])
		if allegiance {
			group.damage += boost
		}
		group.initiative, _ = strconv.Atoi(grep["initiative"])

		if sv, ok := grep["special"]; ok {
			matchImmune := reImmune.FindStringSubmatch(sv)
			if matchImmune != nil {
				IfDebugPrintln("immune ->", matchImmune[1])
				for _, v := range strings.Split(matchImmune[1], ", ") {
					group.receiveMultiplier[v] = 0
				}
			}

			matchWeak := reWeak.FindStringSubmatch(sv)
			if matchWeak != nil {
				IfDebugPrintln("weak ->", matchWeak[1])
				for _, v := range strings.Split(matchWeak[1], ", ") {
					group.receiveMultiplier[v] = 2
				}
			}

			IfDebugPrintln("Parsed group:", group)

			groups = append(groups, &group)
		}
	}
	return groups
}

/*
   --- Day 24: Immune System Simulator 20XX ---

   After a weird buzzing noise, you appear back at the man's cottage. He seems relieved to see his friend, but quickly notices that the little reindeer caught some kind of cold while out exploring.

   The portly man explains that this reindeer's immune system isn't similar to regular reindeer immune systems:

   The immune system and the infection each have an army made up of several groups; each group consists of one or more identical units. The armies repeatedly fight until only one army has units remaining.

   Units within a group all have the same hit points (amount of damage a unit can take before it is destroyed), attack damage (the amount of damage each unit deals), an attack type, an initiative (higher initiative units attack first and win ties), and sometimes weaknesses or immunities. Here is an example group:

   18 units each with 729 hit points (weak to fire; immune to cold, slashing)
    with an attack that does 8 radiation damage at initiative 10

   Each group also has an effective power: the number of units in that group multiplied by their attack damage. The above group has an effective power of 18 * 8 = 144. Groups never have zero or negative units; instead, the group is removed from combat.

   Each fight consists of two phases: target selection and attacking.

   During the target selection phase, each group attempts to choose one target. In decreasing order of effective power, groups choose their targets; in a tie, the group with the higher initiative chooses first. The attacking group chooses to target the group in the enemy army to which it would deal the most damage (after accounting for weaknesses and immunities, but not accounting for whether the defending group has enough units to actually receive all of that damage).

   If an attacking group is considering two defending groups to which it would deal equal damage, it chooses to target the defending group with the largest effective power; if there is still a tie, it chooses the defending group with the highest initiative. If it cannot deal any defending groups damage, it does not choose a target. Defending groups can only be chosen as a target by one attacking group.

   At the end of the target selection phase, each group has selected zero or one groups to attack, and each group is being attacked by zero or one groups.

   During the attacking phase, each group deals damage to the target it selected, if any. Groups attack in decreasing order of initiative, regardless of whether they are part of the infection or the immune system. (If a group contains no units, it cannot attack.)

   The damage an attacking group deals to a defending group depends on the attacking group's attack type and the defending group's immunities and weaknesses. By default, an attacking group would deal damage equal to its effective power to the defending group. However, if the defending group is immune to the attacking group's attack type, the defending group instead takes no damage; if the defending group is weak to the attacking group's attack type, the defending group instead takes double damage.

   The defending group only loses whole units from damage; damage is always dealt in such a way that it kills the most units possible, and any remaining damage to a unit that does not immediately kill it is ignored. For example, if a defending group contains 10 units with 10 hit points each and receives 75 damage, it loses exactly 7 units and is left with 3 units at full health.

   After the fight is over, if both armies still contain units, a new fight begins; combat only ends once one army has lost all of its units.

   For example, consider the following armies:

   Immune System:
   17 units each with 5390 hit points (weak to radiation, bludgeoning) with
    an attack that does 4507 fire damage at initiative 2
   989 units each with 1274 hit points (immune to fire; weak to bludgeoning,
    slashing) with an attack that does 25 slashing damage at initiative 3

   Infection:
   801 units each with 4706 hit points (weak to radiation) with an attack
    that does 116 bludgeoning damage at initiative 1
   4485 units each with 2961 hit points (immune to radiation; weak to fire,
    cold) with an attack that does 12 slashing damage at initiative 4

   If these armies were to enter combat, the following fights, including details during the target selection and attacking phases, would take place:

   Immune System:
   Group 1 contains 17 units
   Group 2 contains 989 units
   Infection:
   Group 1 contains 801 units
   Group 2 contains 4485 units

   Infection group 1 would deal defending group 1 185832 damage
   Infection group 1 would deal defending group 2 185832 damage
   Infection group 2 would deal defending group 2 107640 damage
   Immune System group 1 would deal defending group 1 76619 damage
   Immune System group 1 would deal defending group 2 153238 damage
   Immune System group 2 would deal defending group 1 24725 damage

   Infection group 2 attacks defending group 2, killing 84 units
   Immune System group 2 attacks defending group 1, killing 4 units
   Immune System group 1 attacks defending group 2, killing 51 units
   Infection group 1 attacks defending group 1, killing 17 units

   Immune System:
   Group 2 contains 905 units
   Infection:
   Group 1 contains 797 units
   Group 2 contains 4434 units

   Infection group 1 would deal defending group 2 184904 damage
   Immune System group 2 would deal defending group 1 22625 damage
   Immune System group 2 would deal defending group 2 22625 damage

   Immune System group 2 attacks defending group 1, killing 4 units
   Infection group 1 attacks defending group 2, killing 144 units

   Immune System:
   Group 2 contains 761 units
   Infection:
   Group 1 contains 793 units
   Group 2 contains 4434 units

   Infection group 1 would deal defending group 2 183976 damage
   Immune System group 2 would deal defending group 1 19025 damage
   Immune System group 2 would deal defending group 2 19025 damage

   Immune System group 2 attacks defending group 1, killing 4 units
   Infection group 1 attacks defending group 2, killing 143 units

   Immune System:
   Group 2 contains 618 units
   Infection:
   Group 1 contains 789 units
   Group 2 contains 4434 units

   Infection group 1 would deal defending group 2 183048 damage
   Immune System group 2 would deal defending group 1 15450 damage
   Immune System group 2 would deal defending group 2 15450 damage

   Immune System group 2 attacks defending group 1, killing 3 units
   Infection group 1 attacks defending group 2, killing 143 units

   Immune System:
   Group 2 contains 475 units
   Infection:
   Group 1 contains 786 units
   Group 2 contains 4434 units

   Infection group 1 would deal defending group 2 182352 damage
   Immune System group 2 would deal defending group 1 11875 damage
   Immune System group 2 would deal defending group 2 11875 damage

   Immune System group 2 attacks defending group 1, killing 2 units
   Infection group 1 attacks defending group 2, killing 142 units

   Immune System:
   Group 2 contains 333 units
   Infection:
   Group 1 contains 784 units
   Group 2 contains 4434 units

   Infection group 1 would deal defending group 2 181888 damage
   Immune System group 2 would deal defending group 1 8325 damage
   Immune System group 2 would deal defending group 2 8325 damage

   Immune System group 2 attacks defending group 1, killing 1 unit
   Infection group 1 attacks defending group 2, killing 142 units

   Immune System:
   Group 2 contains 191 units
   Infection:
   Group 1 contains 783 units
   Group 2 contains 4434 units

   Infection group 1 would deal defending group 2 181656 damage
   Immune System group 2 would deal defending group 1 4775 damage
   Immune System group 2 would deal defending group 2 4775 damage

   Immune System group 2 attacks defending group 1, killing 1 unit
   Infection group 1 attacks defending group 2, killing 142 units

   Immune System:
   Group 2 contains 49 units
   Infection:
   Group 1 contains 782 units
   Group 2 contains 4434 units

   Infection group 1 would deal defending group 2 181424 damage
   Immune System group 2 would deal defending group 1 1225 damage
   Immune System group 2 would deal defending group 2 1225 damage

   Immune System group 2 attacks defending group 1, killing 0 units
   Infection group 1 attacks defending group 2, killing 49 units

   Immune System:
   No groups remain.
   Infection:
   Group 1 contains 782 units
   Group 2 contains 4434 units

   In the example above, the winning army ends up with 782 + 4434 = 5216 units.

   You scan the reindeer's condition (your puzzle input); the white-bearded man looks nervous. As it stands now, how many units would the winning army have?

--- Part Two ---

Things aren't looking good for the reindeer. The man asks whether more milk and cookies would help you think.

If only you could give the reindeer's immune system a boost, you might be able to change the outcome of the combat.

A boost is an integer increase in immune system units' attack damage. For example, if you were to boost the above example's immune system's units by 1570, the armies would instead look like this:

Immune System:
17 units each with 5390 hit points (weak to radiation, bludgeoning) with
 an attack that does 6077 fire damage at initiative 2
989 units each with 1274 hit points (immune to fire; weak to bludgeoning,
 slashing) with an attack that does 1595 slashing damage at initiative 3

Infection:
801 units each with 4706 hit points (weak to radiation) with an attack
 that does 116 bludgeoning damage at initiative 1
4485 units each with 2961 hit points (immune to radiation; weak to fire,
 cold) with an attack that does 12 slashing damage at initiative 4

With this boost, the combat proceeds differently:

Immune System:
Group 2 contains 989 units
Group 1 contains 17 units
Infection:
Group 1 contains 801 units
Group 2 contains 4485 units

Infection group 1 would deal defending group 2 185832 damage
Infection group 1 would deal defending group 1 185832 damage
Infection group 2 would deal defending group 1 53820 damage
Immune System group 2 would deal defending group 1 1577455 damage
Immune System group 2 would deal defending group 2 1577455 damage
Immune System group 1 would deal defending group 2 206618 damage

Infection group 2 attacks defending group 1, killing 9 units
Immune System group 2 attacks defending group 1, killing 335 units
Immune System group 1 attacks defending group 2, killing 32 units
Infection group 1 attacks defending group 2, killing 84 units

Immune System:
Group 2 contains 905 units
Group 1 contains 8 units
Infection:
Group 1 contains 466 units
Group 2 contains 4453 units

Infection group 1 would deal defending group 2 108112 damage
Infection group 1 would deal defending group 1 108112 damage
Infection group 2 would deal defending group 1 53436 damage
Immune System group 2 would deal defending group 1 1443475 damage
Immune System group 2 would deal defending group 2 1443475 damage
Immune System group 1 would deal defending group 2 97232 damage

Infection group 2 attacks defending group 1, killing 8 units
Immune System group 2 attacks defending group 1, killing 306 units
Infection group 1 attacks defending group 2, killing 29 units

Immune System:
Group 2 contains 876 units
Infection:
Group 2 contains 4453 units
Group 1 contains 160 units

Infection group 2 would deal defending group 2 106872 damage
Immune System group 2 would deal defending group 2 1397220 damage
Immune System group 2 would deal defending group 1 1397220 damage

Infection group 2 attacks defending group 2, killing 83 units
Immune System group 2 attacks defending group 2, killing 427 units

After a few fights...

Immune System:
Group 2 contains 64 units
Infection:
Group 2 contains 214 units
Group 1 contains 19 units

Infection group 2 would deal defending group 2 5136 damage
Immune System group 2 would deal defending group 2 102080 damage
Immune System group 2 would deal defending group 1 102080 damage

Infection group 2 attacks defending group 2, killing 4 units
Immune System group 2 attacks defending group 2, killing 32 units

Immune System:
Group 2 contains 60 units
Infection:
Group 1 contains 19 units
Group 2 contains 182 units

Infection group 1 would deal defending group 2 4408 damage
Immune System group 2 would deal defending group 1 95700 damage
Immune System group 2 would deal defending group 2 95700 damage

Immune System group 2 attacks defending group 1, killing 19 units

Immune System:
Group 2 contains 60 units
Infection:
Group 2 contains 182 units

Infection group 2 would deal defending group 2 4368 damage
Immune System group 2 would deal defending group 2 95700 damage

Infection group 2 attacks defending group 2, killing 3 units
Immune System group 2 attacks defending group 2, killing 30 units

After a few more fights...

Immune System:
Group 2 contains 51 units
Infection:
Group 2 contains 40 units

Infection group 2 would deal defending group 2 960 damage
Immune System group 2 would deal defending group 2 81345 damage

Infection group 2 attacks defending group 2, killing 0 units
Immune System group 2 attacks defending group 2, killing 27 units

Immune System:
Group 2 contains 51 units
Infection:
Group 2 contains 13 units

Infection group 2 would deal defending group 2 312 damage
Immune System group 2 would deal defending group 2 81345 damage

Infection group 2 attacks defending group 2, killing 0 units
Immune System group 2 attacks defending group 2, killing 13 units

Immune System:
Group 2 contains 51 units
Infection:
No groups remain.

This boost would allow the immune system's armies to win! It would be left with 51 units.

You don't even know how you could boost the reindeer's immune system or what effect it might have, so you need to be cautious and find the smallest boost that would allow the immune system to win.

How many units does the immune system have left after getting the smallest boost it needs to win?
*/
