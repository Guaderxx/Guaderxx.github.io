---
title: rust summary
date: 2024-03-14 17:29:24
tags:
categories:
keywords:
copyright:
copyright_author_href:
copyright_info:
---

## Usage

> Rust 的包叫 `crate`

`cargo new [project-name]` 创建一个项目
`cargo build`  编译
`cargo run` 编译并运行
`cargo check` 不编译，检查项目是否有问题
`Cargo.toml`  在这里增加依赖
`fn main() {}` 主函数
`use [pkg]::func`  引入依赖
`println!` `!` 指示这是个宏 macro ？
`let [mut] variable` 声明变量，无 `mut` 是常量
  - 说是常量又不准确，有 `const` ，且声明时必须标注类型
  - 可以多次声明为不同类型 **这不意味着变量类型可变**
`let var: u32` 声明变量并指明类型 `u32  uint32`
`loop {}` 循环
`match` 好像是 `switch` ，不确定


## Common Programming Concept


### Data Types

Scalar 标量  一个单独的值
整型，浮点型，布尔，字符

- `[i|u][8|16|32|64|128|size]`
  - `isize` `usize` 是根据架构来的，64位架构就是 64, 32位架构就是 32.
  - `_` 可以用来分隔长数字，例如 `123_456_789` 等于 `123456789`。
  - `98_222`, hex `0xff`, octal `0o77`, binary `0b111011`, byte u8 `b'A'`
  - 默认使用的是 `i32`  
- `f[32|64]`
  - 默认是 `f64`
- `true | false`
- `char`
  - 4 bytes, 表示 unicode 值
  - 不是广义的字符，类似 `rune`

Compound 复合，将多个值合为一个类型
元组 tuple, 数组 array

- tuple
  - 长度固定
  -  ```
    // declare a tuple
    let tup: (i32, f64, u8) = (500, 6.4, 1);
    // destructure
    let (x, y, z) = tup;
    println!("The value of y is: {y}");
    // Also can use by index
    println!("The value of y is: {tup.1}")
    ```
  - 不带任何值的元组有个特殊的名称：单元 `unit` 。值和类型都写作 `()` ，表示空值或空的返回类型
- array
  - ```rust
    let a = [1,2,3,4,5];
    let a: [i32; 5] = [1,2,3,4,6];
    let a = [3; 5];  // equal: [3,3,3,3,3]
    ```

### function

Keyword: `fn`
style: `snake case`

区分了 语句 statement 和表达式 expression ？

```rust
// this is a expression, and the return value is 4
// a code block is a expression, caurse the `x + 1` without `;`
let y = {
    let x = 3;
    x + 1
};
```

文档啰哩八索的。

### Comment

`//`


### Control Flow

```rust
// if

let number = 5;
if number < 5 {
    println!("condition is false");
} else if number > 5 {
    println!("condition is false");
} else {
    println!("condition is true")
}
// if 不会显式转换，只接收 bool 值
let cond = true;
let number = if cond { 5 } else { 6 }; // ?


// loop | while | for
let mut counter = 0;

let result = loop {
    counter += 1;
    
    if counter == 20 {
        break counter * 2;  // 上面说不要分号，这里又要了
    }
}

// loop label `'[label]`

let a = [1,2,3,4,5];
for val in a {
    println!("the value is: {val}");
}

for num in (1..4).rev() {
    println!("{num}");
}

// 莫名其妙
```


## Ownership

这是 Rust 最与众不同的特性。
它让 Rust 无需垃圾回收即可保障内存安全。

*所有程序都必须管理其运行时使用计算机内存的方式*
一些语言具有垃圾回收，在程序运行时有规律的寻找不再使用的内存；
另一些语言中，必须亲自手动分配并释放内存。

Rust 选择了第三种方式：通过所有权系统管理内存。


所有权规则：

1. Rust 的每个值都有一个所有者 owner
2. 值在任意时刻有且只有一个所有者
3. 当所有者（变量）离开作用域，这个值将被丢弃。

这里以 `String` 为例：

```rust
// `::` 是运算符，允许将特定的 `from` 函数置于 `String`类型的命名空间（namespace）下。
// 而不需要使用类似 `string_from` 这样的名字
// 我没 get 到他在说什么，我并没有少敲键盘啊
let mut str = String::from("hello");
str.push_str(", world");  // push_str() 在字符串后面追加字面值
println!("{}", str);
```

Rust 的策略是：内存在拥有它的变量离开作用域后就被自动释放。

```rust
{
    let s = String::from("hello"); // 从这里开始，s 有效
    ...
    // 使用 s
} // 作用域结束
  // s 不再有效
```

当变量离开作用域，Rust 自动调用一个特殊的函数。叫做 `drop` 。

> C++ 中，这叫 资源获取即初始化 Resource Acquisition Is Initialization RAII  
> 咋一看挺好，很多问题还不确定

```rust
let x = 5;
let y = x;
// 这两个变量都等于 5, 且都在栈里

let s1 = String::from("hello");
let s2 = s1;

// 首先，String 由三部分组成，ptr, len, capacity
// 所以，这里是相当于浅拷贝的
// 这就有个问题了，当它们离开作用域时，会对相同的底层内存二次释放 double free
// 所以为了确保内存安全，`let s2 = s1;` 后，Rust 认为 s1 不再有效，因此不需要清理
// 所以此时调用 s1: `println!("{}", s1);` 会报错
// ？？？
// 所以这不叫浅拷贝，叫移动 move。
```

如果确实要深度复制 `String` 中堆上的数据，而不仅仅是栈上的数据，可以用 `clone`

```rust
let s1 = String::from("hello");
let s2 = s1.clone();

println!("s1 = {}, s2 = {}", s1, s2);
```

所以：

```rust
fn main() {
    let s = String::from("hello");
    
    takes_ownership(s);   // s 移动到函数
                          // 所以这里开始 s 不再可用
                          // ？？？ get 不到在干嘛
}
```

来回切换所有权有点抽象，可以用元组

```rust
fn main() {
    let s1 = String::from("hello");
    
    let (s2, len) = calculate_length(s1);
    
    println!("The length of '{}' is {}.", s2, len);
}

fn calculate_length(s: String) -> (String, usize) {
    let length = s.len();
    
    (s, length)
}

// 感觉有点蠢
// 算是很纯粹的在关注数据的流动
```

上面这个代码的明显的问题就是 `calculate_length` 居然还要返回 `String` 才能使用。
解决方案是使用引用 reference 。
与指针不同，引用确保指向某个特定类型的有效值。

```rust
fn main() {
    let s = String::from("hello");
    
    let len = calculate_length(&s);
    
    println!("the value '{}' length is {}", s, len);
}

fn calculate_length(s: &String) -> usize {
    s.len()
}
```

以往我理解的引用是为了在函数内部修改变量值。
这里的引用是表明这个值只是暂用，借用，并不拥有。
而且这里的引用不能用来修改值。（默认）不允许修改引用的值。


但也不是不能改。
首先，创建 `String` 时使用 `mut` 标识可修改。
然后，传参时标识可变引用 `&mut s` ，这表明 `s` 是可修改的。

```rust
fn main() {
    let mut s = String::from("hello");
    change(&mut s);
}

fn change(s: &mut String) {
    s.push_str(", world");
}
```

可变引用也有限制：如果你已有一个对该变量的可变引用，就不能再创建对该变量的引用。

```rust
let mut s = String::from("hello");

let r1 = &mut s;
let r2 = &mut s;  // 会报错
```

可以通过 `{}` 来创建局部可变引用，从而有多个可变引用

```rust
let mut s = String::from("hello");

{
    let r1 = &mut s;
}

let r2 = &mut s;
```

且，不可用时有可变引用和不可变引用。

不过不可变引用在使用后就会失效，后续代码就可以用了。

```rust
let mut s = String::from("hello");

let r1 = &s;
let r2 = &s;

println!("{} is equal to {}", r1, r2);
// now, r1 and r2 are unavailable

let r3 = &mut s;
println!("{}", r3);
```

直接返回 `String` 这种堆上分配的数据不可返回引用

```rust
// 这是错的
fn func1() -> &String {
    let tmp = String::from("hello");
    &tmp
}

// 这样直接转移所有权就行
fn func2() -> String {
    let tmp = String::from("hello");
    tmp
}
```


### Slice

slice 允许引用集合中的一段连续的元素，是一种引用。

```rust
// s 是一个以空格分开的字符串
// 返回第一个单词的索引
fn first_work(s: &String) -> usize {
    let bytes = s.as_bytes()
    
    for (i, &item) in bytes.iter().enumerate() {
        if item == b' ' {
            return i;
        }
    }
    
    s.len()
}
```

```rust
// 字符串 slice
let s = String::from("Hello, World");

let hello = &s[0..5];  // [start_index..end_index]
// let hello = &s[..5];  // 如果从 0 开始，可以省略 0
let world = &s[6..11];
// let world = &s[6..];  // 如果包含最后一个字节，可以省略尾部

let len = s.len();
let slice1 = &s[3..len];
// let slice1 = &s[3..];

// 也可以同时省略
let slice2 = &s[..];
```

```rust
// 字符串 slice 写作 &str
fn first_word(s: &String) -> &str {
    let bytes = s.as_bytes();
    
    for (i, &item) in bytes.iter().enumerate() {
        if item == b' ' {
            return &s[..i];
        }
    }
    &s[..]
}
```

```rust
fn first_word(s: &str) -> &str

