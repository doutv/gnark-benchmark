// Objective-C API for talking to gnark-benchmark/ecdsa Go package.
//   gobind -lang=objc gnark-benchmark/ecdsa
//
// File is generated by gobind. Do not edit.

#ifndef __Ecdsa_H__
#define __Ecdsa_H__

@import Foundation;
#include "ref.h"
#include "Universe.objc.h"


FOUNDATION_EXPORT void EcdsaGroth16Prove(NSString* _Nullable fileDir);

FOUNDATION_EXPORT void EcdsaGroth16Setup(NSString* _Nullable fileDir);

FOUNDATION_EXPORT void EcdsaPlonkProve(NSString* _Nullable fileDir);

FOUNDATION_EXPORT void EcdsaPlonkSetup(NSString* _Nullable fileDir);

#endif
