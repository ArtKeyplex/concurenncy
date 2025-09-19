package debounce

import "time"

// Debounce принимает значения и отдаёт только последнее после паузы d.
func Debounce(d time.Duration, in <-chan int) <-chan int {
	// TODO: реализовать дебаунс значений из канала

	out := make(chan int)

	go func() {
		defer close(out)

		var timer *time.Timer
		var lastValue int
		var pending bool

		for {
			select {
			case value, ok := <-in:
				if !ok {
					if pending {
						out <- lastValue
					}
					return
				}

				lastValue = value
				pending = true

				if timer != nil {
					timer.Stop()
				}
				timer = time.AfterFunc(d, func() {
					if pending {
						out <- lastValue
						pending = false
					}
				})
			default:
				time.Sleep(time.Microsecond)
			}
		}
	}()
	return out
}