fn main() {
    let my_str = String::from("hello world");
    
    let word = first_word(&my_str[0..6]);
    let word = first_word(&my_str[..]);
    
    let word = first_word(&my_str);
    
    let my_str_literal = "hello world";  // 这种字面量就是 &str
    
    let word = first_word(&my_str_literal[0..6]);
    let word = first_word(&my_str_literal[..]);
    
    let word = first_word(&my_str_literal);
}
```


```rust
let arr = [1,2,3,4,5];

let slice = &arr[1..3];  // &[i32]

assert_eq!(slice, &[2,3]);
```


## Struct

```rust
struct User {
    active: bool,
    username: String,
    email: String,
    sign_in_count: u64,
}

fn main() {
    let mut user1 = User{
        active: true,
        username: String::from("ua123"),
        email: String::from("ua@mail.com"),
        sign_in_count: 1,
    };
    
    user1.email = String::from("another@mail.com");
    
    // 根据 user1 创建新的 user2
    let user2 = User {
        active: user1.active,
        username: user1.username,
        email: String::from("another@mail.com"),
        sign_in_count: user1.sign_in_count,
    };
    // 这里也可以简写
    let user3 = User{
        email: String::from("third@mail.com"),
        ...user2
    };
}

fn build_user(email: String, username: String) -> User {
    User{
        active: true,
        username: username,
        email: email,
        sign_in_count: 1,
    }
}

// 字段名和参数重复可以简写
fn build_user2(email: String, username: String) -> User {
    User{
        active: true,
        username,
        email,
        sign_in_count: 1,
    }
}
```

也可以定义元组结构体

```rust
struct Color(i32,i32,i32);
struct Point(i32, i32, i32);

fn main() {
    let block = Color(0, 0, 0);
    let origin = Point(0, 0, 0);
}
```

结构体字段类型为 `String` 而不是 `&str` 是有意的，因为希望结构体拥有它的数据，而不是引用。
但也可以，需要加上 生命周期 lifetime


```rust
fn main () {
    let width1 = 30;
    let height1 = 50;
    
    println!("The area of the rectangle is {} square pixels.", area(width1, height1));
}

fn area(width: u32, height: u32) -> u32 {
    width * height
}
// 可读性是较弱的，因为 width, height 是关联的数据，函数签名包括数据结构都没体现出

// 使用元组来指定数据
fn main() {
    let rect1 = (30, 50);
    
    println!("The are of the rectangle is {} square pixels.", area(rect1));
}

fn area(dimensions: (u32, u32)) -> u32 {
    dimensions.0 * dimensions.1
}
// 缺点是没有明确 width, height 的顺序

// 使用结构体再重构
struct Rectangle {
    width: u32,
    height: u32,
};

fn main() {
    let rect1 = Rectangle {
        width: 30,
        height: 50,
    };
    
    println!("The area of rectangle is {} square pixels.", area(&rect1));
}

fn area(rect: &Rectangle) -> u32 {
    rect.width * rect.height
}

// 上面直接在 println! 里输出 Rectangle 是错的
// 因为没有实现 std::fmt::Display 方法 trait?

#[derive(Debug)]      // 增加属性来派生 Debug trait，使用调试格式输出 Rectangle
struct Rectangle { ... };

println!("rect1 is {:?}", rect1); // 也可以是 `{:#?}` ，更易读，格式化了
```


### Method

```rust
#[derive(Debug)]
struct Rectangle {
    width: u32,
    height: u32,
}

// impl implementation
// imple Rectangle 表示这个代码块都是 Rectangle 相关
impl Rectangle {
    fn area(&self) -> u32 {   // 这就是个 方法, 也叫做 关联函数 associated function
        self.width * self.height
    }
    
    fn width(&self) -> {    // 可以和字段名相同
        self.width > 0
    }
    
    fn can_hold(&self, other: &Rectangle) -> bool {
        self.width > other.width && self.height > other.height
    }
    
    // 也可以定义第一个参数不是 `&self` 的方法，这就不是方法了，是 `类` 函数
    // 经常用作构造函数，名称通常为 `new` ，但 `new` 不是关键字
    // `Self` 代值 Rectangle
    fn square(size: u32) -> Self {
        Self {
            width: size,
            height: size,
        }
    }
    // let square1 = Rectangle::square(3); 这样调用
} // 也可以将上述方法分开放置在不同的 impl 块里，是可以的，不过这里没必要

fn main() {
    let rect1 = Rectangle {
        width: 30,
        height: 50,
    };
    
    println!("The area of rectangle is {} square pixels.", rect1.area());
}
```


## Enum 枚举 和 match 模式匹配

```rust
enum IpAddrKind {
    V4,  // 这两个就是枚举的 成员 variants
    V6,
}

// 这两个都是 IpAddrKind 类型
let four = IpAddrKind::V4;
let six = IpAddrKind::V6;

fn route(ip_kind: IpAddrKind) {}

// 可以这样
struct IpAddr {
    kind: IpAddrKind,
    address: String,
}

// 不过可以更简洁
enum IpAddr {
    V4(String),
    V6(String),
}

let home = IpAddr::V4(String::from("127.0.0.1"));
let loopback = IpAddr::V6(String::from("::1"));
// 这样就不需要额外的结构体了

// 或者更方便点
enum IpAddr {
    V4(u8, u8, u8, u8),
    V6(String),
}

let home = IpAddr::V4(127, 0, 0, 1);
let loopback = IpAddr::V6("::1");


// 这是标准库中
struct Ipv4Addr { ... }

struct Ipv6Addr { ... }

enum IpAddr {
    V4(Ipv4Addr),
    V6(Ipv6Addr),
}

// 可以将任意类型的数据放入枚举成员：例如字符串，数字类型，结构体
// 甚至包含另一个枚举

enum Message {
    Quit,   // 没有关联任何数据
    Move { x: i32, y: i32 },  // 类似结构体包含命名字段
    Write(String),
    ChangeColor(i32, i32, i32),
}

// 也可以在枚举上定义方法
impl Message {
    fn call(&self) {
        ...
    }
}

let m = Message::Write(String::from("hello"));
m.call();
```

### Option

Rust 没有 NULL， Option 代表要么有值要么没值

```rust
enum Option<T> {
    None,
    Some(T),
}
```

```rust
let number = Some(5);
let char = Some('e');

let absent_number: Option<i32> = None;
```

```rust
let x: i8 = 5;
let y: Option<i8> = Some(6);

let sum = x + y;  // 不行，会报错
```


### match control flow

```rust
enum Coin {
    Penny,
    Nickel,
    Dime,
    Quarter,
}

fn value_in_cents(coin: Coin) -> u8 {
    match coin {
        Coin::Penny => 1,
        Coin::Nickel => 5,
        Coin::Dime => 10,
        Coin::Quarter => 25,
    }
}
```

```rust
#[derive(Debug)]
enum UsState {
    Alabama,
    Alaska,
    // --snip--
}

enum Coin {
    Penny,
    Nickel,
    Dime,
    Quarter(UsState),
}

fn value_in_cents(coin: Coin) -> u8 {
    match coin {
        Coin::Penny => 1,
        Coin::Nickel => 5,
        Coin::Dime => 10,
        Coin::Quarter(state) => {
            println!("State quarter from {:?}", state)
            25
        }
    }
}
```

用 match 来处理 Option

```rust
fn plus_one(x: Option<i32>) -> Option<i32> {
    match x {
        None => None,
        Some(i) => Some(i + 1),
    }
}

let five = Some(5);
let six = plus_one(five);
let none = plus_one(None);
```

说好听点吧，巧妙的设计。

```rust
fn plus_one(x: Option<i8>) -> Option<i8> {
    match x {
        Some(i) => Some(i + 1),
    }
}
// 会报错，因为 None 没有处理
// match 必须处理所有可能情况
```

这种代码属于覆盖了所有的情况，因为 other 包含了剩余的所有值

```rust
let dice_roll = 9;
match dice_roll {
    3 => add_fancy_hat(),
    7 => remove_fancy_hat(),
    other => move_player(other),
}

// 如果我们不想处理剩下的所有值
match dice_roll {
    3 => add_fancy_hat(),
    7 => remove_fancy_hat(),
    _ => reroll(),   // 这里是忽略了剩下的值
    // Or `_ => (),` 这是什么也不做，一个空元组
}

fn move_player(num_spaces: u8) {}
```



### if let

```rust
let config_max = Some(3u8);
match config_max {
    Some(max) => println!("The maximum is configured to be {}", max),
    _ => (),
}

// 上面这种只关心一个分支的情况，可以用 if let 简写
if let Some(max) = config_max {
    println!("The maximum is configured to be {}", max);
}

// 也可以 if let ... else ...
let mut count = 0;
match coin {
    Coin::Quarter(state) => println!("State Quarter from {:?}", state),
    _ => count += 1;
}

