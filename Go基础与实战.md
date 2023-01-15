# Go语言简介

1. 高性能、高并发

2. 语法简单、学习曲线平缓

3. 丰富的标准库

4. 完善的工具链

5. 静态链接

6. 快速编译

7. 跨平台

8. 垃圾回收

# 基础部分

## 变量声明

```Go

	var a = "initial"

	var b, c int = 1, 2

	var d = true

	var e float64

	f := float32(e)

	g := a + "foo"
```
一般为var或者:=的形式进行命名

## 数组与切片的不同与用法
```Go
var a [5]int
b := [5]int{1, 2, 3, 4, 5}
var twoD [2][3]int
// 数组的声明需要知道数组的长度
s := make([]string, 3)
good := []string{"g", "o", "o", "d"}
// 切片的定义用make，切片一般由数据类型、切片长度和指向切片的指针组成
```
## map
```Go
	m := make(map[string]int)
	m["one"] = 1
	m["two"] = 2
	fmt.Println(m)           // map[one:1 two:2]
	fmt.Println(len(m))      // 2
	fmt.Println(m["one"])    // 1
	fmt.Println(m["unknow"]) // 0

	r, ok := m["unknow"]//ok确定"unkown"是否存在
	fmt.Println(r, ok) // 0 false

	delete(m, "one")//删除map中的键值对

	m2 := map[string]int{"one": 1, "two": 2}
	var m3 = map[string]int{"one": 1, "two": 2}
	fmt.Println(m2, m3)
// map的声明由map[type]type组成字典
```
## 错误处理
```Go
	u, err := findUser([]user{{"wang", "1024"}}, "wang")
	if err != nil {
		fmt.Println(err)
		return
	}
   fmt.Println(u,name) 
    if u, err := findUser([]user{{"wang", "1024"}}, "li"); err != nil {
		fmt.Println(err) // not found
		return
	} else {
		fmt.Println(u.name)
	}
}
```
Go的错误处理，一般用errors.New来生成，处理错误时，可以使用if err!=nil来处理或者if ....;err！=nil{}else{}来处理

## 字符串内置方法
```Go
a := "hello"
	fmt.Println(strings.Contains(a, "ll"))                // true
	fmt.Println(strings.Count(a, "l"))                    // 2
	fmt.Println(strings.HasPrefix(a, "he"))               // true
	fmt.Println(strings.HasSuffix(a, "llo"))              // true
	fmt.Println(strings.Index(a, "ll"))                   // 2
	fmt.Println(strings.Join([]string{"he", "llo"}, "-")) // he-llo
	fmt.Println(strings.Repeat(a, 2))                     // hellohello
	fmt.Println(strings.Replace(a, "e", "E", -1))         // hEllo
	fmt.Println(strings.Split("a-b-c", "-"))              // [a b c]
	fmt.Println(strings.ToLower(a))                       // hello
	fmt.Println(strings.ToUpper(a))                       // HELLO
	fmt.Println(len(a))                                   // 5
	b := "你好"
	fmt.Println(len(b)) // 6
```
## 字符串格式化
一般对于fmt.Printf("%v")可以表示所有类型
fmt.Printf("%+v")可以更清楚地展示字段值
fmt.Printf("%#v")还可以展示构造的具体函数
fmt.Printf("%.2f\n", f)可以打印出浮点数的两位

## json实战
json.Marshal得到结构体序列化结果，然后通过string(buf)得到内容否则会得到数字
json.MarshalIndent以缩进的格式输出，能更清晰地展示
json.Unmarshall来反序列化buf

## 数字解析

```Go
f, _ := strconv.ParseFloat("1.234", 64)
	fmt.Println(f) // 1.234

	n, _ := strconv.ParseInt("111", 10, 64)
	fmt.Println(n) // 111

	n, _ = strconv.ParseInt("0x1000", 0, 64)
	fmt.Println(n) // 4096

	n2, _ := strconv.Atoi("123")
	fmt.Println(n2) // 123

	n2, err := strconv.Atoi("AAA")
	fmt.Println(n2, err) // 0 strconv.Atoi: parsing "AAA": invalid syntax
```
`ParseInt()`参数中第一个为字符串，第二个为转化的进制数，0表示自动表示，第三个表示返回的整型位数

# 实战
## 生出随机数
`rand.Seed(time.Now().UnixNano()`
通过时间戳来生出随机数
`reader := bufio.NewReader(os.Stdin)`来生成一个指向终端输入的一个Reader指针。
`input, err := reader.ReadString('\n')`
来读取终端的输入，直到出现换行符。

## 在线字典
希望通过`go run main.go Hello`将Arg进行翻译，返回翻译结果
使用彩云翻译，检查找到dict和Post
![彩云翻译](https://cdn.jsdelivr.net/gh/richardzhangy26/Pic@main/src202301151352623.png)
点击dict右键，选择`Copy of cURL`获得curl command

[curlconvert](https://curlconverter.com/)网站将curl的结果转换为json
请求的序列话实现
```Go
    request := DictRequest{TransType: "en2zh", Source: word}
	buf, err := json.Marshal(request)
	if err != nil {
		log.Fatal(err)
	}
	var data = bytes.NewReader(buf)// 将buf二进制数组生成reader
    req, err := http.NewRequest("POST", "https://api.interpreter.caiyunai.com/v1/dict", data)//Newrequest生成请求
    //后续通过req.Header.Set来设置reuqest的请求头字段
```
响应返回的是一个比较复杂的字典，需要使用在线网站[jsontogo](https://oktools.net/json2go)
```Go
    resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()//defer保证resp正常关闭
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != 200 {
		log.Fatal("bad StatusCode:", resp.StatusCode, "body", string(bodyText))
	}
	var dictResponse DictResponse
	err = json.Unmarshal(bodyText, &dictResponse)//将响应体反序列化到dictresponse的结构体中

```
## socks5简单实现
Socks5是一种网络协议，它主要用于代理网络连接。它支持TCP和UDP协议，并且可以支持认证，使用者可以使用用户名和密码进行认证。它也可以用于在网络上隐藏真实IP地址，从而增加了网络安全性。