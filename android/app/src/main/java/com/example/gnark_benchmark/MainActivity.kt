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
import dummy1200k.Dummy1200k
import ecdsa.Ecdsa
import eddsa.Eddsa
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.launch
import kotlinx.coroutines.withContext
import kotlinx.serialization.Serializable
import kotlinx.serialization.encodeToString
import kotlinx.serialization.json.Json

// Other imports remain the same

class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        val fileDir =
            getExternalFilesDir(null)?.absolutePath?.toString().orEmpty().removeSuffix("/") + "/"
        setContent {
            GnarkbenchmarkTheme {
                Scaffold(modifier = Modifier.fillMaxSize()) { paddingValues ->
                    BenchmarkComponent(
                        fileDir,
                        modifier = Modifier.padding(paddingValues)
                    )
                }
            }
        }
    }
}

interface SelectableItem {
    val text: String
}

enum class Algorithm(override val text: String) : SelectableItem {
    Dummy1200k("Dummy 1200k"), ECDSA("ECDSA"), EdDSA("EdDSA")
}

enum class System(override val text: String) : SelectableItem {
    Groth16("Groth16"), Plonk("Plonk")
}

@Composable
fun BenchmarkComponent(fileDir: String, modifier: Modifier = Modifier) {
    val setupTime = remember { mutableStateOf("Not started") }
    val proveAndVerifyTime = remember { mutableStateOf("Not started") }
    val selectedAlgorithm = remember { mutableStateOf(Algorithm.entries.first()) }
    val selectedSystem = remember { mutableStateOf(System.entries.first()) }
    val coroutineScope = rememberCoroutineScope()
    val algorithms = Algorithm.entries
    val systems = System.entries

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
                val setupStartTime = java.lang.System.nanoTime()
                when (selectedAlgorithm.value to selectedSystem.value) {
                    Algorithm.Dummy1200k to System.Groth16 ->
                        withContext(Dispatchers.Default) { Dummy1200k.groth16Setup(fileDir) }

                    Algorithm.ECDSA to System.Groth16 ->
                        withContext(Dispatchers.Default) { Ecdsa.groth16Setup(fileDir) }

                    Algorithm.ECDSA to System.Plonk ->
                        withContext(Dispatchers.Default) { Ecdsa.plonkSetup(fileDir) }

                    Algorithm.EdDSA to System.Groth16 ->
                        withContext(Dispatchers.Default) { Eddsa.groth16Setup(fileDir) }

                    Algorithm.EdDSA to System.Groth16 ->
                        withContext(Dispatchers.Default) { Eddsa.plonkSetup(fileDir) }

                    else -> Log.e("BenchmarkComponent", "Invalid selection")
                }
                val setupEndTime = java.lang.System.nanoTime()
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
                val proveAndVerifyStartTime = java.lang.System.nanoTime()
                when (selectedAlgorithm.value to selectedSystem.value) {
                    Algorithm.Dummy1200k to System.Groth16 ->
                        withContext(Dispatchers.Default) { Dummy1200k.groth16Prove(fileDir) }

                    Algorithm.ECDSA to System.Groth16 -> withContext(Dispatchers.Default) {
                        val credential = KycCredential(
                            credential = 12UL,
                            age = 18UL,
                            gender = 1UL,
                            nation = 0b10UL,
                            expireTime = 123UL,
                        )
                        val credentialJson = Json.encodeToString(credential).encodeToByteArray()
                        Ecdsa.groth16Prove(fileDir, credentialJson)
                    }

                    Algorithm.ECDSA to System.Plonk -> withContext(Dispatchers.Default) {
                        val credential = KycCredential(
                            credential = 12UL,
                            age = 18UL,
                            gender = 1UL,
                            nation = 0b10UL,
                            expireTime = 123UL,
                        )
                        val credentialJson = Json.encodeToString(credential).encodeToByteArray()
                        Ecdsa.plonkProve(fileDir, credentialJson)
                    }

                    Algorithm.EdDSA to System.Groth16 -> withContext(Dispatchers.Default) {
                        val attributes = Attributes(intArrayOf(1, 2, 3))
                        val attributesJson = Json.encodeToString(attributes).encodeToByteArray()
                        Eddsa.groth16Prove(fileDir, attributesJson)
                    }

                    Algorithm.EdDSA to System.Groth16 -> withContext(Dispatchers.Default) {
                        val attributes = Attributes(intArrayOf(1, 2, 3))
                        val attributesJson = Json.encodeToString(attributes).encodeToByteArray()
                        Eddsa.plonkProve(fileDir, attributesJson)
                    }

                    else -> Log.e("BenchmarkComponent", "Invalid selection")
                }
                val proveAndVerifyEndTime = java.lang.System.nanoTime()
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
fun <T : SelectableItem> DropdownMenuComponent(
    label: String,
    selectedItem: T,
    items: List<T>,
    onItemSelected: (T) -> Unit
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
            value = selectedItem.text,
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
                    text = { Text(item.text, style = MaterialTheme.typography.bodyLarge) },
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

@Serializable
data class Attributes(val attributes: IntArray) {
    override fun equals(other: Any?): Boolean {
        if (this === other) return true
        if (javaClass != other?.javaClass) return false

        other as Attributes

        return attributes.contentEquals(other.attributes)
    }

    override fun hashCode(): Int {
        return attributes.contentHashCode()
    }
}

@Serializable
data class KycCredential(
    val credential: ULong,
    val age: ULong,
    val gender: ULong,
    val nation: ULong,
    val expireTime: ULong,
)