if let Coin::Quarter(state) = coin {
    println!("state quarter from {:?}", state);
} else {
    count += 1;
}
```



## packages, crates and modules

Rust 有许多功能可以让你管理代码的组织，包括哪些内容可以公开，哪些内容作为私有部分，以及程序每个作用域中的名字。
这些功能，有时被统称为 "模块系统"，包括：

- 包 package   Cargo 的一个功能，允许你构建，测试和分享 crate
- Crates       一个模块的树形结构，它形成了库或二进制项目
- 模块 Modules 和 use  允许你控制作用域和路径的私有性
- 路径 path    一个命名例如结构体，函数或模块等项的方式



### package and crates

Crate 是 Rust 在编译时最小的代码单位。

Crate 有两种形式：二进制项和库。 二进制项可以被编译为可执行文件，比如一个命令行程序或服务器。
库没有 `main` 函数，也不会被编译为可执行文件，这和其他的 `library` 一致。

*crate root* 是一个源文件，Rust 编译器以它为起点，并构成你的 crate 的根模块。

包 package 是提供一系列功能的一个或多个 crate。
一个包会包含一个 `Cargo.toml` 文件。阐述如何去构建这些 crate。

默认： `src/main.rs` 就是一个与包同名的二进制的 crate 的根。
同样的，如果包目录中有 `src/lib.rs` 。。
有点啰嗦，而且没看明白。


### define modules to control scope and privacy

- 从 crate 根节点开始: 当编译一个 crate, 编译器首先在 crate 根文件（通常，对于一个库 crate 而言是src/lib.rs，对于一个二进制 crate 而言是src/main.rs）中寻找需要被编译的代码。
- 声明模块: 在 crate 根文件中，你可以声明一个新模块；比如，你用mod garden;声明了一个叫做garden的模块。编译器会在下列路径中寻找模块代码：
  - 内联，在大括号中，当mod garden后方不是一个分号而是一个大括号
  - 在文件 src/garden.rs
  - 在文件 src/garden/mod.rs
- 声明子模块: 在除了 crate 根节点以外的其他文件中，你可以定义子模块。比如，你可能在src/garden.rs中定义了mod vegetables;。编译器会在以父模块命名的目录中寻找子模块代码：
  - 内联，在大括号中，当mod vegetables后方不是一个分号而是一个大括号
  - 在文件 src/garden/vegetables.rs
  - 在文件 src/garden/vegetables/mod.rs
- 模块中的代码路径: 一旦一个模块是你 crate 的一部分，你可以在隐私规则允许的前提下，从同一个 crate 内的任意地方，通过代码路径引用该模块的代码。举例而言，一个 garden vegetables 模块下的Asparagus类型可以在crate::garden::vegetables::Asparagus被找到。
- 私有 vs 公用: 一个模块里的代码默认对其父模块私有。为了使一个模块公用，应当在声明时使用pub mod替代mod。为了使一个公用模块内部的成员公用，应当在声明前使用pub。
- use 关键字: 在一个作用域内，use关键字创建了一个成员的快捷方式，用来减少长路径的重复。在任何可以引用crate::garden::vegetables::Asparagus的作用域，你可以通过 use crate::garden::vegetables::Asparagus;创建一个快捷方式，然后你就可以在作用域中只写Asparagus来使用该类型。


这里创建一个 `backyard` 的二进制 crate 来说明这些规则。
该 crate 的路径同样命名为 `backyard`

```
backyard
- Cargo.lock
- Cargo.toml
- src
  - garden
    - vegetables.rs
  - garden.rs
  - main.rs
```

```rust
// src/main.rs
use crate::garden::vegetables::Asparagus;

// 告诉编译器应该在 `src/garden.rs` 中发现代码
pub mod garden;

fn main() {
    let plant = Asparagus{};
    println!("I'm growing {:?}", plant);
}

// src/garden.rs
// 告诉编译器 `src/garden/vegetables.rs` 中的代码也应该包括
pub mod vegetables;

// src/garden/vegetables.rs
#[derive(Debug)]
pub struct Asparagus {}
```

开始实操了

通过 `cargo new --lib restaurant` 创建一个新的包

```rust
// src/lib.rs
mod front_of_house {
    mod hosting {
        fn add_to_waitlist() {}
        fn seat_at_table() {}
    }

    mod serving {
        fn take_order() {}
        fn serve_order() {}
        fn take_payment() {}
    }
}
```

创建了字模块 `front_of_house` 及子孙模块 `hosting, serving`
`use crate::front_of_house::hosting`
这是内联的写法。


### path for referring to an item in the module tree

为了调用一个函数，我们需要知道它的路径

- 绝对路径： 以 crate 开头的全路径；对于外部 crate 的代码，是以 crate 名开头的绝对路径，对于当前 crate 的代码，是以字面值 crate 开头
- 相对路径： 从当前模块开始，以 `self, super` 或定义在当前模块中的标识符开头
- 都用 `::` 作为分隔的标识符

```rust
// 这个无法编译的
// 给 `hosting, add_to_waitlist` 加上 `pub`，现在可以了
mod front_of_house {
    pub mod hosting {
        pub fn add_to_waitlist() {}
    }
}

pub fn eat_at_restaurant() {
    // absolute path
    crate::front_of_house::hosting::add_to_waitlist();
    
    // relative path
    front_of_house::hosting::add_to_waitlist();
}
```


```rust
// src/lib.rs
fn deliver_order() {}

mod back_of_house {
    fn fix_incorrect_order() {
        cook_order();
        // relative path
        // 有一说一，让人有点头疼。
        super::deliver_order();
    }
    
    fn cook_order() {}
}
```


```rust
mod back_of_house {
    // 这个结构体是公有的
    // 其中的 toast 也是公有的
    // seasonal_fruit 是私有的
    pub struct Breakfast {
        pub toast: String,
        seasonal_fruit: String,
    }
    
    impl Breakfast {
        pub fn summer(toast: &str) -> Breakfast {
            Breakfast {
                toast: String::from(toast),
                seasonal_fruit: String::from("peaches"),
            }
        }
    }
}

pub fn eat_at_restaurant() {
    //
    let mut meal = back_of_house::Breakfast::summer("Rye");
    // 
    meal.toast = String::from("Wheat");
    println!("I'd like {} toast please.", meal.toast);
    
    // 这个字段不允许修改，编译会失败
    // meal.seasonal_fruit = String::from("blueberries");
}

// 对应的是 enum, 所有字段都会是公有的
mod back_of_house {
    pub enum Appetizer {
        Soup,
        Salad,
    }
}

pub fn eat_at_restaurant() {
    let order1 = back_of_house::Appetizer::Soup;
    let order2 = back_of_house::Appetizer::Salad;
}
```


### bringing paths into scope with the use keyword

```rust
// src/lib.rs
mod front_of_house {
    pub mod hosting {
        pub fn add_to_waitlist() {}
    }
}

use crate::front_of_house::hosting;

pub fn eat_at_restaurant() {
    hosting::add_to_waitlist();
}
```

use 只能创建 use 所在的特定域内的短路径。
如果将其在其他模块内调用，这就是不同于 use 语句的作用域，会导致无法编译。
例如：

```rust
// src/lib.rs
mod front_of_house {
    pub mod hosting {
        pub fn add_to_waitlist() {}
    }
}

use crate::front_of_house::hosting;

mod customer {
    pub fn eat_at_restaurant() {
        hosting::add_to_waitlist();
    }
}

// 可以将 use 移入 customer 模块内
// 或是用 super 调用 `super::hosting::add_to_waitlist();`
```


```rust
// src/main.rs
use std::collections::HashMap;

fn main() {
    let mut map = HashMap::new();
    map.insert(1,2);
}
```


使用 as 关键字提供新的名称

```rust
// src/lib.rs
use std::fmt::Result;
use std::io::Result as IoResult;
```

使用 `pub use` 重导出名称

```rust
mod front_of_house {
    pub mod hosting {
        pub fn add_to_waitlist() {}
    }
}

pub use crate::front_of_house::hosting;

pub fn eat_at_restaurant() {
    hosting::add_to_waitlist();
}
```

使用嵌套路径消除大量的 use 行

```rust
// src/main.rs
use std::cmp::Ordering;
use std::io;

use std::{cmp::Ordering, io};

// ---
use std::io;
use std::io::Write;

use std::io::{self, Write};

// 使用 glob 运算符将所有的公有定义引入
use std::collections::*;
```


### separate modules into different files



## Common collections

三个常见的集合

- vector
- String
- HashMap


### vector

```rust
let v: Vec<i32> = Vec::new();

let v = vec![1,2,3];
```


更新

```rust
let mut v = Vec::new();

v.push(5);
v.push(6);
```


读取

通过索引或 `get`

```rust
let v = vec![1,2,3,4,5];

let third:&i32 = &v[2];
println!("The third element is {third}");

let third: Option<&i32> = v.get(2);
match third {
    Some(num) => println!("The third element is {num}"),
    None => println!("There is no third element."),
}
```

```rust
let v = vec![1,2,3,4,5];

let does_not_exist = &v[100];     // 会 panic
let does_not_exist = v.get(100);  // 会返回 None
```

获取了某个值的引用但未使用，则不可操作 vec

```rust
let mut v = vec![1,2,3,4,5];

let first = &v[0];

v.push(6);

println!("The first element is: {first}");

// 这里会报错，因为 first 有对 v 的不可变引用
// 会在扩容的时候出问题，所以报错
```


遍历

```rust
let v = vec![1,2,3];

for i in &v {
    println!("{i}");
}
```

也可以遍历可变引用从而修改

```rust
let mut v = vec![1,2,3];

for i in &v {
    *i += 50;
}
```

使用枚举来存储多个类型

```rust
enum SpreadsheetCell {
    Int(i32),
    Float(f64),
    Text(String),
}

let row = vec![
    SpreadshellCell::Int(3),
    SpreadshellCell::Text(String::from("blue")),
    SpreadshellCell::Float(10.12),
];
```

类似于其他的 struct ，vector 会在离开作用域时释放



### Strings

String 和 `&str` 并不完全相同

```rust
let mut s = String::new("");  // String

