# Test Drive Dev
### 一.单元测试
单元测试（ 英语：Unit Testing ）又称为模块测试, 是针对程序模块（软件设计的最小单位）来进行正确性检验的测试工作。程序单元是应用的最小可测试部件。在过程化编程中，一个单元就是单个程序、函数、过程等；对于面向对象编程，最小单元就是方法，包括基类（超类）、抽象类、或者派生类（子类）中的方法。
#### 1.golang 单元测试基本使用
Go 语言内置测试框架通过 testing 包以及 go test 命令来提供测试功能。源码中以 _test.go 后缀的文件为测试文件，每个测试函数必须倒入 testing 包，测试函数的名字必须以 Test 开头，可选的后缀名必须以大写字母开头，函数如下：

```
func TestName(t *testing.T) {
    // ...
}
```
其中 t 参数用于报告测试失败和附加的日志信息。完整例子如下：

```
package strings_test

import ( 
    "strings" 
    "testing" 
)

func TestIndex(t *testing.T) { 
    const s, sep, want = "chicken", "ken", 4 
    got := strings.Index(s, sep) 
    if got != want { 
        t.Errorf("Index(%q,%q) = %v; want %v", s, sep, got, want)// 注意原slide中 的 got 和 want 写反了 
    } 
}
```
执行 go test -v 命令，参数 -v 输出详细的执行信息，参数 -run="French|Canal" 对应一个正则表达式，只有测试函数名被它正确匹配的测试函数才会被 go test 测试命令运行：
```
$go test -v strings_test.go 
=== RUN TestIndex 
— PASS: TestIndex (0.00 seconds) 
PASS 
ok      command-line-arguments    0.007s
```
当运行 go test 命令时，go test 会遍历所有的 *_test.go 中符合上述命名规则的函数，然后生成一个临时的 main 包用于调用相应的测试函数，然后构建并运行、报告测试结果，最后清理测试中生成的临时文件。

#### 2.表格驱动测试
表格驱动测试就是定义好测试的输入，输出，然后循环执行被测试的方法，比较输出和期待的输出是否一致，不一致返回错误的描述信息。

```
package word

import "testing"

func TestIsPalindrome(t *testing.T) {
	var tests = []struct {
		input string
		want  bool
	}{
		{"", true},
		{"a", true},
		{"aa", true},
		{"ab", false},
		{"kayak", true},
		{"detartrated", true},
		{"A man, a plan, a canal: Panama", true},
		{"Evil I did dwell; lewd did I live.", true},
		{"Able was I ere I saw Elba", true},
		{"été", true},
		{"Et se resservir, ivresse reste.", true},
		{"palindrome", false}, // non-palindrome
		{"desserts", false},   // semi-palindrome
	}
	for _, test := range tests {
		if got := IsPalindrome(test.input); got != test.want {
			t.Errorf("IsPalindrome(%q) = %v", test.input, got)
		}
	}
}
```
这种表格驱动的测试在 Go 语言中很常见的。我们很容易向表格添加新的测试数据，并且后面的测试逻辑也没有冗余，这样我们可以有更多的精力地完善错误信息。

#### 3.fake & mock
mock 就是在测试过程中，对于某些不容易构造或者不容易获取的对象，用一个虚拟的对象来创建以便测试的测试方法。情景如下：
- 依赖的服务返回不确定的结果，如获取当前时间。
- 依赖的服务返回状态中有的难以重建或复现，比如模拟网络错误。
- 依赖的服务搭建环境代价高，速度慢，需要一定的成本，比如数据库，web 服务
- 依赖的服务行为多变。

使用 mock 可以很好的帮助我们进行单元测试，mock 对象可以被开发人员很灵活的指定传入参数，调用次数，返回值和执行动作，来满足测试的各种情景假设。

每个语言的 mock 实现方式都不一样，golang 的 mock 主要是通过伪对象的方式来实现，通过在代码中使用接口可以更好帮助我们进行实现和测试。

#### 3.1没有使用接口如何 mock
没有使用接口我们通常需要在测试代码中就去实现需要 mock 的对象。假设我们需要测试的函数如下：

这是一个检查配额占比的函数，当超过配额需要发送邮件。

```
func CheckQuota(username string) {
    used := bytesInUse(username)
    const quota = 1000000000 // 1GB
    percent := 100 * used / quota
    if percent < 90 {
        return // OK
    }
    msg := fmt.Sprintf(template, used, percent)
    auth := smtp.PlainAuth("", sender, password, hostname)
    err := smtp.SendMail(hostname+":587", auth, sender,／／发送邮件
        []string{username}, []byte(msg))
    if err != nil {
        log.Printf("smtp.SendMail(%s) failed: %s", username, err)
    }
}
```
我们想测试这个代码，但是我们并不希望发送真实的邮件。因此我们将邮件处理逻辑放到一个私有的 notifyUser 函数中。

