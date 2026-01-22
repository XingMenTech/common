package redis

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func init() {
	err := InitRedisCache(&Config{
		Prefix:   "aaaa",
		Host:     "127.0.0.1:6379",
		Password: "",
		DbNum:    strconv.Itoa(0),
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

	fmt.Println("HSet。。。。。。。。。。。。。。。")
	if err := HSet("test_users", "user_11", &args{
		Name:  "name_11",
		Age:   11,
		Phone: "phone_11",
	}); err != nil {
		t.Error("HSet error:", err)
		return
	}
	fmt.Println("HGet。。。。。。。。。。。。。。。")

	if arg, err := HGet[args]("test_users", "user_11"); err != nil {
		t.Error("HGet error:", err)
		return
	} else {
		t.Log(arg)
	}

	fmt.Println("HMSet。。。。。。。。。。。。。。。")
	data := make(map[string]interface{})
	for i := 0; i < 10; i++ {
		field := fmt.Sprintf("user_%d", i)

		data[field] = &args{
			Name:  fmt.Sprintf("name_%d", i),
			Age:   20,
			Phone: fmt.Sprintf("phone_%d", i),
		}
	}

	if err := HMSet("test_users", data); err != nil {
		t.Error("HSetAll error:", err)
		return
	}

	t.Log("HLen", HLen("test_users"))

	fmt.Println("HKeys。。。。。。。。。。。。。。。")
	if keys, err := HKeys("test_users"); err != nil {
		t.Error("HKeys error:", err)
		return
	} else {
		t.Log(keys)
	}

	fmt.Println("HMGet。。。。。。。。。。。。。。。")
	hmGet := HMGet[args]("test_users", "user_1", "user_2", "user_3")
	for k, v := range hmGet {
		fmt.Println(k, v)
	}

	fmt.Println("HGetAll。。。。。。。。。。。。。。。")
	err, m := HGetAll[args]("test_users")
	if err != nil {
		t.Error("HGetAll error", err)
		return
	}
	for k, v := range m {
		fmt.Println(k, v)
	}

	fmt.Println("HVals。。。。。。。。。。。。。。。。。")
	vals, err := HVals[args]("test_users")
	if err != nil {
		t.Error("HVals error", err)
		return
	}
	for _, v := range vals {
		fmt.Println(v)
	}

}

func TestHSet(t *testing.T) {
	interval := time.Tick(time.Second * 30)
	for range interval {

		fmt.Println(time.Now().Unix())
	}

}

func TestList(t *testing.T) {

	fmt.Println("LPush。。。。。。。。。。。。。。。")
	if err := LPush(listKey, "1", "2", "3", "4", "5", "6", "7", "8", "9", "10"); err != nil {
		t.Error(err)
		return
	}

	fmt.Println("获取列表长度")
	if lenth, err := LLen(listKey); err != nil {
		t.Error(err)
		return
	} else {
		t.Log("列表长度：", lenth)
	}

	fmt.Println("通过索引获取列表中的元素")
	if v, err := LIndex[string](listKey, 2); err != nil {
		t.Error(err)
		return
	} else {
		t.Log(v)
	}

	fmt.Println("LPop。。。。。。。。。。。。。。。")
	if v, err := LPop[string](listKey); err != nil {
		t.Error(err)
		return
	} else {
		t.Log(v)
	}

	// Lindex 4	LINDEX key index 通过索引获取列表中的元素

	// Linsert 5	LINSERT key BEFORE|AFTER pivot value 在列表的元素前或者后插入元素
	fmt.Println("Linsert。。。。。。。。。。。。。。。")
	if err := LInsert(listKey, "BEFORE", 2, 11); err != nil {
		t.Error(err)
		return
	}

	// LPush 8	LPUSH key value1 [value2] 将一个或多个值插入到列表头部
	fmt.Println("LPush。。。。。。。。。。。。。。。")
	if err := LPush(listKey, 20, 21, 22, 23, 24, 25); err != nil {
		t.Error(err)
		return
	}

	// LPushX 9	LPUSHX key value 将一个值插入到已存在的列表头部
	fmt.Println("LPushX。。。。。。。。。。。。。。。")
	if err := LPushX(listKey, 26, 27, 28, 29, 30); err != nil {
		t.Error(err)
		return
	}

	// Lrange 10	LRANGE key start stop 获取列表指定范围内的元素
	fmt.Println("Lrange。。。。。。。。。。。。。。。")
	if v, err := LRange[string](listKey, 0, 10); err != nil {
		t.Error(err)
		return
	} else {
		t.Log(v)
	}

	// LRem 11	LREM key count value 移除列表元素
	fmt.Println("LRem。。。。。。。。。。。。。。。")
	if v, err := LRem(listKey, 2, 27); err != nil {
		t.Error(err)
		return
	} else {
		t.Log(v)
	}

	// Lset 12	LSET key index value 通过索引设置列表元素的值
	fmt.Println("Lset。。。。。。。。。。。。。。。")
	if err := LSet(listKey, 5, 40); err != nil {
		t.Error(err)
		return
	}

	// Ltrim 13	LTRIM key start stop 对一个列表进行修剪(trim)，就是说，让列表只保留指定区间内的元素，不在指定区间之内的元素都将被删除。
	fmt.Println("Ltrim。。。。。。。。。。。。。。。")
	if err := Ltrim(listKey, 5, 10); err != nil {
		t.Error(err)
		return
	}

	// RPop 14	RPOP key 移除列表的最后一个元素，返回值为移除的元素。
	fmt.Println("RPop。。。。。。。。。。。。。。。")
	if v, err := RPop[string](listKey); err != nil {
		t.Error(err)
		return
	} else {
		t.Log(v)
	}

	// RPopLPush 15	RPOPLPUSH source destination 移除列表的最后一个元素，并将该元素添加到另一个列表并返回
	fmt.Println("RPopLPush。。。。。。。。。。。。。。。")
	if v, err := RPopLPush[string](listKey, listKey2); err != nil {
		t.Error(err)
		return
	} else {
		t.Log(v)
	}

	// Rpush 16	RPUSH key value1 [value2] 在列表中添加一个或多个值到列表尾部
	fmt.Println("Rpush。。。。。。。。。。。。。。。")
	if err := RPush(listKey, 1, 2, 3, 4, 5); err != nil {
		t.Error(err)
		return
	}

	// RPushX 17 RPUSHX key value 为已存在的列表添加值
	fmt.Println("RPushX。。。。。。。。。。。。。。。")
	if err := RPushX(listKey, 6, 7, 8, 9); err != nil {
		t.Error(err)
		return
	}

	fmt.Println("获取列表长度")
	if lenth, err := LLen(listKey); err != nil {
		t.Error(err)
		return
	} else {
		res, err := LRange[string](listKey, 0, lenth)
		t.Log("列表长度：", lenth, "列表内容：", res)
		if err != nil {
			t.Error(err)
			return
		}
		t.Log("列表内容：", res)
	}
}

const (
	listKey  = "test_list"
	listKey2 = "test_list2"
)

func TestBLPop(t *testing.T) {
	for {
		k, v, err := BLPop[args](0, listKey, listKey2)
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(k, "====>", v)
	}
}

func TestRPush(t *testing.T) {

	index := 0
	for {

		val := &args{
			Name:  fmt.Sprintf("name_%d", index),
			Age:   20,
			Phone: fmt.Sprintf("phone_%d", index),
		}

		err := RPush(listKey, val)
		if err != nil {
			t.Error(err)
			return
		}
		time.Sleep(10 * time.Second)
		index++
	}
}

func TestRPush2(t *testing.T) {

	index := 0
	for {

		val := &args{
			Name:  fmt.Sprintf("k2_name_%d", index),
			Age:   20,
			Phone: fmt.Sprintf("k2_phone_%d", index),
		}

		err := RPush(listKey2, val)
		if err != nil {
			t.Error(err)
			return
		}
		time.Sleep(5 * time.Second)
		index++
	}
}
