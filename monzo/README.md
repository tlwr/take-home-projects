# monzo

A crawler written as a take home exercise for Monzo.com

## Usage

```
# build

make build
```

```
# get help

./monzo -help
```

```
# crawl my personal website

./monzo       \
  -url 'https://www.toby.codes' \
  -host '*.toby.codes'          \
  -host 'toby.codes'            \
  -host '*.tobys.cloud'

# note that -host is a glob
```

```
# crawler go brrr

./monzo  \
  -host '*.monzo.com'      \
  -host 'monzo.com'        \
  -url 'https://monzo.com' \
  -parallel 64
```

## Development

```
# run the tests
make test

# run the integration tests
make integration
```

```
# check code coverage

make coverage
ginkgo -r -cover -skipPackage integration
Will skip:
  ./integration
  [1612218142] Hostfilter Suite - 3/3 specs ••• SUCCESS! 150.375µs PASS
  coverage: 92.3% of statements
  [1612218142] LinkParser Suite - 7/7 specs ••••••• SUCCESS! 161.833µs PASS
  coverage: 100.0% of statements
  [1612218142] UrlDedupQueue Suite - 3/3 specs ••• SUCCESS! 5.005292ms PASS
  coverage: 100.0% of statements
  [1612218142] WebPageScraper Suite - 4/4 specs •••• SUCCESS! 380.375µs PASS
  coverage: 84.8% of statements

  Ginkgo ran 4 suites in 1.332776334s
  Test Suite Passed
```