```
var notifyUser = func(username, msg string) {
    auth := smtp.PlainAuth("", sender, password, hostname)
    err := smtp.SendMail(hostname+":587", auth, sender,
        []string{username}, []byte(msg))
    if err != nil {
        log.Printf("smtp.SendEmail(%s) failed: %s", username, err)
    }
}

func CheckQuota(username string) {
    used := bytesInUse(username)
    const quota = 1000000000 // 1GB
    percent := 100 * used / quota
    if percent < 90 {
        return // OK
    }
    msg := fmt.Sprintf(template, used, percent)
    notifyUser(username, msg)
}
```
现在我们可以在测试中用伪邮件发送函数替代真实的邮件发送函数。它只是简单记录要通知的用户和邮件的内容。

```
func TestCheckQuotaNotifiesUser(t *testing.T) {
    var notifiedUser, notifiedMsg string
    notifyUser = func(user, msg string) {
        notifiedUser, notifiedMsg = user, msg
    }

    // ...simulate a 980MB-used condition...

    const user = "joe@example.org"
    CheckQuota(user)
    if notifiedUser == "" && notifiedMsg == "" {
        t.Fatalf("notifyUser not called")
    }
    if notifiedUser != user {
        t.Errorf("wrong user (%s) notified, want %s",
            notifiedUser, user)
    }
    const wantSubstring = "98% of your quota"
    if !strings.Contains(notifiedMsg, wantSubstring) {
        t.Errorf("unexpected notification message <<%s>>, "+
            "want substring %q", notifiedMsg, wantSubstring)
    }
}
```
这里有一个问题：当测试函数返回后，CheckQuota 将不能正常工作，因为 notifyUsers 依然使用的是测试函数的伪发送邮件函数（当更新全局对象的时候总会有这种风险）。 我们必须修改测试代码恢复 notifyUsers 原先的状态以便后续其他的测试没有影响，要确保所有的执行路径后都能恢复，包括测试失败或panic异常的情形。在这种情况下，我们建议使用 defer 语句来延后执行处理恢复的代码。

```
func TestCheckQuotaNotifiesUser(t *testing.T) {
    // Save and restore original notifyUser.
    saved := notifyUser
    defer func() { notifyUser = saved }()

    // Install the test's fake notifyUser.
    var notifiedUser, notifiedMsg string
    notifyUser = func(user, msg string) {
        notifiedUser, notifiedMsg = user, msg
    }
    // ...rest of test...
}
```
这种处理模式可以用来暂时保存和恢复所有的全局变量，包括命令行标志参数、调试选项和优化参数；安装和移除导致生产代码产生一些调试信息的钩子函数；还有有些诱导生产代码进入某些重要状态的改变，比如超时、错误，甚至是一些刻意制造的并发行为等因素。

以这种方式使用全局变量是安全的，因为 go test 命令并不会同时并发地执行多个测试。

#### 3.2 使用接口的对象如何 mock

在使用接口的时候，我们去 mock 一个对象只需要去实现这个接口的所有方法,我们可以自己去实现接口中的所有方法，用 mock 对象去替换掉真实对象。

如果每个需要 mock 的对象都要自己一个一个去写，工作量比较大。