let data = "initial contents";  // &str
let s = data.to_string();     // String
```


更新

可以用 `+` 或是 `format!` 宏来拼接

使用 `push` 和 `push_str` 来附加

```rust
let mut s = String::from("foo");

s.push_str("bar");  // push_str 并不获取所有权
s.push('l');        // push 是附加单个字符

// ---
let s1 = String::from("Hello,");
let s2 = String::from(" world.");
let s3 = s1 + &s2;   // 这里 s1 被移动了，不能用了
// `+` 是调用了 `add` 函数
// `fn add(self, s: &str) -> String`

// ---
// 连接多个字符串， `+` 会有点笨
let s1 = String::from("tic");
let s2 = String::from("tac");
let s3 = String::from("toe");

let s = s1 + "-" + &s2 + "-" + &s3;   // tic-tac-toe
// 这里可以用 format!
let s = format!("{s1}-{s2}-{s3}");  // format! 不会获取任何参数的所有权
```


索引字符串

Rust 的字符串不支持索引

因为很容易是无效的数据



字符串 slice

需要小心，也挺容易 panic


遍历字符串

操作字符串的最好方法是明确需要字符还是字节。
对于 Unicode 标量使用 `chars` 

```rust
for c in "中国".chars() {
    println!("{c}");
}

// 中
// 国

// ---
// bytes() 方法返回每一个字节
for b in "中国".bytes() {
    println!("{b}");
}
// 228
// 184
// 173
// 229
// 155
// 189
```


### HashMap

```rust
// 新建并插入一些键值对
use std::collections::HashMap;

let mut scores = HashMap::new();

scores.insert(String::from("Blue"), 10);
scores.insert(String::from("Yellow"), 50);

// 可以通过 `get` 方法取值
let team_name = String::from("Blue");
let score = scores.get(&team_name).copied().unwrap_or(0);
// get 返回 Option(&v)， 没有对应值会返回 None
// 通过 copied 获取一个 Option<i32> 而不是 Option<&i32>
// 再调用 unwrap_or 将 scores 中没有对应项时设置为 0

// ---
// 也可以通过 for in 遍历
for (key, value) in &scores {  // 这是无序的
    println!("{key}: {value}");
}
```


所有权相关

```rust
use std::collections::HashMap;

let field_name = String::from("Favourite Color");
let field_value = String::from("Blue");

let mut map = HashMap::new();
map.insert(field_name, field_value);
// 这里 field_name, field_value 所有权移动到 map 了
// 这两个变量无效了
```

更新

```rust
// 覆盖一个值
use std::collections::HashMap;

let mut scores = HashMap::new();

scores.insert(String::from("Blue"), 10);
scores.insert(String::from("Blue"), 25);

// 只在键没有对应值时插入键值对
// 如果键存在则不操作，不存在则插入
// 为此 map 有一个特定的 API，entry。
let mut scores = HashMap::new();
scores.insert(String::from("Blue"), 10);

scores.entry(String::from("Yellow")).or_insert(50);
scores.entry(String::from("Blue")).or_insert(50);

// ---
// 根据旧值更新
let text = "hello world wonderful world";

let mut map = HashMap::new();

for word in text.split_whitespace() {
    // 如果 word 没有对应的值，赋予 0，然后返回该值（有值则直接返回）
    let count = map.entry(word).or_insert(0);
    *count += 1;
}

println!("{:?}", map);
```

HashMap 使用的是 [SipHash](https://en.wikipedia.org/wiki/SipHash) 的哈希函数。
不算快，不过安全，可以自己切换掉。



## error handling

Rust 将错误分为两大类： 可恢复的 recoverable 和 不可恢复的 unrecoverable。

`Result<T, E>`  用于处理可恢复的错误
还有 `panic!` 宏  用于处理立刻退出的情况


### panic!

当出现 panic 时，程序默认会 展开 unwinding，也就是回溯栈并清理遇到的每个函数的数据，这个回溯并清理的过程有很多工作。
另一个选择是直接 终止 abort。这会不清理直接退出。

通过在 `Cargo.toml` 文件增加：

```toml
[profile.release]
panic = 'abort'
```

```rust
// src/main.rs
fn main() {
    panic!("crush and burn");
}
```


### recoverable errors with result

```rust
enum Result<T, E> {
    Ok(T),
    Err(E),
}
```

```rust
// src/main.rs
use std::fs::File;

fn main () {
    // File::open 的返回值为 Result<std::fs::File, std::io::Error>
    let greeting_file_result = File::open("hello.txt");
}

// 可以用 match 来处理
fn main() {
    let greeting_file_result = File::open("hello.txt");
    
    let greeting_file = match greeting_file_result {
        Ok(file) => file,
        Err(error) => panic!("Problem opening the file: {:?}", error),
    };
}

// 可以根据错误类型做不同的处理
use std::io::ErrorKind;

fn main() {
    let greeting_file_result = File::open("hello.txt");
    
    let greeting_file = match greeting_file_result {
        Ok(file) => file,
        Err(error) => match error.kind() {
            ErrorKind::NotFound => match File::create("hello.txt") {
                Ok(fc) => fc,
                Err(e) => panic!("Problem creating the file: {:?}", e),
            },
            other_error => {
                panic!("Problem opening the file: {:?}", other_error);
            }
        },
    };
}
```

match 太多了，可读性一般。
可以用闭包处理

```rust
use std::io::ErrorKind;
use std::fs:File;

fn main() {
    let greeting_file = File::open("hello.txt").unwrap_or_else(|error| {
        if error.kind() == ErrorKind::NotFound {
            File::create("hello.txt").unwrap_or_else(|error {
                panic!("Problem creating the file: {:?}", error);
            })
        } else {
            panic!("Problem opening the file: {:?}", error);
        }
    });
}
```


失败时 panic 的简写： unwrap 和 expect

如果 Result 值为 Ok, unwrap 会返回 Ok 中的值，
如果 Result 值为 Err, unwrap 会调用 panic!

```rust
let greeting_file = File::open("hello.txt").unwrap();
```

expect 类似不过可以自定义 panic 的消息

```rust
let greeting_file = File::open("hello.txt")
    .expect("hello.txt should be included in this project.");
    // 看起来就更适用
```


传播错误  propagating

```rust
use std::io::{self, Read};
use std::fs::File;

fn read_username_from_file() -> Result<String, io::Error> {
    let username_file_result = File::open("hello.txt");
    
    let mut user_file = match username_file_result {
        Ok(file) => file,
        Err(error) => return Err(error),
    };
    
    let mut username = String::new();
    
    match user_file.read_to_string(&mut username) {
        Ok(_) => Ok(username),
        Err(e) => Err(e),
    }
}
```

这种写法很常见，所以 Rust 提供了 `?` 运算符

```rust
// src/main.rs
use std::fs::File;
use std::io::{self, Read};

fn read_username_from_file() -> Result<String, io::Error> {
    let mut username_file = File::open("hello.txt")?;
    let mut username = String::new();
    username_file.read_to_string(&mut username)?;
    Ok(username)
}

// Result 值后的 `?` 和上面的 match 一样
// 如果值为 Ok, 表达式返回 Ok 的值并继续执行
// 否则，Err 作为整个函数的返回值
```


`?` 运算符消除了大量样板代码。
我们甚至可以在 `?` 后直接适用链式方法调用来进一步缩短代码

```rust
use std::fs::File;
use std::io::{self, Read};

fn read_username_from_file() -> Result<String, io::Error> {
    let mut username = String::new();
    
    File::open("hello.txt")?.read_to_string(&mut username)?;
    
    Ok(username)
}
```

还有更短的写法

```rust
use std::fs;
use std::io;

fn read_username_from_file() -> Result<String, io::Error> {
    fs::read_to_string("hello.txt")
}
```

`?` 只能用在返回值是 `Result` 的情况，不然不兼容。
也可以用在 `Option` 上，

```rust
fn last_char_of_first_line(text: &str) -> Option<char> {
    text.lines().next()?.chars().last()
}
```

当然，不会混搭，`?` 不会将 Result 转为 Option。
可以用 Result 的 Ok 方法或是 Option 的 `ok_or` 方法显式转化

```rust
use std::error::Error;
use std::fs::File;

fn main() -> Result<(), Box<dyn Error>> {
    let greeting_file = File::open("hello.txt")?;
    
    Ok(())
}
```


### To panic or not to panic

首先，Result 是必须处理的，编译器强制要求

```rust
use std::net::IpAddr;

let home: IpAddr = "127.0.0.1"
    .parse()
    .expect("hardcoded IP address should be valid");
    // 这里用 unwrap 也行，原则上这里是不会出错的
```



## Generics

```rust
fn main() {
    let number_list = vec![1,2,3,4,5];
    
    let mut largest = &number_list[0];
    
    for number in &number_list {
        if number > largest {
            largest = number;
        }
    }
    
    println!("The largest number is {largest}");
}

//
fn largest(list: &[i32]) -> &i32 {
    let mut largest = &list[0];
    
    for item in list {
        if item > largest {
            largest = item;
        }
    }
    
    largest
}

//
// fn largest<T>(list: &[T]) -> &T {
fn largest<T: std::cmp::PartialOrd>(list: &[T]) -> &T {
    let mut largest = &list[0];
    
    for item in list {
        if item > largest {
            largest = item;
        }
    }
    
    largest
}
```

也可以定义结构体泛型

```rust
struct Point<T> {
    x: T,
    y: T,
}

