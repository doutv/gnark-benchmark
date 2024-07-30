import SwiftUI
import metamask_ios_sdk

struct ConnectTestView: View {
    @StateObject private var metaMaskSDK = MetaMaskSDK.shared(
        AppMetadata(
            name: "zkkyc2",
            url: "https://zkkyc2.com",
            iconUrl: "https://yourdapp.com/icon.png"
        ),
        transport: .deeplinking(dappScheme: "zkkyc2"),
        sdkOptions: SDKOptions(infuraAPIKey: "0x2d04488d5611460ba7d2c2958e2b7227"))
    
    @State private var connected: Bool = false
    @State private var status: String = "Offline"
    @State private var errorMessage: String = ""
    @State private var showError: Bool = false
    @State private var showProgressView: Bool = false

    var body: some View {
        VStack {
            Text("Status: \(status)")
                .padding()
            
            if showProgressView {
                ProgressView()
                    .scaleEffect(1.5, anchor: .center)
                    .progressViewStyle(CircularProgressViewStyle(tint: .black))
            }
            
            Button(action: {
                Task {
//                    await connectAndCallVerifyFunction()
                    showProgressView = true
                    await metaMaskSDK.connect()
                    showProgressView = false
                }
            }) {
                Text("Connect and Call Verify Function")
                    .padding()
                    .background(Color.blue)
                    .foregroundColor(.white)
                    .cornerRadius(8)
            }
                
            Section {
                Button {
                    metaMaskSDK.clearSession()
                } label: {
                    Text("Clear Session")
                        .modifier(TextButton())
                        .frame(maxWidth: .infinity, maxHeight: 32)
                }
                .modifier(ButtonStyle())

                Button {
                    metaMaskSDK.disconnect()
                } label: {
                    Text("Disconnect")
                        .modifier(TextButton())
                        .frame(maxWidth: .infinity, maxHeight: 32)
                }
                .modifier(ButtonStyle())
            }
            
            .alert(isPresented: $showError) {
                Alert(
                    title: Text("Error"),
                    message: Text(errorMessage),
                    dismissButton: .default(Text("OK"))
                )
            }
        }
        .onReceive(NotificationCenter.default.publisher(for: .Connection)) { notification in
            status = notification.userInfo?["value"] as? String ?? "Offline"
        }
    }
    
    private func sendTx() async {
        showProgressView = true
        
        

       

        let transactionResult = await metaMaskSDK.sendTransaction(from: metaMaskSDK.account, to: "0x74c3e0074dc0ff91252b0485dae9d05ee67145e4", amount: "0x0")
        
        showProgressView = false
        
        switch transactionResult {
        case .success(let result):
            // 处理成功结果
            status = "Function Called: \(result)"
        case .failure(let error):
            // 处理错误
            errorMessage = error.localizedDescription
            showError = true
        }
    }

    private func connectAndCallVerifyFunction() async {
        showProgressView = true
        
        let transaction = Transaction(
            to: "0x74c3e0074dc0ff91252b0485dae9d05ee67145e4",
            from: metaMaskSDK.account, // this is initially empty before connection, will be populated with selected address once connected
            value: "0x0",
            data:"0x8753367f0000000000000000000000000000000000000000000000000000000000003039"
        )

        let parameters: [Transaction] = [transaction]

        let transactionRequest = EthereumRequest(
            method: .ethSendTransaction,
            params: parameters
        )

        let transactionResult = await metaMaskSDK.connectWith(transactionRequest)
        
        showProgressView = false
        
        switch transactionResult {
        case .success(let result):
            // 处理成功结果
            status = "Function Called: \(result)"
        case .failure(let error):
            // 处理错误
            errorMessage = error.localizedDescription
            showError = true
        }
    }
    
    // 编码合约函数调用
    private func encodeFunctionCall(functionName: String, parameters: [String: Any]) -> String {
        // 这里你需要使用ABI编码器将函数名和参数编码为以太坊合约调用的数据格式
        // 你可以使用web3swift或其他库来实现
        // 这是一个简化示例，实际实现可能会有所不同
        return "0xEncodedFunctionCallData"
    }
}

struct ConnectTestView_Previews: PreviewProvider {
    static var previews: some View {
        ConnectView()
    }
}
