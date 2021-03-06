package gocudnn

/*
#include <cudnn.h>


void MakeAlgorithmforBWDFilter(cudnnAlgorithm_t *input,cudnnConvolutionBwdFilterAlgo_t algo ){
	input->algo.convBwdFilterAlgo=algo;
}

*/
import "C"
import (
	"fmt"
	"unsafe"

	"github.com/dereklstinson/cutil"
)

//Algo returns an Algorithm Struct
func (c ConvBwdFiltAlgo) Algo() Algorithm {
	var algorithm C.cudnnAlgorithm_t
	C.MakeAlgorithmforBWDFilter(&algorithm, c.c())
	return Algorithm(algorithm)
}

//GetBackwardFilterAlgorithmMaxCount returns the max number of Algorithm
func (c *ConvolutionD) getBackwardFilterAlgorithmMaxCount(handle *Handle) (int32, error) {

	var count C.int
	x := Status(C.cudnnGetConvolutionBackwardFilterAlgorithmMaxCount(handle.x, &count)).error("GetConvolutionForwardAlgorithmMaxCount")

	return int32(count), x

}

//FindBackwardFilterAlgorithm will find the top performing algoriths and return the best algorithms in accending order they are limited to the number passed in requestedAlgoCount.
//So if 4 is passed through in requestedAlgoCount, then it will return the top 4 performers in the ConvolutionFwdAlgoPerformance struct.  using this could possible give the user cheat level performance :-)
func (c *ConvolutionD) FindBackwardFilterAlgorithm(
	handle *Handle,
	xD *TensorD,
	dyD *TensorD,
	dwD *FilterD,
) ([]ConvBwdFiltAlgoPerformance, error) {
	requestedAlgoCount, err := c.getBackwardFilterAlgorithmMaxCount(handle)
	if err != nil {
		return nil, err
	}
	perfResults := make([]C.cudnnConvolutionBwdFilterAlgoPerf_t, requestedAlgoCount)
	var actualalgocount C.int
	err = Status(C.cudnnFindConvolutionBackwardFilterAlgorithm(
		handle.x,
		xD.descriptor,
		dyD.descriptor,
		c.descriptor,
		dwD.descriptor,
		C.int(requestedAlgoCount),
		&actualalgocount,
		&perfResults[0],
	)).error("FindConvolutionBackwardFilterAlgorithm")

	results := make([]ConvBwdFiltAlgoPerformance, int32(actualalgocount))
	for i := int32(0); i < int32(actualalgocount); i++ {
		results[i] = convertConvBwdFiltAlgoPerformance(perfResults[i])

	}
	return results, err
}

//FindBackwardFilterAlgorithmEx finds some algorithms with memory
func (c *ConvolutionD) FindBackwardFilterAlgorithmEx(
	handle *Handle,
	xD *TensorD, x cutil.Mem,
	dyD *TensorD, dy cutil.Mem,
	dwD *FilterD, dw cutil.Mem,
	wspace cutil.Mem, wspacesize uint) ([]ConvBwdFiltAlgoPerformance, error) {
	reqAlgoCount, err := c.getBackwardFilterAlgorithmMaxCount(handle)
	if err != nil {
		return nil, err
	}
	perfResults := make([]C.cudnnConvolutionBwdFilterAlgoPerf_t, reqAlgoCount)
	var actualalgocount C.int
	err = Status(C.cudnnFindConvolutionBackwardFilterAlgorithmEx(
		handle.x,
		xD.descriptor, x.Ptr(),
		dyD.descriptor, dy.Ptr(),
		c.descriptor,
		dwD.descriptor, dw.Ptr(),
		C.int(reqAlgoCount), &actualalgocount, &perfResults[0], wspace.Ptr(), C.size_t(wspacesize))).error("FindConvolutionBackwardFilterAlgorithmEx")

	results := make([]ConvBwdFiltAlgoPerformance, int32(actualalgocount))
	for i := int32(0); i < int32(actualalgocount); i++ {
		results[i] = convertConvBwdFiltAlgoPerformance(perfResults[i])

	}
	return results, err
}

//FindBackwardFilterAlgorithmExUS is just like FindBackwardFilterAlgorithmEx but uses unsafe.Pointer instead of cutil.Mem
func (c *ConvolutionD) FindBackwardFilterAlgorithmExUS(
	handle *Handle,
	xD *TensorD, x unsafe.Pointer,
	dyD *TensorD, dy unsafe.Pointer,
	dwD *FilterD, dw unsafe.Pointer,
	wspace unsafe.Pointer, wspacesize uint) ([]ConvBwdFiltAlgoPerformance, error) {
	reqAlgoCount, err := c.getBackwardFilterAlgorithmMaxCount(handle)
	if err != nil {
		return nil, err
	}
	perfResults := make([]C.cudnnConvolutionBwdFilterAlgoPerf_t, reqAlgoCount)
	var actualalgocount C.int
	err = Status(C.cudnnFindConvolutionBackwardFilterAlgorithmEx(
		handle.x,
		xD.descriptor, x,
		dyD.descriptor, dy,
		c.descriptor,
		dwD.descriptor, dw,
		C.int(reqAlgoCount), &actualalgocount, &perfResults[0], wspace, C.size_t(wspacesize))).error("FindConvolutionBackwardFilterAlgorithmEx")

	results := make([]ConvBwdFiltAlgoPerformance, int32(actualalgocount))
	for i := int32(0); i < int32(actualalgocount); i++ {
		results[i] = convertConvBwdFiltAlgoPerformance(perfResults[i])

	}
	return results, err
}

