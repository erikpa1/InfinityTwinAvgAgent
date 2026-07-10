package tools

// XRange is an iterator over all the numbers from 0 to the limit.
func XRange(limit int) <-chan int {
	chnl := make(chan int)
	go func() {
		for i := 0; i < limit; i++ {
			chnl <- i
		}

		// Ensure that at the end of the loop we close the channel!
		close(chnl)
	}()
	return chnl
}
