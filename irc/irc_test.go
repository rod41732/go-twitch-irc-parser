package irc_test

import (
	"os"
	"strings"
	"testing"

	"github.com/rod41732/go-twitch-irc-parser/irc"
)

func TestSimpleIRC(t *testing.T) {
	msg := "@foo=bar :user!user@user.tmi.twitch.tv PRIVMSG #pajlada :this is a test"
	parsed := irc.NewIRCMessage(msg)

	expRawTags := "foo=bar"
	if string(parsed.RawTags) != expRawTags {
		t.Errorf("RawTags: expected [%s], got: [%s]", expRawTags, parsed.RawTags)
	}

	expPrefix := "user!user@user.tmi.twitch.tv"
	if string(parsed.Prefix) != expPrefix {
		t.Errorf("Prefix: expected [%s], got: [%s]", expPrefix, parsed.Prefix)
	}

	expCommand := "PRIVMSG"
	if string(parsed.Command) != expCommand {
		t.Errorf("Command: expected [%s], got: [%s]", expCommand, parsed.Command)
	}

	expParams := "#pajlada :this is a test"
	if string(parsed.Params) != expParams {
		t.Errorf("Params: expected [%s], got: [%s]", expParams, parsed.Params)
	}
}

func TestEscape(t *testing.T) {
	msg := "@foo=bar\\sbaz :user!user@user.tmi.twitch.tv PRIVMSG #pajlada :this is a test"
	parsed := irc.NewIRCMessage(msg)
	if parsed.Tag[0].Key != "foo" {
		t.Errorf("Tag key: expected [foo], got: [%v]", parsed.Tag[0])
	}
	if parsed.Tag[0].Value != "bar baz" {
		t.Errorf("Tag value: expected [bar baz], got: [%v]", parsed.Tag[0])
	}

}

func TestDankIRC(t *testing.T) {
	msg := "@foo=bar :user!user@user.tmi.twitch.tv PRIVMSG #pajlada :this is a test"
	parsed := irc.NewIRCMessage(msg)
	if len(parsed.Tag) != 1 {
		t.Error("Expected 1, got", len(parsed.Tag))
	}
	t.Logf("%v", parsed.Tag)

	msg = "@foo=bar;trail :user!user@user.tmi.twitch.tv PRIVMSG #pajlada :this is a test"
	parsed = irc.NewIRCMessage(msg)
	if len(parsed.Tag) != 2 {
		t.Error("Expected 2, got", len(parsed.Tag))
	}
	t.Logf("%v", parsed.Tag)

	msg = "@foo=bar;trail_eq= :user!user@user.tmi.twitch.tv PRIVMSG #pajlada :this is a test"
	parsed = irc.NewIRCMessage(msg)
	if len(parsed.Tag) != 2 {
		t.Error("Expected 2, got", len(parsed.Tag))
	}
	t.Logf("%v", parsed.Tag)

	msg = "@foo=bar;mis;trail_eq= :user!user@user.tmi.twitch.tv PRIVMSG #pajlada :this is a test"
	parsed = irc.NewIRCMessage(msg)
	if len(parsed.Tag) != 3 {
		t.Error("Expected 3, got", len(parsed.Tag))
	}
	t.Logf("%v", parsed.Tag)

	msg = "@foo=bar;empty=;trail_eq= :user!user@user.tmi.twitch.tv PRIVMSG #pajlada :this is a test"
	parsed = irc.NewIRCMessage(msg)
	if len(parsed.Tag) != 3 {
		t.Error("Expected 3, got", len(parsed.Tag))
	}
	t.Logf("%v", parsed.Tag)

	msg = "@foo=bar;empty=;trail_eq=;baz=quux :user!user@user.tmi.twitch.tv PRIVMSG #pajlada :this is a test"
	parsed = irc.NewIRCMessage(msg)
	if len(parsed.Tag) != 4 {
		t.Error("Expected 4, got", len(parsed.Tag))
	}
	t.Logf("%v", parsed.Tag)

	msg = "@mis;key=value :user!user@user.tmi.twitch.tv PRIVMSG #pajlada :this is a test"
	parsed = irc.NewIRCMessage(msg)
	if len(parsed.Tag) != 2 {
		t.Error("Expected 2, got", len(parsed.Tag))
	}
	t.Logf("%v", parsed.Tag)

	msg = "@mis;mis2 :user!user@user.tmi.twitch.tv PRIVMSG #pajlada :this is a test"
	parsed = irc.NewIRCMessage(msg)
	if len(parsed.Tag) != 2 {
		t.Error("Expected 2, got", len(parsed.Tag))
	}
	t.Logf("%v", parsed.Tag)

	msg = "@mis=;key=value :user!user@user.tmi.twitch.tv PRIVMSG #pajlada :this is a test"
	parsed = irc.NewIRCMessage(msg)
	if len(parsed.Tag) != 2 {
		t.Error("Expected 2, got", len(parsed.Tag))
	}
	t.Logf("%v", parsed.Tag)
}

func BenchmarkParsingSingleMessage(b *testing.B) {
	for i := 0; i < b.N; i++ {
		m := irc.NewIRCMessage("@badge-info=subscriber/22;badges=subscriber/3012;color=#FFFF00;display-name=FELYP8;emote-only=1;emotes=521050:0-6,8-14,16-22,24-30,32-38,40-46,48-54,56-62,64-70,72-78,80-86,88-94,96-102,104-110,148-154,156-162,164-170,172-178,180-186,188-194,196-202,204-210,212-218,220-226,228-234,236-242,244-250,252-258,260-266/302827730:112-119/302827734:121-128/302827735:130-137/302827737:139-146;first-msg=0;flags=;id=1844235a-c24e-4e18-937b-805d6601aebe;mod=0;returning-chatter=0;room-id=22484632;subscriber=1;tmi-sent-ts=1685664001040;turbo=0;user-id=162760707;user-type= :felyp8!felyp8@felyp8.tmi.twitch.tv PRIVMSG #forsen :forsenE forsenE forsenE forsenE forsenE forsenE forsenE forsenE forsenE forsenE forsenE forsenE forsenE forsenE forsenE1 forsenE2 forsenE3 forsenE4 forsenE forsenE forsenE forsenE forsenE forsenE forsenE forsenE forsenE forsenE forsenE forsenE forsenE forsenE forsenE")
		_ = m
		// m.ParseTags()
	}
}

// around 1.47..153
func BenchmarkParsing1000Messages(b *testing.B) {
	f, err := os.ReadFile("../data.txt")
	if err != nil {
		b.Fatalf("Read file failed %s", err)
	}

	lines := strings.Split(string(f), "\n")
	if lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	if len(lines) != 1000 {
		b.Fatalf("Not 1000 lines")
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, line := range lines {
			m := irc.NewIRCMessage(line)
			_ = m
			// m.ParseTags()
		}
	}
}
