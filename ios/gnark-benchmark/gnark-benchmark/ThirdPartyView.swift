import SwiftUI

struct ThirdPartyView: View {
    @Binding var selectedTab: Int
    // 0: age
    // 1: gender
    // 2: nationality
    @Binding var attribute : Int
    // 0: equal
    // 1: not equal
    // 2: less than
    // 3: greater than
    @Binding var op :Int
    @Binding var value :Int
    
    
    
    var body: some View {
        VStack {
            Text("Criteria: age > 18")
                .padding()
            
            Button("Prove") {
                attribute = 0
                op = 3
                value = 18
                selectedTab = 0
            }
            .padding()
            .background(Color.blue)
            .foregroundColor(.white)
            .cornerRadius(8)
            
            
        }
    }
}

struct ThirdPartyView_Previews: PreviewProvider {
    @State static var selectedTab = 1
    @State static var attribute = 0
    @State static var op = 3
    @State static var value = 18
    
    
    static var previews: some View {
        ThirdPartyView(selectedTab: $selectedTab,attribute:$attribute,op:$op,value:$value)
    }
}