//GetBackwardFilterAlgorithmV7 will find the top performing algoriths and return the best algorithms in accending order they are limited to the number passed in requestedAlgoCount.
//So if 4 is passed through in requestedAlgoCount, then it will return the top 4 performers in the ConvolutionFwdAlgoPerformance struct.  using this could possible give the user cheat level performance :-)
func (c *ConvolutionD) GetBackwardFilterAlgorithmV7(
	handle *Handle,
	xD *TensorD,
	dyD *TensorD,
	dwD *FilterD,
) ([]ConvBwdFiltAlgoPerformance, error) {
	requestedAlgoCount, err := c.getBackwardFilterAlgorithmMaxCount(handle)
	if err != nil {
		return nil, err
	}
	perfResults := make([]C.cudnnConvolutionBwdFilterAlgoPerf_t, requestedAlgoCount)
	var actualalgocount C.int
	err = Status(C.cudnnGetConvolutionBackwardFilterAlgorithm_v7(
		handle.x,
		xD.descriptor,
		dyD.descriptor,
		c.descriptor,
		dwD.descriptor,
		C.int(requestedAlgoCount),
		&actualalgocount,
		&perfResults[0])).error("GetConvolutionBackwardFilterAlgorithm_v7")
	results := make([]ConvBwdFiltAlgoPerformance, int32(actualalgocount))

	for i := int32(0); i < int32(actualalgocount); i++ {
		results[i] = convertConvBwdFiltAlgoPerformance(perfResults[i])

	}
	return results, err
}

//GetBackwardFilterAlgorithm gives a good algo with the limits given to it
func (c *ConvolutionD) GetBackwardFilterAlgorithm(
	handle *Handle,
	xD *TensorD,
	dyD *TensorD,
	dwD *FilterD,
	pref ConvBwdFilterPref, wsmemlimit uint) (ConvBwdFiltAlgo, error) {
	var algo C.cudnnConvolutionBwdFilterAlgo_t
	err := Status(C.cudnnGetConvolutionBackwardFilterAlgorithm(
		handle.x,
		xD.descriptor,
		dyD.descriptor,
		c.descriptor,
		dwD.descriptor,
		pref.c(), C.size_t(wsmemlimit), &algo)).error("GetConvolutionBackwardFilterAlgorithm")

	return ConvBwdFiltAlgo(algo), err
}

func (c ConvBwdFiltAlgo) print() {
	switch c {
	case ConvBwdFiltAlgo(C.CUDNN_CONVOLUTION_BWD_FILTER_ALGO_0):
		fmt.Println("ConvBwdFiltAlgo0")
	case ConvBwdFiltAlgo(C.CUDNN_CONVOLUTION_BWD_FILTER_ALGO_1):
		fmt.Println("ConvBwdFiltAlgo1")
	case ConvBwdFiltAlgo(C.CUDNN_CONVOLUTION_BWD_FILTER_ALGO_FFT):
		fmt.Println("ConvBwdFiltAlgoFFT")
	case ConvBwdFiltAlgo(C.CUDNN_CONVOLUTION_BWD_FILTER_ALGO_3):
		fmt.Println("ConvBwdFiltAlgo3")
	case ConvBwdFiltAlgo(C.CUDNN_CONVOLUTION_BWD_FILTER_ALGO_WINOGRAD):
		fmt.Println("ConvBwdFiltAlgoWinGrad")
	case ConvBwdFiltAlgo(C.CUDNN_CONVOLUTION_BWD_FILTER_ALGO_WINOGRAD_NONFUSED):
		fmt.Println("ConvBwdFiltAlgoNonFused")
	case ConvBwdFiltAlgo(C.CUDNN_CONVOLUTION_BWD_FILTER_ALGO_FFT_TILING):
		fmt.Println("ConvBwdFiltAlgoFFTTiling")
	case ConvBwdFiltAlgo(C.CUDNN_CONVOLUTION_BWD_FILTER_ALGO_COUNT):
		fmt.Println("ConvBwdFiltAlgoCount")
	default:
		fmt.Println("Not supported")
	}
}

//ConvBwdFiltAlgoPerformance is the return struct in the finding algorithm funcs
type ConvBwdFiltAlgoPerformance struct {
	Algo        ConvBwdFiltAlgo `json:"algo,omitempty"`
	Status      Status          `json:"status,omitempty"`
	Time        float32         `json:"time,omitempty"`
	Memory      uint            `json:"memory,omitempty"`
	Determinism Determinism     `json:"determinism,omitempty"`
	MathType    MathType        `json:"math_type,omitempty"`
}

func convertConvBwdFiltAlgoPerformance(input C.cudnnConvolutionBwdFilterAlgoPerf_t) ConvBwdFiltAlgoPerformance {
	var x ConvBwdFiltAlgoPerformance
	x.Algo = ConvBwdFiltAlgo(input.algo)
	x.Status = Status(input.status)
	x.Time = float32(input.time)
	x.Memory = uint(input.memory)
	x.Determinism = Determinism(input.determinism)
	x.MathType = MathType(input.mathType)
	return x
}

//Print prints a human readable copy of the algorithm
func (cb ConvBwdFiltAlgoPerformance) Print() {
	fmt.Println("Convolution Backward Filter Algorithm Performance")
	fmt.Println("-------------------------------------------------")
	ConvBwdFiltAlgo(cb.Algo).print()
	fmt.Println("Status:", Status(cb.Algo).GetErrorString())
	fmt.Println("Time:", cb.Time)
	fmt.Println("Memory:", cb.Memory)
	fmt.Println("Determinism:", cb.Determinism)
	fmt.Println("MathType:", cb.MathType)
}
