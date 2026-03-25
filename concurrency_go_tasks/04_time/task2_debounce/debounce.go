package debounce

import "time"

// Debounce принимает значения и отдаёт только последнее после паузы d.
func Debounce(d time.Duration, in <-chan int) <-chan int {
	// TODO: реализовать дебаунс значений из канала
	out := make(chan int)
	go func() {
		defer close(out)

		var (
			timer *time.Timer
			last int
		)
		
		for {
			var timerCh <-chan time.Time
			if timer != nil {
				timerCh = timer.C
			}
			select {
			case v, ok := <-in:
				if !ok {
					if timer != nil {
						<-timer.C
						out <- last
					}
					return
				}

				last = v

				if timer == nil {
					timer = time.NewTimer(d)
				} else {
					if !timer.Stop(){
						<-timer.C
					}
					timer.Reset(d)
				}
			case <-timerCh:
				out <- last
				timer = nil
			}
		}
	}()
	return out
}
