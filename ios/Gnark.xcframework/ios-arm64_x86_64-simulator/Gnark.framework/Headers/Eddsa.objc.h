// Objective-C API for talking to gnark-benchmark/eddsa Go package.
//   gobind -lang=objc gnark-benchmark/eddsa
//
// File is generated by gobind. Do not edit.

#ifndef __Eddsa_H__
#define __Eddsa_H__

@import Foundation;
#include "ref.h"
#include "Universe.objc.h"


@class EddsaDummyCircuit;

@interface EddsaDummyCircuit : NSObject <goSeqRefInterface> {
}
@property(strong, readonly) _Nonnull id _ref;

- (nonnull instancetype)initWithRef:(_Nonnull id)ref;
- (nonnull instancetype)init;
// skipped field DummyCircuit.A with unsupported type: github.com/consensys/gnark/frontend.Variable

// skipped field DummyCircuit.C with unsupported type: github.com/consensys/gnark/frontend.Variable

// skipped method DummyCircuit.Define with unsupported parameter or return types

@end

FOUNDATION_EXPORT void EddsaGroth16Prove(NSString* _Nullable fileDir);

FOUNDATION_EXPORT void EddsaGroth16Setup(NSString* _Nullable fileDir);

FOUNDATION_EXPORT void EddsaPlonkSetup(NSString* _Nullable fileDir);

#endif
