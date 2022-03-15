## Fast or First routine (FFR)

### Situation:

there are times when you need to get information and write this value to some buffer within various streams.

### for example:
it is necessary to obtain information from the first worked out correction search function and this functionality will help.


### use:

```go

func main() {
	fg := firstroutine.New()
	
	var adj any

    fg.Go(func(m *Mutable) error {
		adjFromFirstApi, err := //// example api 
		if err != nil {
		    return err	
        }   
		
        m.Set(&adj, adjFromFirstApi)
        return nil
    })

    fg.Go(func(m *Mutable) error {
        adjFromSecondApi, err := //// example api
        if err != nil {
            return err
        }
        
        m.Set(&adj, adjFromSecondApi)
        return nil
    })
	
	if err := fg.Wait(); err != nil {
		panic(err) 
    }

	// if first api worked first
	// adj = adjFromFirstApi
	// if second api worked first
	// adj = adjFromSecondApi
	fmt.Println(adj)
}

```
