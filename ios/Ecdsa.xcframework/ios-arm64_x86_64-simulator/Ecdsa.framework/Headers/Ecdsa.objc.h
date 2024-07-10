// Objective-C API for talking to gnark-benchmark/ecdsa Go package.
//   gobind -lang=objc gnark-benchmark/ecdsa
//
// File is generated by gobind. Do not edit.

#ifndef __Ecdsa_H__
#define __Ecdsa_H__

@import Foundation;
#include "ref.h"
#include "Universe.objc.h"


FOUNDATION_EXPORT void EcdsaProveAndVerify(void);

FOUNDATION_EXPORT void EcdsaSetup(void);

#endif
