package com.freelanceplatform.android.ui.profile

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.navigation.NavController
import com.freelanceplatform.android.data.api.RetrofitClient
import com.freelanceplatform.android.data.local.TokenManager
import kotlinx.coroutines.launch

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun ProfileScreen(navController: NavController) {
    var user by remember { mutableStateOf<com.freelanceplatform.android.data.model.User?>(null) }
    val scope = rememberCoroutineScope()

    LaunchedEffect(Unit) {
        scope.launch {
            try {
                val res = RetrofitClient.apiService.getMe()
                if (res.code == 0) user = res.data?.user
            } catch (_: Exception) {}
        }
    }

    Scaffold(topBar = {
        TopAppBar(title = { Text("我的") }, colors = TopAppBarDefaults.topAppBarColors(containerColor = MaterialTheme.colorScheme.primary, titleContentColor = androidx.compose.ui.graphics.Color.White))
    }) { padding ->
        LazyColumn(modifier = Modifier.fillMaxSize().padding(padding).padding(16.dp), horizontalAlignment = Alignment.CenterHorizontally) {
            item {
                Surface(modifier = Modifier.size(80.dp), shape = androidx.compose.foundation.shape.CircleShape, color = MaterialTheme.colorScheme.primaryContainer) {}
                Spacer(Modifier.height(12.dp))
                Text(user?.nickname ?: "未登录", style = MaterialTheme.typography.headlineMedium, fontWeight = FontWeight.Bold)
                Text(when (user?.userType) { "developer" -> "开发者"; "client" -> "需求方"; "both" -> "双重身份"; else -> "" }, style = MaterialTheme.typography.bodyMedium, color = MaterialTheme.colorScheme.onSurfaceVariant)
                Spacer(Modifier.height(24.dp))
            }
            item {
                val menuItems = listOf("编辑资料" to "edit-profile", "我的项目" to "projects", "我的钱包" to "wallet", "设置" to "")
                menuItems.forEach { (label, route) ->
                    TextButton(onClick = { if (route.isNotBlank()) navController.navigate(route) }, modifier = Modifier.fillMaxWidth()) {
                        Text(label, modifier = Modifier.weight(1f), style = MaterialTheme.typography.bodyLarge)
                        Text("›", style = MaterialTheme.typography.bodyLarge, color = MaterialTheme.colorScheme.onSurfaceVariant)
                    }
                    HorizontalDivider()
                }
            }
            item {
                Spacer(Modifier.height(24.dp))
                OutlinedButton(onClick = {
                    TokenManager(navController.context).clearTokens()
                    navController.navigate("login") { popUpTo(0) }
                }, modifier = Modifier.fillMaxWidth(), colors = ButtonDefaults.outlinedButtonColors(contentColor = MaterialTheme.colorScheme.error)) { Text("退出登录") }
            }
        }
    }
}

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun EditProfileScreen(navController: NavController) {
    var nickname by remember { mutableStateOf("") }
    var title by remember { mutableStateOf("") }
    var bio by remember { mutableStateOf("") }

    Scaffold(topBar = {
        TopAppBar(title = { Text("编辑资料") }, colors = TopAppBarDefaults.topAppBarColors(containerColor = MaterialTheme.colorScheme.primary, titleContentColor = androidx.compose.ui.graphics.Color.White))
    }) { padding ->
        Column(modifier = Modifier.fillMaxSize().padding(padding).padding(16.dp)) {
            OutlinedTextField(value = nickname, onValueChange = { nickname = it }, label = { Text("昵称") }, modifier = Modifier.fillMaxWidth())
            Spacer(Modifier.height(12.dp))
            OutlinedTextField(value = title, onValueChange = { title = it }, label = { Text("职位头衔") }, modifier = Modifier.fillMaxWidth())
            Spacer(Modifier.height(12.dp))
            OutlinedTextField(value = bio, onValueChange = { bio = it }, label = { Text("简介") }, modifier = Modifier.fillMaxWidth().height(120.dp), maxLines = 5)
            Spacer(Modifier.height(24.dp))
            Button(onClick = { navController.popBackStack() }, modifier = Modifier.fillMaxWidth()) { Text("保存") }
        }
    }
}
