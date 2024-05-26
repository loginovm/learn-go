package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if done != nil {
		return applyStagesWithCancellation(in, done, stages)
	}

	return applyStages(in, stages)
}

func applyStages(in In, stages []Stage) Out {
	if len(stages) == 0 {
		return in
	}

	return applyStages(stages[0](in), stages[1:])
}

func applyStagesWithCancellation(in In, done In, stages []Stage) Out {
	in2 := make(Bi)

	go func() {
		defer close(in2)

		for {
			select {
			case m, ok := <-in:
				if !ok {
					return
				}
				in2 <- m
			case <-done:
				return
			}
		}
	}()

	if len(stages) == 0 {
		return in2
	}

	return applyStagesWithCancellation(stages[0](in2), done, stages[1:])
}
