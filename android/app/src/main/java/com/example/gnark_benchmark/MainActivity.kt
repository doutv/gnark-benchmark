package com.example.gnark_benchmark

import android.annotation.SuppressLint
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.tooling.preview.Preview
import com.example.gnark_benchmark.ui.theme.GnarkbenchmarkTheme
import java.io.BufferedReader
import java.io.File
import java.io.InputStreamReader
import ecdsa.Ecdsa
import eddsa.Eddsa
import android.os.Bundle
import androidx.activity.ComponentActivity
import android.util.Log
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.height
import androidx.compose.material3.Button
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.rememberCoroutineScope
import androidx.compose.ui.Alignment
import androidx.compose.ui.unit.dp
import kotlinx.coroutines.launch
import android.content.Context
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.wrapContentSize
import androidx.compose.material3.DropdownMenu
import androidx.compose.material3.DropdownMenuItem
import androidx.compose.runtime.getValue
import androidx.compose.runtime.setValue

// Other imports remain the same

class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContent {
            GnarkbenchmarkTheme {
                Scaffold(modifier = Modifier.fillMaxSize()) {
                    BenchmarkComponent(applicationContext.filesDir.toString())
                }
            }
        }
    }
}

@Composable
fun BenchmarkComponent(fileDir: String) {
    val setupTime = remember { mutableStateOf("Not started") }
    val proveAndVerifyTime = remember { mutableStateOf("Not started") }
    val selectedAlgorithm = remember { mutableStateOf("EdDSA") }
    val selectedSystem = remember { mutableStateOf("Groth16") }
    val coroutineScope = rememberCoroutineScope()
    val algorithms = listOf("EdDSA", "ECDSA")
    val systems = listOf("Groth16", "Plonk")

    Column(modifier = Modifier
        .padding(top = 200.dp)
        .fillMaxSize(), horizontalAlignment = Alignment.CenterHorizontally) {
        // Algorithm Selector
        DropdownMenuComponent("Select Algorithm", selectedAlgorithm.value, algorithms) { selected ->
            selectedAlgorithm.value = selected
        }
        // System Selector
        DropdownMenuComponent("Select System", selectedSystem.value, systems) { selected ->
            selectedSystem.value = selected
        }
        Button(onClick = {
            setupTime.value = "Setting up..."
            coroutineScope.launch {
                val setupStartTime = System.nanoTime()
                when (selectedAlgorithm.value to selectedSystem.value) {
                    "ECDSA" to "Groth16" -> Ecdsa.groth16Setup(fileDir)
                    "ECDSA" to "Plonk" -> Ecdsa.plonkSetup(fileDir)
                    "EdDSA" to "Groth16" -> Eddsa.groth16Setup(fileDir)
                    "EdDSA" to "Plonk" -> Eddsa.plonkSetup(fileDir)
                    else -> Log.e("BenchmarkComponent", "Invalid selection")
                }
                val setupEndTime = System.nanoTime()
                setupTime.value = "Setup: ${(setupEndTime - setupStartTime) / 1_000_000} ms"
            }
        }) {
            Text("Setup")
        }
        Text(text = setupTime.value)
        Spacer(modifier = Modifier.height(8.dp))
        Button(onClick = {
            proveAndVerifyTime.value = "Proving..."
            coroutineScope.launch {
                val proveAndVerifyStartTime = System.nanoTime()
                when (selectedAlgorithm.value to selectedSystem.value) {
                    "ECDSA" to "Groth16" -> Ecdsa.groth16ProveAndVerify(fileDir)
                    "ECDSA" to "Plonk" -> Ecdsa.plonkProveAndVerify(fileDir)
                    "EdDSA" to "Groth16" -> Eddsa.groth16ProveAndVerify(fileDir)
                    "EdDSA" to "Plonk" -> Eddsa.plonkProveAndVerify(fileDir)
                    else -> Log.e("BenchmarkComponent", "Invalid selection")
                }
                val proveAndVerifyEndTime = System.nanoTime()
                proveAndVerifyTime.value = "Prove and Verify: ${(proveAndVerifyEndTime - proveAndVerifyStartTime) / 1_000_000} ms"
            }
        }) {
            Text("Prove and Verify")
        }
        Spacer(modifier = Modifier.height(16.dp))
        Text(text = proveAndVerifyTime.value)
    }
}

@Composable
fun DropdownMenuComponent(label: String, selectedItem: String, items: List<String>, onItemSelected: (String) -> Unit) {
    var expanded: Boolean by remember { mutableStateOf(false) }
    Box(modifier = Modifier
        .fillMaxWidth()
        .wrapContentSize(Alignment.TopStart)) {
        Text(selectedItem, modifier = Modifier
            .fillMaxWidth()
            .clickable { expanded = true })
        DropdownMenu(
            expanded = expanded,
            onDismissRequest = { expanded = false },
            modifier = Modifier.fillMaxWidth()
        ) {
            items.forEach { label ->
                DropdownMenuItem(text = { Text(text = label) }, onClick = {
                        onItemSelected(label)
                        expanded = false
                    })
            }
        }
    }
}