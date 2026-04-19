package system

import (
	"testing"

	initservice "github.com/flipped-aurora/gin-vue-admin/server/service/system"
)

func TestInitOrderMenuAuthorityRunsAfterLabSeeds(t *testing.T) {
	// Lab seeds currently occupy InitOrderSystem + 26~28.
	// Menu-authority binding must run later so role-menu mapping includes lab menus on first boot.
	if initOrderMenuAuthority <= initservice.InitOrderSystem+28 {
		t.Fatalf(
			"initOrderMenuAuthority=%d must be greater than %d",
			initOrderMenuAuthority,
			initservice.InitOrderSystem+28,
		)
	}
}