struct Point2<T, R> {
    x: T,
    y: R,
}

fn main() {
    let integer = Point{x: 5, y: 10};
    let float = Point{x: 1.0, y: 2.0};
    let int_and_float = Point2{x: 4, y: 2.0};
}
```


在方法上也可以实现泛型

```rust
struct Point<T> {
    x: T,
    y: T,
}

impl<T> Point<T> {
    fn x(&self) -> &T {
        &self.x
    }
}

// 也可以为具体类型指定方法
impl Point<f32> {
    fn distance_from_origin(&self) -> f32 {
        (self.x.powi(2) + self.y.powi(2)).sqrt()
    }
}

impl<X1,Y1> Point<X1, Y1> {
    // 方法依然可以指定其他类型
    fn mixup<X2, Y2>(self, other: Point<X2, Y2>) -> Point<X1, Y2> {
        Point{
            x: self.x,
            y: other.y,
        }
    }
}

fn main() {
    let p = Point{x: 5, y: 10};
    println!("p.x = {}", p.x());
}
```


### Trait 定义共同行为

> 类似 interface ，但有一些不同

```rust
// src/lib.rs
// 定义 trait Summary
pub trait Summary {
    fn summarize(&self) -> &String;
}

pub struct NewsArticle {
    pub headline: String,
    pub location: String,
    pub author: String,
    pub content: String,
}

// 实现方法时要注明实现的 trait
impl Summary for NewsArticle {
    fn summarize(&self) -> String {
        format!("{}, by {} ({})", self.headline, self.author, self.location)
    }
}

pub struct Tweet {
    pub username: String,
    pub content: String,
    pub reply: bool,
    pub retweet: bool,
}

impl Summary for Tweet {
    fn Summarize(&self) -> String {
        format!("{}: {}", self.username, self.content)
    }
}

// 使用时也要引入 trait

// src/main.rs
use aggregator::{Summary, Tweet};

fn main() {
    let tweet = Tweet {
        username: String::from("horse_ebooks"),
        content: String::from("of course, as you probably already know, people"),
        reply: false,
        retweet: false,
    };
    
    println!("1 new tweet: {}", tweet.Summarize());
}
```

不能为外部类型实现外部 trait


默认实现。
这样可以在没有实现覆盖方法时有默认调用。

```rust
// src/lib.rs
pub trait Summary {
    fn summarize(&self) -> String {
        String::from("read more...")
    }
}
```

无法从相同的重载实现中调用默认方法



*trait 作为参数*

```rust
pub fn notify(item: &impl Summary) {
    println!("Breaking news {}", item.summarize());
}
// 这是个 trait bound 的语法糖

pub fn notify<T: Summary>(item: &T);

pub fn notify(item1: &impl Summary, item2: &impl Summary);
// 这里两个参数可以是不同类型
// 或者通过泛型定义为相同类型
pub fn notify<T: Summary>(item1: &T, item: &T);

// 指定多个 trait
pub fn notify(item: &(impl Summary + Display));

// 参数长了会比较难读，所以可以用 where 从句 (离谱)
fn some_fn<T: Display + Clone, U: Clone + Debug>(t: &T, u: &u) -> i32 {};
// 这样
fn some_fn<T, U>(t: &T, u: &U) -> i32
where
    T: Display + Clone,
    U: Clone + Debug,
{};
```


也可以用在返回值

```rust
fn return_summarizable() -> impl Summary {
    Tweet{
        // ...
    }
}
```


### lifetime

Rust 中的每一个引用都有其 生命周期。

```rust
fn main() {
    let str1 = String::from("abc");
    let str2 = "xyz";
    
    let result = longest(str1.as_str(), str2);
    println!("The longest string is {}", result);
}
// 编译不通过，会失败
fn longest(x: &str, y: &str) -> &str {
    if x.len() >= y.len() {
        x
    } else {
        y
    }
}
// 因为不确定该返回那个的引用
```

生命周期注解：通常以 `'` 开头，其名称通常全是小写。

```rust
&i32        // 引用
&'a i32     // 带有显式生命周期的引用
&'a mut i32 // 带有显式生命周期的可变引用
```

单个的是没有意义的，这个理解多个如何联系的。
假如有个生命周期 `'a` 的 i32 的引用参数 `first`，还有另一个 `'a` 的 i32 的引用参数 `second` 。
这两个生命周期注解意味着引用 `first, second` 必须和这泛型生命周期存在的一样久。

*函数签名中的生命周期注解*

这和泛型类似，如下的意思是：这两个参数和返回的引用存活的一样久。

```rust
fn longest<'a>(x: &'a str, y: &'a str) -> &'a str {
    if x.len() >= y.len() {
        x
    } else {
        y
    }
}
```

实际含义是： `longest` 返回的参数的生命周期与引用的参数的生命周期的较小者一致。


*结构体定义中的生命周期*

```rust
// src/main.rs
struct ImportantExcerpt<'a> {
    part: &'a str,
}
// 这个注解意味着： ImportantExcerpt 的实例不能比 part 中的引用存在的更久

fn main() {
    let novel = String::from("Call me Judy. Some years age...");
    let first_sentence = novel.split(".").next().expect("Cound not find a '.");
    let i = ImportantExcerpt {
        part: first_sentence,
    };
}
```


*生命周期省略 Lifetime Elision*

```rust
fn first_word(s: &str) -> &str {
    let bytes = s.as_bytes();
    
    for (i, &item) in bytes.iter().enumerate() {
        if item == b' ' {
            return &s[..i];
        }
    }
    
    &s[..]
}
```


*静态生命周期*

有一个特殊的生命周期： `static` ，其生命周期能存活在整个程序期间。
所有的字符串字面值都有 `static` 生命周期，也可以手动标注。

```rust
let s: &'static str = "I have a static lifetime.";
```


```rust
use std::fmt::Display;

fn longest_with_an_announcement<'a, T>(x: &'a str, y: &'a str, ann: T) ->&'a str
where
    T: Display,
{
    println!("Announcement. {}", ann);
    if x.len() >= y.len() {
        x
    } else {
        y
    }
}
```


## testing

### Writing tests

`assert!`, `assert_eq!`, `assert_ne!`, `should_panic`


### running tests



## functional features

- 闭包 Closures
- 迭代器 Iterators


### Closures

```rust
#[derive(Debug, PartialEq, Copy, Clone)]
enum ShirtColor {
    Red,
    Blue,
}

struct Inventory {
    shirts: Vec<ShirtColor>,
}

impl Inventory {
    fn giveaway(&self, user_preference: Option<ShirtColor>) -> ShirtColor {
        user_preference.unwrap_or_else(|| self.most_stocked())
    }
    
    fn most_stocked(&self) -> ShirtColor {
        let mut num_red = 0;
        let mut num_blue = 0;
        
        for color in &self.shirts {
            match {
            ShirtColor::Red => num_red += 1,
            ShirtColor::Blue => num_blue += 1,
            }
        }
        
        if num_red > num_blue {
            ShirtColor::Red
        } else {
            ShirtColor::Blue
        }
    }
}

fn main() {
    let store = Inventory{
        shirts: vec![ShirtColor::Blue, ShirtColor::Red, ShirtColor::Blue],
    };
    
    let user_pref1 = Some(ShirtColor::Red);
    let giveaway1 = store.giveaway(user_pref1);
    println!(
        "The user will preference {:?} gets {:?}",
        user_pref1, giveaway1
    );
    
    let user_pref2 = None;
    let giveaway2 = store.giveaway(user_pref2);
    println!(
        "The user will preference {:?} gets {:?}",
        user_pref2, giveaway2
    );
}
```

一般类似根据上下文或是推导，也可以增加类型标注

```rust
let expensive_closure = |num: u32| -> u32 {
    println!("calculating slowly...");
    thread::sleep(Duration::from_secs(2));
    num
}
```

一些对比

```rust
fn  add_one_v1   (x: u32) -> u32 { x + 1 }
let add_one_v2 = |x: u32| -> u32 { x + 1 };
let add_one_v3 = |x|             { x + 1 };
let add_one_v4 = |x|               x + 1 ;
```


*捕获引用或者移动所有权*

```rust
fn main() {
    let list = vec![1,2,3];
    println!("before defining closure: {:?}", list);
    
    let only_borrows = || println!("From closure: {:?}", list);
    
    println!("before calling closure: {:?}", list);
    only_borrows();
    println!("after calling closure: {:?}", list);
}
```


## Smart Pointers

- `Box<T>`
- `Rc<T>`
- `Ref<T>`, `RefMut<T>`, `RefCell<T>`


### Drop trait

```rust

struct CustomSmartPointer {
    data: String,
}

impl Drop for CustomSmartPointer {
    fn drop(&mut self) {
        println!("Dropping CustomSmartPointer with data: `{}`", self.data);
    }
}

fn main() {
    let  c = CustomSmartPointer{
        data: String::from("my stuff"),
    };
    let d = CustomSmartPointer{
        data: String::from("other stuff"),
    };
    println!("CustomSmartPointers created.");
}

// CustomSmartPointers created.
// Dropping CustomSmartPointer with data: `other stuff`
// Dropping CustomSmartPointer with data: `my stuff`
```


可以通过 `std::mem::drop` 提前释放



### rc

`RC<T>` 只能用于单线程场景

