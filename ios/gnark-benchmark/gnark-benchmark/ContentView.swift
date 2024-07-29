import SwiftUI
import Gnark
import Foundation

struct ContentView: View {
    @State private var isRunning = false
    @State private var setupMessage = ""
    @State private var proveMessage = ""
    @State private var algorithmSelection = "EdDSA"
    @State private var systemSelection = "Groth16"
    var directory: URL {
        let documentDirectory = FileManager().urls(for: .documentDirectory, in: .userDomainMask).first!
        return documentDirectory.appendingPathComponent("depositStorage")
    }
    var body: some View {
        VStack {
            
            
            Button("Setup") {
                isRunning = true
                let startTime = Date()
                
                DispatchQueue.global().async {
                    
                    EddsaGroth16Setup(directory.filePath)
                    
                    let endTime = Date()
                    DispatchQueue.main.async {
                        setupMessage = "Setup Time: \(endTime.timeIntervalSince(startTime)) seconds"
                        isRunning = false
                    }
                }
            }
            .padding()
            .disabled(isRunning)
            
            Text(setupMessage)
            Button("Prove") {
                isRunning = true
                let startTime = Date()
                EddsaGroth16Prove(directory.filePath)
                
                let endTime = Date()
                DispatchQueue.main.async {
                    proveMessage = "Prove Time: \(endTime.timeIntervalSince(startTime)) seconds"
                    isRunning = false
                }
            }
            .disabled(isRunning)
            Text(proveMessage)
        }
        .padding()
    }
}
// Preview
struct ContentView_Previews: PreviewProvider {
    static var previews: some View {
        ContentView()
    }
}
extension URL {
    var filePath: String {
        absoluteString.replacingOccurrences(of: "file://", with: "")
    }
}


