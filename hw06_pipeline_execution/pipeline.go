package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	for _, stage := range stages {
		in = applyStage(in, done, stage)
	}

	return in
}

func applyStage(in In, done In, stage Stage) Out {
	if done == nil {
		return stage(in)
	}

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

	return stage(in2)
}
