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
	}, "your consul")
	cl.GetServiceXX(true)
	time.Sleep(5 * time.Second)

	fmt.Println(cl.GetMembers())
	cl.UpateLocalRoute("k1", member_manager.NewMember("testnode", "127.0.0.1", 11))
	k1member := cl.GetMember("k1")
	fmt.Println("k1member", k1member)

	cl.UpateLocalRoute("k1", member_manager.NewMember("testnode", "127.0.0.1", 11))
	k1member = cl.GetMember("k1")
	fmt.Println("k1member22222", k1member)
	//routerdb := cl.GetRouter()
	//routerdb.View(func(tx *buntdb.Tx) error {

	//	val, err := tx.Get("k1")
	//	fmt.Println(string(val), err)
	//	return nil
	//})

}
