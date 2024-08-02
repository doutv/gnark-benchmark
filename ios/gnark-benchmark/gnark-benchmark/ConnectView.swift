//
//  ConnectView.swift
//  metamask-ios-sdk_Example
//

import SwiftUI
import metamask_ios_sdk
extension Notification.Name {
    static let Event = Notification.Name("event")
    static let Connection = Notification.Name("connection")
}

private let DAPP_SCHEME = "zkkyc"

@MainActor
struct ConnectView: View {
//    @State var selectedTransport: Transport = .deeplinking(dappScheme: DAPP_SCHEME)
    @State private var dappScheme: String = DAPP_SCHEME

    // We recommend adding support for Infura API for read-only RPCs (direct calls) via SDKOptions
    @ObservedObject var metaMaskSDK = MetaMaskSDK.shared(
        AppMetadata(
            name: "zkkyc app",
            url: "https://zkkyc.com",
            iconUrl: "https://cdn.sstatic.net/Sites/stackoverflow/Img/apple-touch-icon.png"
        ),
        transport: .socket,
        sdkOptions: nil)

    @State private var connected: Bool = false
    @State private var status: Bool = false

    @State private var errorMessage = ""
    @State private var showError = false

    @State private var connectAndSignResult = ""
    @State private var isConnect = true
    @State private var isConnectAndSign = false
    @State private var isConnectWith = false

    @State private var showProgressView = false
    
    @State private var selectedTab = 0
    @State private var attribute = -1
    @State private var op = -1
    @State private var value = -1
    @State private var proof:Data = Data()
    @State private var proofGenerated = false
    @State private var claimClicked = false
    @State private var calldata :String = ""
    @State private var publicwitnessCount = 0

    var body: some View {
        
        TabView(selection: $selectedTab) {
            OkxView(selectedTab:$selectedTab,proofGenerated: $proofGenerated, proof:$proof, attribute:$attribute,op: $op,value:$value,calldata: $calldata)
                .tabItem {
                    Label("Okx", systemImage: "1.circle")
                }
                .tag(0)

            
            if !proofGenerated {
                ThirdPartyView(selectedTab: $selectedTab,attribute:$attribute,op: $op,value:$value,proofGenerated: $proofGenerated,claimClicked: $claimClicked)
                    .tabItem {
                        Label("ThirdParty", systemImage: "2.circle")
                    }
                    .tag(1)
            }
            
            if proofGenerated {
                // 读取proof
                
                VStack {
                        if status{
                            Text("Claim Success!")
                        }
                            Section {

                                if !status{
                                    Button{
                                        Task{
                                            await connectAndCallVerifyFunction()
                                        }
                                    }label:{
                                        Text("Claim")
                                    }
                                }
                                
                                if showProgressView {
                                    ProgressView()
                                        .scaleEffect(1.5, anchor: .center)
                                        .progressViewStyle(CircularProgressViewStyle(tint: .black))
                                }
                                
                            } footer: {
                                Text(connectAndSignResult)
                                    .modifier(TextCaption())
                            }
                    
                    .font(.body)
                    .onAppear {
                        showProgressView = false
                    }
                }
                    .tabItem {
                        Label("ThirdParty", systemImage: "2.circle")
                    }
                    .tag(1)
            }
        }
    }

    func connectSDK() async {
        showProgressView = true
        let result = await metaMaskSDK.connect()
        showProgressView = false

        switch result {
        case .success:
            status = true
        case let .failure(error):
            errorMessage = error.localizedDescription
            showError = true
        }
    }
    func connectAndCallVerifyFunction() async {
        showProgressView = true
       
        
        let transaction = Transaction(
            to: "0x73fdb44133ead4b38aba350df187de1f74626f7d",
            from: metaMaskSDK.account,
            value: "0x0",
            data:calldata
        )

        let parameters: [Transaction] = [transaction]

        let transactionRequest = EthereumRequest(
            method: .ethSendTransaction,
            params: parameters
        )

        let transactionResult = await metaMaskSDK.connectWith(transactionRequest)
        
        showProgressView = false
        print("transactionResult:\(transactionResult)")
        switch transactionResult {
        case .success(_):
            // 处理成功结果
            status = true
        case .failure(let error):
            // 处理错误
            errorMessage = error.localizedDescription
            showError = true
        }
    }
}


struct ConnectView_Previews: PreviewProvider {
    static var previews: some View {
        ConnectView()
    }
}

func readFileFromDocumentsDirectory(fileName: String) -> Data? {
        let fileManager = FileManager.default
        let documentDirectory = fileManager.urls(for: .documentDirectory, in: .userDomainMask).first!
        let fileURL = documentDirectory.appendingPathComponent(fileName)
        
        do {
            let fileData = try Data(contentsOf: fileURL)
            return fileData
        } catch {
            print("无法读取文件 \(fileName): \(error.localizedDescription)")
            return nil
        }
    }



