package com.freelanceplatform.android.ui.chat

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextOverflow
import androidx.compose.ui.unit.dp
import androidx.navigation.NavController
import com.freelanceplatform.android.data.model.Conversation
import com.freelanceplatform.android.data.model.ChatMessage
import com.freelanceplatform.android.data.api.RetrofitClient
import kotlinx.coroutines.launch

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun ChatListScreen(navController: NavController) {
    var conversations by remember { mutableStateOf<List<Conversation>>(emptyList()) }
    var isLoading by remember { mutableStateOf(true) }
    val scope = rememberCoroutineScope()

    LaunchedEffect(Unit) {
        scope.launch {
            try {
                val res = RetrofitClient.apiService.listConversations(page = 1, pageSize = 50)
                if (res.code == 0) conversations = res.data ?: emptyList()
            } catch (_: Exception) {}
            isLoading = false
        }
    }

    Scaffold(topBar = {
        TopAppBar(title = { Text("消息") }, colors = TopAppBarDefaults.topAppBarColors(containerColor = MaterialTheme.colorScheme.primary, titleContentColor = androidx.compose.ui.graphics.Color.White))
    }) { padding ->
        if (isLoading) {
            Box(Modifier.fillMaxSize().padding(padding), contentAlignment = Alignment.Center) { CircularProgressIndicator() }
        } else {
            LazyColumn(modifier = Modifier.fillMaxSize().padding(padding)) {
                items(conversations) { conv ->
                    Card(modifier = Modifier.fillMaxWidth().padding(horizontal = 16.dp, vertical = 4.dp), shape = RoundedCornerShape(12.dp), onClick = { navController.navigate("chat/${conv.id}") }) {
                        Row(modifier = Modifier.padding(16.dp), verticalAlignment = Alignment.CenterVertically) {
                            Column(modifier = Modifier.weight(1f)) {
                                Text("会话", style = MaterialTheme.typography.titleMedium, fontWeight = FontWeight.SemiBold)
                                conv.lastMessage?.content?.let { Text(it, style = MaterialTheme.typography.bodyMedium, maxLines = 1, overflow = TextOverflow.Ellipsis, color = MaterialTheme.colorScheme.onSurfaceVariant) }
                            }
                        }
                    }
                }
            }
        }
    }
}

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun ChatDetailScreen(navController: NavController, conversationId: String) {
    var messages by remember { mutableStateOf<List<ChatMessage>>(emptyList()) }
    var inputText by remember { mutableStateOf("") }
    val scope = rememberCoroutineScope()

    LaunchedEffect(conversationId) {
        scope.launch {
            try {
                val res = RetrofitClient.apiService.getMessages(conversationId, page = 1, pageSize = 50)
                if (res.code == 0) messages = res.data ?: emptyList()
            } catch (_: Exception) {}
        }
    }

    Scaffold(topBar = {
        TopAppBar(title = { Text("聊天") }, colors = TopAppBarDefaults.topAppBarColors(containerColor = MaterialTheme.colorScheme.primary, titleContentColor = androidx.compose.ui.graphics.Color.White))
    }) { padding ->
        Column(modifier = Modifier.fillMaxSize().padding(padding)) {
            LazyColumn(modifier = Modifier.weight(1f).padding(horizontal = 16.dp)) {
                items(messages) { msg ->
                    val isSelf = msg.senderId == "self"
                    Row(modifier = Modifier.fillMaxWidth().padding(vertical = 4.dp), horizontalArrangement = if (isSelf) Arrangement.End else Arrangement.Start) {
                        Surface(color = if (isSelf) MaterialTheme.colorScheme.primary else MaterialTheme.colorScheme.surfaceVariant, shape = RoundedCornerShape(12.dp)) {
                            Text(msg.content ?: "", modifier = Modifier.padding(12.dp), color = if (isSelf) androidx.compose.ui.graphics.Color.White else MaterialTheme.colorScheme.onSurface, style = MaterialTheme.typography.bodyLarge)
                        }
                    }
                }
            }
            Row(modifier = Modifier.fillMaxWidth().padding(16.dp), verticalAlignment = Alignment.CenterVertically) {
                OutlinedTextField(value = inputText, onValueChange = { inputText = it }, modifier = Modifier.weight(1f), placeholder = { Text("输入消息") }, shape = RoundedCornerShape(24.dp))
                Spacer(Modifier.width(8.dp))
                Button(onClick = {
                    if (inputText.isNotBlank()) {
                        val newMsg = ChatMessage(id = System.currentTimeMillis().toString(), conversationId = conversationId, senderId = "self", content = inputText, messageType = "text", createdAt = "")
                        messages = messages + newMsg
                        scope.launch { try { RetrofitClient.apiService.sendMessage(conversationId, content = inputText, messageType = "text") } catch (_: Exception) {} }
                        inputText = ""
                    }
                }) { Text("发送") }
            }
        }
    }
}
