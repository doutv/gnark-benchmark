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
    @State private var status: String = "Offline"

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
            NavigationView {
                List {
                    Section {
                        Group {

                            HStack {
                                Text("Account")
                                    .bold()
                                    .modifier(TextCallout())
                                Spacer()
                                Text(metaMaskSDK.account)
                                    .modifier(TextCaption())
                            }
                        }
                    }


                        Section {
                            Group {
                               

                                NavigationLink("Transact") {
                                    TransactionView().environmentObject(metaMaskSDK)
                                }

                            }
                        }
                    

                    if metaMaskSDK.account.isEmpty {
                        Section {
//                            Button {
//                                isConnectWith = true
//                            } label: {
//                                Text("Connect With Request")
//                                    .modifier(TextButton())
//                                    .frame(maxWidth: .infinity, maxHeight: 32)
//                            }
//                            .sheet(isPresented: $isConnectWith, onDismiss: {
//                                isConnectWith = false
//                            }) {
//                                TransactionView(isConnectWith: true)
//                                    .environmentObject(metaMaskSDK)
//                            }
//                            .modifier(ButtonStyle())

                            
                            
                            Button {
                                Task {
                                    await connectSDK()
                                }
                            } label: {
                                Text("Connect to MetaMask")
                                    .modifier(TextButton())
                                    .frame(maxWidth: .infinity, maxHeight: 32)
                            }
                            .modifier(ButtonStyle())

                            
                            
                            Button{
                                Task{
                                    await connectAndCallVerifyFunction()
                                }
                            }label:{
                                Text("Connect and verify")
                                    .modifier(TextButton())
                                    .frame(maxWidth: .infinity, maxHeight: 32)
                            }
                            .modifier(ButtonStyle())
                            if showProgressView {
                                ProgressView()
                                    .scaleEffect(1.5, anchor: .center)
                                    .progressViewStyle(CircularProgressViewStyle(tint: .black))
                            }
                            
                            
                        } footer: {
                            Text(connectAndSignResult)
                                .modifier(TextCaption())
                        }
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
                    
                    
                        
                    
                }
                .font(.body)
                .onReceive(NotificationCenter.default.publisher(for: .Connection)) { notification in
                    status = notification.userInfo?["value"] as? String ?? "Offline"
                }
                .navigationTitle("ZKkyc")
                .onAppear {
                    showProgressView = false
                }
            }
                .tabItem {
                    Label("Connect", systemImage: "link.circle")
                }
                .tag(2)
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
}

struct ConnectView_Previews: PreviewProvider {
    static var previews: some View {
        ConnectView()
    }
}