```rust
enum List {
    Cons(i32, Rc<List>),
    Nil,
}

use crate::List::{Cons, Nil};
use std::rc::Rc;

fn main() {
    let a = Rc::new(Cons(5, Rc::new(Cons(10, Rc::new(Nil)))));
    println!("count after creating a = {}", Rc::strong_count(&a));
    let b = Cons(3, Rc::clone(&a));
    println!("count after creating b = {}", Rc::strong_count(&a));
    {
        let c = Cons(4, Rc::clone(&a));
        println!("count after creating c = {}", Rc::strong_count(&a));
    }
    println!("count after c goes out of scope = {}", Rc::strong_count(&a));
}

// 1
// 2
// 3
// 2
```


### Interior mutability



## Concurrency



## OOP

关于面向对象编程有很多相互矛盾的定义。
在一些定义下，Rust 是面向对象的；在其他定义下，Rust 不是。




## Pattern

模式 Pattern 是 Rust 的特殊语法。
结合使用模式和 `match` 表达式以及其他结构可以提供更多对程序控制流的支配权。
模式由如下一些内容组成：

- 字面值
- 解构的数组，枚举，结构体或者元组
- 变量
- 通配符
- 占位符

refutable 和 irrefutable 模式的区别


### All the places for patterns

match 分支

```rust
match VALUE {
    PATTERN => EXPRESSION,
    PATTERN => EXPRESSION,
    PATTERN => EXPRESSION,
}

// e.g.
match x {
    None => None,
    Some(i) => Some(i + 1),
}
```


if let 表达式

```rust
fn main() {
    let favourite_color: Option<&str> = None;
    let is_tuesday = false;
    let age: Result<u8, _> = "34".parse();
    
    if let Some(color) = favourite_color {
        println!("Using your favourite color, {color}, as the background");
    } else if is_tuesday {
        println!("Tuesday is green day.");
    } else if let Ok(age) = age {
        if age > 30 {
            println!("Using purple as the background color");
        } else {
            println!("Using orange as the background color");
        }
    } else {
        println!("Using Blue as the background color");
    }
}
```



while let 循环

```rust
let mut stack = Vec::new();

stack.push(1);
stack.push(2);
stack.push(3);

while let Some(top) = stack.pop() {
    println!("{}", top);
}
```


for 循环

```rust
let v = vec!['a', 'b', 'c'];

for (index, val) in v.iter().enumerate() {
    println!("{} is at index {}", val, index);
}
```


let 语句

```rust
let x = 5;

let PATTERN = EXPRESSION;

let (x, y, z) = (1,2,3);
```



函数参数



### Refutability 可反驳性：模式是否会匹配失效


### Pattern Syntax

匹配字面值

```rust
let x = 1;

match x {
    1 => println!("one"),
    2 => println!("two"),
    3 => println!("three"),
    _ => println!("anything"),
}
```


匹配命名变量

```rust
let x = Some(5);
let y = 10;

match x {
    Some(50) => println!("Got 50"),
    // 会匹配到这，因为 Result 只有 Some 和 None，所以肯定是这，y 只是形参，和上面那个 y 毫无关系
    Some(y) => println!("Matched, y = {y}"),
    _ => println!("Default case, x = {:?}", x),
}

println!("at the end, x = {:?}, y= {y}", x);
```


多个模式

```rust
let x = 1;

match x {
    1 | 2 => println!("one or two"),
    3 => println!("three"),
    _ => println!("anything"),
}
```



通过 ..= 匹配值的范围

```rust
let x = 5;

match x {
    1..=5 => println!("one through 5"),
    _ => println!("something else"),
}


// 范围只允许用于数字或 char 值
let x = 'c';

match x {
    'a'..='j' => println!("early ASCII letter"),
    'k'..='z' => println!("late ASCII letter"),
    _ => println!("something else"),
}
```


解构并分解值

```rust
// 解构结构体
struct Point {
    x: i32,
    y: i32,
}

fn main() {
    let p = Point{ x: 0, y: 7};
    
    let Point{x: a, y: b} = p;
    // Or let Point{x, y} = p;
    assert_eq!(0, a);
    assert_eq!(7, b);
}

//
fn main() {
    let p = Point{ x: 0, y: 7};
    
    match p {
        Point{x, y: 0} => println!("On the x axis at {x}"),
        Point{x: 0, y} => println!("On the y axis at {y}"),
        Point{ x, y } => {
            println!("On neither axis: ({x}, {y})");
        }
    }
}


// 解构枚举
enum Message {
    Quit,
    Move{x: i32, y: i32},
    Write(String),
    ChangeColor(i32, i32, i32),
}

fn main() {
    let msg = Message::ChangeColor(0, 160, 255);
    
    match msg {
        Message::Quit => {
            println!("The Quit variant has no data to destructure");
        }
        Message::Move{x, y} => {
            println!("Move in the x direction {x} and in the y direction {y}");
        }
        Message::Write(text) => {
            println!("Text message: {text}");
        }
        Message::ChangeColor(r,g,b) => {
            println!("Change the color to red {r}, green {g}, and blue {b}");
        }
    }
}


// 解构嵌套的结构体和枚举
enum Color {
    Rgb(i32, i32, i32),
    Hsv(i32, i32, i32),
}

enum Message {
    Quit,
    Move{x: i32, y: i32},
    Write(String),
    ChangeColor(Color),
}

fn main() {
    let msg = Message::ChangeColor(Color::Hsv(0, 160, 255));
    
    match msg {
        Message::ChangeColor(Color::Rgb(r, g, b)) => {
            println!("Change color to red {r}, green {g}, and blue {b}");
        }
        Message::ChangeColor(Color::Hsv(h, s, v)) => {
            println!("Change color to hue {h}, saturation {s}, value {v}");
        }
        _ => (),
    }
}


// 解构结构体和元组
let ((feet, inches), Point{x, y}) = ((3, 19), Point{x: 3, y: -19});


// 忽略模式中的值
fn foo(_: i32, y: i32) {
    println!("This code only uses the y parameter: {}", y);
}

fn main() {
    foo(3,4);
}
```

可以通过在变量名前加上 `_` 来忽略未使用的变量。

```rust
fn main() {
    let _x = 5;
    let y = 10;
}
```

用 `..` 忽略剩余值

```rust
struct Point {
    x: i32,
    y: i32,
    z: i32,
}

let origin = Point{ x: 0, y: 0, z: 0 };

match origin {
    Point{ x, .. } => println!("x is {}", x),
    _ => (),
}

//
fn main() {
    let numbers = (2, 4, 6, 8, 10);
    
    match numbers {
        (first, .., last) => {
            println!("some numbers: {first}, {last}");
        }
    }
}
```


匹配守卫提供的额外条件

匹配守卫 match guard.
用于表达比单独的模式所能允许的更为复杂的情况

```rust
let num = Some(4);

match num {
    Some(x) if x % 2 == 0 => println!("The number {} is even", x),
    Some(x) => println!("The number {} is odd", x),
    None => (),
}
```



`@` 绑定

```rust
enum Message {
    Hello{ id: i32 },
}

let msg = Message::Hello{ id: 5 };

match msg {
    Message::Hello{
        // 通过 @ 将具体的 id 绑定到 id_variable，从而可以在后续代码中使用
        // 这里一般直接绑定到 id
        id: id_variable @ 3..=7,
    } => println!("Found an id in range {}", id_variable),
    Message::Hello{
        id: 10..=12,
    } => println!("Found an id in another range"),
    Message::Hello{ id } => println!("Found some other id: {}", id),
}
```



## Advanced Features

- 不安全 Rust ：用于当需要舍弃 Rust 的某些保证并负责手动维持这些保证
- 高级 Trait ：
- 高级类型
- 高级函数和闭包
- 宏



### 解引用裸指针

裸指针是不可变和可变的： `* const T` 和 `*mut T`

裸指针与引用和智能指针的区别在与：

- 允许忽略借用规则，可以同时拥有不可变和可变的指针，或多个指向相同位置的可变指针
- 不保证指向有效的内存
- 允许为空
- 不能实现任何自动清理功能

```rust
let mut num = 5;

let r1 = &num as *const i32;
let r2 = &mut num as *mut i32;

unsafe {
    println!("r1 is {}", *r1);
    println!("r2 is {}", *r2);
}
```

```rust
let address = 0x012345usize;
let r = address as *const i32;
```


调用不安全函数

```rust
unsafe fn dangerous() {}

unsafe {
    dangerous();
}
```


创建不安全代码的安全抽象

```rust
let mut v = vec![1,2,3,4,5,6];

let r = &mut v[..];

let (a, b) = r.split_at_mut(3);

assert_eq!(a, &mut [1,2,3]);
assert_eq!(b, &mut [4,5,6]);
```


访问或修改可变静态变量

全局变量被称为 静态 staic 变量。

```rust
static HELLO_WORLD: &str = "Hello, World.";

fn main() {
    println!("name is: {}", HELLO_WORLD);
}
```


读取和修改一个可变静态变量是不安全的

```rust
static mut COUNTER: u32 = 0;

fn add_to_count(inc: u32) {
    unsafe {
        COUNTER += inc;
    }
}

fn main() {
    add_to_count(3);
    
    unsafe {
        println!("COUNTER: {}", COUNTER);
    }
}
```


实现不安全 trait

```rust
unsafe trait Foo {}

unsafe impl Foo for i32 {}
```

访问联合体的字段


### Advanced Trait

