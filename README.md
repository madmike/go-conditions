# Smart Conditions Library for Go

The `smart-conditions` library is a powerful and flexible tool designed to evaluate conditions and logical expressions dynamically in Go applications. It supports a wide range of operators for simple, common, and logical conditions, enabling complex decision-making scenarios based on runtime data.

## Features

- **Simple Operators:** Includes checks for null, defined, undefined, exist, empty, blank, truly, and falsy values.
- **Common Operators:** Supports common comparison and matching operations, such as equals, not equals, less than, greater than, regex match, in list, starts with, ends with, and more.
- **Logical Operators:** Facilitates complex logical expressions using OR, XOR, AND, and NOT operators.
- **Dynamic Value Evaluation:** Dynamically evaluates conditions against runtime data, allowing for flexible and powerful logic constructs within your Go applications.

## Installation

To install the `go-conditions` library, use the following `go get` command:

```bash
go get -u github.com/madmike/go-conditions
```

## Usage

### Creating a New Conditions Object

Start by creating a new instance of `Conditions`:

```go
import "github.com/madmike/go-conditions"

cond := conditions.NewConditions()
```

### Evaluating Conditions

You can evaluate conditions using the Check method. This method takes two parameters: the instance to check against and the condition to evaluate.

```go
instance := map[string]any{
    "name": "John Doe",
    "age": 30,
}

condition := map[string]any{
    "$gt": map[string]any{
        "age": 18,
    },
}

result := cond.Check(instance, condition)
fmt.Println(result) // Output: true
```

## Supported Operators

### Simple Operators

- **NULL**: `$null`
- **DEFINED**: `$defined`
- **UNDEFINED**: `$undefined`
- **EXIST**: `$exist`
- **EMPTY**: `$empty`
- **BLANK**: `$blank`
- **TRULY**: `$truly`
- **FALSY**: `$falsy`

### Common Operators

- **EQ (Equal)**: `$eq`
- **NE (Not Equal)**: `$ne`
- **LT (Less Than)**: `$lt`
- **GT (Greater Than)**: `$gt`
- **LTE (Less Than or Equal To)**: `$lte`
- **GTE (Greater Than or Equal To)**: `$gte`
- **RE (Regex)**: `$re`
- **IN (In List)**: `$in`
- **NI (Not In List)**: `$ni`
- **SW (Starts With)**: `$sw`
- **EW (Ends With)**: `$ew`
- **INCL (Includes)**: `$incl`
- **EXCL (Excludes)**: `$excl`
- **HAS (Has Property)**: `$has`
- **POWER (Bitwise Power)**: `$power`
- **BETWEEN (Between)**: `$between`
- **SOME (Some)**: `$some`
- **EVERY (Every)**: `$every`
- **NOONE (No One)**: `$noone`

### Logical Operators

- **OR**: `$or`
- **XOR**: `$xor`
- **AND**: `$and`
- **NOT**: `$not`

## Testing and Coverage

Testing is an essential part of software development that ensures your application behaves as expected. The `go-conditions` library supports extensive testing capabilities to help maintain high-quality code.

### Running Tests

To run tests, execute the following command in the terminal at the root of your project:

```bash
go test .
```

This command will run all tests in the current project and its subdirectories. To see detailed output for each test, including successful tests, use the -v (verbose) flag:

```bash
go test -v .
```

### Viewing Test Coverage

Go provides built-in support for analyzing test coverage. To generate a test coverage report, run:

```bash
go test -coverprofile=coverage.out .
```

Then, to view the coverage report in your terminal, use:

```bash
go tool cover -func=coverage.out
```

For a more visual representation, generate an HTML coverage report:

```bash
go tool cover -html=coverage.out
```

This command opens the generated report in your default web browser, showing which lines of code were executed during tests.

### Continuous Improvement

As you add new features or discover bugs, writing tests should be a parallel effort to coding. Aim to improve test coverage iteratively, especially in areas identified as under-tested in your coverage reports.

## Contributing

We welcome contributions to the smart-conditions library! Please feel free to submit issues, pull requests, or enhancements to improve the library's functionality and usability.

## License

This library is licensed under the MIT License. Feel free to use it, modify it, and distribute it as you see fit.
