# Contributing to XQB

Thank you for your interest in contributing to **XQB**, a powerful and flexible Sql query builder for Go. Whether you're fixing a bug, improving documentation, or adding a new feature — your contribution is highly appreciated!

---

## 📋 Before You Start

* Familiarize yourself with the [README](./README.md).
* Ensure Go is installed and your environment is set up.
* Read through the codebase and existing issues to avoid duplication.

---

## 🛠 How to Contribute

### 1. Fork the Repository

```bash
git clone https://github.com/your-username/xqb.git
cd xqb
```

### 2. Create a Feature or Fix Branch

```bash
git checkout -b feature/short-description
```

### 3. Make Your Changes

* Write clean, idiomatic Go.
* Keep the fluent query-building API consistent.
* Add unit tests for new logic.

### 4. Run Tests

```bash
go test ./...
```

> All tests must pass before submitting a pull request.

### 5. Format Your Code

```bash
go fmt ./...
```

### 6. Push and Open a Pull Request

```bash
git push origin feature/short-description
```

Then open a PR against the `main` branch with a clear title and description.

---

## 🔪 Tests

All contributions must include tests under `xqb_test.go` or relevant test files.

Test any of the following if affected:

* Sql compilation logic
* Parameter bindings
* Transactions
* Connection behavior
* Fluent chaining

---

## 🛉 Style Guide

* Use `go fmt` to auto-format.
* Name functions and variables clearly and consistently.
* Avoid breaking changes unless necessary — backward compatibility matters.

---

## 🐛 Bug Reports

If you find a bug, open an issue with:

* A clear title and description.
* Steps to reproduce.
* Expected vs actual behavior.
* Relevant code snippets or Sql output.

---

## 💡 Feature Requests

* Describe the use case.
* Include an example query or fluent chain.
* Mention if this feature exists in Laravel’s query builder.

---

## 💬 Discussions

Use GitHub Discussions or open an issue if you're unsure how to implement something or want feedback before writing code.

---

## 📄 License

By contributing, you agree that your contributions will be licensed under the [MIT License](./LICENSE).

---

## 🙏 Thanks

Thanks for helping improve **XQB**! Your time and ideas make this project better.
