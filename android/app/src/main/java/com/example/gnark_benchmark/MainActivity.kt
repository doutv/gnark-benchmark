package com.example.gnark_benchmark

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
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
// Other imports remain the same

class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        enableEdgeToEdge()
        setContent {
            val setupTime = remember { mutableStateOf("setuping...") }
            val proveAndVerifyTime = remember { mutableStateOf("proving...") }

            GnarkbenchmarkTheme {
                Scaffold(modifier = Modifier.fillMaxSize()) { innerPadding ->
                    val setupStartTime = System.nanoTime()
                    Ecdsa.setup()
                    val setupEndTime = System.nanoTime()
                    setupTime.value = "Setup: ${(setupEndTime - setupStartTime) / 1_000_000} ms"

                    val proveAndVerifyStartTime = System.nanoTime()
                    Ecdsa.proveAndVerify()
                    val proveAndVerifyEndTime = System.nanoTime()
                    proveAndVerifyTime.value = "Prove and Verify: ${(proveAndVerifyEndTime - proveAndVerifyStartTime) / 1_000_000} ms"

                    Column(modifier = Modifier.padding(innerPadding).fillMaxSize()) {
                        Text(text = setupTime.value)
                        Text(text = proveAndVerifyTime.value)
                    }
                }
            }
        }
    }
}