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

// Other imports remain the same

class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContent {
            GnarkbenchmarkTheme {
                Scaffold(modifier = Modifier.fillMaxSize()) {
                    EcdsaComponent(applicationContext.filesDir.toString())
                }
            }
        }
    }
}

@Composable
fun EcdsaComponent(fileDir: String) {
    val setupTime = remember { mutableStateOf("Not started") }
    val proveAndVerifyTime = remember { mutableStateOf("Not started") }
    val coroutineScope = rememberCoroutineScope()

    Column(modifier = Modifier.padding(top = 200.dp).fillMaxSize(), horizontalAlignment = Alignment.CenterHorizontally) {
        Button(onClick = {
            setupTime.value = "Setuping..."
            coroutineScope.launch {
                val setupStartTime = System.nanoTime()
                Ecdsa.setup(fileDir)
                val setupEndTime = System.nanoTime()
                setupTime.value = "Setup: ${(setupEndTime - setupStartTime) / 1_000_000} ms"
            }
        }) {
            Text("Setup")
        }
        Text(text = setupTime.value)
        Spacer(modifier = Modifier.height(8.dp)) // Add space between buttons
        Button(onClick = {
            proveAndVerifyTime.value = "Proving..."
            coroutineScope.launch {
                val proveAndVerifyStartTime = System.nanoTime()
                Ecdsa.proveAndVerify(fileDir)
                val proveAndVerifyEndTime = System.nanoTime()
                proveAndVerifyTime.value = "Prove and Verify: ${(proveAndVerifyEndTime - proveAndVerifyStartTime) / 1_000_000} ms"
            }
        }) {
            Text("Prove and Verify")
        }
        Spacer(modifier = Modifier.height(16.dp)) // Add space between button and text
        Text(text = proveAndVerifyTime.value)
    }
}