```rust
pub trait Iterator {
    type Item;
    
    fn next(&mut self) -> Option<Self::Item>;
}

impl Iterator for Counter {
    type Item = u32;
    
    fn next(&mut self) -> Option<Self::Item>
}
```


默认泛型类型参数和运算符重载

```rust
use std::ops::Add;

#[derive(Debug, Copy, Clone, PartialEq)]
struct Point {
    x: i32,
    y: i32,
}

impl Add for Point {
    type Output = Point;
    
    fn add(self, other: Point) -> Point {
        Point{
            x: self.x + other.x,
            y: self.y + other.y,
        }
    }
}

fn main() {
    assert_eq!(
        Point{x:1, y:0} + Point{x:2, y: 3},
        Point{x: 3, y: 3}
    )
}
```

定义

```rust
trait Add<Rhs=Self> {
    type Output;
    
    fn add(self, rhs: Rhs) -> Self::Output;
}
```

Rhs right hand side 泛型类型参数。
`Rhs=Self` 默认类型参数 default type parameters.


下面是一个实现 `add trait` 时希望自定义 `Rhs` 类型而不是使用默认类型的例子。

```rust
use std::ops::Add;

struct Millimeters(u32);
struct Meters(u32);

// 在 Millimeters 上实现 Add, 能够将 Millimeters 与 Meters 相加
impl Add<Meters> for Millimeters {
    type Output = Millimeters;
    
    fn add(self, other: Meters) -> Millimeters {
        Millimeters(self.0 + (other.0 * 1000))
    }
}
```



完全限定语法与消歧义：调用相同名称的方法

```rust
trait Pilot {
    fn fly(&self);
}

trait Wizard {
    fn fly(&self);
}

struct Human;

impl Pilot for Human {
    fn fly(&self) {
        println!("This is your captain speaking");
    }
}

impl Wizard for Human {
    fn fly(&self) {
        println!("up");
    }
}

impl Human {
    fn fly(&self) {
        println!("*waving arms furiously*");
    }
}

fn main() {
    let person = Human;
    person.fly();      // 这里会默认调用直接实现在类型上的方法
    
    Pilot::fly(&person);  // 显式调用
    Wizard::fly(&person); 
}
```

然而，不是方法的关联函数没有 `self` 参数。
当存在多个类型或者 trait 定义了相同函数名的非方法函数时，Rust 就不总是能计算出我们期望的是哪一个类型，除非使用 **完全限定语法** fully qualified syntax.

```rust
trait Animal {
    fn baby_name() -> String;
}

struct Dog;

impl Dog {
    fn baby_name() -> String {
        string::from("Spot")
    }
}

impl Animal for Dog {
    fn baby_name() -> String {
        string::from("Puppy")
    }
}

fn main() {
    // 这里还是会默认调用直接实现的，也就是 Spot
    println!("A baby dog is called a {}", Dog::baby_name());
    // 因为 baby_name 没有 self 参数，所以不能 Animal::baby_name(&dog);
    // 而 Animal::baby_name() 又不能确定是要调用 Dog
    // 这里就需要 完全限定语法
    println!("A baby dog is called a {}", <Dog as Animal>::baby_name());
}
```

通常，完全限定语法为 `<Type as Trait>::function(receiver_if_method, next_arg, ...);`

对于不是方法的关联函数，其没有一个 `receiver` ，故而只有其他参数的列表。



父 trait 用于在另一个 trait 中使用某 trait 的功能

有时我们可能需要编写一个依赖另一个 trait 的 trait 定义：
对于一个实现了第一个 trait 的类型，你希望要求这个类型也实现了第二个 trait 。
如此就可使 trait 定义使用第二个 trait 的关联项。
这个所需的 trait 是我们实现的 trait 的父 trait superTrait.


例如我们希望实现一个带有 `outline_print` 方法的 trait `OutlinePrint` ，它会将给定的值格式化为带有星号框。
也就是说，给定一个实现了标准库 Display trait 并返回 `(x,y)` 的 Point，调用 `outline_print` 会输出：

```
**********
* (1, 3) *
**********
```

在 `outline_print` 的实现的，因为希望能够使用 Display trait 的功能，需要说明 `OutlinePrint` 只能用于同时实现了 Display 并提供了 `OutlinePrint` 需要的功能的类型。可以通过在 trait 定义中指定 `OutlinePrint: Display` 来实现。
类似于为 trait 增加 trait bound。

...



newtype 模式用以在外部类型上实现外部 trait

例如，如果想在 `Vec<T>` 上实现 Display，孤儿规则会阻止我们直接这样做，因为 Display trait 和 `Vec<T>` 都定义在我们的 trait 以外。
可以包裹一下。

```rust
use std::fmt;

struct Wrapper(Vec<String>);

impl fmt::Display for Wrapper {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f, "[{}]", self.0.join(", "))
    }
}
```

缺点是，Wrapper 是个新类型，没有自己的方法；必须直接在 Wrapper 上实现 `Vec<T>` 的所有方法。
如果希望新类型拥有其内部类型的每一个方法，为封装类型实现 `Deref` trait 并返回其内部类型是一个解决方案。
如果不希望封装类型用于所有内部类型的方法，则必须自己实现所需的方法。



### Advanced Types

类型别名

```rust
// 创建 i32 的类型别名 Kilometers
// 类型别名意思是这不是个 新的，单独的 类型。
// Kilometers 类型的值会被完全当作 i32
type Kilometers = i32;
```


主要用途是减少重复。例如，会有这样很长的类型

```rust
Box<dyn Fn() + Send + 'static>
```

在函数签名或注解中每次写这么长，即枯燥且容易出错，例如：

```rust
let f: Box<dyn Fn() + Send + 'static> = Box::new(|| println!("hi"));

fn takes_long_type(f: Box<dyn Fn() + Send + 'static>) {}

fn returns_long_type() -> Box<dyn Fn() + Send + 'static> {}
```

我们可以为这个冗长的类型引入一个别名，比如 `Thunk` 

```rust
type Thunk = Box<dyn Fn() + Send + 'static>;

let f: Thunk = Box::new(|| println!("hi"));

fn takes_long_type(f: Thunk){}

fn returns_long_type() -> Thunk {}
```


类型别名也经常与 `Result<T, E>` 结合使用来减少重复。
例如标准库的 `std::io` ，IO 操作通常返回一个 `Result<T, E>` ，因为操作可能会失败。
标准库中的 `std::io::Error` 结构体代表了所有可能的 IO 错误。
所以 `Result<..., Error` 可以有个类型别名

```rust
type Result<T> = std::result::Result<T, std::io::Error>;
```



从不返回的 never type

Rust 有一个叫做 `!` 的特殊类型。在类型理论术语中，称为 empty type。
我们更倾向于称为 never type。在函数从不返回的时候充当返回值。

```rust
fn bar() -> ! {}
```

从不返回的函数称为 发散函数 diverging functions。

```rust
// e.g.
let guess: u32 = match guess.trim().parse() {
    Ok(num) => num,
    Err(_) => continue,
};
// match 的分支必须返回相同的类型，而这里是 u32 和 continue
// continue 的返回值就是 !
```

never type 的另一个用途是 `panic!` 。

```rust
impl<T> Option<T> {
    pub fn unwrap(self) -> T {
        match self {
            Some(val) => val,
            None => panic!("called `Option::unwrap()` on a `None` value"),
        }
    }
}
```

这里的 `panic!` 返回值也是 `!` 类型。


最后一个有 `!` 类型的是 loop

```rust
print!("forever ");

loop {
    print!("and ever");
}
```



动态大小类型和 Sized trait

动态大小类型 dynamically sized types。
有时称为 DST 或 unsized types。
这种类型允许我们只有在运行时才知道大小的类型。

比如 `&str`。
这是 Rust 中动态大小类型的常规用法：它们有一些额外的元信息来储存动态信息的大小。
这引出了动态大小类型的黄金法则：必须将动态大小类型的值置于某种指针之后。

为了处理 DST，Rust 提供了 Sized trait 来决定一个类型的大小是否在编译时可知。
这个 trait 自动为编译器在编译时就知道大小的类型实现。
此外，Rust 隐式的为每个泛型函数增加了 Sized bound。
也就是说，对于以下泛型：

```rust
fn generic<T>(t: T){}
```

实际上被当作如下处理：

```rust
fn generic<T: Sized>(t: T){}
```

泛型函数默认只能用于在编译时已知大小的类型。然而可以用如下特殊语法来放宽这个限制：

```rust
fn generic<T: ?Sized>(t: &T){}
```

`?Sized` 意味着 T 有可能不是 Sized。
巴拉巴拉，复杂且难理解。



### Advanced functions and closures

函数满足 `fn` ，不要与闭包 trait 的 `Fn` 相混淆。
`fn` 称为 函数指针 function pointer。
通过函数指针允许我们使用函数作为另一个函数的参数。

```rust
fn add_one(x: i32) -> i32 {
    x + 1
}

fn do_twice(f: fn(i32) -> i32, arg: i32) i32 {
    f(arg) + f(arg)
}

fn main() {
    let answer = do_twice(add_one, 5);
    
    println!("the answer is {}", answer);   // 12
}
```

函数指针实现了所有三个闭包 trait `Fn`, `FnMut`，`FnOnce` ，所以总是可以在调用参数为闭包的函数时传入函数指针。
倾向于编写使用泛型和闭包 trait 的函数，这样它就能接收函数或闭包作为参数。

