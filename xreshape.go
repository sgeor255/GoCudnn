package gocudnn

import "errors"

//Reshape is failing ---- changed all to private
//Probably have to move to CPU
//FindSegmentedOutputTensor creates a tensordescriptor for the segmeented size
func (xt Xtra) findSegmentedOutputTensor(descX *TensorD, h, w int32) (*TensorD, error) {
	dtype, dims, _, err := descX.GetDescrptor()
	if err != nil {
		return nil, err
	}
	if dims[0] > 1 {
		return nil, errors.New("The N dim shouldn't be larger than 1")
	}
	var frmt TensorFormatFlag
	n1 := int32(divideandroundup(dims[2], h))
	n2 := int32(divideandroundup(dims[3], w))
	return Tensor{}.NewTensor4dDescriptor(dtype, frmt.NCHW(), []int32{n1 * n2, dims[1], h, w})
}

//SegmentedBatches1CHWtoNCHWForward only works on a 1CHW to NCHW where the batches are the size of the window. If the segmented
func (xt Xtra) segmentedBatches1CHWtoNCHWForward(handle *XHandle, xDesc *TensorD, x *Malloced, yDesc *TensorD, y *Malloced) error {
	datatypex, dimsx, _, err := xDesc.GetDescrptor()
	if err != nil {
		return err
	}
	datatypey, dimsy, _, err := xDesc.GetDescrptor()
	if err != nil {
		return err
	}
	n1 := divideandroundup(dimsx[2], dimsy[2])

	n2 := divideandroundup(dimsx[3], dimsy[3])
	n3 := n1 * n2

	//originalarea := dimsx[2] * dimsx[3]
	var dtflg DataTypeFlag
	if datatypex != datatypey || datatypex != dtflg.Float() {
		return errors.New("Datatypes Don't Match or Datatype is not float")
	}
	if uint32(dimsy[0]) != n3 {
		return errors.New("SegmentedBatches Tensors don't match up")

	}
	if dimsx[1] != dimsy[1] {
		return errors.New("Channels Need to be same size")
	}
	var cu Cuda
	kern, err := cu.MakeKernel("NCHWsegmentfrom1CHWfloat", handle.mod)
	if err != nil {
		return err
	}
	blocksize := int32(16)
	bs := uint32(blocksize)
	gx := divideandroundup(dimsy[2], blocksize)
	gy := divideandroundup(dimsy[3], blocksize)
	//gz := divideandroundup(dimsy[3], blocksize)
	//for channel := int32(0); channel < dimsx[1]; channel++ {
	//	for ah, h := uint32(0), int32(0); ah < n1; ah, h = ah+1, h+dimsy[2] {

	//		for aw, w := uint32(0), int32(0); aw < n2; aw, w = aw+1, w+dimsy[3] {

	//			index := ah*n2 + aw
	err = kern.Launch(gx, gy, 1, bs, bs, 1, 0, handle.s, int32(dimsx[1]), int32(dimsx[2]), int32(dimsx[3]), x, y)

	if err != nil {
		return err
	}

	//	}

	/*
	   const int BatchIndex,
	   const int ChannelIndex,
	   const int ChannelLength,
	   const int OriginalStartX,
	   const int OriginalStartY,
	   const int OriginalSizeX,
	   const int OriginalSizeY,
	   float *oMem,
	   float *nMem){
	*/

	//	}

	//}

	return nil
}

//SegmentedBatches1CHWtoNCHWBackward only works on a 1CHW to NCHW where the batches are the size of the window. If the segmented
func (xt Xtra) segmentedBatches1CHWtoNCHWBackward(handle *XHandle, xDesc *TensorD, x *Malloced, yDesc *TensorD, y *Malloced) error {
	datatypex, dimsx, _, err := xDesc.GetDescrptor()
	if err != nil {
		return err
	}
	datatypey, dimsy, _, err := xDesc.GetDescrptor()
	if err != nil {
		return err
	}
	n1 := divideandroundup(dimsx[2], dimsy[2])

	n2 := divideandroundup(dimsx[3], dimsy[3])
	var dtflg DataTypeFlag
	if datatypex != datatypey || datatypex != dtflg.Float() {
		return errors.New("Datatypes Don't Match or Datatype is not float")
	}
	if uint32(dimsy[0]) != n1*n2 {
		return errors.New("SegmentedBatches Tensors don't match up")

	}
	if dimsx[1] != dimsy[1] {
		return errors.New("Channels Need to be same size")
	}
	var cu Cuda
	kern, err := cu.MakeKernel("CHWfromSegmentedNCHWfloat", handle.mod)
	if err != nil {
		return err
	}
	blocksize := int32(4)
	bs := uint32(blocksize)
	gx := divideandroundup(dimsy[1], blocksize)
	gy := divideandroundup(dimsy[2], blocksize)
	gz := divideandroundup(dimsy[3], blocksize)
	OriginalTotalVol := findvol(dimsx)
	for i := uint32(0); i < n1; i++ {
		for k := uint32(0); k < n2; k++ {
			index := int32(i*n2 + k)
			err = kern.Launch(gx, gy, gz, bs, bs, bs, 0, handle.s, index, 1, int32(i), int32(k), 1, n1, n2, OriginalTotalVol, x.ptr, y.ptr)
			if err != nil {
				return err
			}
		}

	}

	return nil
}
func findvol(dims []int32) int32 {
	mult := int32(1)
	for i := 0; i < len(dims); i++ {
		mult *= dims[i]
	}
	return mult
}
func divideandroundup(den, num int32) uint32 {

	return uint32(((den - 1) / num) + 1)

}
