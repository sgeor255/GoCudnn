package gocudnn

/*
#include <cudnn.h>
*/
import "C"

//LRNmode is used for the flags in LRNmode
type LRNmode C.cudnnLRNMode_t

func (l LRNmode) c() C.cudnnLRNMode_t { return C.cudnnLRNMode_t(l) }

//LRNCrossChanelDim1 is the only flag I guess for now for LRNmode
const LRNCrossChanelDim1 LRNmode = C.CUDNN_LRN_CROSS_CHANNEL_DIM1

//DivNormMode is usde for C.cudnnDivNormMode_t flags
type DivNormMode C.cudnnDivNormMode_t

//DivnormPrecomputedMeans flag for divNomMode
const DivnormPrecomputedMeans DivNormMode = C.CUDNN_DIVNORM_PRECOMPUTED_MEANS

func (d DivNormMode) c() C.cudnnDivNormMode_t { return C.cudnnDivNormMode_t(d) }

/* LRN functions: output = alpha * normalize(x) + beta * old_y */

//LRNCrossChannelForward  LRN cross-channel forward computation. Double parameters cast to tensor data type
func (handle *Handle) LRNCrossChannelForward(
	norm *LRND,
	mode LRNmode,
	alpha CScaler,
	xD *TensorD,
	x Memer,
	beta CScaler,
	yD *TensorD,
	y Memer,
) error {
	return Status(C.cudnnLRNCrossChannelForward(
		handle.x,
		norm.descriptor,
		mode.c(),
		alpha.CPtr(),
		xD.descriptor,
		x.Ptr(),
		beta.CPtr(),
		yD.descriptor,
		y.Ptr(),
	)).error("LRNCrossChannelForward")
}

//LRNCrossChannelBackward  LRN cross-channel backward computation. Double parameters cast to tensor data type
func (handle *Handle) LRNCrossChannelBackward(
	norm *LRND,
	mode LRNmode,
	alpha CScaler,
	yD *TensorD,
	y Memer,
	dyD *TensorD,
	dy Memer,
	xD *TensorD,
	x Memer,
	beta CScaler,
	dxD *TensorD,
	dx Memer,
) error {
	return Status(C.cudnnLRNCrossChannelBackward(
		handle.x,
		norm.descriptor,
		mode.c(),
		alpha.CPtr(),
		yD.descriptor,
		y.Ptr(),
		dyD.descriptor,
		dy.Ptr(),
		xD.descriptor,
		x.Ptr(),
		beta.CPtr(),
		dxD.descriptor,
		dx.Ptr(),
	)).error("LRNCrossChannelForward")
}

//DivisiveNormalizationForward   LCN/divisive normalization functions: y = alpha * normalize(x) + beta * y
func (handle *Handle) DivisiveNormalizationForward(
	norm LRND,
	mode DivNormMode,
	alpha CScaler,
	xD TensorD, /* same desc for means, temp, temp2 */
	x Memer,
	means Memer, /* if NULL, means are assumed to be zero */
	temp Memer,
	temp2 Memer,
	beta CScaler,
	yD TensorD,
	y Memer,
) error {
	return Status(C.cudnnDivisiveNormalizationForward(
		handle.x,
		norm.descriptor,
		mode.c(),
		alpha.CPtr(),
		xD.descriptor,
		x.Ptr(),
		means.Ptr(),
		temp.Ptr(),
		temp2.Ptr(),
		beta.CPtr(),
		yD.descriptor,
		y.Ptr(),
	)).error("DivisiveNormalizationForward")
}

//DivisiveNormalizationBackward  LRN cross-channel backward computation. Double parameters cast to tensor data type
func (handle *Handle) DivisiveNormalizationBackward(
	norm *LRND,
	mode DivNormMode,
	alpha CScaler,
	xD *TensorD, /* same desc for x, means, dy, temp, temp2 */
	x Memer,
	means Memer, /* if NULL, means are assumed to be zero */
	dy Memer,
	temp Memer,
	temp2 Memer,
	beta CScaler,
	dXdMeansDesc *TensorD, /* same desc for dx, dMeans */
	dx Memer, /* output x differential */
	dMeans Memer, /* output means differential, can be NULL */
) error {
	return Status(C.cudnnDivisiveNormalizationBackward(
		handle.x,
		norm.descriptor,
		mode.c(),
		alpha.CPtr(),
		xD.descriptor,
		x.Ptr(),
		means.Ptr(),
		dy.Ptr(),
		temp.Ptr(),
		temp2.Ptr(),
		beta.CPtr(),
		dXdMeansDesc.descriptor,
		dx.Ptr(),
		dMeans.Ptr(),
	)).error("DivisiveNormalizationBackward")
}