一个只期望接收 `fn` 而不接收闭包的例子是与不存在闭包的外部代码交互时：C 代码函数可以接收函数，但 C 里没闭包。

一个既可以用内联定义的闭包又可以用命名函数的例如，`map`

```rust
let list_of_numbers = vec![1,2,3];
let list_of_strings: Vec<String> =
    list_of_numbers.iter().map(|i| i.to_string()).collect();
```

或者将函数作为 `map` 的参数来代替闭包：

```rust
let list_of_numbers = vec![1,2,3];
let list_of_strings: Vec<String> =
    list_of_numbers.iter().map(ToString::to_string).collect();
```


```rust
enum Status {
    Value(u32),
    Stop,
}

let list_of_statuses: Vec<Status> = (0u32..20).map(Status::Value).collect();
```

这里创建了 `Status::Value` 实例，通过 `map` 用范围的每一个 `u32` 值调用 `Status::Value` 的初始化函数。
一些人倾向于函数风格，一些人喜欢闭包。这两种形式最终会产生同样的代码.



返回闭包

闭包表现为 trait，这意味着不能直接返回闭包。对于大部分需要返回 trait 的情况，可以使用实现了期望返回的 trait 的具体类型来代替函数的返回值。但是这不能用于闭包，因为它们没有一个可返回的具体类型；例如不允许使用函数指针 fn 作为返回值类型。

```rust
// 这段代码不能编译
fn returns_closure() -> dyn Fn(i32) -> i32 {
    |x| x + 1
}
```

错误会指向 Sized trait，Rust 并不知道需要多少空间来储存闭包。
不过我们上一部分见过了：可以用 trait 对象

```rust
fn returns_closure() -> Box<dyn Fn(i32) -> i32> {
    Box::new(|x| x + 1)
}
```



### Macros

我们已经用过像是 `println!` 这样的宏了。
宏 macro 指的是 Rust 中一系列的功能：使用 `macro_rules!` 的声明 declarative 宏，和三种 过程 procedural 宏

- 自定义 `#[derive]` 宏在结构体和枚举上指定通过 `derive` 属性添加的代码
- 类属性 Attribute-like 宏定义可用于任意项的自定义属性
- 类函数宏看起来像函数不过作用于作为参数传递的 token


宏和函数的区别

本质上说，宏是一种为写其他代码而写代码的方式，即所谓的 元编程 meta programming。
所有的这些宏以 展开 的方式来生成比你手写出的更多的代码。

元编程对于减少大量编写和维护的代码是非常有用的，它也扮演了函数扮演的角色。但宏有一些函数没有的附加功能。

一个函数签名必须声明函数参数个数和类型。相比之下，宏能接收不同数量的参数。
而且，宏可以在编译器翻译代码前展开，例如，宏可以在给定类型上实现 trait，而函数不行，因为函数是运行时被调用，而 trait 需要在编译时实现。

实现宏不如实现函数的一面是 宏定义比函数定义复杂，因为你编写的是生成 Rust 代码的 Rust 代码。
由于这样的间接性，宏定义通常比函数定义更难阅读，理解及维护。

宏和函数的最后一个重要的区别是：在一个文件里调用宏 之前 必须定义它，或将其引入作用域，而函数可以在任何地方定义和调用。




使用 `macro_rules!` 的声明宏用于元编程

Rust 最常用的宏形式是 声明宏 declarative macros。有时也称为 macros by example，`macro_rules!` 宏或就是 macros。
核心概念是，声明宏允许我们编写一些类似 Rust match 表达式的代码。
宏也将一个值和包含相关代码的模式进行比较；此种情况下，该值是传递给宏的 Rust 源代码字面值，模式用于与该字面值进行比较，
每个模式的相关代码会替换传递给宏的代码。所有这一切发生在编译时。

```rust
// 这个宏用三个整数创建一个 vector
let v: Vec<u32> = vec![1, 2, 3];
```

也可以用 `vec!` 宏来构造两个整数的 vector 或五个字符串的 vector。
但无法用函数做相同的事，因为我们无法预先知道参数值的数量和类型

下面是个 `vec!` 稍微简化的定义

```rust
#[macro_export]
macro_rules! vec {
    ( $( $x:expr ), *) => {
        {
            let mut temp_vec = Vec::new();
            $(
                temp_vec.push($x);
            )*
            temp_vec
        }
    };
}
```

`#[macro_export]` 注解表明只要引入了定义该宏的 crate，该宏就应该是可用的。
如果没有该注解，这个宏不能被引入作用域。

接着用 `macro_rules!` 和宏名称开始定义宏，且定义的宏 **不带** `!` 。
名字后跟大括号表示宏定义体，在该例宏名称为 `vec`

`vec!` 宏的结构和 match 类似。
这个是分支模式 `( $( $x:expr ),* )` ，后跟 `=>` 以及和模式相关的代码块。

首先，一对 `()` 包含整个模式。
使用 `$` 在宏系统声明一个变量来包含匹配该模式的 Rust 代码。
`$` 符号明确表明这是个宏变量而不是普通 Rust 变量。
之后是一对 `()` ，捕获了符合括号内模式的值以在替代代码中使用。
`$()` 内则是 `$x:expr` ，匹配 Rust 的任意表达式，并将该表达式命名为 `$x` 。
`$()` 后的逗号说明一个可有可无的逗号分隔符可以出现在 `$()` 所匹配的代码后。
紧随其后的 `*` 说明该模式匹配 0 个或多个 `*` 之前的任何模式。

当以 `vec![1,2,3];` 调用时， `$x` 模式与三个表达式 `1` , `2` , `3` 进行了三次匹配。

匹配到模式中的 `$()` 的每一部分，都会在 `=>` 右边 `$()*` 里生成 `temp_vec.push($x)` ，生成几次取决于匹配了几次。
`$x` 与匹配的表达式替换。
当以 `vec![1,2,3];` 调用时，会替换成如下：

```rust
{
    let mut temp_vec = Vec::new();
    temp_vec.push(1);
    temp_vec.push(2);
    temp_vec.push(3);
    temp_vec
}
```



用于从属性生成代码的过程宏

第二种形式的宏称为 过程宏 procedural macros，因为它们更像函数，一种过程类型。
过程宏接收 Rust 代码作为输入，在这些代码上进行操作，然后产生另一些代码作为输出，而非像声明式宏那样匹配对应模式然后以另一部分代码替换当前代码。
有三种类型的过程宏 自定义派生 derive，类属性和类函数，不过它们的工作方式都类似。

创建过程宏时，其定义必须驻留在它们自己具有特殊 crate 类型的 crate 中。

```rust
// 一个定义过程宏的例子
use proc_macro;

#[some_attribute]
pub fn some_name(input: TokenStream) -> TokenStream {}
```

实例：

创建一个 `hello_macro` 的 crate，包含名为 `HelloMacro` 的 trait 和关联函数 `hello_macro` 。
不同于让用户为其每一个类型实现 `HelloMacro` trait。
我们会提供一个过程宏以便用户可以使用 `#[derive(HelloMacro)]` 注解其类型从而获得 `hello_macro` 函数的默认实现。

1. 创建一个库 crate
   `cargo new hello_macro --lib`
2. 定义 `HelloMacro` trait 及其关联函数
   ```rust
   pub trait HelloMacro {
       fn hello_macro();
   }
   ```
   
   此时，用户可以实现该 trait 以达到期望的功能：
   ```rust
   use hello_macro::HelloMacro;
   
   struct Pancakes;
   
   impl HelloMacro for Pancakes {
       fn hello_macro() {
           println!("Hello, Macro. My name is Pancakes.");
       }
   }
   
   fn main() {
       Pancakes::hello_macro();
   }
   ```
3. 定义过程宏。编写该文档时，过程宏必须定义在自己的 crate 内。该限制是否取消待定。
   惯例是：对于一个 `foo` 包，一个自定义的派生过程宏被称为 `foo_derive`。
   因此，在 `hello_macro` 项目中新建 crate `hello_macro_derive`
   `cargo new hello_macro_derive --lib`
   我们需要声明 `hello_macro_derive` crate 是过程宏 proc-macro crate。
   还需要 `syn` 和 `quote` crate 中的功能，需要将它们添加到依赖中。
   ```toml
   // hello_macro_derive/Cargo.toml
   [lib]
   proc-macro = true
   
   [dependencies]
   syn = "1.0"
   quote = "1.0"
   ```
   
   ```rust
   // hello_macro_derive/src/lib.rs
   use proc_macro::TokenStream;
   use quote::quote;
   use syn;
   
   #[proc_macro_derive(HelloMacro)]
   pub fn hello_macro_derive(input: TokenStream) -> TokenStream {
       // 
       let ast = syn::parse(input).unwrap();
       
       // Build the trait implementation
       impl_hello_macro(&ast);
   }
   ```
   
   注意我们将代码分为了 `hello_macro_derive` 和 `impl_hello_macro` 两个函数。
   前者负责解析 TokenStream ，后者负责转换语法树：这使得编写过程宏更方便。
   
*有点复杂了，看原文可能好点*

可以这样指定依赖： `hello_macro = { path = "../hello_macro" }`



## Last Project  -- web server

改善性能的方法：线程池。
这是改善 web server 吞吐量的方法之一。
其他的比如：

- `fork/join model`
- 单线程异步I/O模型 `single-threaded async I/O model`
- 多线程异步I/O模型 `multi-threaded async I/O model`
