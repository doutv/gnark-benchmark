import SwiftUI
import Gnark
import Foundation
import metamask_ios_sdk



struct ContentView: View {
    @State private var selectedTab = 0
    @State private var attribute = -1
    @State private var op = -1
    @State private var value = -1
    @State private var proof = -1
     
    
    var body: some View {
        TabView(selection: $selectedTab) {
            OkxView(selectedTab:$selectedTab, attribute:$attribute,op: $op,value:$value)
                .tabItem {
                    Label("Okx", systemImage: "1.circle")
                }
                .tag(0)

            ThirdPartyView(selectedTab: $selectedTab,attribute:$attribute,op: $op,value:$value)
                .tabItem {
                    Label("ThirdParty", systemImage: "2.circle")
                }
                .tag(1)
            ConnectView()
                .tabItem {
                    Label("Connect", systemImage: "link.circle")
                }
                .tag(2)
        }
    }
}

// Preview
struct ContentView_Previews: PreviewProvider {
    
    static var previews: some View {
        ContentView()
    }
}
