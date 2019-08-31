package mass

import (
	"github.com/niwho/mass/member_manager"
	"testing"
)

func TestNewMemberManager(t *testing.T) {
	member_manager.NewMemberManager("nodename11", "test_mass_unit", "127.0.0.1", 3111, map[string]string{
		"c_host": "127.0.0.1",
		"c_port": "9111",
	}, "")
}