import SwiftUI
import metamask_ios_sdk
import Foundation

struct ThirdPartyView: View {
    @Binding var selectedTab: Int
    @Binding var attribute: Int
    @Binding var op: Int
    @Binding var value: Int
    
    
    
    // 创建 MetaMaskSDK 实例
    @ObservedObject var metaMaskSDK = MetaMaskSDK.shared(
        AppMetadata(
            name: "Your App Name",
            url: "https://yourapp.com",
            iconUrl: "https://yourapp.com/icon.png"
        ),
        transport: .deeplinking(dappScheme: "yourappscheme"),
        sdkOptions: nil
    )
    
    @State private var status: String = "Offline"
    @State private var errorMessage = ""
    @State private var showError = false
    @State private var showProgressView = false
    
    var body: some View {
        VStack {
            Text("Criteria: nationality not equal China")
                .padding()
            
            Button("Prove") {
                attribute = 2
                op = 1
                value = 1
                selectedTab = 0
            }
            .padding()
            .background(Color.blue)
            .foregroundColor(.white)
            .cornerRadius(8)
            
            Button("Connect to MetaMask") {
                Task {
                    await connectSDK()
                }
            }
            .padding()
            .background(Color.green)
            .foregroundColor(.white)
            .cornerRadius(8)
            
            if showProgressView {
                ProgressView()
                    .scaleEffect(1.5, anchor: .center)
                    .progressViewStyle(CircularProgressViewStyle(tint: .black))
            }
            
            if showError {
                Text("Error: \(errorMessage)")
                    .foregroundColor(.red)
            }
            
            Text("Status: \(status)")
                .padding()
        }
    }
    
    func connectSDK() async {
        showProgressView = true
        let result = await metaMaskSDK.connect()
        showProgressView = false
        switch result {
        case .success:
            status = "Online"
        case let .failure(error):
            errorMessage = error.localizedDescription
            showError = true
        }
    }
}

struct ThirdPartyView_Previews: PreviewProvider {
    @State static var selectedTab = 1
    @State static var attribute = 0
    @State static var op = 3
    @State static var value = 18
    
    static var previews: some View {
        ThirdPartyView(selectedTab: $selectedTab, attribute: $attribute, op: $op, value: $value)
    }
}
