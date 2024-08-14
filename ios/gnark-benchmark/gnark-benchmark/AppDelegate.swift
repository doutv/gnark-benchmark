//
//  AppDelegate.swift
//  metamask-ios-sdk
//

import UIKit
import metamask_ios_sdk

@UIApplicationMain
class AppDelegate: UIResponder, UIApplicationDelegate {
    var window: UIWindow?

    

    func application(_ app: UIApplication, open url: URL, options: [UIApplication.OpenURLOptionsKey: Any] = [:]) -> Bool {
        if URLComponents(url: url, resolvingAgainstBaseURL: true)?.host == "mmsdk" {
            MetaMaskSDK.sharedInstance?.handleUrl(url)
        } else {
            // handle other deeplinks
        }
        return true
    }
    
    func application(_ application: UIApplication, didFinishLaunchingWithOptions launchOptions: [UIApplication.LaunchOptionsKey: Any]?) -> Bool {
            copyFilesToDocumentsDirectory()
            return true
        }
    
    func copyFilesToDocumentsDirectory() {
            let fileManager = FileManager.default
            let documentDirectory = fileManager.urls(for: .documentDirectory, in: .userDomainMask).first!

            let files = ["dummy.r1cs", "dummy.zkey", "dummy.vkey"]
            for file in files {
                if let sourceURL = Bundle.main.url(forResource: file, withExtension: nil) {
                    let targetURL = documentDirectory.appendingPathComponent(file)
                    do {
                        if !fileManager.fileExists(atPath: targetURL.path) {
                            try fileManager.copyItem(at: sourceURL, to: targetURL)
                            print("\(file) 已成功复制到文档目录")
                        } else {
                            print("\(file) 已存在于文档目录")
                        }
                    } catch {
                        print("无法复制文件 \(file): \(error.localizedDescription)")
                    }
                }
            }
        }
        

    func applicationWillResignActive(_: UIApplication) {
        // Sent when the application is about to move from active to inactive state. This can occur for certain types of temporary interruptions (such as an incoming phone call or SMS message) or when the user quits the application and it begins the transition to the background state.
        // Use this method to pause ongoing tasks, disable timers, and throttle down OpenGL ES frame rates. Games should use this method to pause the game.
    }

    func applicationDidEnterBackground(_: UIApplication) {
        // Use this method to release shared resources, save user data, invalidate timers, and store enough application state information to restore your application to its current state in case it is terminated later.
        // If your application supports background execution, this method is called instead of applicationWillTerminate: when the user quits.
    }

    func applicationWillEnterForeground(_: UIApplication) {
        // Called as part of the transition from the background to the inactive state; here you can undo many of the changes made on entering the background.
    }

    func applicationDidBecomeActive(_: UIApplication) {
        // Restart any tasks that were paused (or not yet started) while the application was inactive. If the application was previously in the background, optionally refresh the user interface.
    }

    func applicationWillTerminate(_: UIApplication) {
        // Called when the application is about to terminate. Save data if appropriate. See also applicationDidEnterBackground:.
    }
}