package job

import (
	"context"
	"encoding/base64"
	"strings"
	"sync"

	"github.com/pkg/errors"
	xlog "github.com/qiniu/xlog.v1"
	outer "qiniu.com/argus/dbstorage/outer_service"
	"qiniu.com/argus/dbstorage/proto"
	"qiniu.com/argus/dbstorage/util"
)

type FaceJob struct {
	index        int
	ctx          context.Context
	distributor  *Distributor
	faceGroup    outer.IFaceGroup
	imageContent []byte
	imageId      proto.ImageId
	imageURI     proto.ImageURI
	tag          proto.ImageTag
	desc         proto.ImageDesc
	process      proto.ImageProcess
	preCheckErr  error
	wg           *sync.WaitGroup
}

func NewFaceJob(ctx context.Context, index int, distributor *Distributor, faceGroup outer.IFaceGroup,
	imageContent []byte, imageId proto.ImageId, imageURI proto.ImageURI, tag proto.ImageTag, desc proto.ImageDesc, process proto.ImageProcess, preCheckErr error, wg *sync.WaitGroup) *FaceJob {
	return &FaceJob{
		index:        index,
		ctx:          ctx,
		distributor:  distributor,
		faceGroup:    faceGroup,
		imageContent: imageContent,
		imageId:      imageId,
		imageURI:     imageURI,
		tag:          tag,
		desc:         desc,
		process:      process,
		preCheckErr:  preCheckErr,
		wg:           wg,
	}
}

func (fj *FaceJob) execute(workerIndex int) {
	defer func() {
		fj.wg.Done()
		IncrementBar(fj.distributor.bar)
	}()

	fj.distributor.Lock()
	if fj.distributor.cancelled {
		fj.distributor.Unlock()
		return
	}
	fj.distributor.Unlock()

	xl := xlog.FromContextSafe(fj.ctx)

	handled := true
	defer func() {
		if handled {
			if err := fj.distributor.UpdateCount(); err != nil {
				xl.Errorf("distributor.UpdateCount fail: %s", err)
			}
		}
	}()

	//when start to handle a face, we save its index
	if err := fj.distributor.UpdateProcess(workerIndex, fj.index); err != nil {
		xl.Errorf("distributor.UpdateProcess fail: %s", err)
		return
	}

	if fj.preCheckErr != nil {
		fj.distributor.UpdateErrorLog(string(fj.imageURI), fj.preCheckErr)
		return
	}

	var err error

	hash := string(fj.imageURI)
	if fj.imageContent != nil {
		hash = util.GetSha1(fj.imageContent)
	}

	//check if this hash exist
	existed := fj.distributor.IsHashExist(hash)
	if existed {
		if fj.process == proto.NOT_HANDLED {
			fj.distributor.UpdateErrorLog(string(fj.imageURI), proto.ErrDupImage)
		}
		return
	}

	defer func() {
		if handled {
			if err := fj.distributor.UpdateHash(hash); err != nil {
				xl.Errorf("distributor.UpdateHash fail: %s", err)
			}
		}
	}()

	//call group_add service to store the image
	url := fj.imageURI
	if fj.imageContent != nil {
		imgBase64 := proto.BASE64_PREFIX + base64.StdEncoding.EncodeToString(fj.imageContent)
		url = proto.ImageURI(imgBase64)
	}
	_, err = fj.faceGroup.Add(fj.ctx, fj.distributor.task.Config, fj.distributor.task.GroupName, fj.imageId, url, fj.tag, fj.desc)

	if err == nil {
		//success
		handled = true
	} else if err == context.Canceled {
		//canceled by user
		handled = false
	} else if strings.Contains(err.Error(), proto.ErrMsgGroupNotExistFG) || strings.Contains(err.Error(), proto.ErrMsgGroupNotExistFGP) {
		//fatal error, group not exists, stop immediately
		err2 := errors.Errorf("group %s not exist", fj.distributor.task.GroupName)
		xl.Error(err2)
		handled = false
		_ = fj.distributor.UpdateLastError(fj.ctx, fj.distributor.task.TaskId, proto.TaskError(err2.Error()))
		fj.distributor.Stop()
	} else if strings.Contains(err.Error(), proto.ErrMsgFeatureExistFG) || strings.Contains(err.Error(), proto.ErrMsgFeatureExistFGP) {
		//duplicated
		xl.Error("id already exists", fj.imageId, fj.imageURI)
		fj.distributor.UpdateErrorLog(string(fj.imageURI), proto.ErrIdExist)
		handled = true
	} else {
		//real error, check it
		handled = true
		if strings.Contains(err.Error(), proto.ErrMsgNoFaceFG) || strings.Contains(err.Error(), proto.ErrMsgNoFaceFGP) {
			//no face detected
			xl.Errorf("no face detected in image : %s", fj.imageURI)
			fj.distributor.UpdateErrorLog(string(fj.imageURI), proto.ErrNoFace)
		} else if strings.Contains(err.Error(), proto.ErrMsgMultiFaceFG) || strings.Contains(err.Error(), proto.ErrMsgMultiFaceFGP) {
			//multi face detected
			xl.Errorf("multi face detected in image : %s", fj.imageURI)
			fj.distributor.UpdateErrorLog(string(fj.imageURI), proto.ErrMultiFace)
		} else if strings.Contains(err.Error(), proto.ErrMsgFaceSizeFG) {
			//face too small
			xl.Errorf("face too small in image : %s", fj.imageURI)
			fj.distributor.UpdateErrorLog(string(fj.imageURI), proto.ErrFaceTooSmall)
		} else {
			xl.Errorf("fail to add %s to group %s : %s", fj.imageURI, fj.distributor.task.GroupName, err)
			fj.distributor.UpdateErrorLog(string(fj.imageURI), err)
		}
	}
}
