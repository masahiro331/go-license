# go-license
License Parser for License Expression

The library parse license expression like the bellow string.

> Public Domain AND ( GPLv2+ or AFL ) AND LGPLv2+ with distribution exceptions

Parse to like bellow structure
```

"Public Domain" AND 
( 
    "GPLV2+" OR 
    "AFL" 
) AND
"LGPLv2+ with distribution exceptions"
```


## Quick Start

```
func main() {
	input := "Public Domain AND (GPLv2+ or AFL) and LGPLv2+ with distribution exceptions"
	l := lexer.New(input)
	p := New(l)
	
	expression, err := p.Parse()
	fmt.Println(p.Normalize(expression))
}
```

Output: 
```
Public Domain AND ( GPLv2+ OR AFL ) AND LGPLv2+ with distribution exceptions
```


### normilize callback 
You can register callback function when normalize license expression.

```
func main() {
	input := "Public Domain AND ( GPLv2+ or AFL ) AND LGPLv2+ with distribution exceptions"
	l := lexer.New(input)
	p := New(l).RegisterNormalizeFunc(func(n string) string {
	    return strings.Replace(n, " ", "_", -1)
	})
	
	expression, err := p.Parse()
	fmt.Println(p.Normalize(expression))
}
```

Output:
```
Public_Domain AND ( GPLv2+ OR AFL ) AND LGPLv2+_with_distribution_exceptions
```
