package mass

import (
	"fmt"
	"github.com/niwho/mass/member_manager"
	"testing"
	"time"
)

func TestNewMemberManager(t *testing.T) {
	cl, _ := member_manager.NewMemberManager("nodename11", "mem_dist", "127.0.0.1", 3111, map[string]string{
		"c_host": "127.0.0.1",
		"c_port": "9111",
	}, "")
	time.Sleep(5 * time.Second)
	fmt.Println(cl.GetMembers())
}
