package com.example.gnark_benchmark

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
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

class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        enableEdgeToEdge()
        setContent {
            GnarkbenchmarkTheme {
                Scaffold(modifier = Modifier.fillMaxSize()) { innerPadding ->
                    exec("gnark", "")
                    // val log = runBinaryAndCaptureLog()
                    // Greeting(
                    //     name = log,
                    //     modifier = Modifier.padding(innerPadding)
                    // )
                }
            }
        }
    }

    private fun runBinaryAndCaptureLog(): String {
        try {
            val command = "/data/local/tmp/gnark"
            val process = Runtime.getRuntime().exec(command)
            val reader = BufferedReader(InputStreamReader(process.inputStream))
            val output = StringBuilder()
            var line: String?
            while (reader.readLine().also { line = it } != null) {
                output.append(line).append("\n")
            }
            process.waitFor()
            return output.toString()
        } catch (e: Exception) {
            e.printStackTrace()
            return "Error running binary: ${e.message}"
        }
    }

    private fun exec(command: String, params: String): String {
        try {
            val process = ProcessBuilder()
                .directory(File(filesDir.parentFile!!, "lib"))
                .command(command, params)
                .redirectErrorStream(true)
                .start()
            val reader = BufferedReader(
                InputStreamReader(process.inputStream)
            )
            val text = reader.readText()
            reader.close()
            process.waitFor()
            return text
        } catch (e: Exception) {
            return e.message ?: "IOException"
        }
    }
}

@Composable
fun Greeting(name: String, modifier: Modifier = Modifier) {
    Text(
        text = "Log: $name",
        modifier = modifier
    )
}

@Preview(showBackground = true)
@Composable
fun GreetingPreview() {
    GnarkbenchmarkTheme {
        Greeting("Android")
    }
}