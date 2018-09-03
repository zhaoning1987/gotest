package job

import pb "gopkg.in/cheggaaa/pb.v1"

func InitBar(count int) *pb.ProgressBar {
	bar := pb.StartNew(count)
	bar.ShowSpeed = true
	return bar
}

func IncrementBar(bar *pb.ProgressBar) {
	if bar != nil {
		bar.Increment()
	}
}

func CompleteBar(bar *pb.ProgressBar, msg string) {
	if bar != nil {
		bar.FinishPrint(msg)
	}
}
