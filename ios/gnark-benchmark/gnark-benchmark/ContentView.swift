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
            Picker("Select Algorithm", selection: $algorithmSelection) {
                Text("EdDSA").tag("EdDSA")
                Text("ECDSA").tag("ECDSA")
            }
            .pickerStyle(SegmentedPickerStyle())
            .padding()
            Picker("Select System", selection: $systemSelection) {
                Text("Groth16").tag("Groth16")
                Text("Plonk").tag("Plonk")
            }
            .pickerStyle(SegmentedPickerStyle())
            .padding()
            Button("Setup") {
                isRunning = true
                let startTime = Date()
                if algorithmSelection == "ECDSA" && systemSelection == "Groth16" {
                    EcdsaGroth16Setup(directory.filePath)
                } else if algorithmSelection == "EdDSA" && systemSelection == "Groth16" {
                    EddsaGroth16Setup(directory.filePath)
                } else if algorithmSelection == "ECDSA" && systemSelection == "Plonk" {
                    EcdsaPlonkSetup(directory.filePath) //
                } else if algorithmSelection == "EdDSA" && systemSelection == "Plonk" {
                    EddsaPlonkSetup(directory.filePath)
                }
                let endTime = Date()
                setupMessage = "Setup Time: \(endTime.timeIntervalSince(startTime)) seconds"
                isRunning = false
            }
            .disabled(isRunning)
            Text(setupMessage)
            Button("Prove and Verify") {
                isRunning = true
                let startTime = Date()
                if algorithmSelection == "ECDSA" && systemSelection == "Groth16" {
                    EcdsaGroth16ProveAndVerify(directory.filePath)
                } else if algorithmSelection == "EdDSA" && systemSelection == "Groth16" {
                    EddsaGroth16ProveAndVerify(directory.filePath)
                } else if algorithmSelection == "ECDSA" && systemSelection == "Plonk" {
                    EcdsaPlonkProveAndVerify(directory.filePath)
                } else if algorithmSelection == "EdDSA" && systemSelection == "Plonk" {
                    EddsaPlonkProveAndVerify(directory.filePath)
                }
                let endTime = Date()
                proveMessage = "Prove and Verify Time: \(endTime.timeIntervalSince(startTime)) seconds"
                isRunning = false
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
