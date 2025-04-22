package redis

import (
	"fmt"
	"testing"
)

func init() {
	err := InitRedisCache(&Config{
		Prefix:   "aaaa",
		Host:     "127.0.0.1:6379",
		Password: "",
		DbNum:    0,
	})
	if err != nil {
		panic(err)
	}
}

type args struct {
	Name  string `json:"name,omitempty"`
	Age   int    `json:"age,omitempty"`
	Phone string `json:"phone,omitempty"`
}

func TestHGetAll(t *testing.T) {
	exist2 := IsExist("test_users2")
	fmt.Println("test_users2 IsExist", exist2)

	if err := Delete("test_users"); err != nil {
		t.Error(err)
		return
	}

	exist := IsExist("test_users")
	fmt.Println("test_users IsExist", exist)
	if exist {
		t.Log("HLen", HLen("test_users"))
	}

	data := make(map[string]interface{})
	for i := 0; i < 10; i++ {
		field := fmt.Sprintf("user_%d", i)

		data[field] = &args{
			Name:  fmt.Sprintf("name_%d", i),
			Age:   20,
			Phone: fmt.Sprintf("phone_%d", i),
		}
		//err := HSet("test_users", fmt.Sprintf("user_%d", i), &args{
		//	Name:  fmt.Sprintf("name_%d", i),
		//	Age:   20,
		//	Phone: fmt.Sprintf("phone_%d", i),
		//})
		//if err != nil {
		//	panic(err)
		//}
	}

	err := HMSet("test_users", data)
	if err != nil {
		t.Error("HSetAll error:", err)
		return
	}

	t.Log("HLen", HLen("test_users"))

	err, m := HGetAll("test_users", args{})
	if err != nil {
		t.Error(err)
	}
	for k, v := range m {
		fmt.Println(k, v.(*args))
	}
	t.Log(m)
}
