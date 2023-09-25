package cnf

import (
	"net"

	"github.com/superstes/calamary/proc/meta"
)

type RuleRaw struct {
	Match  MatchRaw `yaml:"match"`
	Action string   `yaml:"action"`
}

type MatchRaw struct {
	SrcNet   YamlStringArray `yaml:"src"`
	DestNet  YamlStringArray `yaml:"dest"`
	SrcPort  YamlStringArray `yaml:"sport"`
	DestPort YamlStringArray `yaml:"port"`
	ProtoL3  YamlStringArray `yaml:"protoL3"`
	ProtoL4  YamlStringArray `yaml:"protoL4"`
	ProtoL5  YamlStringArray `yaml:"protoL5"`
	Domains  YamlStringArray `yaml:"dns"`
}

type Rule struct {
	Match  Match
	Action meta.Action
}

type Match struct {
	SrcNet    []*net.IPNet
	SrcNetN   []*net.IPNet
	DestNet   []*net.IPNet
	DestNetN  []*net.IPNet
	SrcPort   []uint16
	SrcPortN  []uint16
	DestPort  []uint16
	DestPortN []uint16
	ProtoL3   []meta.Proto
	ProtoL3N  []meta.Proto
	ProtoL4   []meta.Proto
	ProtoL4N  []meta.Proto
	ProtoL5   []meta.Proto
	ProtoL5N  []meta.Proto
	Domains   []string
}
