import SwiftUI
import Gnark
import Foundation

struct ContentView: View {
    @State private var selectedTab = 0
    @State private var attribute = -1
    @State private var op = -1
    @State private var value = -1
    
    var body: some View {
        TabView(selection: $selectedTab) {
            OkxView(attribute:$attribute,op: $op,value:$value)
                .tabItem {
                    Label("Okx", systemImage: "1.circle")
                }
                .tag(0)

            ThirdPartyView(selectedTab: $selectedTab,attribute:$attribute,op: $op,value:$value)
                .tabItem {
                    Label("ThirdParty", systemImage: "2.circle")
                }
                .tag(1)
        }
    }
}

// Preview
struct ContentView_Previews: PreviewProvider {
    
    static var previews: some View {
        ContentView()
    }
}
