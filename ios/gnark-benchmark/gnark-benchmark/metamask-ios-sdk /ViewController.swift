//
//  ViewController.swift
//  metamask-ios-sdk
//

import UIKit
import SwiftUI

class ViewController: UIViewController {
    let connectView = ContentView()

    override func viewDidLoad() {
        super.viewDidLoad()
        let childView = UIHostingController(rootView: connectView)
        addChild(childView)
        childView.view.frame = view.bounds
        view.addSubview(childView.view)
        childView.didMove(toParent: self)
    }
}
