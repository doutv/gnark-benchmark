//
//  ContentView.swift
//  gnark-benchmark
//
//  Created by Jason HUANG on 8/7/2024.
//

import SwiftUI
import Ecdsa

struct ContentView: View {
    @State private var isRunning = false
    @State private var setupMessage = ""
    @State private var proveMessage = ""
    
    var body: some View {
        VStack {
            Button("Setup") {
                isRunning = true
                let startTime = Date()
                EcdsaSetup()
                let endTime = Date()
                setupMessage = "Setup Time: \(endTime.timeIntervalSince(startTime)) seconds"
                isRunning = false
            }
            .disabled(isRunning)
            Text(setupMessage)
            
            Button("Prove and Verify") {
                isRunning = true
                let startTime = Date()
                EcdsaProveAndVerify()
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
