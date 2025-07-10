package pipeline

func Run(nums []int) int {
    squareCh := make(chan int)
    go func() {
        defer close(squareCh)
        for _, n := range nums {
            squareCh <- n * n
        }
    }()

    doubleCh := make(chan int)
    go func() {
        defer close(doubleCh)
        for n := range squareCh {
            doubleCh <- n * 2
        }
    }()

    sum := 0
    for n := range doubleCh {
        sum += n
    }

    return sum
}