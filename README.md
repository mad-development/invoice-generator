# invoice-generator


## Description

Basic execution:

```bash
go run main.go -day=18 \
               -month=7 \
               -year=2023 \
               -client="John Doe" \
               -id="INV1234" \
               -services="Service 1:100.0,Service 2:75.0,Service 3:25.5" \
               -originName="My Company" \
               -originAddress="123 Main St, City, Country"
```
