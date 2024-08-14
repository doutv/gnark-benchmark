package com.example.gnark_benchmark

import android.os.Bundle
import android.util.Log
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Button
import androidx.compose.material3.DropdownMenuItem
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.ExposedDropdownMenuBox
import androidx.compose.material3.ExposedDropdownMenuDefaults
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Text
import androidx.compose.material3.TextField
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.rememberCoroutineScope
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import com.example.gnark_benchmark.ui.theme.GnarkbenchmarkTheme
import ecdsa.Ecdsa
import eddsa.Eddsa
import kotlinx.coroutines.launch

// Other imports remain the same

class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContent {
            GnarkbenchmarkTheme {
                Scaffold(modifier = Modifier.fillMaxSize()) { paddingValues ->
                    BenchmarkComponent(
                        applicationContext.filesDir.toString(),
                        modifier = Modifier.padding(paddingValues)
                    )
                }
            }
        }
    }
}

@Composable
fun BenchmarkComponent(fileDir: String, modifier: Modifier = Modifier) {
    val setupTime = remember { mutableStateOf("Not started") }
    val proveAndVerifyTime = remember { mutableStateOf("Not started") }
    val selectedAlgorithm = remember { mutableStateOf("EdDSA") }
    val selectedSystem = remember { mutableStateOf("Groth16") }
    val coroutineScope = rememberCoroutineScope()
    val algorithms = listOf("EdDSA", "ECDSA")
    val systems = listOf("Groth16", "Plonk")

    Column(
        modifier = modifier
            .padding(top = 200.dp)
            .fillMaxSize(), horizontalAlignment = Alignment.CenterHorizontally
    ) {
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
                proveAndVerifyTime.value =
                    "Prove and Verify: ${(proveAndVerifyEndTime - proveAndVerifyStartTime) / 1_000_000} ms"
            }
        }) {
            Text("Prove and Verify")
        }
        Spacer(modifier = Modifier.height(16.dp))
        Text(text = proveAndVerifyTime.value)
    }
}

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun DropdownMenuComponent(
    label: String,
    selectedItem: String,
    items: List<String>,
    onItemSelected: (String) -> Unit
) {
    var expanded: Boolean by remember { mutableStateOf(false) }
    ExposedDropdownMenuBox(
        expanded = expanded,
        onExpandedChange = { expanded = it }
    ) {
        TextField(
            // The `menuAnchor` modifier must be passed to the text field to handle
            // expanding/collapsing the menu on click. A read-only text field has
            // the anchor type `PrimaryNotEditable`.
            modifier = Modifier.menuAnchor(),
            value = selectedItem,
            onValueChange = {},
            readOnly = true,
            singleLine = true,
            label = { Text(text = label) },
            trailingIcon = { ExposedDropdownMenuDefaults.TrailingIcon(expanded = expanded) },
            colors = ExposedDropdownMenuDefaults.textFieldColors()
        )
        ExposedDropdownMenu(
            expanded = expanded,
            onDismissRequest = { expanded = false }
        ) {
            items.forEach { item ->
                DropdownMenuItem(
                    text = { Text(item, style = MaterialTheme.typography.bodyLarge) },
                    onClick = {
                        onItemSelected(item)
                        expanded = false
                    },
                    contentPadding = ExposedDropdownMenuDefaults.ItemContentPadding
                )
            }
        }
    }
}