常用的 mock 框架主要有：
- [testify](https://github.com/stretchr/testify)
- [gomock](https://github.com/golang/mock)

我们使用 testify 框架，原因：

1. 提供了 assertion，request 等功能。
    - assertion：测试的结果和预期的值做比较的一个工具
    - request： 功能和 assertion 类似，只是如果第一个期待的值错误就直接退出函数。
2. mock 出来的对象代码结构更清晰，使用简单。
3. fock 数更高 testify 2000+，gomock 400+

下面提供一个简单的例子来描述 testify 如何去使用：

```
package calculator

type Random interface {
	Random(limit int) int
}

type Calculator interface {
	Add(x, y int) int
	Subtract(x, y int) int
	Multiply(x, y int) int
	Divide(x, y int) int
	Random() int
}

func newCalculator(rnd Random) Calculator {
	return calc{
		rnd: rnd,
	}
}

type calc struct {
	rnd Random
}

func (c calc) Add(x, y int) int {
	return x + y
}

func (c calc) Subtract(x, y int) int {
	return x - y
}

func (c calc) Multiply(x, y int) int {
	return x * y
}

func (c calc) Divide(x, y int) int {
	return x / y
}

func (c calc) Random() int {
	return c.rnd.Random(100)
}
```
上面是一个简单的计算器功能的程序，加减乘除的功能都已经实现了，但是随机函数的功能我们并没有去实现，这里是我们需要去 mock 的地方：

```
package calculator

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type randomMock struct {
	mock.Mock
}

func (o randomMock) Random(limit int) int {
  args := o.Called(limit)
  return args.Int(0)
}

func TestAdd(t *testing.T) {
	calc := newCalculator(nil)
	assert.Equal(t, 9, calc.Add(5, 4))
}

func TestSubtract(t *testing.T) {
	calc := newCalculator(nil)
	assert.Equal(t, 1, calc.Subtract(5, 4))
}

func TestMultiply(t *testing.T) {
	calc := newCalculator(nil)
	assert.Equal(t, 20, calc.Multiply(5, 4))
}

func TestDivide(t *testing.T) {
	calc := newCalculator(nil)
	assert.Equal(t, 5, calc.Divide(20, 4))
}

func TestRandom(t *testing.T) {
	rnd := new(randomMock)
	rnd.On("Random", 100).Return(7)
	calc := newCalculator(rnd)
	assert.Equal(t, 7, calc.Random())
}
```
上述的代码中我们使用了 assertion 去做基本的加减乘除功能的测试，然后 mock 随机函数，输入为100，返回7，在实际代码的开发中我们可以自由的去设置这些值，来达到我们测试的需求。

使用接口的好处是，我们的 mock 对象和设置输入输出是可以完全分开的。比如上述的代码我们可以把 mock 对象的代码放在一个单独的 mock 文件夹中，来给其他的人来使用，而使用的人可以自由的去设置输入输出。上述代码 mock 对象可以完全移出：

```
var Mock Random = new(Mock)

type Mock struct {
	mock.Mock
}

func (o Mock) Random(limit int) int {
  args := o.Called(limit)／／调用参数
  return args.Int(0)／／args就是返回值
}
```
我们可以看到上面的代码其实都是相似的，我们可以使用[tool](https://github.com/vektra/mockery)来自动生成。测试的人只需要调用如下方法
```
Mock.On("接口名",输入).ruturn(返回值)
```
4.测试覆盖率

单元测试可以很好的去测试函数的功能，控制产品质量，也可以帮助我们去对某个代码块进行阅读，而测试覆盖率则是对测试代码的一种检查，golang 的测试框架提供的 cover 的功能，使用也很简单。主要 command 如下：

```
go test -cover  pkgs                             //可以指定多个包名，结果输出在终端
go test -coverprofile=a.out pkg                  //只能制定一个包名，结果保存在文件中
go test -coverprofile=a.out -covermode=count pkg //除了覆盖率的结果，还会得到代码块执行次数
go tool cover -func=all.out                      //解析保存的文件以函数形式展示
go tool cover -html=all.out                      //解析保存的文件以html展示，会直接打开你的默认浏览器
go tool cover -html=all.out -o coverage.html     //解析保存的文件以html展示，持久化成一个文件
```
####
由于 go test -coverprofile=a.out pkg 只能指定一个包名， golang 的测试框架代码覆盖率是通过 package 计算测试覆盖率，所以下面写了一个 Makefile 脚本遍历 package，然后聚合文件，最后得到整个项目的覆盖率。
Makefile:

```
SHELL:=/bin/bash

test:
	go test -cover ./...

#采集测试覆盖率，是否覆盖
collect_bool:
	rm all.out coverage.out
	for pkg in $$(go list ./...);do\
				go test -coverprofile=coverage.out $${pkg} || exit $$?;\
				if [ -f coverage.out ] ; then \
				sed -i '1d' coverage.out ;\
				cat coverage.out >> all.out ;\
				fi ; \
	done;\
	sed -i '1 i\mode: set' all.out 

#采集测试覆盖率，显示次数，可以得到代码块的调用次数
collect_count:
	rm all.out coverage.out
	for pkg in $$(go list ./...);do\
				go test -coverprofile=coverage.out -covermode=count $${pkg} || exit $$?;\
				if [ -f coverage.out ] ; then \
				sed -i '1d' coverage.out ;\
				cat coverage.out >> all.out ;\
				fi ; \
	done;\
	sed -i '1 i\mode: count' all.out

#将数据用 html 显示，持久化成文件形式 
html:
	go tool cover -html=all.out -o coverage.html

#将数据用 html 显示，持久化成文件形式 
html_open:
	go tool cover -html=all.out	

#函数形式显示
func:
	go tool cover -func=all.out

.PHONY: collect
```
执行make test:

```
go test -cover ./...
ok  	TestDriveDev/mock/calculator	0.016s	coverage: 100.0% of statements
ok  	TestDriveDev/mock/example	0.012s	coverage: 66.7% of statements
ok  	TestDriveDev/tableTriveTest/word	0.007s	coverage: 100.0% of statements
```
执行make collect_bool && make func:

```
TestDriveDev/mock/calculator/calculator.go:15:	newCalculator				100.0%
TestDriveDev/mock/calculator/calculator.go:25:	Add					100.0%
TestDriveDev/mock/calculator/calculator.go:29:	Subtract				100.0%
TestDriveDev/mock/calculator/calculator.go:33:	Multiply				100.0%
TestDriveDev/mock/calculator/calculator.go:37:	Divide					100.0%
TestDriveDev/mock/calculator/calculator.go:41:	Random					100.0%
TestDriveDev/mock/example/example.go:20:	DoSomething				0.0%
TestDriveDev/mock/example/example.go:27:	DoSomething2				0.0%
TestDriveDev/mock/example/example.go:31:	targetFuncThatDoesSomethingWithObj	100.0%
TestDriveDev/mock/example/example.go:36:	targetFuncThatDoesSomethingWithObj2	100.0%
TestDriveDev/mock/example/example_mock.go:11:	DoSomething				100.0%
TestDriveDev/mock/example/example_mock.go:17:	DoSomething2				100.0%
TestDriveDev/tableTriveTest/word/word.go:8:	IsPalindrome				100.0%
total:						(statements)				84.6%
```

### 二：e2e 测试
E2E ,也就是 End To End,就是所谓的“用户真实场景”。

一种测试分类的方法是基于测试者是否需要了解被测试对象的内部工作原理。黑盒测试只需要测试包公开的文档和 API 行为，内部实现对测试代码是透明的。所以 e2e 测试也是黑盒测试的一种。

e2e 测试其实不需要什么特别的框架，只要能够去调用 http api，由于我们的接口都是 golang 写的，为了更简单的去解析返回值（很好可以用现成的），我们的框架会使用下面两种：

- [ginkgo](https://github.com/onsi/ginkgo)：原来项目中已经使用
- [fortest](https://github.com/emicklei/forest)：封装好了 http 请求和解析返回值，go-restful 的作者。

e2e 测试的例子：template_test.go

```
var _ = Describe("apptempalte", func() {
	var (
		ns1 string
		ns2 string
	)

	BeforeSuite(func() {
		// Wait kubernetes admin to start.
		util.WaitComponents()
		// Create partition
		ns1 = fmt.Sprintf("template-test-ns-%s", 
		。。。
	})

	AfterSuite(func() {
		_, err := client.DeletePartitionWithUidAndCid(ns1)
	    。。。
	})

	It("should be a available.", func() {
		Expect(util.IsAvailable()).To(Equal(true))
	})

	It("should be able to create, get and delete an app template.", func() {
		// Create an application
		app1 := util.CreateApplicationV0_2(ns1, "template-test-1")
		。。。

		// Create a template
		templateName := "template-test"
		template1Req := util.CreateTemplate(templateName, ns1, "template-test-1")
        。。。

		// List all templates
		resp, err = clientv2.ListAllAppTemplate()
		。。。

		// Get a template
		resp, err = clientv2.GetAppTemplate(templateName)
		。。。

		/*
			template := v0_2.AppTemplate{}
			err = jsonutil.GetJsonPayloadFromResponse(resp, &template)
			fmt.Println(template)
			Expect(err).To(BeNil())

			// Update a template
			template.Description = "description-after-updated"
			resp, err = clientv2.UpdateAppTemplate(&template)
		
			// Get tempalte again
			resp, err = clientv2.GetAppTemplate(templateName)
        。。。
			// Verify Update succeeded
			Expect(template2.Description).To(Equal("description-after-updated"))
		*/
	})
})
```
上面的例子就是模拟用户对 template 进行操作：


```
创建分区-->创建应用-->创建模版-->显示模版列表-->显示某一个模版-->删除模版-->删除分区
```
真实的开发中，e2e 测试代码较长，这里取了最短的一个。

最后执行一个 e2e 脚本，遍历执行所有的测试：

```
# Run tests.
testcase=(
proxy
partition
application
job
template
tag
volume
solution
networkpolicy
)

for case in ${testcase[@]}
do
    K8S_HOST="${K8S_HOST}" \
    K8S_TOKEN="${K8S_TOKEN}" \
    K8S_UID="" \
    K8S_CID="" \
    go test -v ${K8S_ADMIN_ROOT}/tests/${case}
    if [ $? -ne 0 ]
    then
        exit -1
    fi
done
```





























