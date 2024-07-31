import SwiftUI
import Gnark
import Foundation

struct OkxView: View {
    @State private var isRunning = false
    @State private var setupMessage = ""
    @State private var proveMessage = ""
    
    
    
    
    @Binding var selectedTab: Int
    @Binding var proofGenerated: Bool

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
    var contryMap:[Int:String] = [0:"America",1:"China"]
    
    var directory: URL {
        let documentDirectory = FileManager().urls(for: .documentDirectory, in: .userDomainMask).first!
        return documentDirectory.appendingPathComponent("depositStorage")
    }

    var body: some View {
        VStack {
            
            if  attribute != -1, op != -1, value != -1 {
                Text("Requirement: \(attributesMap[attribute] ?? "") \(opMap[op] ?? "") \(contryMap[value] ?? "")")
            }
            
            Button("Prove") {
                isRunning = true
                print("directory.filePath:\(directory.filePath)")
                let setupStartTime = Date()
                
                DispatchQueue.global().async {
                    EddsaGroth16Setup(directory.filePath)
                    
                    let setupEndTime = Date()
                    DispatchQueue.main.async {
                        setupMessage = "Setup Time: \(setupEndTime.timeIntervalSince(setupStartTime)) seconds"
                    }
                    
                    
                    
                    let proveStartTime = Date()
                    
                    EddsaGroth16Prove(directory.filePath,Int64(attribute) ,Int64(op) ,Int64(value))

                  
                    let proveEndTime = Date()
                    proofGenerated = true
                    DispatchQueue.main.async {
                        proveMessage = "Prove Time: \(proveEndTime.timeIntervalSince(proveStartTime)) seconds"
                        isRunning = false
                    }
                }
            }
            .padding()
            .disabled(isRunning)
            
            
            Text("\(setupMessage) \n \(proveMessage)")
            Button("Go back to third party app"){
                selectedTab = 1
            }
            .disabled(!proofGenerated)
        }
        .padding()
        .navigationTitle("Okx")
        
    }
}

struct OkxView_Previews: PreviewProvider {
    @State static var selectedTab = 0
    @State static var attribute = -1
    @State static var op = -1
    @State static var value = -1
    @State static var proofGenerated = false
    static var previews: some View {
        OkxView(selectedTab:$selectedTab, proofGenerated: $proofGenerated, attribute: $attribute, op: $op, value: $value)
    }
}

extension URL {
    var filePath: String {
        absoluteString.replacingOccurrences(of: "file://", with: "")
    }
}
