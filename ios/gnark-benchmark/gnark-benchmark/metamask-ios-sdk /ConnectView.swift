//
//  ConnectView.swift
//  metamask-ios-sdk_Example
//

import SwiftUI
import metamask_ios_sdk
import CryptoKit
import BigInt
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
                    
                   
//                    if !proof.isEmpty {
//                        ScrollView {
//                                    Text("Proof data: \(proof.map { String($0) }.joined(separator: " "))")
//                                        .padding()
//                                }
//                        } else {
//                            Text("No proof data available")
//                        }
                   
                    
                    if calldata != ""{
                        Text(calldata)
                    }

                        if status{
                            Text("Claim Success!")
                        }
                        
                        
                        
                            Section {

                                if !status{
                                    Button{
                                        Task{
    //                                        metaMaskSDK.clearSession()
    //                                        metaMaskSDK.disconnect()
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
                        

                        
//                            Section {
//                                Button {
//
//                                } label: {
//                                    Text("Clear Session and disconnect")
//                                        .modifier(TextButton())
//                                        .frame(maxWidth: .infinity, maxHeight: 32)
//                                }
//                                .modifier(ButtonStyle())
//
//                            }
                        
                        
                            
                        
                    
                    .font(.body)
//                    .onReceive(NotificationCenter.default.publisher(for: .Connection)) { notification in
//                        status = notification.userInfo?["value"] as? String ?? "Offline"
//                    }
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
        
//        var calldata : String
//        构造calldata
//        if let publicWitness = readFileFromDocumentsDirectory(fileName: "public_witness.bin"){
//            publicwitnessCount = publicWitness.count
//            var tempCalldata = Data()
//            let functionSelector = "0x94e4398a"
//           
//            
//
//            
//
////            calldata = tempCalldata.map { String(format: "%02x", $0) }.joined()
//            
//
//            
//        }else{
//            fatalError()
//        }
       
        

        
        
        let transaction = Transaction(
            to: "0x73fdb44133ead4b38aba350df187de1f74626f7d",
            from: metaMaskSDK.account, // this is initially empty before connection, will be populated with selected address once connected
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
        case .success(let result):
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

func readFileFromBundle(resourceName: String, fileExtension: String) -> Data? {
    // 获取文件的 URL
    guard let fileURL = Bundle.main.url(forResource: resourceName, withExtension: fileExtension) else {
        print("文件未找到：\(resourceName).\(fileExtension)")
        return nil
    }
    
    do {
        // 读取文件内容
        let fileData = try Data(contentsOf: fileURL)
        return fileData
    } catch {
        print("无法读取文件：\(error.localizedDescription)")
        return nil
    }
}

//func copyProofToDocumentsDirectory() {
//        let fileManager = FileManager.default
//        let documentDirectory = fileManager.urls(for: .documentDirectory, in: .userDomainMask).first!
//
//        let files = ["eddsa.proof"]
//        for file in files {
//            if let sourceURL = Bundle.main.url(forResource: file, withExtension: nil) {
//                let targetURL = documentDirectory.appendingPathComponent(file)
//                do {
//                    if !fileManager.fileExists(atPath: targetURL.path) {
//                        try fileManager.copyItem(at: sourceURL, to: targetURL)
//                        print("\(file) 已成功复制到文档目录")
//                    } else {
//                        print("\(file) 已存在于文档目录")
//                    }
//                } catch {
//                    print("无法复制文件 \(file): \(error.localizedDescription)")
//                    fatalError()
//                }
//            }
//        }
//    }
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

func dataToUint256Array(_ data: Data, count: Int) -> [UInt256] {
    var array: [UInt256] = []
    for i in 0..<count {
        let start = i * 32
        let end = start + 32
        let chunk = data[start..<end]
        let uint256 = UInt256(data: chunk)
        array.append(uint256)
    }
    return array
}

// uint256 类型的表示
struct UInt256 {
    var value: UInt64

    init(data: Data) {
        self.value = data.withUnsafeBytes { $0.load(as: UInt64.self) }
    }

    var data: Data {
        var bigEndianValue = value.bigEndian
        return Data(bytes: &bigEndianValue, count: MemoryLayout<UInt64>.size).prefix(32)
    }
}
func encodeUint256Array(_ array: [UInt256]) -> Data {
    var encoded = Data()
    for uint256 in array {
        encoded.append(uint256.data)
    }
    return encoded
}
