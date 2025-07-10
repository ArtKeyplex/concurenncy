package pipelinectx

import "context"


func Run(ctx context.Context, nums []int) (int, error) {
    doubleCh := make(chan int)
    go func() {
        defer close(doubleCh)
        for _, n := range nums {
            select {
            case <-ctx.Done():
                return
            case doubleCh <- n * 2:
            }
        }
    }()

    sum := 0
    for n := range doubleCh {
        select {
        case <-ctx.Done():
            return 0, ctx.Err()
        default:
            sum += n
        }
    }

    return sum, nil
}
