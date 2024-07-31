import SwiftUI
import metamask_ios_sdk
import Foundation

struct ThirdPartyView: View {
    @Binding var selectedTab: Int
    @Binding var attribute: Int
    @Binding var op: Int
    @Binding var value: Int
    @Binding var proofGenerated:Bool
    @Binding var claimClicked:Bool
    
    
    
    
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
            
            
            
            
            
//            if proofGenerated{
//                Button("Claim") {
//                    claimClicked = true
//                    
//                }
//                .padding()
//                .background(Color.blue)
//                .foregroundColor(.white)
//                .cornerRadius(8)
//            }
            
            
            
            
            
//            Button("Connect to MetaMask") {
//                Task {
//                    await connectSDK()
//                }
//            }
//            .padding()
//            .background(Color.green)
//            .foregroundColor(.white)
//            .cornerRadius(8)
            
            if showProgressView {
                ProgressView()
                    .scaleEffect(1.5, anchor: .center)
                    .progressViewStyle(CircularProgressViewStyle(tint: .black))
            }
            
            if showError {
                Text("Error: \(errorMessage)")
                    .foregroundColor(.red)
            }
            
            
        }
    }
    
    
}

struct ThirdPartyView_Previews: PreviewProvider {
    @State static var selectedTab = 1
    @State static var attribute = 0
    @State static var op = 3
    @State static var value = 18
    @State static var proofGenerated = false;
    @State static var claimClicked = false
    
    static var previews: some View {
        ThirdPartyView(selectedTab: $selectedTab, attribute: $attribute, op: $op, value: $value,proofGenerated:$proofGenerated,claimClicked: $claimClicked)
    }
}
