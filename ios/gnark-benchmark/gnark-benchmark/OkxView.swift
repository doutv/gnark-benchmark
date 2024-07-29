import SwiftUI
import Gnark
import Foundation

struct OkxView: View {
    @State private var isRunning = false
    @State private var setupMessage = ""
    @State private var proveMessage = ""
    
    // 0: age
    // 1: gender
    // 2: nationality
    @Binding var attribute : Int
    // 0: equal
    // 1: not equal
    // 2: less than
    // 3: greater than
    @Binding var op :Int
    @Binding var value :Int
    
    var attributesMap : [Int:String] = [0:"age",1:"gender",2:"nationality"]
    var opMap  : [Int:String] = [0:"equal",1:"not equal",2:"less than",3:"greater than"]
    
    var directory: URL {
        let documentDirectory = FileManager().urls(for: .documentDirectory, in: .userDomainMask).first!
        return documentDirectory.appendingPathComponent("depositStorage")
    }

    var body: some View {
        VStack {
            
            if  attribute != -1, op != -1, value != -1 {
                Text("Requirement: \(attributesMap[attribute] ?? "") \(opMap[op] ?? "") \(value)")
                Text("attribute:\(attribute) op:\(op) value:\(value)")
            }
            
            Button("Setup and Prove") {
                isRunning = true
                let setupStartTime = Date()
                
                DispatchQueue.global().async {
                    EddsaGroth16Setup(directory.filePath)
                    
                    let setupEndTime = Date()
                    DispatchQueue.main.async {
                        setupMessage = "Setup Time: \(setupEndTime.timeIntervalSince(setupStartTime)) seconds"
                    }
                    
                    
                    
                    let proveStartTime = Date()
                    
                    var res: Void =  EddsaGroth16Prove(directory.filePath,Int64(attribute) ,Int64(op) ,Int64(value))

                    if res == Void {
                        print("res:\(res)")
                    }
                    
                    let proveEndTime = Date()
                    DispatchQueue.main.async {
                        proveMessage = "Prove Time: \(proveEndTime.timeIntervalSince(proveStartTime)) seconds"
                        isRunning = false
                    }
                }
            }
            .padding()
            .disabled(isRunning)
            
            Text("\(setupMessage) \n \(proveMessage)")
        }
        .padding()
        .navigationTitle("Okx")
        
    }
}

struct OkxView_Previews: PreviewProvider {
    @State static var attribute = -1
    @State static var op = -1
    @State static var value = -1
    static var previews: some View {
        OkxView(attribute: $attribute, op: $op, value: $value)
    }
}

extension URL {
    var filePath: String {
        absoluteString.replacingOccurrences(of: "file://", with: "")
    }
}
