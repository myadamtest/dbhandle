package main

import "github.com/myadamtest/dbhandle/minecache"

type Product struct {
	Code  string `json:"code"`
	Price uint   `toml:"f"`
}

type Article struct {
	Id     int64  `json:"id" orm:"pk"`
	UserId int64  `json:"userId"`
	Title  string `json:"title"`
}

func main() {

	store := minecache.NewStore("user:passw@tcp(172.0.0.1:3306)/db?charset=utf8", "172.0.0.1:6379")
	art := &Article{UserId: 1, Title: "mf"}
	store.Add(art)

	//var f *string
	//f = new(string)
	//*f = "ad"
	//
	//_,k := whereToKey(1,3,4,"f",f)
	//fmt.Println(k)
}

//func whereToKey(where ...interface{}) (bool,string) {
//	str := ""
//	for i:=0;i<len(where);i++ {
//		tempVal := where[i]
//		t := reflect.TypeOf(where[i])
//		if t.Kind() == reflect.Ptr {//指针类型处理和转换
//			t = t.Elem()
//			tempVal = reflect.ValueOf(where[i]).Elem()
//		}
//
//		switch t.Kind() {
//			case reflect.Int,reflect.Int8,reflect.Int16,reflect.Int32,reflect.Int64,reflect.String,reflect.Bool,
//				reflect.Float32,reflect.Float64,reflect.Uint,reflect.Uint8,reflect.Uint16,reflect.Uint32,reflect.Uint64:
//				str += fmt.Sprintf("/%v",tempVal)
//				continue
//			default:
//				return false,""
//		}
//	}
//	return true,str
//}
//
//type MyStruct struct {
//	UserId int
//	Id int
//	Name string `orm:"-"json:"name"`
//	M2 M2
//	Age int64
//	Income *float64
//}
//
//type M2 struct {
//	UserId int
//	Count int
//}
//
//func group(is []MyStruct) {
//	mp := make(map[int]*M2,0)
//	for i:=0;i < len(is);i++ {
//		if mp[is[i].UserId] == nil {
//			mp[is[i].UserId] = &M2{}
//		}
//
//		mp[is[i].UserId].Count = mp[is[i].UserId].Count +1
//	}
//}
