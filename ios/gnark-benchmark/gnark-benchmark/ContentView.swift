//
//  ContentView.swift
//  gnark-benchmark
//
//  Created by Jason HUANG on 8/7/2024.
//

import SwiftUI
import Ecdsa

struct ContentView: View {
    var body: some View {
        VStack {
            Image(systemName: "globe")
                .imageScale(.large)
                .foregroundStyle(.tint)
            Text("Hello, world!")
        }
        .padding()
        .onAppear {
            print("Log")
            
        }
        
    }
}

#Preview {
    ContentView()
}
